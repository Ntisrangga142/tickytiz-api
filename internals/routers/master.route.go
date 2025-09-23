package routers

import (
	"github.com/Ntisrangga142/API_tickytiz/internals/handlers"
	"github.com/Ntisrangga142/API_tickytiz/internals/middlewares"
	"github.com/Ntisrangga142/API_tickytiz/internals/repositories"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func InitMasterRoute(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	repo := repositories.NewMasterRepo(db)
	handler := handlers.NewMasterHandler(repo)

	master := router.Group("/master")
	{
		master.GET("/directors", middlewares.Authentication, middlewares.Authorization("admin"), handler.GetDirectors)
		master.GET("/actors", middlewares.Authentication, middlewares.Authorization("admin"), handler.GetActors)
		master.GET("/locations", middlewares.Authentication, middlewares.Authorization("admin"), handler.GetLocations)
		master.GET("/times", middlewares.Authentication, middlewares.Authorization("admin"), handler.GetTimes)
		master.GET("/cinemas", middlewares.Authentication, middlewares.Authorization("admin"), handler.GetCinemas)
	}
}
