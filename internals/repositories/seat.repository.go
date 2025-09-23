package repositories

import (
	"context"

	"github.com/Ntisrangga142/API_tickytiz/internals/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SeatRepository struct {
	DB *pgxpool.Pool
}

func NewSeatRepository(db *pgxpool.Pool) *SeatRepository {
	return &SeatRepository{DB: db}
}

func (r *SeatRepository) GetSoldSeats(ctx context.Context, scheduleID int) (models.Seat, error) {
	query := `
		SELECT od.id_seat
		FROM orders o
		JOIN orderdetails od ON o.id = od.id_order
		WHERE o.id_schedule = $1;
	`

	rows, err := r.DB.Query(ctx, query, scheduleID)
	if err != nil {
		return models.Seat{}, err
	}
	defer rows.Close()

	seats := models.Seat{ID: []string{}}
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			return models.Seat{}, err
		}
		seats.ID = append(seats.ID, s)
	}

	return seats, nil
}
