package models

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func (m *DBModel) GetLocationByID(id int) (*Location, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `Select id, location_name, created_at, updated_at from location where id = ?`

	row := m.DB.QueryRowContext(ctx, query, id)
	var location Location
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

func (m *DBModel) GetAllLocations() ([]*Location, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := fmt.Sprint(`select id, location_name, created_at, updated_at from location`)
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var locations []*Location
	for rows.Next() {
		var location Location
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

func (m *DBModel) InsertLocation(location Location) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		insert into location (location_name) value (?)
	`

	_, err := m.DB.ExecContext(ctx, stmt, location.LocationName)
	if err != nil {
		return err
	}
	return nil
}

func (m *DBModel) UpdateLocation(location Location) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		Update location set location_name = ?, updated_at = ? where id = ?
	`

	_, err := m.DB.ExecContext(ctx, stmt, location.LocationName, location.UpdatedAt, location.ID)
	if err != nil {
		return err
	}

	return nil
}

func (m *DBModel) DeleteLocation(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `Delete from location where id = ?`
	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *DBModel) DeleteLocations(idList []int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

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

	_, err := m.DB.ExecContext(ctx, sb.String())
	if err != nil {
		return err
	}

	return nil
}
