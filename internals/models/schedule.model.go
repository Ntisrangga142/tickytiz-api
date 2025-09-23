package models

import "time"

type Schedule struct {
	ID        int       `json:"id"`
	Date      time.Time `json:"date"`
	Cinema    string    `json:"cinema"`
	CinemaIMG string    `json:"cinema_img"`
	Price     int       `json:"price"`
	Location  string    `json:"location"`
	ShowTime  string    `json:"show_time"`
}

type ScheduleResponse struct {
	MovieID  int        `json:"movie_id" example:"101"`
	Schedule []Schedule `json:"schedule"`
}

