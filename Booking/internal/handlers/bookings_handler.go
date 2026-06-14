
package handlers

import (
	"BookingSystem/Booking/internal/dtos"
	"BookingSystem/Booking/internal/errorcodes"
	"BookingSystem/Booking/internal/loggers"
	"BookingSystem/Booking/internal/services"
	"BookingSystem/Booking/internal/utils"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type UserHandler struct {
	userService services.UserServiceInterface
	logger      *loggers.Logger
}

type UserInterface interface {

	// Reserve
	 CreateBooking(w http.ResponseWriter, r *http.Request)
	 GetAvaliableRooms(w http.ResponseWriter, r *http.Request)
	 UpdateBooking(w http.ResponseWriter,r *http.Request)
	 DeleteBooking(w http.ResponseWriter,r *http.Request)
	 GetBookingByID(w http.ResponseWriter, r *http.Request)
}

func NewUserHandler(userService services.UserServiceInterface, logger *loggers.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger, 
	}
}

func (h *UserHandler) CreateBooking(w http.ResponseWriter, r *http.Request) {

	reqID := uuid.NewString()
	h.logger.Info("Received request to book a room")

	var req dtos.BookingRequest
	//  Json Validations
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode request body")
		apiResponse := &dtos.APIResponse{
			Status:  "Error",
			Code:    http.StatusBadRequest,
			Message: "Invalid request format",
			Errors: []dtos.Error{{
				Field:   "body",
				Message: "Invalid JSON format",
				Code:    errorcodes.ErrorCodeStatus[errorcodes.InvalidJSONFormatErrorCode]}},
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
		utils.WriteResponse(w, http.StatusBadRequest, apiResponse)
		return
	}

	// Validation
	validationErrors := utils.ValidateStruct(&req)
	if len(validationErrors) > 0 {
		utils.WriteResponse(w, http.StatusBadRequest, &dtos.APIResponse{
			Status:    "Error",
			RequestID: reqID,
			Message:   "Validation failed",
			Errors:    validationErrors,
			Timestamp: time.Now().UTC().Format(time.RFC3339),

		})
		return
	}

	// Service call
	response := h.userService.CreateBooking(r.Context(), req, reqID)

	utils.WriteResponse(w, response.Code, response)
}

// func (h *UserHandler) GetAvaliableRooms(w http.ResponseWriter , r *http.Request) {

// 	// reqID := uuid.NewString()
// 	h.logger.Info("Received request get Avalaible Rooms... [ Fetching Available Rooms]")

// }

func (h *UserHandler) GetBookingByID(w http.ResponseWriter, r *http.Request) {

	reqID := uuid.NewString()
	h.logger.Info("Received request to get booking by ID---------[GetBookingByID---HANDLER]")

	vars := mux.Vars(r)
	bookingID := vars["id"]

	if bookingID == "" {
		utils.WriteResponse(w, http.StatusBadRequest, &dtos.APIResponse{
			Status:    "Error",
			Code:      http.StatusBadRequest,
			RequestID: reqID,
			Message:   "Booking ID is required",
			Timestamp: time.Now().UTC().Format(time.RFC3339),

		})
		return
	}

	response := h.userService.GetBookingByID(r.Context(), bookingID, reqID)

	// response := utils.MapErrorCode(apiResponse)
	utils.WriteResponse(w, response.Code, response)
}

func (h *UserHandler) DeleteBooking(w http.ResponseWriter,r *http.Request,) {
	h.logger.Info("Received request to Cancel Request.. [Handler]")
	reqID := uuid.NewString()

	bookingID := mux.Vars(r)["id"]
	
	if bookingID == "" {
		utils.WriteResponse(w, http.StatusBadRequest, &dtos.APIResponse{
			Status:    "Error",
			Code:      http.StatusBadRequest,
			RequestID: reqID,
			Message:   "Booking ID is required",
			Timestamp: time.Now().UTC().Format(time.RFC3339),

		})
		return
	}
	response := h.userService.DeleteBooking(r.Context(),bookingID,reqID,)
	// response := utils.MapErrorCode(response_cancel)
	utils.WriteResponse(w, response.Code, response)
}

func (h *UserHandler) UpdateBooking(w http.ResponseWriter,r *http.Request,) {
	h.logger.Info("Received request to UpdateBooking Request.. [Handler]")
	reqID := uuid.NewString()

	var req dtos.BookingRequest
	//  Json Validations.
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode request body")
		apiResponse := &dtos.APIResponse{
			Status:  "Error",
			Code:    http.StatusUnprocessableEntity,
			Message: "Invalid request format",
			Errors: []dtos.Error{{
				Field:   "body",
				Message: "Invalid JSON format",
				Code:    errorcodes.ErrorCodeStatus[errorcodes.InvalidJSONFormatErrorCode]}},
			RequestID: reqID,
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		}
		utils.WriteResponse(w, http.StatusUnprocessableEntity, apiResponse)
		return
	}

	validationErrors := utils.ValidateStruct(&req)
	if len(validationErrors) > 0 {
		utils.WriteResponse(w, http.StatusUnprocessableEntity, &dtos.APIResponse{
			Code:  http.StatusUnprocessableEntity,
			Status:    "Error",
			RequestID: reqID,
			Message:   "Validation failed",
			Errors:    validationErrors,
			Timestamp: time.Now().UTC().Format(time.RFC3339),

		})
		return
	}

	bookingID := mux.Vars(r)["id"]
	if bookingID == "" {
		utils.WriteResponse(w, http.StatusBadRequest, &dtos.APIResponse{
			Status:    "Error",
			Code:      http.StatusBadRequest,
			RequestID: reqID,
			Message:   "Booking ID is required",
			Timestamp: time.Now().UTC().Format(time.RFC3339),

		})
		return
	}
	response := h.userService.UpdateBooking(r.Context(),bookingID,req,reqID,)

	utils.WriteResponse(w, response.Code, response)
}

// func (h *UserHandler) GetAvaliableRooms(
// 	w http.ResponseWriter,
// 	r *http.Request,
// ) {

// 	reqID := uuid.NewString()

// 	date := r.URL.Query().Get("date")

// 	if date == "" {

// 		utils.WriteResponse(w, http.StatusBadRequest, &dtos.APIResponse{
// 			Status:    "Error",
// 			Code:      http.StatusBadRequest,
// 			RequestID: reqID,
// 			Message:   "date query param is required",
// 			Timestamp: time.Now().UTC().Format(time.RFC3339),

// 		})

// 		return
// 	}

// 	response := h.userService.GetAvailableRooms(
// 		r.Context(),
// 		date,
// 		reqID,
// 	)

// 	utils.WriteResponse(w, response.Code, response)
// }


