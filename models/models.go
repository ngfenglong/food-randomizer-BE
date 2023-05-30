package models

import (
	"database/sql"
	"time"
)

type Models struct {
	DB DBModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		DB: DBModel{
			DB: db,
		},
	}
}

type Place struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	Category     string    `json:"category"`
	IsHalal      bool      `json:"is_halal"`
	IsVegetarian bool      `json:"is_vegetarian"`
	Location     string    `json:"location"`
	Lat          string    `json:"lat"`
	Lon          string    `json:"lon"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Category struct {
	ID           int       `json:"id"`
	CategoryName string    `json:"category_name"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
}

type PlaceCategory struct {
	ID         int       `json:"-"`
	PlaceID    int       `json:"-"`
	CategoryID int       `json:"-"`
	Category   Category  `json:"category"`
	CreatedAt  time.Time `json:"-"`
	UpdatedAt  time.Time `json:"-"`
}

type Location struct {
	ID           int       `json:"id"`
	LocationName string    `json:"location_name"`
	StreetName   string    `json:"street_name"`
	Lat          string    `json:"lat"`
	Lon          string    `json:"lon"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
}

type PlaceLocation struct {
	ID         int       `json:"-"`
	PlaceID    int       `json:"-"`
	LocationID int       `json:"-"`
	Location   Location  `json:"location"`
	CreatedAt  time.Time `json:"-"`
	UpdatedAt  time.Time `json:"-"`
}

type User struct {
	ID        int       `json:"id"`
	UserName  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Role      int       `json:"role"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

type Token struct {
	ID     int
	UserID int
	Token  string
	Expiry time.Time
}
