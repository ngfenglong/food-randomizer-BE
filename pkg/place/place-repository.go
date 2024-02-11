package place

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/ngfenglong/food-randomizer-BE/pkg/models"
)

var _ PlaceRepository = &SQLPlaceRepository{}

type PlaceRepository interface {
	// Add more methods as needed
	GetPlaceByID(ctx context.Context, id int) (*models.Place, error)
	GetAllPlaces(ctx context.Context, category ...string) ([]*models.Place, error)
	GetAllPlacesWithFilter(ctx context.Context, isHalal, isVegetarian bool) ([]*models.Place, error)
	InsertPlace(ctx context.Context, place models.Place) error
	UpdatePlace(ctx context.Context, place models.Place) error
	DeletePlace(ctx context.Context, id int) error
	DeletePlaces(ctx context.Context, idList []int) error
}

type SQLPlaceRepository struct {
	db *sql.DB
}

func NewSQLPlaceRepository(db *sql.DB) *SQLPlaceRepository {
	return &SQLPlaceRepository{db: db}
}

func (r *SQLPlaceRepository) GetPlaceByID(ctx context.Context, id int) (*models.Place, error) {
	query := `select id, name, description, is_halal, is_vegetarian, location, lat, lon, created_at, updated_at, category from place where id = ?`

	row := r.db.QueryRowContext(ctx, query, id)
	var place models.Place
	err := row.Scan(
		&place.ID,
		&place.Name,
		&place.Description,
		&place.IsHalal,
		&place.IsVegetarian,
		&place.Location,
		&place.Lat,
		&place.Lon,
		&place.CreatedAt,
		&place.UpdatedAt,
		&place.Category,
	)
	if err != nil {
		return nil, err
	}

	return &place, nil
}

func (r *SQLPlaceRepository) GetAllPlaces(ctx context.Context, category ...string) ([]*models.Place, error) {
	where := ""
	if len(category) > 0 {
		println(category)
		where = fmt.Sprintf("where category =  %s)", category)
	}

	query := fmt.Sprintf(`select id, name, description, is_halal, is_vegetarian, location, lat, lon, created_at, updated_at, category from place %s order by name`, where)
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var places []*models.Place
	for rows.Next() {
		var place models.Place
		err := rows.Scan(
			&place.ID,
			&place.Name,
			&place.Description,
			&place.IsHalal,
			&place.IsVegetarian,
			&place.Location,
			&place.Lat,
			&place.Lon,
			&place.CreatedAt,
			&place.UpdatedAt,
			&place.Category,
		)
		if err != nil {
			return nil, err
		}

		places = append(places, &place)
	}
	return places, nil
}

func (r *SQLPlaceRepository) GetAllPlacesWithFilter(ctx context.Context, isHalal, isVegetarian bool) ([]*models.Place, error) {
	where := ""
	if isHalal {
		where = fmt.Sprintf("where is_halal = 1")
	}

	if isVegetarian {
		if len(where) > 0 {
			where += fmt.Sprintf(" and is_vegetarian = 1")
		} else {
			where = fmt.Sprintf("where is_vegetarian = 1")
		}
	}

	query := fmt.Sprintf(`select id, name, description, is_halal, is_vegetarian, location, lat, lon, created_at, updated_at, category from place %s order by name`, where)
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var places []*models.Place
	for rows.Next() {
		var place models.Place
		err := rows.Scan(
			&place.ID,
			&place.Name,
			&place.Description,
			&place.IsHalal,
			&place.IsVegetarian,
			&place.Location,
			&place.Lat,
			&place.Lon,
			&place.CreatedAt,
			&place.UpdatedAt,
			&place.Category,
		)
		if err != nil {
			return nil, err
		}

		places = append(places, &place)
	}
	return places, nil
}

func (r *SQLPlaceRepository) InsertPlace(ctx context.Context, place models.Place) error {
	stmt := `
		insert into place 
		(name, description, is_halal, is_vegetarian, location, lat, lon, created_at, updated_at, category) 
		values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := r.db.ExecContext(ctx, stmt,
		place.Name,
		place.Description,
		place.IsHalal,
		place.IsVegetarian,
		place.Location,
		place.Lat,
		place.Lon,
		place.CreatedAt,
		place.UpdatedAt,
		place.Category,
	)

	if err != nil {
		return err
	}
	return nil
}

func (r *SQLPlaceRepository) UpdatePlace(ctx context.Context, place models.Place) error {
	stmt := `Update place set name = ?, description = ?, is_halal = ?, is_vegetarian = ?, location = ?, lat = ?, lon = ?, created_at = ? , updated_at = ? , category = ? where id = ?`

	_, err := r.db.ExecContext(ctx, stmt,
		place.Name,
		place.Description,
		place.IsHalal,
		place.IsVegetarian,
		place.Location,
		place.Lat,
		place.Lon,
		place.CreatedAt,
		place.UpdatedAt,
		place.Category,
		place.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *SQLPlaceRepository) DeletePlace(ctx context.Context, id int) error {
	stmt := "Delete from place where id = ?"

	_, err := r.db.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *SQLPlaceRepository) DeletePlaces(ctx context.Context, idList []int) error {
	if len(idList) == 0 {
		return errors.New("The list is empty")
	}

	var sb strings.Builder
	sb.WriteString("Delete from place where id in (")
	sb.WriteString(strconv.Itoa(idList[0]))
	for i := 1; i < len(idList); i++ {
		sb.WriteString(",")
		sb.WriteString(strconv.Itoa(idList[i]))
	}
	sb.WriteString(")")

	_, err := r.db.ExecContext(ctx, sb.String())
	if err != nil {
		return err
	}

	return nil
}
