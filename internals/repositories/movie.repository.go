package repositories

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Ntisrangga142/API_tickytiz/internals/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MovieRepo struct {
	DB *pgxpool.Pool
}

func NewMovieRepo(db *pgxpool.Pool) *MovieRepo {
	return &MovieRepo{DB: db}
}

func (r *MovieRepo) GetUpcoming(ctx context.Context, page int) ([]models.Movie, error) {
	const limit = 4
	offset := (page - 1) * limit

	query := `
		SELECT m.id, m.title, m.poster, m.backdrop, m.release_date, 
		       m.duration, m.synopsis, m.rating, array_agg(g.name) AS genres
		FROM movies m
		JOIN movies_genres mg ON mg.id_movie = m.id
		JOIN genres g ON g.id = mg.id_genre
		WHERE m.release_date > CURRENT_DATE AND m.delete_at IS NULL
		GROUP BY m.id
		ORDER BY m.release_date ASC 
		LIMIT $1 OFFSET $2
	`

	rows, err := r.DB.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie
		var releaseDate time.Time
		if err := rows.Scan(
			&m.ID, &m.Title, &m.Poster, &m.Backdrop, &releaseDate,
			&m.Duration, &m.Synopsis, &m.Rating, &m.Genres,
		); err != nil {
			return nil, err
		}
		m.ReleaseDate = models.DateOnly(releaseDate)
		movies = append(movies, m)
	}

	return movies, nil
}

func (r *MovieRepo) GetPopular(ctx context.Context) ([]models.Movie, error) {
	query := `
		SELECT m.id, m.title, m.poster, m.backdrop, m.release_date, m.duration, m.synopsis, m.rating, array_agg(g.name) AS genres
		FROM movies m
		JOIN movies_genres mg ON mg.id_movie = m.id
		JOIN genres g ON g.id = mg.id_genre
		WHERE m.popularity > 60.0 AND delete_at IS NULL
		GROUP BY m.id
		ORDER BY m.popularity ASC LIMIT 4;
	`
	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movie
	var releaseDate time.Time
	for rows.Next() {
		var m models.Movie
		if err := rows.Scan(&m.ID, &m.Title, &m.Poster, &m.Backdrop, &releaseDate, &m.Duration, &m.Synopsis, &m.Rating, &m.Genres); err != nil {
			return nil, err
		}
		m.ReleaseDate = models.DateOnly(releaseDate)
		movies = append(movies, m)
	}

	return movies, nil
}

func (r *MovieRepo) GetFilteredMovies(ctx context.Context, title string, genres []string, page int) ([]models.Movie, int, error) {
	limit := 12
	offset := (page - 1) * limit

	baseQuery := `
		FROM movies m
		LEFT JOIN movies_genres mg ON mg.id_movie = m.id
		LEFT JOIN genres g ON g.id = mg.id_genre
		WHERE m.delete_at IS NULL
	`

	args := []any{}
	argID := 1

	// filter title
	if title != "" {
		baseQuery += fmt.Sprintf(" AND m.title ILIKE $%d", argID)
		args = append(args, "%"+title+"%")
		argID++
	}

	// filter genres
	groupBy := " GROUP BY m.id"
	if len(genres) > 0 {
		placeholders := []string{}
		for _, g := range genres {
			placeholders = append(placeholders, fmt.Sprintf("$%d", argID))
			args = append(args, g)
			argID++
		}
		baseQuery += fmt.Sprintf(" AND g.name IN (%s)", strings.Join(placeholders, ","))
		groupBy += fmt.Sprintf(" HAVING COUNT(DISTINCT g.id) = %d", len(genres))
	}

	// 🔹 Hitung total data
	countQuery := "SELECT COUNT(DISTINCT m.id) " + baseQuery
	var totalItems int
	if err := r.DB.QueryRow(ctx, countQuery, args...).Scan(&totalItems); err != nil {
		return nil, 0, err
	}

	// 🔹 Ambil data dengan pagination
	dataQuery := `
		SELECT m.id, m.title, m.poster, m.backdrop, m.release_date, m.duration,
		       m.synopsis, m.rating,
		       COALESCE(ARRAY_AGG(DISTINCT g.name), '{}') AS genres
		` + baseQuery + groupBy +
		fmt.Sprintf(" ORDER BY m.id ASC LIMIT %d OFFSET %d", limit, offset)

	rows, err := r.DB.Query(ctx, dataQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie
		var genresArr []*string
		var releaseDate time.Time
		if err := rows.Scan(&m.ID, &m.Title, &m.Poster, &m.Backdrop, &releaseDate,
			&m.Duration, &m.Synopsis, &m.Rating, &genresArr); err != nil {
			return nil, 0, err
		}
		m.ReleaseDate = models.DateOnly(releaseDate)
		m.Genres = genresArr
		movies = append(movies, m)
	}

	return movies, totalItems, nil
}

func (r *MovieRepo) GetMovieDetail(ctx context.Context, movieID int) (*models.MovieDetail, error) {
	query := `
		SELECT
			m.id,
			m.title,
			m.poster,
			m.backdrop,
			m.release_date,
			m.duration,
			m.synopsis,
			m.rating,
			COALESCE(ARRAY_AGG(DISTINCT g.name), '{}') AS genres,
			d.name AS director,
			COALESCE(ARRAY_AGG(DISTINCT a.name), '{}') AS casts
		FROM movies m
		JOIN directors d ON m.id_director = d.id
		LEFT JOIN movies_genres mg ON m.id = mg.id_movie
		LEFT JOIN genres g ON mg.id_genre = g.id
		LEFT JOIN movies_actors ma ON m.id = ma.id_movie
		LEFT JOIN actors a ON ma.id_actor = a.id
		WHERE m.id = $1
		GROUP BY m.id, d.name
	`

	row := r.DB.QueryRow(ctx, query, movieID)

	var movie models.MovieDetail
	var releaseDate time.Time
	err := row.Scan(
		&movie.ID,
		&movie.Title,
		&movie.Poster,
		&movie.Backdrop,
		&releaseDate,
		&movie.Duration,
		&movie.Synopsis,
		&movie.Rating,
		&movie.Genres,
		&movie.Director,
		&movie.Casts,
	)
	if err != nil {
		return nil, err
	}
	movie.ReleaseDate = models.DateOnly(releaseDate)

	return &movie, nil
}
