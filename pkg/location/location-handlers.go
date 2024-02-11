package location

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/ngfenglong/food-randomizer-BE/pkg/models"
	"github.com/ngfenglong/food-randomizer-BE/pkg/utils"
)

type LocationDto struct {
	ID           int    `json:"id"`
	LocationName string `json:"location_name"`
}

func GetLocationByID(repo LocationRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := httprouter.ParamsFromContext(r.Context())
		id, err := strconv.Atoi(params.ByName("id"))
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		location, err := repo.GetLocationByID(ctx, id)
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		err = utils.WriteJSON(w, http.StatusOK, location, "location")
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}
	}
}

func GetAllLocations(repo LocationRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		locations, err := repo.GetAllLocations(ctx)
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		err = utils.WriteJSON(w, http.StatusOK, locations, "locations")
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}
	}
}

func EditLocation(repo LocationRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload LocationDto
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		var location models.Location
		if payload.ID != 0 {
			m, err := repo.GetLocationByID(ctx, payload.ID)
			if err != nil {
				utils.ErrorJSON(w, err)
				return
			}
			location = *m
		}

		location.ID = payload.ID
		location.LocationName = payload.LocationName
		location.UpdatedAt = time.Now()

		if location.ID == 0 {
			err = repo.InsertLocation(ctx, location)
			if err != nil {
				utils.ErrorJSON(w, err)
				return
			}
		} else {
			err = repo.UpdateLocation(ctx, location)
			if err != nil {
				utils.ErrorJSON(w, err)
				return
			}
		}

		err = utils.WriteJSON(w, http.StatusOK, nil, "response")
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}
	}
}

func DeleteLocation(repo LocationRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := httprouter.ParamsFromContext(r.Context())
		id, err := strconv.Atoi(params.ByName("id"))
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		err = repo.DeleteLocation(ctx, id)
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		err = utils.WriteJSON(w, http.StatusOK, "Deleted Successfully", "response")
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}
	}
}

func DeleteLocations(repo LocationRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var idList []int

		err := json.NewDecoder(r.Body).Decode(&idList)
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		if len(idList) == 0 {
			utils.ErrorJSON(w, errors.New("the ID list is empty"))
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		err = repo.DeleteLocations(ctx, idList)
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		err = utils.WriteJSON(w, http.StatusOK, "Deleted Successfully", "response")
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}
	}
}
