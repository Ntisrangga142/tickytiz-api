package routers

import (
	"github.com/Ntisrangga142/API_tickytiz/internals/handlers"
	"github.com/Ntisrangga142/API_tickytiz/internals/middlewares"
	repo "github.com/Ntisrangga142/API_tickytiz/internals/repositories"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func InitOrderRoute(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	repo := repo.NewOrderRepo(db)
	handler := handlers.NewOrderHandler(repo, rdb)

	order := router.Group("/order")
	order.POST("", middlewares.Authentication, middlewares.Authorization("user"), handler.CreateOrder)

}
