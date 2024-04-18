package main

import (
	"flag"
	"hypha/api/internal/config"
	"hypha/api/internal/db"
	"hypha/api/internal/db/tables"
	"hypha/api/internal/http"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	configPath := flag.String("config", "config.yaml", "Path to configuration file")
	flag.Parse()

	cfg, err := config.ReadConfig(*configPath)
	if err != nil {
		log.Fatal("Could not read config: ", err)
	}

	dbConn, err := db.Connect(cfg)
	if err != nil {
		log.Fatal("Could not connect to database: ", err)
	}

	if err := db.AutoMigrate(dbConn, &tables.Product{}, &tables.Repo{}); err != nil {
		log.Fatalf("Could not migrate tables: %v", err)
	}

	router := gin.Default()
	dbConnWrapper := &db.DBConnWrapper{DB: dbConn}
	http.SetupRoutes(router, dbConnWrapper)

	router.Run("localhost:5052")
}
