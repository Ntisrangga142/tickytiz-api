package routers

import (
	"github.com/Ntisrangga142/API_tickytiz/internals/handlers"
	repo "github.com/Ntisrangga142/API_tickytiz/internals/repositories"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func InitPayment(r *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {

	repo := repo.NewPaymentMethodRepository(db)
	handler := handlers.NewPaymentMethodHandler(repo)

	group := r.Group("/payment-methods")
	{
		group.GET("/", handler.GetAll)
	}
}
