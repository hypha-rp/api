package http

import (
	"hypha/api/internal/http/db/product"
	"hypha/api/internal/http/db/utils"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

func InitRoutes(router *gin.Engine, dbOperations utils.DatabaseOperations) {
	log.Info().Msg("Initializing routes")
	dbGroup := router.Group("/db")
	product.InitProductRoutes(dbGroup, dbOperations)
}
