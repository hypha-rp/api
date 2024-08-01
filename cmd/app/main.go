package main

import (
	"flag"
	"fmt"
	"hypha/api/internal/config"
	"hypha/api/internal/db"
	"hypha/api/internal/http"
	"hypha/api/internal/utils/logging"
	"hypha/api/internal/utils/router"
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

	router, err := router.InitRouter(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize router")
	}

	dbConnWrapper := &db.DBConnWrapper{DB: dbConn}
	http.InitRoutes(router, dbConnWrapper)

	port := cfg.Http.Port
	if err := router.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatal().Err(err).Msg("Failed to run server")
	}
}
