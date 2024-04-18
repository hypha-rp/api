package productsrepo

import (
	"hypha/api/internal/db/tables"
	"hypha/api/internal/http/db/utils"

	"github.com/gin-gonic/gin"
)

func SetupProductsRepoRoutes(router *gin.RouterGroup, dbOperations utils.DatabaseOperations) {
	router.POST("/products_repo", func(context *gin.Context) {
		CreateProductsRepo(dbOperations, context)
	})
	router.GET("/products_repo/:id", func(context *gin.Context) {
		GetProductsRepo(dbOperations, context)
	})
}

func CreateProductsRepo(dbOperations utils.DatabaseOperations, context *gin.Context) {
	var newProductsRepo tables.ProductsRepo
	utils.CreateResource(dbOperations, context, &newProductsRepo)
}

func GetProductsRepo(dbOperations utils.DatabaseOperations, context *gin.Context) {
	var existingProductsRepo tables.ProductsRepo
	utils.GetResource(dbOperations, context, &existingProductsRepo, "id", "ProductsRepo")
}
