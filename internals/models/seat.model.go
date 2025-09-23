package models

type Seat struct {
	ID []string `json:"seat_id"`
}

type SeatResponse struct {
	ScheduleID int      `json:"schedule_id" example:"12"`
	Seat       []string `json:"seat_id" example:"A1,A2,A3"`
}
