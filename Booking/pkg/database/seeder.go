package database

import (
	"BookingSystem/Booking/internal/models"
	"time"

	// "github.com/google/uuid"
	"gorm.io/gorm"
)

func SeedDatabase(db *gorm.DB) error {
	// 1. Seed 3 Rooms
	rooms := []models.Room{
		{RoomUUID: "c50fe215-347f-46a0-bdc0-51479f96d451", RoomName: "CR-1", Capacity: 5},
		{RoomUUID: "c50fe215-347f-46a0-bdc0-51479f96d452", RoomName: "SeminarHall", Capacity: 50},
		{RoomUUID: "c50fe215-347f-46a0-bdc0-51479f96d453", RoomName: "MR-1", Capacity: 15},
	}

	for i := range rooms {
		// Use RoomName to check for existence before creating
		if err := db.Where(models.Room{RoomName: rooms[i].RoomName}).FirstOrCreate(&rooms[i]).Error; err != nil {
			return err
		}
	}

	// 2. Seed 3 Customers
	customers := []models.Customer{
		{
			CustomerUUID:    "e71260ef-4b14-4b99-9ef3-eba0ddfd48b1",
			CustomerName:    "Arun",
			CustomerEmail:   "Arun@continental.com",
			CustomerPhone:   "6380001985",
			CustomerAddress: "Chennai",
		},
		{
			CustomerUUID:    "e71260ef-4b14-4b99-9ef3-eba0ddfd48b2",
			CustomerName:    "dhanush",
			CustomerEmail:   "dhanush@221b.com",
			CustomerPhone:   "6380001858",
			CustomerAddress: "Chennai , Avadi",
		},
		{
			CustomerUUID:    "e71260ef-4b14-4b99-9ef3-eba0ddfd48b3",
			CustomerName:    "Tomy",
			CustomerEmail:   "tomy@starkintl.com",
			CustomerPhone:   "1234567890",
			CustomerAddress: "Chennai, perambalur",
		},
	}

	for i := range customers {
		// Use CustomerEmail to check for existence
		if err := db.Where(models.Customer{CustomerEmail: customers[i].CustomerEmail}).FirstOrCreate(&customers[i]).Error; err != nil {
			return err
		}
	}

	// 3. Seed 3 Bookings
	// Note: We use the ID (uint) assigned by GORM after FirstOrCreate above
	bookings := []models.Booking{
		{
			BookingUUID:  "94841ade-468c-480f-8b69-ee911e6fcbdb",
			CustomerID:   customers[0].ID,
			CustomerUUID: customers[0].CustomerUUID, // Populates structural string field
			RoomID:       rooms[0].ID,
			RoomUUID:     rooms[0].RoomUUID,         // Populates structural string field
			BookingDate:  time.Now(),
			StartTime:    time.Now().Add(time.Hour * 2),
			EndTime:      time.Now().Add(time.Hour * 4),
		},
		{
			BookingUUID:  "94841ade-468c-480f-8b69-ee911e6fcbdc",
			CustomerID:   customers[1].ID,
			CustomerUUID: customers[1].CustomerUUID, 
			RoomID:       rooms[1].ID,
			RoomUUID:     rooms[1].RoomUUID,         
			BookingDate:  time.Now().AddDate(0, 0, 1),
			StartTime:    time.Now().Add(time.Hour * 24),
			EndTime:      time.Now().Add(time.Hour * 26),
		},
		{
			BookingUUID:  "94841ade-468c-480f-8b69-ee911e6fcbdd",
			CustomerID:   customers[2].ID,
			CustomerUUID: customers[2].CustomerUUID, 
			RoomID:       rooms[2].ID,
			RoomUUID:     rooms[2].RoomUUID,         
			BookingDate:  time.Now().AddDate(0, 0, 2),
			StartTime:    time.Now().Add(time.Hour * 48),
			EndTime:      time.Now().Add(time.Hour * 50),
		},
	}

	for i := range bookings {
		// Use BookingUUID to check for existence
		if err := db.Where(models.Booking{BookingUUID: bookings[i].BookingUUID}).FirstOrCreate(&bookings[i]).Error; err != nil {
			return err
		}
	}

	return nil
}
