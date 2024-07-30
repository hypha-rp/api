package http

import (
	"hypha/api/internal/db/ops"
	"hypha/api/internal/http/db/integration"
	"hypha/api/internal/http/db/product"
	"hypha/api/internal/http/report"
	"hypha/api/internal/utils/logging"

	"github.com/gin-gonic/gin"
)

var log = logging.Logger

func InitRoutes(router *gin.Engine, dbOperations ops.DatabaseOperations) {
	log.Info().Msg("Initializing routes")

	dbGroup := router.Group("/db")
	product.InitProductRoutes(dbGroup, dbOperations)
	integration.InitIntegrationRoutes(dbGroup, dbOperations)

	reportGroup := router.Group("/report")
	report.InitReportRoutes(reportGroup, dbOperations)

	log.Info().Msg("Routes initialized")
}
