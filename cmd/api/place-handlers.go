package main

import (
	"backend/models"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

// Dto
type PlaceDto struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Category     string `json:"category"`
	IsHalal      bool   `json:"is_halal"`
	IsVegetarian bool   `json:"is_vegetarian"`
	Location     string `json:"location"`
}

func (app *application) generatePlace(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	isHalalParam, isHalalExists := queryParams["is_halal"]
	isVegetarianParam, isVegetarianExists := queryParams["is_vegetarian"]

	var isHalal, isVegetarian bool
	var err error

	if isHalalExists && len(isHalalParam) > 0 {
		isHalal, err = strconv.ParseBool(isHalalParam[0])
		if err != nil {
			app.errorJSON(w, err)
			return
		}
	}

	if isVegetarianExists && len(isVegetarianParam) > 0 {
		isVegetarian, err = strconv.ParseBool(isVegetarianParam[0])
		if err != nil {
			app.errorJSON(w, err)
			return
		}
	}

	var places []*models.Place
	if isHalalExists || isVegetarianExists {
		places, err = app.models.DB.GetAllPlacesWithFilter(isHalal, isVegetarian)
	} else {
		places, err = app.models.DB.GetAllPlaces()
	}

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	idx := rand.Intn(len(places))
	err = app.WriteJSON(w, http.StatusOK, places[idx], "place")
	if err != nil {
		app.errorJSON(w, err)
	}
}

func (app *application) getAllPlaces(w http.ResponseWriter, r *http.Request) {
	places, err := app.models.DB.GetAllPlaces()

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.WriteJSON(w, http.StatusOK, places, "places")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

func (app *application) getPlaceByID(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.logger.Print(errors.New("invalid id parameter"))
		app.errorJSON(w, err)
		return
	}

	place, err := app.models.DB.GetPlaceByID(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}
	err = app.WriteJSON(w, http.StatusOK, place, "place")

	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

func (app *application) deletePlace(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	// Check if ID exists
	place, err := app.models.DB.GetPlaceByID(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if place == nil {
		app.errorJSON(w, errors.New("ID does not exists"))
	}

	err = app.models.DB.DeletePlace(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.WriteJSON(w, http.StatusOK, nil, "response")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

func (app *application) deletePlaces(w http.ResponseWriter, r *http.Request) {
	var idList []int

	err := json.NewDecoder(r.Body).Decode(&idList)
	if err != nil {
		app.errorJSON(w, err)
	}

	if len(idList) == 0 {
		app.errorJSON(w, errors.New("the ID list is empty"))
	}

	err = app.models.DB.DeletePlaces(idList)
	if err != nil {
		app.errorJSON(w, err)
	}

	err = app.WriteJSON(w, http.StatusOK, "Deleted Successfully", "response")
	if err != nil {
		app.errorJSON(w, err)
	}
}

func (app *application) editPlace(w http.ResponseWriter, r *http.Request) {
	var payload PlaceDto

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		log.Println("error", err)
		app.errorJSON(w, err)
		return
	}

	var place models.Place

	if payload.ID != 0 {
		m, err := app.models.DB.GetPlaceByID(payload.ID)

		if err != nil {
			if err == sql.ErrNoRows {
				payload.ID = 0
			} else {
				app.errorJSON(w, err)
				return
			}
		} else {
			place = *m
		}
	}

	place.ID = payload.ID
	place.Name = payload.Name
	place.Description = payload.Description
	place.Category = payload.Category
	place.IsHalal = payload.IsHalal
	place.IsVegetarian = payload.IsVegetarian
	place.Location = payload.Location
	place.Lat = " "
	place.Lon = " "
	place.CreatedAt = time.Now()
	place.UpdatedAt = time.Now()

	if place.ID == 0 {
		err = app.models.DB.InsertPlace(place)
		if err != nil {
			app.errorJSON(w, err)
			return
		}
	} else {
		err = app.models.DB.UpdatePlace(place)
		if err != nil {
			app.errorJSON(w, err)
			return
		}
	}

	err = app.WriteJSON(w, http.StatusOK, nil, "response")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}
