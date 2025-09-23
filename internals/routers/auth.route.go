package routers

import (
	"github.com/Ntisrangga142/API_tickytiz/internals/handlers"
	"github.com/Ntisrangga142/API_tickytiz/internals/repositories"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func InitAuthRoutes(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	repo := repositories.NewAuthRepo(db, rdb)
	handler := handlers.NewAuthHandler(repo)

	auth := router.Group("/auth")
	auth.POST("", handler.Login)
	auth.DELETE("", handler.Logout)
	auth.POST("/register", handler.Register)
}
