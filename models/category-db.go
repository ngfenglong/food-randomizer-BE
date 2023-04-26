package models

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
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

func (m *DBModel) InsertCategory(category Category) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		insert into category (category_name) value (?)
	`

	_, err := m.DB.ExecContext(ctx, stmt, category.CategoryName)
	if err != nil {
		return err
	}

	return nil
}

func (m *DBModel) UpdateCategory(category Category) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		Update category set category_name = ? where id = ?
	`

	_, err := m.DB.ExecContext(ctx, stmt, category.CategoryName, category.ID)
	if err != nil {
		return err
	}

	return nil
}

func (m *DBModel) DeleteCategory(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := "Delete from category where id = ?"

	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *DBModel) DeleteCategories(idList []int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if len(idList) == 0 {
		return errors.New("the list is empty")
	}

	var sb strings.Builder
	sb.WriteString("Delete from category where id in (")
	sb.WriteString(strconv.Itoa(idList[0]))
	for i := 1; i < len(idList); i++ {
		sb.WriteString(",")
		sb.WriteString(strconv.Itoa(idList[i]))
	}
	sb.WriteString(")")

	_, err := m.DB.ExecContext(ctx, sb.String())
	if err != nil {
		return err
	}

	return nil
}
