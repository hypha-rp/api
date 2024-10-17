package db

import (
	"bytes"
	"encoding/json"
	"hypha/api/internal/db"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func InitIntegrationRoutes(router *gin.RouterGroup, dbOperations db.DatabaseOperations) {
	router.POST("/integration", func(context *gin.Context) {
		CreateIntegration(dbOperations, context)
	})
	router.GET("/integration/:id", func(context *gin.Context) {
		GetIntegration(dbOperations, context)
	})
}

func CreateIntegration(dbOperations db.DatabaseOperations, context *gin.Context) {
	var requestBody map[string]string

	body, err := io.ReadAll(context.Request.Body)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	context.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	if err := json.Unmarshal(body, &requestBody); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	productID1 := requestBody["productID1"]
	productID2 := requestBody["productID2"]

	if productID1 == "" || productID2 == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Product IDs cannot be empty"})
		return
	}

	if productID1 == productID2 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Cannot create integration for the same product"})
		return
	}

	var newIntegration db.Integration
	newIntegration.ID = db.GenerateUniqueID()
	db.CreateResource(dbOperations, context, &newIntegration)
}

func GetIntegration(dbOperations db.DatabaseOperations, context *gin.Context) {
	var existingIntegration db.Integration
	if err := dbOperations.Connection().
		Preload("Product1").
		Preload("Product2").
		Where("id = ?", context.Param("id")).
		First(&existingIntegration).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, existingIntegration)
}
