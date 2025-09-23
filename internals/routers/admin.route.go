package routers

import (
	"github.com/Ntisrangga142/API_tickytiz/internals/handlers"
	"github.com/Ntisrangga142/API_tickytiz/internals/middlewares"
	"github.com/Ntisrangga142/API_tickytiz/internals/repositories"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func InitAdminRoute(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	repo := repositories.NewAdminRepo(db)
	handler := handlers.NewAdminHandler(repo, rdb)

	middlewares.InitRedis(rdb)

	admin := router.Group("/admin")
	admin.GET("", middlewares.Authentication, middlewares.Authorization("admin"), handler.GetAllMovie)
	admin.GET(":id", middlewares.Authentication, middlewares.Authorization("admin"), handler.GetMovieDetailUpdate)
	admin.POST("", middlewares.Authentication, middlewares.Authorization("admin"), handler.CreateMovie)
	admin.PATCH(":id", middlewares.Authentication, middlewares.Authorization("admin"), handler.UpdateMovieAdmin)
	admin.DELETE(":id", middlewares.Authentication, middlewares.Authorization("admin"), handler.DeleteMovie)
}
