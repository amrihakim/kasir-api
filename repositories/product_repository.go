package repositories

import (
	"database/sql"
	"errors"
	"kasir-api/models"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetAll(name string) ([]models.Product, error) {
	query := "SELECT id, name, price, stock, category_id FROM products"

	var args []interface{}
	if name != "" {
		query += " WHERE name ILIKE $1"
		args = append(args, "%"+name+"%")
	}

	rows, err := r.db.Query(query, args...)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return products, nil
}

func (r *ProductRepository) Create(newProduct *models.Product) (models.Product, error) {
	var createdProduct models.Product
	err := r.db.QueryRow(
		"INSERT INTO products (name, price, stock, category_id) VALUES ($1, $2, $3, $4) RETURNING id, name, price, stock, category_id",
		newProduct.Name, newProduct.Price, newProduct.Stock, newProduct.CategoryID,
	).Scan(&createdProduct.ID, &createdProduct.Name, &createdProduct.Price, &createdProduct.Stock, &createdProduct.CategoryID)
	if err != nil {
		return models.Product{}, err
	}
	return createdProduct, nil
}

func (r *ProductRepository) GetByID(id int) (*models.ProductWithCategory, error) {
	var p models.ProductWithCategory
	err := r.db.QueryRow(
		"SELECT p.id, p.name, p.price, p.stock, p.category_id, c.id, c.name, c.description FROM products p JOIN categories c ON p.category_id = c.id WHERE p.id = $1", id,
	).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &p.CategoryID, &p.Category.ID, &p.Category.Name, &p.Category.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Product Not Found!")
		}
		return nil, err
	}
	return &p, nil
}

func (r *ProductRepository) Update(product *models.Product) (*models.ProductWithCategory, error) {
	query := "UPDATE products SET name = $1, price = $2, stock = $3 WHERE id = $4"
	result, err := r.db.Exec(query, product.Name, product.Price, product.Stock, product.ID)
	if err != nil {
		return nil, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rows == 0 {
		return nil, errors.New("Product Not Found!")
	}

	return r.GetByID(product.ID)
}

func (r *ProductRepository) Delete(id int) error {
	result, err := r.db.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("Product Not Found!")
	}

	return nil
}
