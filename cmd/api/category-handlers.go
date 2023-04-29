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

type CategoryDto struct {
	ID           int    `json:"id"`
	CategoryName string `json:"category_name"`
}

func (app *application) getCategoryByID(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.errorJSON(w, err)
	}

	category, err := app.models.DB.GetCategoryByID(id)
	if err != nil {
		app.errorJSON(w, err)
	}

	err = app.WriteJSON(w, http.StatusOK, category, "category")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

func (app *application) getAllCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := app.models.DB.GetAllCategories()

	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.WriteJSON(w, http.StatusOK, categories, "categories")
	if err != nil {
		app.errorJSON(w, err)
		return
	}
}

func (app *application) editCategory(w http.ResponseWriter, r *http.Request) {
	var payload CategoryDto
	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {
		app.errorJSON(w, err)
	}

	var category models.Category
	if payload.ID != 0 {
		m, err := app.models.DB.GetCategoryByID(payload.ID)
		if err != nil {
			app.errorJSON(w, err)
			return
		}
		category = *m
	}

	category.ID = payload.ID
	category.CategoryName = payload.CategoryName
	category.UpdatedAt = time.Now()

	if category.ID == 0 {
		err = app.models.DB.InsertCategory(category)
		if err != nil {
			app.errorJSON(w, err)
			return
		}
	} else {
		err = app.models.DB.UpdateCategory(category)
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

func (app *application) deleteCategory(w http.ResponseWriter, r *http.Request) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		app.errorJSON(w, err)
		return
	}

	err = app.models.DB.DeleteCategory(id)
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

func (app *application) deleteCategories(w http.ResponseWriter, r *http.Request) {
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

	err = app.models.DB.DeleteCategories(idList)
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
