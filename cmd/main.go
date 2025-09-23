package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	_ "github.com/Ntisrangga142/API_tickytiz/docs"
	"github.com/Ntisrangga142/API_tickytiz/internals/configs"
	"github.com/Ntisrangga142/API_tickytiz/internals/routers"
	"github.com/joho/godotenv"
)

// @title Tickytiz API
// @version 1.0
// @description API untuk sistem tiket bioskop

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	// Load Environment Variable
	if err := godotenv.Load(); err != nil {
		log.Println("Failed to load env\nCause:", err.Error())
	}

	// Init Database
	db, err := configs.InitDB()
	if err != nil {
		log.Fatal("Failed to connect DB:", err.Error())
	}
	defer db.Close()

	rdb := configs.InitRedis()

	// Init Router
	router := routers.InitRouter(db, rdb)
	if runtime.GOOS == "windows" {
		link := fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
		router.Run(link)
	} else {
		router.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))
	}
}
