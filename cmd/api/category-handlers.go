package main

import "net/http"

func (app *application) getCategoryByID(w http.ResponseWriter, r *http.Request) {

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

}

func (app *application) deleteCategory(w http.ResponseWriter, r *http.Request) {

}
