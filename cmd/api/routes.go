package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	// Places
	router.HandlerFunc(http.MethodGet, "/v1/places", app.getAllPlaces)
	router.HandlerFunc(http.MethodGet, "/v1/places/:id", app.getPlaceByID)
	router.HandlerFunc(http.MethodPost, "/v1/admin/updatePlace", app.editPlace)
	router.HandlerFunc(http.MethodGet, "/v1/admin/deletePlace/:id", app.deletePlace)
	router.HandlerFunc(http.MethodPost, "/v1/admin/deletePlaces", app.deletePlaces)
	router.HandlerFunc(http.MethodGet, "/v1/generatePlace", app.generatePlace)

	// Categories
	router.HandlerFunc(http.MethodGet, "/v1/admin/categories", app.getAllCategories)
	router.HandlerFunc(http.MethodGet, "/v1/admin/categories/:id", app.getCategoryByID)
	router.HandlerFunc(http.MethodPost, "/v1/admin/updateCategory", app.editCategory)
	router.HandlerFunc(http.MethodGet, "/v1/admin/deleteCategory/:id", app.deleteCategory)
	router.HandlerFunc(http.MethodPost, "/v1/admin/deleteCategories", app.deleteCategories)

	// Location
	router.HandlerFunc(http.MethodGet, "/v1/admin/locations", app.getAllLocations)
	router.HandlerFunc(http.MethodGet, "/v1/admin/locations/:id", app.getLocationByID)
	router.HandlerFunc(http.MethodPost, "/v1/admin/updateLocation", app.editLocation)
	router.HandlerFunc(http.MethodGet, "/v1/admin/deleteLocation/:id", app.deleteLocation)
	router.HandlerFunc(http.MethodPost, "/v1/admin/deleteLocations", app.deleteLocations)

	router.HandlerFunc(http.MethodPost, "/v1/auth/login", app.login)
	router.HandlerFunc(http.MethodPost, "/v1/auth/logout", app.logout)
	router.HandlerFunc(http.MethodPost, "/v1/auth/register", app.register)
	router.HandlerFunc(http.MethodPost, "/v1/auth/forget-password", app.forgetPassword)

	return app.enableCORS(router)
}
