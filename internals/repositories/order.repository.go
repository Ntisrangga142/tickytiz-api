package repositories

import (
	"context"

	"github.com/Ntisrangga142/API_tickytiz/internals/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepo struct {
	DB *pgxpool.Pool
}

func NewOrderRepo(db *pgxpool.Pool) *OrderRepo {
	return &OrderRepo{DB: db}
}

func (r *OrderRepo) CreateOrder(ctx context.Context, req models.OrderRequest, userID int) (*models.OrderResponse, error) {
	query := `
		INSERT INTO orders (ispaid, total_price, qrcode, name, email, phone, id_schedule, id_payment_method, id_user)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, name, email, phone, qrcode;
	`

	var id int
	var name, email, phone, qrcode string
	err := r.DB.QueryRow(ctx, query,
		req.IsPaid,
		req.TotalPrice,
		req.QRCode,
		req.Name,
		req.Email,
		req.Phone,
		req.ScheduleID,
		req.PaymentMethodID,
		userID,
	).Scan(&id, &name, &email, &phone, &qrcode)

	if err != nil {
		return nil, err
	}

	return &models.OrderResponse{ID: id, Name: name, Email: email, Phone: phone, QRCode: qrcode, Seat: []string{}}, nil
}

func (r *OrderRepo) CreateOrderDetails(ctx context.Context, req models.OrderRequest, orderID int) ([]string, error) {
	query := `INSERT INTO orderdetails (id_order, id_seat) VALUES ($1, $2) RETURNING id_seat`

	var insertedSeats []string

	for _, seat := range req.Seat {
		var seatID string
		err := r.DB.QueryRow(ctx, query, orderID, seat).Scan(&seatID)
		if err != nil {
			return nil, err
		}
		insertedSeats = append(insertedSeats, seatID)
	}

	return insertedSeats, nil
}
