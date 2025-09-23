package repositories

import (
	"context"
	"fmt"

	"github.com/Ntisrangga142/API_tickytiz/internals/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ScheduleRepo struct {
	DB *pgxpool.Pool
}

func NewScheduleRepo(db *pgxpool.Pool) *ScheduleRepo {
	return &ScheduleRepo{DB: db}
}

// Schedule fetch schedules with optional filters: date, location, showTime
func (r *ScheduleRepo) Schedule(ctx context.Context, movieID int, date, location, showTime *string) ([]models.Schedule, error) {
	query := `
		SELECT 
			s.id AS schedule_id,
			s.date,
			c.name AS cinema,
			c.logo AS cinema_img,
			c.price,
			l.name AS location,
			t.time AS show_time
		FROM schedule s
		JOIN cinema c   ON s.id_cinema = c.id
		JOIN location l ON s.id_location = l.id
		JOIN time t     ON s.id_time = t.id
		WHERE s.id_movie = $1
		  AND s.delete_at IS NULL
	`
	args := []any{movieID}
	argIdx := 2

	if date != nil && *date != "" {
		query += fmt.Sprintf(" AND s.date::date = $%d", argIdx) // cast ke date agar cocok filter
		args = append(args, *date)
		argIdx++
	}

	if location != nil && *location != "" {
		query += fmt.Sprintf(" AND l.name = $%d", argIdx)
		args = append(args, *location)
		argIdx++
	}

	if showTime != nil && *showTime != "" {
		query += fmt.Sprintf(" AND t.time = $%d", argIdx)
		args = append(args, *showTime)
		argIdx++
	}

	query += " ORDER BY s.date, t.time;"

	rows, err := r.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to load schedule: %w", err)
	}
	defer rows.Close()

	var schedules []models.Schedule
	for rows.Next() {
		var s models.Schedule
		if err := rows.Scan(
			&s.ID,
			&s.Date,
			&s.Cinema,
			&s.CinemaIMG,
			&s.Price,
			&s.Location,
			&s.ShowTime,
		); err != nil {
			return nil, fmt.Errorf("failed to scan schedule row: %w", err)
		}
		schedules = append(schedules, s)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return schedules, nil
}
