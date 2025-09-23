package routers

import (
	"github.com/Ntisrangga142/API_tickytiz/internals/handlers"
	"github.com/Ntisrangga142/API_tickytiz/internals/repositories"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func InitScheduleRoute(router *gin.Engine, db *pgxpool.Pool, rdb *redis.Client) {
	repoSchedule := repositories.NewScheduleRepo(db)
	handlerSchedule := handlers.NewScheduleHandler(repoSchedule, rdb)

	repoSeat := repositories.NewSeatRepository(db)
	handlerSeat := handlers.NewSeatHandler(repoSeat)

	schedule := router.Group("/schedule")
	schedule.GET("/:id", handlerSchedule.ScheduleMovie)
	schedule.GET("/seat/:id", handlerSeat.GetSoldSeats)
}
