package models

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type DBModel struct {
	DB *sql.DB
}

func (m *DBModel) Get(id int) (*Place, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `select id, name, description, is_halal, is_vegetarian, location, lat, long, created_at, updated_at, category from place where id = $1`

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
		&place.Long,
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
	// 			pc.place_id = $1
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

func (m *DBModel) All(category ...string) ([]*Place, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	where := ""
	if len(category) > 0 {
		println(category)
		where = fmt.Sprintf("where category =  %s)", category)
	}

	query := fmt.Sprintf(`select id, name, description, is_halal, is_vegetarian, location, lat, long, created_at, updated_at, category from place %s order by name`, where)
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
			&place.Long,
			&place.CreatedAt,
			&place.UpdatedAt,
			&place.Category,
		)
		if err != nil {
			return nil, err
		}
		// 	categoryQuery := `Select
		// 				pc.id, pc.place_id, pc.category_id, c.category_name
		// 			From
		// 				place_category pc
		// 				left join category c on (c.id = pc.category_id)
		// 			Where
		// 				pc.place_id = $1
		// 			`

		// 	categoryRows, _ := m.DB.QueryContext(ctx, categoryQuery, place.ID)

		// 	category := make(map[int]string)
		// 	for categoryRows.Next() {
		// 		var pc PlaceCategory
		// 		err := categoryRows.Scan(
		// 			&pc.ID,
		// 			&pc.MovieID,
		// 			&pc.CategoryID,
		// 			&pc.Category.CategoryName,
		// 		)
		// 		if err != nil {
		// 			return nil, err
		// 		}
		// 		category[pc.ID] = pc.Category.CategoryName
		// 	}
		// 	categoryRows.Close()
		// 	place.Category = category
		places = append(places, &place)
	}
	return places, nil
}

func (m *DBModel) InsertPlace(place Place) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `
		insert into place 
		(name, description, is_halal, is_vegetarian, location, lat, long, created_at, updated_at, category) 
		values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := m.DB.ExecContext(ctx, stmt,
		place.Name,
		place.Description,
		place.IsHalal,
		place.IsVegetarian,
		place.Location,
		place.Lat,
		place.Long,
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

	stmt := `Update place set name = $1, description = $2, is_halal = $3, is_vegetarian = $4, location = $5, lat = $6, long = $7, created_at = $8 , updated_at = $9 , category = $10 where id = $11`

	_, err := m.DB.ExecContext(ctx, stmt,
		place.Name,
		place.Description,
		place.IsHalal,
		place.IsVegetarian,
		place.Location,
		place.Lat,
		place.Long,
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

	stmt := "Delete from place where id = $1"
	println(id)
	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {
		return err
	}

	return nil
}
