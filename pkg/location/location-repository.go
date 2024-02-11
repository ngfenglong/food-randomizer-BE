package location

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ngfenglong/food-randomizer-BE/pkg/models"
)

type LocationRepository interface {
	GetLocationByID(ctx context.Context, id int) (*models.Location, error)
	GetAllLocations(ctx context.Context) ([]*models.Location, error)
	InsertLocation(ctx context.Context, location models.Location) error
	UpdateLocation(ctx context.Context, location models.Location) error
	DeleteLocation(ctx context.Context, id int) error
	DeleteLocations(ctx context.Context, idList []int) error
}

type SQLLocationRepository struct {
	db *sql.DB
}

var _ LocationRepository = &SQLLocationRepository{}

// Constructor function for SQLLocationRepository
func NewSQLLocationRepository(db *sql.DB) *SQLLocationRepository {
	return &SQLLocationRepository{db: db}
}

func (r *SQLLocationRepository) GetLocationByID(ctx context.Context, id int) (*models.Location, error) {
	query := `Select id, location_name, created_at, updated_at from location where id = ?`

	row := r.db.QueryRowContext(ctx, query, id)
	var location models.Location
	err := row.Scan(
		&location.ID,
		&location.LocationName,
		&location.CreatedAt,
		&location.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &location, nil
}

func (r *SQLLocationRepository) GetAllLocations(ctx context.Context) ([]*models.Location, error) {
	query := fmt.Sprint(`select id, location_name, created_at, updated_at from location`)
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var locations []*models.Location
	for rows.Next() {
		var location models.Location
		err := rows.Scan(
			&location.ID,
			&location.LocationName,
			&location.CreatedAt,
			&location.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		locations = append(locations, &location)
	}
	return locations, nil
}

func (r *SQLLocationRepository) InsertLocation(ctx context.Context, location models.Location) error {

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		insert into location (location_name) value (?)
	`

	_, err := r.db.ExecContext(ctx, stmt, location.LocationName)
	if err != nil {
		return err
	}
	return nil
}

func (r *SQLLocationRepository) UpdateLocation(ctx context.Context, location models.Location) error {
	stmt := `
		Update location set location_name = ?, updated_at = ? where id = ?
	`

	_, err := r.db.ExecContext(ctx, stmt, location.LocationName, location.UpdatedAt, location.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *SQLLocationRepository) DeleteLocation(ctx context.Context, id int) error {
	stmt := `Delete from location where id = ?`
	_, err := r.db.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *SQLLocationRepository) DeleteLocations(ctx context.Context, idList []int) error {
	if len(idList) == 0 {
		return errors.New("the list is empty")
	}

	var sb strings.Builder
	sb.WriteString("Delete from location where id in (")
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
