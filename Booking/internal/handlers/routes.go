
package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"BookingSystem/Booking/internal/config"
	"BookingSystem/Booking/internal/loggers"
	"BookingSystem/Booking/internal/repository"
	"BookingSystem/Booking/internal/services"
	"BookingSystem/Booking/pkg/database"

	"github.com/gorilla/mux"
)

func SetupRoutes(router *mux.Router, db *database.Db, logger *loggers.Logger, cfg *config.Config) {
	//  Wiring Takes place here
	// Initialize repositories
	userRepo := repository.NewUserRepository(db, logger)

	// Initialize services
	userService := services.NewUserService(userRepo, logger)

	// Initialize handlers
	userHandler := NewUserHandler(userService, logger)

	// API routes
	api := router.PathPrefix("/api/v1").Subrouter()
// POST Create 
	api.HandleFunc("/booking", userHandler.CreateBooking).Methods(http.MethodPost)

//  Get booking ID 
	api.HandleFunc("/bookings/{id}", userHandler.GetBookingByID).Methods(http.MethodGet)
// PUT
	api.HandleFunc("/bookings/{id}",userHandler.UpdateBooking,).Methods(http.MethodPut)
// Soft Delete
	api.HandleFunc("/bookings/{id}",userHandler.DeleteBooking,).Methods(http.MethodDelete)

	// Get Avaialable Bookings 
	// api.HandleFunc( "/rooms"  , userHandler.GetAvaliableRooms).Methods(http.MethodGet)

	// Health check endpoint
	router.HandleFunc("/health", func  (w http.ResponseWriter, r *http.Request) {
		logger.Info("Received request to Health Request..------------------------------------->>>>> [Reoutes]")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "healthy", "service": "training-service"}`))
	}).Methods("GET")

	// Define the /docs/openapi.yaml route first to avoid being shadowed
	router.HandleFunc("/docs/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		openAPIPath := filepath.Join(cfg.ProjectRoot, "docs", "openapi.yaml")
		// Read the openAPI.yaml file
		yamlContent, err := os.ReadFile(openAPIPath)
		if err != nil {
			fmt.Println("Error reading openAPI.yaml:", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		// Replace the placeholder with the actual base URL
		updatedContent := strings.ReplaceAll(string(yamlContent), "{BASE_URL}", cfg.BaseUrl)
		// Serve the updated OpenAPI documentation
		w.Header().Set("Content-Type", "application/x-yaml")

		if _, writeErr := w.Write([]byte(updatedContent)); writeErr != nil {
			logger.Info("Error writing response")
		}
	}).Methods(http.MethodGet)

	// Serve the Swagger UI files on the /docs path
	swaggerUIPath := http.Dir(filepath.Join(cfg.ProjectRoot, "docs", "swaggerui"))
	fs := http.FileServer(swaggerUIPath)
	router.PathPrefix("/docs/").Handler(http.StripPrefix("/docs/", fs))

	logger.Info("Routes configured successfully")
}



