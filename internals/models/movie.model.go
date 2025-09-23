package models

type Movie struct {
	ID          int       `json:"id" example:"1"`
	Title       string    `json:"title" example:"Inception"`
	Poster      string    `json:"poster" example:"inception.jpg"`
	Backdrop    string    `json:"backdrop" example:"inception_backdrop.jpg"`
	ReleaseDate DateOnly  `json:"release_date" example:"2010-07-16"`
	Duration    int       `json:"duration" example:"148"`
	Synopsis    string    `json:"synopsis" example:"A thief who enters the dreams of others to steal secrets must pull off the ultimate heist."`
	Rating      float64   `json:"rating" example:"8.8"`
	Genres      []*string `json:"genres" example:"Action,Sci-Fi,Thriller"`
}

type MovieDetail struct {
	ID          int       `json:"id" example:"1"`
	Title       string    `json:"title" example:"Inception"`
	Poster      string    `json:"poster" example:"inception.jpg"`
	Backdrop    string    `json:"backdrop" example:"inception_backdrop.jpg"`
	ReleaseDate DateOnly  `json:"release_date" example:"2010-07-16T00:00:00Z"`
	Duration    int       `json:"duration" example:"148"`
	Synopsis    string    `json:"synopsis" example:"A thief who enters the dreams of others to steal secrets must pull off the ultimate heist."`
	Rating      float64   `json:"rating" example:"8.8"`
	Genres      []*string `json:"genres" example:"Action,Sci-Fi,Thriller"`
	Director    string    `json:"director" example:"Christopher Nolan"`
	Casts       []string  `json:"casts" example:"Leonardo DiCaprio, Joseph Gordon-Levitt, Ellen Page"`
}

