package product

import (
	"hypha/api/internal/db/tables"
	"hypha/api/internal/http/db/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetupProductRoutes(router *gin.RouterGroup, dbOperations utils.DatabaseOperations) {
	router.POST("/product", func(context *gin.Context) {
		CreateProduct(dbOperations, context)
	})
	router.GET("/product/:id", func(context *gin.Context) {
		GetProduct(dbOperations, context)
	})
}

func CreateProduct(dbOperations utils.DatabaseOperations, context *gin.Context) {
	var newProduct tables.Product
	if err := context.ShouldBindJSON(&newProduct); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := dbOperations.Create(&newProduct); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, newProduct)
}

func GetProduct(dbOperations utils.DatabaseOperations, context *gin.Context) {
	var existingProduct tables.Product
	productID := context.Param("id")

	if err := dbOperations.First(&existingProduct, "id = ?", productID); err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}

	context.JSON(http.StatusOK, existingProduct)
}
