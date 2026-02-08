package repositories

import (
	"database/sql"
	"fmt"
	"kasir-api/models"
	"strings"
	"time"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	for _, item := range items {
		var productPrice, stock int
		var productName string

		err := tx.QueryRow("SELECT name, price, stock FROM products WHERE id = $1", item.ProductID).Scan(&productName, &productPrice, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}

		subtotal := productPrice * item.Quantity
		totalAmount += subtotal

		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	var transactionID int
	var createdAt sql.NullTime
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id, created_at", totalAmount).Scan(&transactionID, &createdAt)
	if err != nil {
		return nil, err
	}

	insertQuery := "INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES "
	params := []interface{}{}
	valueStrings := []string{}
	paramIndex := 1

	for _, d := range details {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d)", paramIndex, paramIndex+1, paramIndex+2, paramIndex+3))
		params = append(params, transactionID, d.ProductID, d.Quantity, d.Subtotal)
		paramIndex += 4
	}
	insertQuery += strings.Join(valueStrings, ",")
	insertQuery += " RETURNING id"

	rows, err := tx.Query(insertQuery, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	i := 0
	for rows.Next() {
		var detailID int
		if err := rows.Scan(&detailID); err != nil {
			return nil, err
		}
		details[i].ID = detailID
		details[i].TransactionID = transactionID
		i++
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		CreatedAt:   createdAt.Time,
		Details:     details,
	}, nil
}

func (repo *TransactionRepository) GetReport(startDate string, endDate string) (*models.ReportResponse, error) {
	var totalRevenue, totalTransactions int
	var bestSellingProductName sql.NullString
	var bestSellingProductQty sql.NullInt64

	err := repo.db.QueryRow(`SELECT COALESCE(SUM(t.total_amount), 0) AS total_revenue,
		COUNT(t.id) AS total_transactions
		FROM transactions t
		WHERE DATE(t.created_at) BETWEEN $1 AND $2`, startDate, endDate).Scan(&totalRevenue, &totalTransactions)

	if err != nil {
		return nil, err
	}

	err = repo.db.QueryRow(`SELECT p.name, SUM(td.quantity) AS total_quantity
		FROM transaction_details td
		JOIN products p ON td.product_id = p.id
		JOIN transactions t ON td.transaction_id = t.id
		WHERE DATE(t.created_at) BETWEEN $1 AND $2
		GROUP BY p.name
		ORDER BY total_quantity DESC
		LIMIT 1`, startDate, endDate).Scan(&bestSellingProductName, &bestSellingProductQty)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	report := &models.ReportResponse{
		TotalRevenue:   totalRevenue,
		TotalTransaksi: totalTransactions,
	}
	if bestSellingProductName.Valid && bestSellingProductQty.Valid {
		report.ProdukTerlaris.Nama = bestSellingProductName.String
		report.ProdukTerlaris.QtyTerjual = int(bestSellingProductQty.Int64)
	}

	return report, nil
}

func (repo *TransactionRepository) GetReportToday() (*models.ReportResponse, error) {
	today := time.Now().Format("2006-01-02")
	return repo.GetReport(today, today)
}
