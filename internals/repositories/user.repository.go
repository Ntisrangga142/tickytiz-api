package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/Ntisrangga142/API_tickytiz/internals/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return UserRepository{db: db}
}

func (r *UserRepository) GetProfileByID(ctx context.Context, id int) (*models.UserProfile, error) {
	query := `
	SELECT 
		u.id,
		u.firstname,
		u.lastname,
		u.phone,
		u.profileimg,
		u.virtual_account,
		u.point,
		a.email,
		a.role
	FROM users u
	JOIN account a ON u.id = a.id
	WHERE u.id = $1;
	`

	row := r.db.QueryRow(ctx, query, id)
	var user models.UserProfile
	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&user.ProfileImg,
		&user.VirtualAccount,
		&user.Point,
		&user.Email,
		&user.Role,
	)

	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetHistoryByUserID(ctx context.Context, userID int) (models.OrderHistoryResponse, error) {
	query := `
	SELECT 
		o.id AS order_id,
		o.ispaid,
		o.total_price,
		o.qrcode,
		o.name AS order_name,
		o.email AS order_email,
		o.phone AS order_phone,
		pm.name AS payment_method,
		ns.date AS show_date,
		t.time AS show_time,
		c.name AS cinema_name,
		c.logo AS cinema_logo,
		l.name AS location_name,
		m.title AS movie_title,
		m.poster AS movie_poster,
		m.backdrop AS movie_backdrop,
		m.duration,
		m.rating,
		ARRAY_AGG(od.id_seat) AS seats
	FROM orders o
	JOIN payment_method pm ON o.id_payment_method = pm.id
	JOIN schedule ns ON o.id_schedule = ns.id
	JOIN movies m ON ns.id_movie = m.id
	JOIN cinema c ON ns.id_cinema = c.id
	JOIN location l ON ns.id_location = l.id
	JOIN time t ON ns.id_time = t.id
	LEFT JOIN orderdetails od ON o.id = od.id_order
	WHERE o.id_user = $1
	GROUP BY 
		o.id, o.ispaid, o.total_price, o.qrcode, o.name, o.email, o.phone,
		pm.name, ns.date, t.time, c.name, l.name, m.title, m.poster, m.backdrop, m.duration, m.rating, c.logo
	ORDER BY o.id DESC;
	`

	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return models.OrderHistoryResponse{}, err
	}
	defer rows.Close()

	var histories []models.OrderHistory
	for rows.Next() {
		var history models.OrderHistory
		err := rows.Scan(
			&history.OrderID,
			&history.IsPaid,
			&history.TotalPrice,
			&history.QRCode,
			&history.OrderName,
			&history.OrderEmail,
			&history.OrderPhone,
			&history.PaymentMethod,
			&history.ShowDate,
			&history.ShowTime,
			&history.CinemaName,
			&history.CinemaLogo,
			&history.LocationName,
			&history.MovieTitle,
			&history.MoviePoster,
			&history.MovieBackdrop,
			&history.Duration,
			&history.Rating,
			&history.Seats,
		)
		if err != nil {
			return models.OrderHistoryResponse{}, err
		}
		histories = append(histories, history)
	}

	var res models.OrderHistoryResponse
	res.UserID = userID
	res.ListHistory = histories

	return res, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, id int, data models.UpdateProfile) error {
	setClauses := []string{}
	args := []any{}
	argID := 1

	if data.ProfileImg != nil && *data.ProfileImg != "" {
		setClauses = append(setClauses, fmt.Sprintf("profileimg=$%d", argID))
		args = append(args, *data.ProfileImg)
		argID++
	}
	if data.FirstName != nil && *data.FirstName != "" {
		setClauses = append(setClauses, fmt.Sprintf("firstname=$%d", argID))
		args = append(args, *data.FirstName)
		argID++
	}
	if data.LastName != nil && *data.LastName != "" {
		setClauses = append(setClauses, fmt.Sprintf("lastname=$%d", argID))
		args = append(args, *data.LastName)
		argID++
	}
	if data.Phone != nil && *data.Phone != "" {
		setClauses = append(setClauses, fmt.Sprintf("phone=$%d", argID))
		args = append(args, *data.Phone)
		argID++
	}
	if data.VirtualAccount != nil && *data.VirtualAccount != "" {
		setClauses = append(setClauses, fmt.Sprintf("virtual_account=$%d", argID))
		args = append(args, *data.VirtualAccount)
		argID++
	}

	if len(setClauses) == 0 {
		return nil
	}

	setClauses = append(setClauses, "update_at=NOW()")

	query := fmt.Sprintf("UPDATE users SET %s WHERE id=$%d", strings.Join(setClauses, ", "), argID)
	args = append(args, id)

	_, err := r.db.Exec(ctx, query, args...)
	return err
}

func (r *UserRepository) GetVirtualAccountByID(ctx context.Context, userID int) (models.UserVA, error) {
	var user models.UserVA

	row := r.db.QueryRow(ctx, "SELECT virtual_account FROM users WHERE id=$1 LIMIT 1", userID)
	if err := row.Scan(&user.VirtualAccount); err != nil {
		if err == pgx.ErrNoRows {
			return user, nil
		}
		return user, err
	}

	return user, nil
}

func (r *UserRepository) GetPasswordByID(ctx context.Context, userID int) (string, error) {
	var password string
	err := r.db.QueryRow(ctx, `SELECT password FROM account WHERE id = $1`, userID).Scan(&password)
	if err != nil {
		return "", err
	}
	return password, nil
}

func (r *UserRepository) UpdatePassword(ctx context.Context, userID int, newPassword string) error {
	cmd, err := r.db.Exec(ctx, `UPDATE account SET password = $1 WHERE id = $2`, newPassword, userID)
	if err != nil {
		return err
	}
	if cmd.RowsAffected() == 0 {
		return errors.New("no rows updated")
	}
	return nil
}
