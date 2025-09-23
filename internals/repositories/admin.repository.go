package repositories

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Ntisrangga142/API_tickytiz/internals/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AdminRepo struct {
	DB *pgxpool.Pool
}

func NewAdminRepo(db *pgxpool.Pool) *AdminRepo {
	return &AdminRepo{DB: db}
}

func (r *AdminRepo) GetAllMovie(ctx context.Context) ([]models.AdminMovie, error) {
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
			ARRAY_AGG(DISTINCT g.name) AS genres,
			d.name AS director,
			ARRAY_AGG(DISTINCT a.name) AS casts
		FROM movies m
		JOIN movies_genres mg ON m.id = mg.id_movie
		JOIN genres g ON mg.id_genre = g.id
		JOIN directors d ON d.id = m.id_director
		JOIN movies_actors ma ON ma.id_movie = m.id
		JOIN actors a ON a.id = ma.id_actor
		WHERE m.delete_at IS NULL
		GROUP BY m.id, m.title, d.name
	`

	rows, err := r.DB.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.AdminMovie

	for rows.Next() {
		var m models.AdminMovie
		err := rows.Scan(
			&m.ID,
			&m.Title,
			&m.Poster,
			&m.Backdrop,
			&m.ReleaseDate,
			&m.Duration,
			&m.Synopsis,
			&m.Rating,
			&m.Genres,
			&m.Director,
			&m.Casts,
		)
		if err != nil {
			return nil, err
		}
		movies = append(movies, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil
}

func (r *AdminRepo) UpdateMovie(ctx context.Context, req models.AdminUpdate, id int) error {
	setClauses := []string{}
	args := []any{}
	argID := 1

	// --- String fields ---
	if req.Title != nil && *req.Title != "" {
		setClauses = append(setClauses, fmt.Sprintf("title = $%d", argID))
		args = append(args, *req.Title)
		argID++
	}
	if req.Synopsis != nil && *req.Synopsis != "" {
		setClauses = append(setClauses, fmt.Sprintf("synopsis = $%d", argID))
		args = append(args, *req.Synopsis)
		argID++
	}

	// --- Time field ---
	if req.ReleaseDate != nil && !req.ReleaseDate.IsZero() {
		setClauses = append(setClauses, fmt.Sprintf("release_date = $%d", argID))
		args = append(args, *req.ReleaseDate)
		argID++
	}

	// --- File fields ---
	if req.Poster != nil {
		setClauses = append(setClauses, fmt.Sprintf("poster = $%d", argID))
		args = append(args, *req.Poster) // bisa diganti path setelah upload
		argID++
	}
	if req.Backdrop != nil {
		setClauses = append(setClauses, fmt.Sprintf("backdrop = $%d", argID))
		args = append(args, *req.Backdrop)
		argID++
	}

	// --- Numeric pointer fields ---
	if req.Duration != nil {
		setClauses = append(setClauses, fmt.Sprintf("duration = $%d", argID))
		args = append(args, *req.Duration)
		argID++
	}
	if req.Rating != nil {
		setClauses = append(setClauses, fmt.Sprintf("rating = $%d", argID))
		args = append(args, *req.Rating)
		argID++
	}
	if req.IDDirector != nil && *req.IDDirector < 0 {
		setClauses = append(setClauses, fmt.Sprintf("id_director = $%d", argID))
		args = append(args, *req.IDDirector)
		argID++
	}

	// --- Jika tidak ada field untuk diupdate ---
	if len(setClauses) == 0 {
		return fmt.Errorf("no fields to update")
	}

	// --- Buat query ---
	query := fmt.Sprintf(`
		UPDATE movies
		SET %s
		WHERE id = $%d
	`, strings.Join(setClauses, ", "), argID)
	args = append(args, id)

	// --- Eksekusi ---
	_, err := r.DB.Exec(ctx, query, args...)
	return err
}

func (r *AdminRepo) DeleteMovie(ctx context.Context, id int) (*models.AdminDelete, error) {
	query := `
		UPDATE movies
		SET delete_at = NOW()
		WHERE id = $1
		RETURNING id
	`

	var deletedID int
	err := r.DB.QueryRow(ctx, query, id).Scan(&deletedID)
	if err != nil {
		return nil, err
	}

	return &models.AdminDelete{
		ID:      deletedID,
		Message: fmt.Sprintf("Movie with ID %d has been soft deleted", deletedID),
	}, nil
}

func (r *AdminRepo) InsertMovie(ctx context.Context, req models.AdminInsertMovie) (int, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	// --- Director ---
	var directorID int
	err = tx.QueryRow(ctx, "SELECT id FROM directors WHERE name=$1", req.DirectorName).Scan(&directorID)
	if err != nil {
		// insert jika tidak ada
		err = tx.QueryRow(ctx, "INSERT INTO directors(name) VALUES($1) RETURNING id", req.DirectorName).Scan(&directorID)
		if err != nil {
			return 0, err
		}
	}

	// --- Insert Movie ---
	var movieID int
	releaseDate, _ := time.Parse("2006-01-02", req.ReleaseDate)
	err = tx.QueryRow(ctx,
		`INSERT INTO movies(title, poster, backdrop, release_date, duration, synopsis, id_director)
		 VALUES($1,$2,$3,$4,$5,$6,$7) RETURNING id`,
		req.Title, req.Poster.Filename, req.Backdrop.Filename, releaseDate, req.Duration, "synopsis placeholder", directorID,
	).Scan(&movieID)
	if err != nil {
		return 0, err
	}

	// --- Genres ---
	genres := strings.Split(req.Genres, ",")
	for _, g := range genres {
		var genreID int
		err = tx.QueryRow(ctx, "SELECT id FROM genres WHERE name=$1", g).Scan(&genreID)
		if err != nil {
			err = tx.QueryRow(ctx, "INSERT INTO genres(name) VALUES($1) RETURNING id", g).Scan(&genreID)
			if err != nil {
				return 0, err
			}
		}
		_, err = tx.Exec(ctx, "INSERT INTO movies_genres(id_movie, id_genre) VALUES($1,$2)", movieID, genreID)
		if err != nil {
			return 0, err
		}
	}

	// --- Actors ---
	actors := strings.Split(req.CastName, ",")
	for _, a := range actors {
		var actorID int
		err = tx.QueryRow(ctx, "SELECT id FROM actors WHERE name=$1", a).Scan(&actorID)
		if err != nil {
			err = tx.QueryRow(ctx, "INSERT INTO actors(name) VALUES($1) RETURNING id", a).Scan(&actorID)
			if err != nil {
				return 0, err
			}
		}
		_, err = tx.Exec(ctx, "INSERT INTO movies_actors(id_movie, id_actor) VALUES($1,$2)", movieID, actorID)
		if err != nil {
			return 0, err
		}
	}

	// --- Schedule ---
	locations := strings.Split(req.Locations, ",")
	dates := strings.Split(req.Dates, ",")
	times := strings.Split(req.Times, ",")

	for _, loc := range locations {
		var locationID int
		err = tx.QueryRow(ctx, "SELECT id FROM location WHERE name=$1", loc).Scan(&locationID)
		if err == pgx.ErrNoRows {
			err = tx.QueryRow(ctx, "INSERT INTO location(name) VALUES($1) RETURNING id", loc).Scan(&locationID)
			if err != nil {
				return 0, err
			}
		}

		for _, d := range dates {
			dParsed, _ := time.Parse("2006-01-02", d)

			for _, t := range times {
				tParsed, _ := time.Parse("15:04", t) // parse string time

				var timeID int
				err = tx.QueryRow(ctx, "SELECT id FROM time WHERE time=$1", tParsed.Format("15:04:00")).Scan(&timeID)
				if err == pgx.ErrNoRows {
					err = tx.QueryRow(ctx, "INSERT INTO time(time) VALUES($1) RETURNING id", tParsed.Format("15:04:00")).Scan(&timeID)
					if err != nil {
						return 0, err
					}
				}

				// insert schedule
				_, err = tx.Exec(ctx, "INSERT INTO schedule(date, id_movie, id_cinema, id_location, id_time) VALUES($1,$2,$3,$4,$5)",
					dParsed, movieID, 1, locationID, timeID)
				if err != nil {
					return 0, err
				}
			}
		}
	}

	return movieID, nil
}

func (r *AdminRepo) CreateMovieWithRelations(ctx context.Context, m *models.MovieAdminInsert, genreIDs []int, actorIDs []int, combos []models.ScheduleComboAdminInsert) (int, error) {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(ctx)

	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now

	var movieID int
	err = tx.QueryRow(ctx,
		`INSERT INTO movies (title, release_date, duration, synopsis, id_director, create_at, update_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		m.Title, m.ReleaseDate, m.Duration, m.Synopsis, m.IdDirector, m.CreatedAt, m.UpdatedAt,
	).Scan(&movieID)
	if err != nil {
		return 0, fmt.Errorf("insert movie: %w", err)
	}

	// genres
	for _, gid := range genreIDs {
		_, _ = tx.Exec(ctx, `INSERT INTO movies_genres (id_movie, id_genre) VALUES ($1,$2)`, movieID, gid)
	}
	// actors
	for _, aid := range actorIDs {
		_, _ = tx.Exec(ctx, `INSERT INTO movies_actors (id_movie, id_actor) VALUES ($1,$2)`, movieID, aid)
	}
	// schedule
	for _, c := range combos {
		_, _ = tx.Exec(ctx, `INSERT INTO schedule (date, id_movie, id_cinema, id_location, id_time, update_at)
		 VALUES ($1,$2,$3,$4,$5,$6)`, c.Date, movieID, c.IdCinema, c.IdLocation, c.IdTime, now)
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, err
	}

	return movieID, nil
}

