
package repository

import (
	"context"
	"errors"

	// "fmt"
	"time"

	"BookingSystem/Booking/internal/dtos"
	"BookingSystem/Booking/internal/loggers"
	"BookingSystem/Booking/internal/models"
	"BookingSystem/Booking/pkg/database"

	// "go.uber.org/zap"
	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	GetCustomerByUUID(ctx context.Context, uuid string) (*models.Customer, *dtos.Error)
	GetRoomByUUID(ctx context.Context, uuid string) (*models.Room, *dtos.Error)
	GetBookingByUUID(ctx context.Context, uuid string) (*models.Booking, *dtos.Error)
	// GetAvailableRooms(ctx context.Context, date time.Time) ([]models.Room, *dtos.Error)
	CheckRoomAvailability(ctx context.Context, roomID uint, start, end time.Time) (bool , string )
	CreateBooking(ctx context.Context, booking *models.Booking) *dtos.Error
	SoftDeleteBooking(ctx context.Context, bookingUUID string) *dtos.Error
// update
	UpdateBooking(ctx context.Context, booking *models.Booking) *dtos.Error
	CheckRoomAvailabilityForUpdate(ctx context.Context, roomID uint, start, end time.Time, bookingUUID string) bool

}


// In repositories 2 Operations That is Read and Write   
type UserRepository struct {
	db     *database.Db
	logger *loggers.Logger
}

func NewUserRepository(db *database.Db, logger *loggers.Logger) UserRepositoryInterface {
	return &UserRepository{
		db:     db,
		logger: logger,
	}
}

func (r *UserRepository) GetCustomerByUUID(ctx context.Context, uuid string) (*models.Customer, *dtos.Error) {
		r.logger.Info("Received request to ---------[GetCustomerByUUID---REPO]")

	var customer models.Customer
	if err := r.db.Gorm.WithContext(ctx).Where("customer_uuid = ?", uuid).First(&customer).Error; err != nil {
		return nil, &dtos.Error{Message: "Customer not found"}
	}
	return &customer, nil
}

func (r *UserRepository) GetRoomByUUID(ctx context.Context, uuid string) (*models.Room, *dtos.Error) {
		r.logger.Info("Received request to --------[GetRoomByUUID ---REPO]")

	var room models.Room
	if err := r.db.Gorm.WithContext(ctx).Where("room_uuid = ?", uuid).First(&room).Error; err != nil {
		return nil, &dtos.Error{Message: "Room not found"}
	}
	return &room, nil
}

func (r *UserRepository) CheckRoomAvailability(ctx context.Context, roomID uint, start, end time.Time) (bool , string)  {

	r.logger.Info("Received request to --------[CheckRoomAvailability ---REPO]")

	var conflictingBooking models.Booking
	
	// Query to find any overlapping booking exists
	err := r.db.Gorm.WithContext(ctx).Model(&models.Booking{}).
		Select("booking_uuid").
		Where("room_id = ? AND start_time < ? AND end_time > ? AND status != ?", roomID, end, start , "CANCELLED").
		Limit(1).
		Find(&conflictingBooking).Error

	// If a record was found, a conflict exists For the debugin purpose  returning the  UUID
	if err == nil && conflictingBooking.BookingUUID != "" {
		return true, conflictingBooking.BookingUUID
	}
	
	return false, ""
}



func (r *UserRepository) CreateBooking(ctx context.Context, booking *models.Booking) *dtos.Error {
		r.logger.Info("Received request to --------[CreateBooking ---REPO]")

	if err := r.db.Gorm.WithContext(ctx).Create(booking).Error; err != nil {
		return &dtos.Error{Message: "Failed to create booking"}
	}
	return nil
}

func (r *UserRepository) GetBookingByUUID(ctx context.Context, uuid string) (*models.Booking, *dtos.Error) {
		r.logger.Info("Received request to --------[ GetBookingByUUID ---REPO]")

	var booking models.Booking
	err := r.db.Gorm.WithContext(ctx).
		Where("booking_uuid = ? AND is_active = ? AND status != ?", uuid, true, "CANCELLED").
		First(&booking).Error

	// fmt.Println("--- DATABASE QUERY DATA RESULT BELOW ---")
	// fmt.Printf("Booking Struct Content: %+v\n", booking) 
	// fmt.Println("----------------------------------------")



	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &dtos.Error{
				Message: "Booking not found",
			}
		}
		return nil, &dtos.Error{
			Message: "Failed to fetch booking",
		}
	}
	return &booking, nil
}



func (r *UserRepository) SoftDeleteBooking(ctx context.Context, bookingUUID string) *dtos.Error {
	r.logger.Info("Received request to --------[SoftDeleteBooking ---REPO]")
	// r.logger.Info("Received request to SoftDeleteBooking REPO", zap.String("booking_uuid", bookingUUID))

	var booking models.Booking

	err := r.db.Gorm.WithContext(ctx).
		Select("status"). 
		Where("booking_uuid = ?", bookingUUID).
		First(&booking).Error

	// Handle case where booking doesn't exist at all
	if err != nil {
		return &dtos.Error{
			Message: "Booking not found",
		}
	}

	// Check if it is already cancelled
	if booking.Status == "CANCELLED" {
		return &dtos.Error{
			Message: "Booking Already Cancelled...!!",
		}
	}

	// Update it if it is active [ here we are not making it as inactive ]
	result := r.db.Gorm.WithContext(ctx).
		Model(&models.Booking{}).
		Where("booking_uuid = ?", bookingUUID).
		Updates(map[string]interface{}{
			"status":    "CANCELLED",
			"is_active": false,
		})

	if result.Error != nil {
		return &dtos.Error{
			Message: "Failed to delete booking",
		}
	}

	return nil
}


func (r *UserRepository) CheckRoomAvailabilityForUpdate(ctx context.Context,roomID uint,start, end time.Time,bookingUUID string,) bool {
	var count int64
	r.logger.Info("Received request to --------[ CheckRoomAvailabilityForUpdate  ---REPO]")

	// Check for the Any Conflict for the updating it. 
	r.db.Gorm.WithContext(ctx).Model(&models.Booking{}).
		Where("room_id = ? AND booking_uuid != ? AND start_time < ? AND end_time > ? AND status != ?", roomID, bookingUUID, end, start , "CANCELLED").
		Count(&count)

	return count > 0
}

func (r *UserRepository) UpdateBooking(ctx context.Context, booking *models.Booking) *dtos.Error {
			r.logger.Info("Received request to --------[UpdateBooking ---REPO]")

	// Save it
	if err := r.db.Gorm.WithContext(ctx).Save(booking).Error; err != nil {
		return &dtos.Error{
			Message: "Failed to update booking in database",
		}
	}
	return nil
}












// func (r *UserRepository) GetAvailableRooms(ctx context.Context, date time.Time) ([]models.Room, *dtos.Error) {
// 	var rooms []models.Room

// 	// Calculate bounds for the target date
// 	dayStart := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC)
// 	dayEnd := dayStart.Add(24 * time.Hour)

// 	// Fetch rooms that are active and not tied to an active booking within this time block
// 	err := r.db.Gorm.WithContext(ctx).
// 		Where("is_active = ? AND id NOT IN (?)", true,
// 			r.db.Gorm.Model(&models.Booking{}).
// 				Select("room_id").
// 				Where("start_time < ? AND end_time > ?", dayEnd, dayStart),
// 		).Find(&rooms).Error

// 	if err != nil {
// 		return nil, &dtos.Error{
// 			Message: "Failed to fetch available rooms",
// 		}
// 	}

// 	return rooms, nil
// }