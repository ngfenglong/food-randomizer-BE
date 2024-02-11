package place

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/ngfenglong/food-randomizer-BE/pkg/models"
	"github.com/ngfenglong/food-randomizer-BE/pkg/utils"
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

func GeneratePlace(repo PlaceRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		queryParams := r.URL.Query()
		isHalalParam, isHalalExists := queryParams["is_halal"]
		isVegetarianParam, isVegetarianExists := queryParams["is_vegetarian"]

		var isHalal, isVegetarian bool
		var err error

		if isHalalExists && len(isHalalParam) > 0 {
			isHalal, err = strconv.ParseBool(isHalalParam[0])
			if err != nil {
				utils.ErrorJSON(w, err)
				return
			}
		}

		if isVegetarianExists && len(isVegetarianParam) > 0 {
			isVegetarian, err = strconv.ParseBool(isVegetarianParam[0])
			if err != nil {
				utils.ErrorJSON(w, err)
				return
			}
		}

		var places []*models.Place
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if isHalalExists || isVegetarianExists {
			places, err = repo.GetAllPlacesWithFilter(ctx, isHalal, isVegetarian)
		} else {
			places, err = repo.GetAllPlaces(ctx)
		}

		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		idx := rand.Intn(len(places))
		err = utils.WriteJSON(w, http.StatusOK, places[idx], "place")
		if err != nil {
			utils.ErrorJSON(w, err)
		}
	}
}

func GetAllPlaces(repo PlaceRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		places, err := repo.GetAllPlaces(ctx)

		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		err = utils.WriteJSON(w, http.StatusOK, places, "places")
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}
	}
}

func GetPlaceByID(repo PlaceRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := httprouter.ParamsFromContext(r.Context())
		id, err := strconv.Atoi(params.ByName("id"))
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		place, err := repo.GetPlaceByID(ctx, id)
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}
		err = utils.WriteJSON(w, http.StatusOK, place, "place")

		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}
	}
}

func DeletePlace(repo PlaceRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := httprouter.ParamsFromContext(r.Context())
		id, err := strconv.Atoi(params.ByName("id"))
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		// Check if ID exists
		place, err := repo.GetPlaceByID(ctx, id)
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		if place == nil {
			utils.ErrorJSON(w, errors.New("ID does not exists"))
		}

		err = repo.DeletePlace(ctx, id)
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		err = utils.WriteJSON(w, http.StatusOK, nil, "response")
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}
	}
}

func DeletePlaces(repo PlaceRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var idList []int

		err := json.NewDecoder(r.Body).Decode(&idList)
		if err != nil {
			utils.ErrorJSON(w, err)
		}

		if len(idList) == 0 {
			utils.ErrorJSON(w, errors.New("the ID list is empty"))
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		err = repo.DeletePlaces(ctx, idList)
		if err != nil {
			utils.ErrorJSON(w, err)
		}

		err = utils.WriteJSON(w, http.StatusOK, "Deleted Successfully", "response")
		if err != nil {
			utils.ErrorJSON(w, err)
		}
	}
}

func EditPlace(repo PlaceRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload PlaceDto

		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			log.Println("error", err)
			utils.ErrorJSON(w, err)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		var place models.Place
		if payload.ID != 0 {
			m, err := repo.GetPlaceByID(ctx, payload.ID)

			if err != nil {
				if err == sql.ErrNoRows {
					payload.ID = 0
				} else {
					utils.ErrorJSON(w, err)
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
			err = repo.InsertPlace(ctx, place)
			if err != nil {
				utils.ErrorJSON(w, err)
				return
			}
		} else {
			err = repo.UpdatePlace(ctx, place)
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
