package main

import (
	"backend/models"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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

func (app *application) getAllPlaces(w http.ResponseWriter, r *http.Request) {
	fmt.Print("TYesting")
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

	err = app.models.DB.DeletePlace(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.WriteJSON(w, http.StatusOK, nil, "response")
	if err != nil {
		app.errorJSON(w, err)
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
	}

	var place models.Place

	if payload.ID != 0 {
		m, _ := app.models.DB.GetPlaceByID(payload.ID)
		place = *m
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
