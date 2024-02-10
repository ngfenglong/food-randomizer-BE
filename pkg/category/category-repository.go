package category

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/ngfenglong/food-randomizer-BE/pkg/models"
)

var _ CategoryRepository = &SQLCategoryRepository{}

type CategoryRepository interface {
	GetCategoryByID(ctx context.Context, id int) (*models.Category, error)
	GetAllCategories(ctx context.Context) ([]*models.Category, error)
	InsertCategory(ctx context.Context, category models.Category) error
	UpdateCategory(ctx context.Context, category models.Category) error
	DeleteCategory(ctx context.Context, id int) error
	DeleteCategories(ctx context.Context, idList []int) error
}

type SQLCategoryRepository struct {
	db *sql.DB
}

func NewSQLCategoryRepostory(db *sql.DB) *SQLCategoryRepository {
	return &SQLCategoryRepository{db: db}
}

func (repo *SQLCategoryRepository) GetCategoryByID(ctx context.Context, id int) (*models.Category, error) {
	query := `select id, category_name, created_at, updated_at from category where id = ?`

	row := repo.db.QueryRowContext(ctx, query, id)
	var category models.Category
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

func (repo *SQLCategoryRepository) GetAllCategories(ctx context.Context) ([]*models.Category, error) {
	query := fmt.Sprintf(`select id, category_name, created_at, updated_at from category`)
	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var categories []*models.Category
	for rows.Next() {
		var category models.Category
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

func (repo *SQLCategoryRepository) InsertCategory(ctx context.Context, category models.Category) error {
	stmt := `
		insert into category (category_name) value (?)
	`

	_, err := repo.db.ExecContext(ctx, stmt, category.CategoryName)
	if err != nil {
		return err
	}

	return nil
}

func (repo *SQLCategoryRepository) UpdateCategory(ctx context.Context, category models.Category) error {
	stmt := `
		Update category set category_name = ? where id = ?
	`

	_, err := repo.db.ExecContext(ctx, stmt, category.CategoryName, category.ID)
	if err != nil {
		return err
	}

	return nil
}

func (repo *SQLCategoryRepository) DeleteCategory(ctx context.Context, id int) error {
	stmt := "Delete from category where id = ?"

	_, err := repo.db.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

func (repo *SQLCategoryRepository) DeleteCategories(ctx context.Context, idList []int) error {
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

	_, err := repo.db.ExecContext(ctx, sb.String())
	if err != nil {
		return err
	}

	return nil
}
