package repositories

import (
	"database/sql"
	"errors"
	"kasir-api/models"
)

type CategoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) GetAll() ([]models.Category, error) {
	rows, err := r.db.Query("SELECT id, name, description FROM categories ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}

func (r *CategoryRepository) Create(newCategory *models.Category) (models.Category, error) {
	var createdCategory models.Category
	err := r.db.QueryRow(
		"INSERT INTO categories (name, description) VALUES ($1, $2) RETURNING id, name, description",
		newCategory.Name, newCategory.Description,
	).Scan(&createdCategory.ID, &createdCategory.Name, &createdCategory.Description)
	if err != nil {
		return models.Category{}, err
	}
	return createdCategory, nil
}

func (r *CategoryRepository) GetByID(id int) (*models.Category, error) {
	var c models.Category
	err := r.db.QueryRow(
		"SELECT id, name, description FROM categories WHERE id = $1", id,
	).Scan(&c.ID, &c.Name, &c.Description)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Category Not Found!")
		}
		return nil, err
	}
	return &c, nil
}

func (r *CategoryRepository) Update(category *models.Category) (*models.Category, error) {
	query := "UPDATE categories SET name = $1, description = $2 WHERE id = $3"
	result, err := r.db.Exec(query, category.Name, category.Description, category.ID)
	if err != nil {
		return nil, err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}

	if rows == 0 {
		return nil, errors.New("Category Not Found!")
	}

	return r.GetByID(category.ID)
}

func (r *CategoryRepository) Delete(id int) error {
	result, err := r.db.Exec("DELETE FROM categories WHERE id = $1", id)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("Category Not Found!")
	}

	return nil
}
