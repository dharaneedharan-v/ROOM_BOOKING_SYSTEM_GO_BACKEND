

package models

import (
	"time"

	"gorm.io/gorm"
)

type Room struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"-"`            // Internal PK
	RoomUUID  string    `gorm:"type:char(36);uniqueIndex;not null" json:"id"` // Public ID
	RoomName  string    `gorm:"type:varchar(200);not null" json:"room_name"`
	Capacity  int       `gorm:"not null" json:"capacity"`

	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedBy string    `gorm:"type:varchar(100);default:SYSTEM" json:"created_by"`
	UpdatedBy string    `gorm:"type:varchar(100);default:SYSTEM" json:"updated_by"`
	IsActive  bool      `gorm:"default:true;not null" json:"is_active"`
}

type Customer struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"-"`
	CustomerUUID    string    `gorm:"type:char(36);uniqueIndex;not null" json:"id"`
	CustomerName    string    `gorm:"type:varchar(200);not null" json:"customer_name"`
	CustomerPhone   string    `gorm:"type:varchar(20);uniqueIndex;not null" json:"customer_phone"`
	CustomerAddress string    `gorm:"type:text" json:"customer_address"`
	CustomerEmail   string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"customer_email"`

	CreatedAt       time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedBy       string    `gorm:"type:varchar(100);default:SYSTEM" json:"created_by"`
	UpdatedBy       string    `gorm:"type:varchar(100);default:SYSTEM" json:"updated_by"`
	IsActive        bool      `gorm:"default:true;not null" json:"is_active"`
}

type Booking struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"-"`
	BookingUUID string    `gorm:"type:char(36);uniqueIndex;not null" json:"id"`

	//  FK using the UUIDS
	CustomerUUID string   `gorm:"type:char(36);not null" json:"customer_uuid"`
	RoomUUID     string   `gorm:"type:char(36);not null" json:"room_uuid"`

	// --- Internal Database Links ---
	CustomerID  uint      `gorm:"not null" json:"-"`
	RoomID      uint      `gorm:"not null" json:"-"`

	// --- GORM Relations ---
	Customer    Customer  `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Room        Room      `gorm:"foreignKey:RoomID" json:"room,omitempty"`

	// --- Booking Details ---
	BookingDate time.Time `gorm:"type:date;not null" json:"booking_date"`
	StartTime   time.Time `gorm:"not null" json:"start_time"` // Removed type:time for better Go compatibility
	EndTime     time.Time `gorm:"not null" json:"end_time"`

	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	CreatedBy   string    `gorm:"type:varchar(100);default:SYSTEM" json:"created_by"`
	UpdatedBy   string    `gorm:"type:varchar(100);default:SYSTEM" json:"updated_by"`
	IsActive    bool      `gorm:"default:true;not null" json:"is_active"`

	// ---- Soft Delete -------

	Status      string  `gorm:"type:varchar(50);default:CONFIRMED;not null" json:"status"`
	DeletedAt gorm.DeletedAt   `gorm:"index" json:"-"`
	 
}