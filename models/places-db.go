package models

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type DBModel struct {
	DB *sql.DB
}

func (m *DBModel) GetPlaceByID(id int) (*Place, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, name, description, is_halal, is_vegetarian, location, lat, lon, created_at, updated_at, category from place where id = ?`

	row := m.DB.QueryRowContext(ctx, query, id)
	var place Place
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

	// query = `Select
	// 			pc.id, pc.place_id, pc.category_id, c.category_name
	// 		From
	// 			place_category pc
	// 			left join category c on (c.id = pc.category_id)
	// 		Where
	// 			pc.place_id =?
	// 		`
	// rows, _ := m.DB.QueryContext(ctx, query, id)

	// defer rows.Close()

	// category := make(map[int]string)
	// for rows.Next() {
	// 	var pc PlaceCategory
	// 	err := rows.Scan(
	// 		&pc.ID,
	// 		&pc.MovieID,
	// 		&pc.CategoryID,
	// 		&pc.Category,
	// 	)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	category[pc.ID] = pc.Category.CategoryName
	// }

	// place.Category = category
	return &place, nil
}

func (m *DBModel) GetAllPlaces(category ...string) ([]*Place, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	where := ""
	if len(category) > 0 {
		println(category)
		where = fmt.Sprintf("where category =  %s)", category)
	}

	query := fmt.Sprintf(`select id, name, description, is_halal, is_vegetarian, location, lat, lon, created_at, updated_at, category from place %s order by name`, where)
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var places []*Place
	for rows.Next() {
		var place Place
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

func (m *DBModel) GetAllPlacesWithFilter(isHalal, isVegetarian bool) ([]*Place, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

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
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var places []*Place
	for rows.Next() {
		var place Place
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

func (m *DBModel) InsertPlace(place Place) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		insert into place 
		(name, description, is_halal, is_vegetarian, location, lat, lon, created_at, updated_at, category) 
		values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := m.DB.ExecContext(ctx, stmt,
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

func (m *DBModel) UpdatePlace(place Place) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `Update place set name = ?, description = ?, is_halal = ?, is_vegetarian = ?, location = ?, lat = ?, lon = ?, created_at = ? , updated_at = ? , category = ? where id = ?`

	_, err := m.DB.ExecContext(ctx, stmt,
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

func (m *DBModel) DeletePlace(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := "Delete from place where id = ?"

	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *DBModel) DeletePlaces(idList []int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

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

	_, err := m.DB.ExecContext(ctx, sb.String())
	if err != nil {
		return err
	}

	return nil
}