func (r *AdminRepo) UpdateMovieFiles(ctx context.Context, movieID int, poster, backdrop string) error {
	_, err := r.DB.Exec(ctx,
		`UPDATE movies SET poster=$1, backdrop=$2 WHERE id=$3`, poster, backdrop, movieID)
	return err
}

func (r *AdminRepo) GetMovieDetail(ctx context.Context, movieID int) (*models.GetMovieDetailUpdate, error) {
	movie := models.GetMovieDetailUpdate{}

	// ambil movie utama
	err := r.DB.QueryRow(ctx, `
		SELECT m.id, m.title, m.poster, m.backdrop,
		       to_char(m.release_date, 'YYYY-MM-DD'),
		       m.duration, m.synopsis, m.id_director, d.name
		FROM movies m
		JOIN directors d ON d.id = m.id_director
		WHERE m.id=$1
	`, movieID).Scan(
		&movie.ID, &movie.Title, &movie.Poster, &movie.Backdrop,
		&movie.ReleaseDate, &movie.Duration, &movie.Synopsis,
		&movie.IdDirector, &movie.Director,
	)
	if err != nil {
		return nil, err
	}

	// genres
	rows, err := r.DB.Query(ctx, `
		SELECT g.id, g.name 
		FROM movies_genres mg 
		JOIN genres g ON g.id = mg.id_genre
		WHERE mg.id_movie=$1
	`, movieID)
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var g models.Genre
			if err := rows.Scan(&g.ID, &g.Name); err == nil {
				movie.Genres = append(movie.Genres, g)
			}
		}
	}

	// actors
	rows2, err := r.DB.Query(ctx, `
		SELECT a.id, a.name 
		FROM movies_actors ma 
		JOIN actors a ON a.id = ma.id_actor
		WHERE ma.id_movie=$1
	`, movieID)
	if err == nil {
		defer rows2.Close()
		for rows2.Next() {
			var a models.Actor
			if err := rows2.Scan(&a.ID, &a.Name); err == nil {
				movie.Actors = append(movie.Actors, a)
			}
		}
	}

	// schedules
	rows3, err := r.DB.Query(ctx, `
		SELECT s.id, to_char(s.date,'YYYY-MM-DD'),
		       c.id, c.name,
		       l.id, l.name,
		       t.id, to_char(t.time,'HH24:MI')
		FROM schedule s
		JOIN cinema c ON c.id = s.id_cinema
		JOIN location l ON l.id = s.id_location
		JOIN time t ON t.id = s.id_time
		WHERE s.id_movie=$1
	`, movieID)
	if err == nil {
		defer rows3.Close()
		for rows3.Next() {
			var sc models.ScheduleDetail
			if err := rows3.Scan(&sc.ID, &sc.Date,
				&sc.CinemaID, &sc.Cinema,
				&sc.LocationID, &sc.Location,
				&sc.TimeID, &sc.Time); err == nil {
				movie.Schedules = append(movie.Schedules, sc)
			}
		}
	}

	return &movie, nil
}

