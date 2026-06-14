package dtos


type BookingRequest struct {
	CustomerUUID string `json:"customer_uuid" validate:"required,uuid"`
	RoomUUID     string `json:"room_uuid" validate:"required,uuid"`
	BookingDate  string `json:"booking_date" validate:"required"` // YYYY-MM-DD
	StartTime    string `json:"start_time" validate:"required"`   // HH:MM
	EndTime      string `json:"end_time" validate:"required"`     // HH:MM
}

type BookingResponse struct {
	BookingUUID string `json:"booking_id"`
	CustomerUUID string `json:"customer_uuid"`
	RoomUUID     string `json:"room_uuid"`
	BookingDate  string `json:"booking_date"`
	StartTime    string `json:"start_time"`
	EndTime      string `json:"end_time"`
}



type AvailableRoomResponse struct {
	RoomUUID string `json:"room_uuid"`
	RoomName string `json:"room_name"`
	Capacity int    `json:"capacity"`
	Status   string `json:"status"`
}
