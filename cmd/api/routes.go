package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	// secure := alice.New(app.checkToken)
	router.HandlerFunc(http.MethodGet, "/v1/places/:id", app.getPlaceByID)
	router.HandlerFunc(http.MethodGet, "/v1/places", app.getAllPlaces)
	router.HandlerFunc(http.MethodPost, "/v1/admin/editPlace", app.editPlace)
	router.HandlerFunc(http.MethodGet, "/v1/admin/deletePlace/:id", app.deletePlace)
	// router.HandlerFunc(http.MethodPost, "/status", app.statusHandler)

	return app.enableCORS(router)
}
