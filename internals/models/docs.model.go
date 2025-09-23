package models

// Login
type LoginDocs struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Login successful"`
	Data    string `json:"data" example:"token"`
}

// Register
type RegisterDocs struct {
	Success bool   `json:"success" example:"true"`
	Message string `json:"message" example:"Register successful"`
	Data    string `json:"data" example:""`
}

type AdminMovieResponse struct {
	Success bool       `json:"success" example:"true"`
	Message string     `json:"message" example:"Movie fetched successfully"`
	Data    AdminMovie `json:"data"`
}

type AdminMovieUpdateResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Movie fetched successfully"`
	Data    AdminUpdate `json:"data"`
}

type AdminMovieDeleteResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Movie fetched successfully"`
	Data    AdminDelete `json:"data"`
}

type ResponseMovies struct {
	Success bool    `json:"success" example:"true"`
	Message string  `json:"message" example:"Success Load Movies"`
	Data    []Movie `json:"data"`
}

type ResponseMovieDetail struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Success Load Movies"`
	Data    MovieDetail `json:"data"`
}

type ResponseOrders struct {
	Success bool          `json:"success" example:"true"`
	Message string        `json:"message" example:"Success Load Movies"`
	Data    OrderResponse `json:"data"`
}

type ResponseSchedule struct {
	Success bool             `json:"success" example:"true"`
	Message string           `json:"message" example:"Success Load Movies"`
	Data    ScheduleResponse `json:"data"`
}

type ResponseSeats struct {
	Success bool         `json:"success" example:"true"`
	Message string       `json:"message" example:"Success Load Movies"`
	Data    SeatResponse `json:"data"`
}

type ResponseUserProfile struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Success Load Movies"`
	Data    UserProfile `json:"data"`
}

type ResponseOrderHistory struct {
	Success bool         `json:"success" example:"true"`
	Message string       `json:"message" example:"Success Load Movies"`
	Data    OrderHistory `json:"data"`
}

type ResponseUpdateProfile struct {
	Success bool          `json:"success" example:"true"`
	Message string        `json:"message" example:"Success Load Movies"`
	Data    UpdateProfile `json:"data"`
}

type ResponseCreateMovie struct {
	Success bool                     `json:"success" example:"true"`
	Message string                   `json:"message" example:"Success Load Movies"`
	Data    AdminInsertMovieResponse `json:"data"`
}
