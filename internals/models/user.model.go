package models

import "time"

type UserVA struct {
	ID             int    `json:"id"`
	VirtualAccount string `json:"virtual_account"`
}

type UserProfile struct {
	ID             int     `json:"id" example:"101"`
	FirstName      *string `json:"firstname" example:"Rangga"`
	LastName       *string `json:"lastname" example:"Saputra"`
	Phone          *string `json:"phone" example:"+628123456789"`
	ProfileImg     *string `json:"profileimg" example:"https://example.com/uploads/profile_101.png"`
	VirtualAccount *string `json:"virtual_account" example:"1234567890123456"`
	Point          int     `json:"point" example:"250"`
	Email          string  `json:"email" example:"rangga@example.com"`
	Role           string  `json:"role" example:"user"`
}

type UpdateProfile struct {
	ProfileImg     *string `json:"profileimg,omitempty" example:"https://example.com/uploads/profile_101.png"` // akan diisi path setelah upload file
	FirstName      *string `form:"firstname" example:"Rangga"`
	LastName       *string `form:"lastname" example:"Saputra"`
	Phone          *string `form:"phone" example:"+628123456789"`
	VirtualAccount *string `form:"virtual_account" example:"1234567890123456"`
}

type OrderHistory struct {
	OrderID       int       `json:"order_id" example:"501"`
	IsPaid        bool      `json:"ispaid" example:"true"`
	TotalPrice    int       `json:"total_price" example:"150000"`
	QRCode        string    `json:"qrcode" example:"https://example.com/qrcode/501.png"`
	OrderName     string    `json:"order_name" example:"Rangga Saputra"`
	OrderEmail    string    `json:"order_email" example:"rangga@example.com"`
	OrderPhone    string    `json:"order_phone" example:"+628123456789"`
	PaymentMethod string    `json:"payment_method" example:"Credit Card"`
	ShowDate      time.Time `json:"show_date" example:"2025-09-20T19:30:00Z"`
	ShowTime      string    `json:"show_time" example:"19:30"`
	CinemaName    string    `json:"cinema_name" example:"XXI Plaza Indonesia"`
	CinemaLogo    string    `json:"cinema_logo" example:"XXI.jpg"`
	LocationName  string    `json:"location_name" example:"Jakarta"`
	MovieTitle    string    `json:"movie_title" example:"Avengers: Endgame"`
	MoviePoster   string    `json:"movie_poster" example:"https://example.com/posters/avengers.jpg"`
	MovieBackdrop string    `json:"movie_backdrop" example:"https://example.com/backdrops/avengers-bg.jpg"`
	Duration      int       `json:"duration" example:"180"`
	Rating        float32   `json:"rating" example:"8.5"`
	Seats         []*string `json:"seats" example:"[\"A1\",\"A2\",\"A3\"]"`
}

type OrderHistoryResponse struct {
	UserID      int            `json:"user_id" example:"101"`
	ListHistory []OrderHistory `json:"history"`
}

type UserPassword struct {
	ID       int    `json:"id"`
	Password string `json:"password"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}
