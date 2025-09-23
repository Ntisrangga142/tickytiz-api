package repositories

import (
	"context"
	"fmt"

	"github.com/Ntisrangga142/API_tickytiz/internals/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type GenreRepo struct {
	DB *pgxpool.Pool
}

func NewGenreRepo(db *pgxpool.Pool) *GenreRepo {
	return &GenreRepo{DB: db}
}

// GetGenres ambil daftar genre
func (r *GenreRepo) GetGenres(ctx context.Context) ([]models.Genre, error) {
	query := `SELECT id, name FROM genres`

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed query genres: %w", err)
	}
	defer rows.Close()

	var genres []models.Genre
	for rows.Next() {
		var g models.Genre
		if err := rows.Scan(&g.ID, &g.Name); err != nil {
			return nil, err
		}
		genres = append(genres, g)
	}

	return genres, nil
}
