package models

import (
	"mime/multipart"
	"time"
)

type AdminMovie struct {
	ID          int       `json:"id" example:"1"`
	Title       string    `json:"title" example:"Avengers: Endgame"`
	Poster      string    `json:"poster" example:"https://image.example.com/posters/avengers_endgame.jpg"`
	Backdrop    string    `json:"backdrop" example:"https://image.example.com/backdrops/avengers_endgame.jpg"`
	ReleaseDate time.Time `json:"release_date" example:"2019-04-26T00:00:00Z"`
	Duration    int       `json:"duration" example:"181"`
	Synopsis    string    `json:"synopsis" example:"Setelah pertempuran melawan Thanos, para Avengers bersatu untuk mengembalikan keseimbangan di alam semesta."`
	Rating      float64   `json:"rating" example:"8.4"`
	Genres      []string  `json:"genres" example:"[\"Action\", \"Adventure\", \"Sci-Fi\"]"`
	Director    string    `json:"director" example:"Anthony Russo, Joe Russo"`
	Casts       []string  `json:"casts" example:"[\"Robert Downey Jr.\", \"Chris Evans\", \"Mark Ruffalo\", \"Chris Hemsworth\"]"`
}

type AdminUpdate struct {
	Title       *string    `form:"title" example:"Avengers: Endgame - Director's Cut"`
	Poster      *string    `form:"poster" example:"https://image.example.com/posters/avengers_endgame_cut.jpg"`
	Backdrop    *string    `form:"backdrop" example:"https://image.example.com/backdrops/avengers_endgame_cut.jpg"`
	ReleaseDate *time.Time `form:"release_date" example:"2019-05-01T00:00:00Z"`
	Duration    *int       `form:"duration" example:"183"`
	Synopsis    *string    `form:"synopsis" example:"Versi director's cut dari Avengers: Endgame."`
	Rating      *float32   `form:"rating" example:"8.5"`
	IDDirector  *int       `form:"id_director" example:"12"`
}

type AdminDelete struct {
	ID      int    `json:"id" example:"1"`
	Message string `json:"message" example:"Movie successfully deleted"`
}

type AdminInsertMovie struct {
	Title        string                `form:"title" binding:"required"`
	Genres       string                `form:"genres" binding:"required"`
	ReleaseDate  string                `form:"release_date" binding:"required"`
	Duration     int                   `form:"duration" binding:"required"`
	DirectorName string                `form:"director_name" binding:"required"`
	CastName     string                `form:"cast_name" binding:"required"`
	Locations    string                `form:"location" binding:"required"`
	Dates        string                `form:"dates" binding:"required"`
	Times        string                `form:"times" binding:"required"`
	Poster       *multipart.FileHeader `form:"poster" binding:"required"`
	Backdrop     *multipart.FileHeader `form:"backdrop" binding:"required"`
}

type AdminInsertMovieResponse struct {
	ID          int      `json:"id"`
	Title       string   `json:"title"`
	Poster      string   `json:"poster"`
	Backdrop    string   `json:"backdrop"`
	ReleaseDate string   `json:"release_date"`
	Duration    int      `json:"duration"`
	Synopsis    string   `json:"synopsis"`
	Rating      float64  `json:"rating"`
	Genres      []string `json:"genres"`
	Director    string   `json:"director"`
	Casts       []string `json:"casts"`
}

type MovieAdminInsert struct {
	Title       string
	ReleaseDate time.Time
	Duration    int
	Synopsis    string
	IdDirector  int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type ScheduleComboAdminInsert struct {
	Date       string `json:"date"`
	IdCinema   int    `json:"id_cinema"`
	IdLocation int    `json:"id_location"`
	IdTime     int    `json:"id_time"`
}

type GetMovieDetailUpdate struct {
	ID          int              `json:"id"`
	Title       string           `json:"title"`
	Poster      string           `json:"poster"`
	Backdrop    string           `json:"backdrop"`
	ReleaseDate string           `json:"release_date"`
	Duration    int              `json:"duration"`
	Synopsis    string           `json:"synopsis"`
	IdDirector  int              `json:"id_director"`
	Director    string           `json:"director"`
	Genres      []Genre          `json:"genres"`
	Actors      []Actor          `json:"actors"`
	Schedules   []ScheduleDetail `json:"schedules"`
}

type Actor struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ScheduleDetail struct {
	ID         int    `json:"id"`
	Date       string `json:"date"`
	CinemaID   int    `json:"cinema_id"`
	Cinema     string `json:"cinema"`
	LocationID int    `json:"location_id"`
	Location   string `json:"location"`
	TimeID     int    `json:"time_id"`
	Time       string `json:"time"`
}

type MovieUpdateAdmin struct {
	ID           int
	Title        string
	ReleaseDate  time.Time
	Duration     int
	DirectorID   int
	Synopsis     string
	PosterPath   string
	BackdropPath string
	Schedules    []ScheduleUpdate
}

// ScheduleUpdate untuk update schedule
type ScheduleUpdate struct {
	Date       string `json:"date"`
	CinemaID   int    `json:"id_cinema"`
	LocationID int    `json:"id_location"`
	TimeID     int    `json:"id_time"`
}
