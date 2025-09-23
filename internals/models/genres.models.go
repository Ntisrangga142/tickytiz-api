package models

// Genre model sesuai tabel genres
type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Response wrapper untuk daftar genre
type GenreResponse struct {
	Genres []Genre `json:"genres"`
}
