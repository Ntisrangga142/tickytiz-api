package routers

import (
	"github.com/Ntisrangga142/API_tickytiz/internals/handlers"
	repo "github.com/Ntisrangga142/API_tickytiz/internals/repositories"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func InitMovieRoutes(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	repo := repo.NewMovieRepo(db)
	handler := handlers.NewMovieHandler(repo, rdb)

	movie := router.Group("/movies")

	movie.GET("/", handler.FilteredMovies)
	movie.GET("/:id", handler.MovieDetail)
	movie.GET("/popular", handler.PopularMovies)
	movie.GET("/upcoming", handler.UpcomingMovies)
}
