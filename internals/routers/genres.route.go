package routers

import (
	"github.com/Ntisrangga142/API_tickytiz/internals/handlers"
	"github.com/Ntisrangga142/API_tickytiz/internals/repositories"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func InitRouteGenres(r *gin.Engine, db *pgxpool.Pool) {
	genreRepo := repositories.NewGenreRepo(db)
	genreHandler := handlers.NewGenreHandler(genreRepo)

	genres := r.Group("/genres")
	{
		genres.GET("", genreHandler.GetGenres)
	}
}
