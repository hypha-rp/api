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
	if err := dbOperations.Connection().
		Preload("Product1").
		Preload("Product2").
		Where("id = ?", context.Param("id")).
		First(&existingIntegration).Error; err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}
	context.JSON(200, existingIntegration)
}
