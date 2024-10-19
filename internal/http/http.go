package http

import (
	"hypha/api/internal/db"
	"hypha/api/internal/http/routes"
	"hypha/api/internal/utils/logging"

	"github.com/gin-gonic/gin"
)

var log = logging.Logger

func InitRoutes(router *gin.Engine, dbOps db.DatabaseOperations) {
	log.Info().Msg("Initializing routes")

	dbGroup := router.Group("/db")
	routes.InitProductRoutes(dbGroup, dbOps)
	routes.InitRelationRoutes(dbGroup, dbOps)
	routes.InitRuleRoutes(dbGroup, dbOps)

	resultsGroup := router.Group("/results")
	routes.InitResultsRoutes(resultsGroup, dbOps)

	log.Info().Msg("Routes initialized")
}