// Update movie
// UpdateMovieAdmin update movie, poster/backdrop path, genre & actor relations, schedules
func (r *AdminRepo) UpdateMovieAdmin(
	ctx context.Context,
	movieID int,
	m models.MovieUpdateAdmin,
	addedGenres, removedGenres, addedActors, removedActors []int,
) error {
	tx, err := r.DB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Update basic movie info
	query := `UPDATE movies 
		SET title=$1, release_date=$2, duration=$3, id_director=$4, synopsis=$5, poster=$6, backdrop=$7 
		WHERE id=$8`
	_, err = tx.Exec(ctx, query,
		m.Title,
		m.ReleaseDate,
		m.Duration,
		m.DirectorID,
		m.Synopsis,
		m.PosterPath,
		m.BackdropPath,
		movieID,
	)
	if err != nil {
		return err
	}

	// Update genres
	if len(addedGenres) > 0 {
		vals := []string{}
		for _, id := range addedGenres {
			vals = append(vals, fmt.Sprintf("(%d,%d)", movieID, id))
		}
		q := fmt.Sprintf("INSERT INTO movies_genres (id_movie, id_genre) VALUES %s ON CONFLICT DO NOTHING", strings.Join(vals, ","))
		_, err = tx.Exec(ctx, q)
		if err != nil {
			return err
		}
	}

	if len(removedGenres) > 0 {
		q := `DELETE FROM movies_genres WHERE id_movie=$1 AND id_genre=ANY($2::int[])`
		_, err = tx.Exec(ctx, q, movieID, removedGenres)
		if err != nil {
			return err
		}
	}

	// Update actors
	if len(addedActors) > 0 {
		vals := []string{}
		for _, id := range addedActors {
			vals = append(vals, fmt.Sprintf("(%d,%d)", movieID, id))
		}
		q := fmt.Sprintf("INSERT INTO movies_actors (id_movie, id_actor) VALUES %s ON CONFLICT DO NOTHING", strings.Join(vals, ","))
		_, err = tx.Exec(ctx, q)
		if err != nil {
			return err
		}
	}

	if len(removedActors) > 0 {
		q := `DELETE FROM movies_actors WHERE id_movie=$1 AND id_actor=ANY($2::int[])`
		_, err = tx.Exec(ctx, q, movieID, removedActors)
		if err != nil {
			return err
		}
	}

	// Update schedules (hapus semua lama + insert baru)
	if len(m.Schedules) > 0 {
		_, err = tx.Exec(ctx, "DELETE FROM schedule WHERE id_movie=$1", movieID)
		if err != nil {
			return err
		}
		for _, s := range m.Schedules {
			_, err = tx.Exec(ctx,
				"INSERT INTO schedule (id_movie, date, id_cinema, id_location, id_time) VALUES ($1,$2,$3,$4,$5)",
				movieID, s.Date, s.CinemaID, s.LocationID, s.TimeID,
			)
			if err != nil {
				return err
			}
		}
	}

	return tx.Commit(ctx)
}
