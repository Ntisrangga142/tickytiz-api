package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"

	docs "github.com/Ntisrangga142/API_tickytiz/docs"
	"github.com/Ntisrangga142/API_tickytiz/internals/middlewares"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter(db *pgxpool.Pool, rdb *redis.Client) *gin.Engine {
	router := gin.Default()
	router.Use(middlewares.CORSMiddleware())

	router.StaticFS("/profile", gin.Dir("./public/profiles", false))
	router.StaticFS("/poster", gin.Dir("./public/movie/poster", false))
	router.StaticFS("/backdrop", gin.Dir("./public/movie/backdrop", false))
	router.StaticFS("/payment-logo", gin.Dir("./public/payment_method", false))
	router.StaticFS("/cinema", gin.Dir("./public/cinema", false))

	InitAuthRoutes(router, db, rdb)
	InitMovieRoutes(router, db, rdb)
	InitScheduleRoute(router, db, rdb)
	InitOrderRoute(router, db, rdb)
	InitUserRoute(router, db, rdb)
	InitAdminRoute(router, db, rdb)
	InitRouteGenres(router, db)
	InitPayment(router, db, rdb)
	InitMasterRoute(router, db, rdb)

	docs.SwaggerInfo.BasePath = "/"
	router.GET("/tickytiz/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
