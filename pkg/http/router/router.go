package router

import (
	"database/sql"

	"github.com/gorilla/mux"

	"github.com/ngfenglong/food-randomizer-BE/pkg/auth"
	"github.com/ngfenglong/food-randomizer-BE/pkg/category"
	"github.com/ngfenglong/food-randomizer-BE/pkg/config"
	"github.com/ngfenglong/food-randomizer-BE/pkg/location"
	"github.com/ngfenglong/food-randomizer-BE/pkg/middleware"
	"github.com/ngfenglong/food-randomizer-BE/pkg/place"
)

func NewRouter(cfg *config.Config, db *sql.DB) *mux.Router {
	r := mux.NewRouter()

	r.Use(middleware.EnableCORS)

	// Create Repo
	authRepo := auth.NewSQLAuthRepository(db)
	categoryRepo := category.NewSQLCategoryRepostory(db)
	locationRepo := location.NewSQLLocationRepository(db)
	placeRepo := place.NewSQLPlaceRepository(db)

	// Handle API
	api := r.PathPrefix("/v1").Subrouter()
	// api.HandleFunc("/places", auth.PlaceHandler(db)).Methods("GET")

	// Places
	api.HandleFunc("/places", place.GetAllPlaces(placeRepo)).Methods("GET")
	api.HandleFunc("/places/:id", place.GetPlaceByID(placeRepo)).Methods("GET")
	api.HandleFunc("/admin/updatePlace", place.EditPlace(placeRepo)).Methods("PUT")
	api.HandleFunc("/admin/deletePlace/:id", place.DeletePlace(placeRepo)).Methods("DELETE")
	api.HandleFunc("/admin/deletePlaces", place.DeletePlaces(placeRepo)).Methods("POST")
	api.HandleFunc("/generatePlace", place.GeneratePlace(placeRepo)).Methods("GET")

	// Categories
	api.HandleFunc("/admin/categories", category.GetAllCategories(categoryRepo)).Methods("GET")
	api.HandleFunc("/admin/categories/:id", category.GetCategoryByID(categoryRepo)).Methods("GET")
	api.HandleFunc("/admin/updateCategory", category.EditCategory(categoryRepo)).Methods("PUT")
	api.HandleFunc("/v1/admin/deleteCategory/:id", category.DeleteCategory(categoryRepo)).Methods("DELETE")
	api.HandleFunc("/v1/admin/deleteCategories", category.DeleteCategories(categoryRepo)).Methods("POST")

	// Location
	api.HandleFunc("/v1/admin/locations", location.GetAllLocations(locationRepo)).Methods("GET")
	api.HandleFunc("/v1/admin/locations/:id", location.GetLocationByID(locationRepo)).Methods("GET")
	api.HandleFunc("/v1/admin/updateLocation", location.EditLocation(locationRepo)).Methods("PUT")
	api.HandleFunc("/v1/admin/deleteLocation/:id", location.DeleteLocation(locationRepo)).Methods("DELETE")
	api.HandleFunc("/v1/admin/deleteLocations", location.DeleteLocations(locationRepo)).Methods("POST")

	api.HandleFunc("/v1/auth/login", auth.Login(authRepo)).Methods("POST")
	api.HandleFunc("/v1/auth/logout", auth.Logout(authRepo)).Methods("POST")
	api.HandleFunc("/v1/auth/register", auth.Register(authRepo, cfg)).Methods("POST")
	api.HandleFunc("/v1/auth/forget-password", auth.ForgetPassword(authRepo)).Methods("POST")

	// Telegram_Access
	api.HandleFunc("/v1/auth/requestAccess", auth.Request_access(authRepo)).Methods("POST")

	// return app.enableCORS(router)

	return r
}
