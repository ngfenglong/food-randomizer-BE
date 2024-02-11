package category

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

type CategoryDto struct {
	ID           int    `json:"id"`
	CategoryName string `json:"category_name"`
}

func GetCategoryByID(repo CategoryRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := httprouter.ParamsFromContext(r.Context())
		id, err := strconv.Atoi(params.ByName("id"))
		if err != nil {
			utils.ErrorJSON(w, err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		category, err := repo.GetCategoryByID(ctx, id)
		if err != nil {
			utils.ErrorJSON(w, err)
		}

		err = utils.WriteJSON(w, http.StatusOK, category, "category")
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}
	}
}

func GetAllCategories(repo CategoryRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		categories, err := repo.GetAllCategories(ctx)

		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}

		err = utils.WriteJSON(w, http.StatusOK, categories, "categories")
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}
	}
}

func EditCategory(repo CategoryRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var payload CategoryDto
		err := json.NewDecoder(r.Body).Decode(&payload)

		if err != nil {
			utils.ErrorJSON(w, err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		var category models.Category
		if payload.ID != 0 {
			m, err := repo.GetCategoryByID(ctx, payload.ID)
			if err != nil {
				utils.ErrorJSON(w, err)
				return
			}
			category = *m
		}

		category.ID = payload.ID
		category.CategoryName = payload.CategoryName
		category.UpdatedAt = time.Now()

		if category.ID == 0 {
			err = repo.InsertCategory(ctx, category)
			if err != nil {
				utils.ErrorJSON(w, err)
				return
			}
		} else {
			err = repo.UpdateCategory(ctx, category)
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

func DeleteCategory(repo CategoryRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := httprouter.ParamsFromContext(r.Context())
		id, err := strconv.Atoi(params.ByName("id"))
		if err != nil {
			utils.ErrorJSON(w, err)
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		err = repo.DeleteCategory(ctx, id)
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

func DeleteCategories(repo CategoryRepository) http.HandlerFunc {
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

		err = repo.DeleteCategories(ctx, idList)
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
