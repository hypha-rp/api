package main

import (
	"flag"
	"hypha/api/internal/config"
	"hypha/api/internal/db"
	"hypha/api/internal/http"
	"hypha/api/internal/utils/logging"

	"github.com/gin-gonic/gin"
)

var log = logging.Logger

func main() {
	configPath := flag.String("config", "config.yaml", "Path to configuration file")
	flag.Parse()

	cfg, err := config.ReadConfig(*configPath)
	if err != nil {
		log.Fatal().Msgf("Could not read config: %v", err)
	}

	dbConn, err := db.Connect(cfg)
	if err != nil {
		log.Fatal().Msgf("Could not connect to database: %v", err)
	}

	if err := db.AutoMigrate(dbConn); err != nil {
		log.Fatal().Msgf("Could not migrate tables: %v", err)
	}

	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(logging.GinLogger())

	dbConnWrapper := &db.DBConnWrapper{DB: dbConn}
	http.InitRoutes(router, dbConnWrapper)

	router.Run(":8081")
}
