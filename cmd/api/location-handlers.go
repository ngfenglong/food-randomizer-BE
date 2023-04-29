package main

import (
	"backend/models"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

type LocationDto struct {
	ID           int    `json:"id"`
	LocationName string `json:"location_name"`
}

func (app *application) getLocationByID(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	location, err := app.models.DB.GetLocationByID(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.WriteJSON(w, http.StatusOK, location, "location")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}
func (app *application) getAllLocations(w http.ResponseWriter, r *http.Request) {
	locations, err := app.models.DB.GetAllLocations()
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.WriteJSON(w, http.StatusOK, locations, "locations")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

func (app *application) editLocation(w http.ResponseWriter, r *http.Request) {
	var payload LocationDto
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	var location models.Location
	if payload.ID != 0 {
		m, err := app.models.DB.GetLocationByID(payload.ID)
		if err != nil {
			app.errorJSON(w, err)
			return
		}
		location = *m
	}

	location.ID = payload.ID
	location.LocationName = payload.LocationName
	location.UpdatedAt = time.Now()

	if location.ID == 0 {
		err = app.models.DB.InsertLocation(location)
		if err != nil {
			app.errorJSON(w, err)
			return
		}
	} else {
		err = app.models.DB.UpdateLocation(location)
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

func (app *application) deleteLocation(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.models.DB.DeleteLocation(id)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.WriteJSON(w, http.StatusOK, "Deleted Successfully", "response")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

func (app *application) deleteLocations(w http.ResponseWriter, r *http.Request) {
	var idList []int

	err := json.NewDecoder(r.Body).Decode(&idList)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	if len(idList) == 0 {
		app.errorJSON(w, errors.New("the ID list is empty"))
		return
	}

	err = app.models.DB.DeleteLocations(idList)
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.WriteJSON(w, http.StatusOK, "Deleted Successfully", "response")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}
