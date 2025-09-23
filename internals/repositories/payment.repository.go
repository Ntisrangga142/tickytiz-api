package repositories

import (
	"context"
	"log"

	"github.com/Ntisrangga142/API_tickytiz/internals/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PaymentMethodRepository struct {
	db *pgxpool.Pool
}

func NewPaymentMethodRepository(db *pgxpool.Pool) *PaymentMethodRepository {
	return &PaymentMethodRepository{db: db}
}

func (r *PaymentMethodRepository) GetAll(ctx context.Context) ([]models.PaymentMethod, error) {
	query := `SELECT id, name, logo FROM payment_method`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		log.Println("Failed to fetch payment methods:", err)
		return nil, err
	}
	defer rows.Close()

	var methods []models.PaymentMethod
	for rows.Next() {
		var pm models.PaymentMethod
		if err := rows.Scan(&pm.ID, &pm.Name, &pm.Logo); err != nil {
			log.Println("Failed to scan payment method:", err)
			return nil, err
		}
		methods = append(methods, pm)
	}

	return methods, nil
}
