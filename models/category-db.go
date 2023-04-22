package models

import (
	"context"
	"fmt"
	"time"
)

func (m *DBModel) GetCategoryByID(id int) (*Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, category_name, created_at, updated_at from category where id = ?`

	row := m.DB.QueryRowContext(ctx, query, id)
	var category Category
	err := row.Scan(
		&category.ID,
		&category.CategoryName,
		&category.CreatedAt,
		&category.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (m *DBModel) GetAllCategories() ([]*Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := fmt.Sprintf(`select id, category_name, created_at, updated_at from category`)
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var categories []*Category
	for rows.Next() {
		var category Category
		err := rows.Scan(
			&category.ID,
			&category.CategoryName,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		categories = append(categories, &category)
	}
	return categories, nil

}
