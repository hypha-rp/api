package http

import (
	"hypha/api/internal/db"
	db_group "hypha/api/internal/http/db"
	"hypha/api/internal/http/report"
	"hypha/api/internal/http/results"
	"hypha/api/internal/utils/logging"

	"github.com/gin-gonic/gin"
)

var log = logging.Logger

func InitRoutes(router *gin.Engine, dbOperations db.DatabaseOperations) {
	log.Info().Msg("Initializing routes")

	dbGroup := router.Group("/db")
	db_group.InitProductRoutes(dbGroup, dbOperations)
	db_group.InitIntegrationRoutes(dbGroup, dbOperations)

	resultsGroup := router.Group("/results")
	results.InitResultsRoutes(resultsGroup, dbOperations)

	reportGroup := router.Group("/report")
	report.InitReportRoutes(reportGroup, dbOperations)

	log.Info().Msg("Routes initialized")
}
