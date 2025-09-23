package routers

import (
	"github.com/Ntisrangga142/API_tickytiz/internals/handlers"
	"github.com/Ntisrangga142/API_tickytiz/internals/middlewares"
	"github.com/Ntisrangga142/API_tickytiz/internals/repositories"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func InitUserRoute(r *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	userRepo := repositories.NewUserRepository(db)
	userHandler := handlers.NewUserHandler(&userRepo, rdb)

	userGroup := r.Group("/user")
	userGroup.GET("", middlewares.Authentication, middlewares.Authorization("user"), userHandler.GetProfile)
	userGroup.PATCH("", middlewares.Authentication, middlewares.Authorization("user"), userHandler.UpdateProfile)
	userGroup.PATCH("/password", middlewares.Authentication, middlewares.Authorization("user"), userHandler.ChangePassword)
	userGroup.GET("/va", middlewares.Authentication, middlewares.Authorization("user"), userHandler.GetVirtualAccountHandler)
	userGroup.GET("/history", middlewares.Authentication, middlewares.Authorization("user"), userHandler.GetHistory)

}
