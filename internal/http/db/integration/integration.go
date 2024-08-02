package integration

import (
	"hypha/api/internal/db/ops"
	"hypha/api/internal/db/tables"

	"github.com/gin-gonic/gin"
)

func InitIntegrationRoutes(router *gin.RouterGroup, dbOperations ops.DatabaseOperations) {
	router.POST("/integration", func(context *gin.Context) {
		CreateIntegration(dbOperations, context)
	})
	router.GET("/integration/:id", func(context *gin.Context) {
		GetIntegration(dbOperations, context)
	})
}

func CreateIntegration(dbOperations ops.DatabaseOperations, context *gin.Context) {
	var newIntegration tables.Integration
	newIntegration.ID = ops.GenerateUniqueID()
	ops.CreateResource(dbOperations, context, &newIntegration)
}

func GetIntegration(dbOperations ops.DatabaseOperations, context *gin.Context) {
	var existingIntegration tables.Integration
	ops.GetResource(dbOperations, context, &existingIntegration, "id", "Product")
}
