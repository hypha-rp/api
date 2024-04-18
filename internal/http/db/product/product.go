package product

import (
	"hypha/api/internal/db/tables"
	"hypha/api/internal/http/db/utils"

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
	utils.CreateResource(dbOperations, context, &newProduct)
}

func GetProduct(dbOperations utils.DatabaseOperations, context *gin.Context) {
	var existingProduct tables.Product
	utils.GetResource(dbOperations, context, &existingProduct, "id", "Product")
}
