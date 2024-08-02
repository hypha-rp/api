package product

import (
	"hypha/api/internal/db/ops"
	"hypha/api/internal/db/tables"

	"github.com/gin-gonic/gin"
)

func InitProductRoutes(router *gin.RouterGroup, dbOperations ops.DatabaseOperations) {
	router.POST("/product", func(context *gin.Context) {
		CreateProduct(dbOperations, context)
	})
	router.GET("/product/:id", func(context *gin.Context) {
		GetProduct(dbOperations, context)
	})
	router.GET("/product/:id/integrations", func(context *gin.Context) {
		GetProductIntegrations(dbOperations, context)
	})
	router.GET("/products", func(context *gin.Context) {
		GetAllProducts(dbOperations, context)
	})
}

func CreateProduct(dbOperations ops.DatabaseOperations, context *gin.Context) {
	var newProduct tables.Product
	newProduct.ID = ops.GenerateUniqueID()
	ops.CreateResource(dbOperations, context, &newProduct)
}

func GetProduct(dbOperations ops.DatabaseOperations, context *gin.Context) {
	var existingProduct tables.Product
	ops.GetResource(dbOperations, context, &existingProduct, "id", "Product")
}

func GetProductIntegrations(dbOperations ops.DatabaseOperations, context *gin.Context) {
	var integrations []tables.Integration
	productID := context.Param("id")
	if err := dbOperations.Connection().
		Where("product_id1 = ? OR product_id2 = ?", productID, productID).
		Preload("Product1").
		Preload("Product2").
		Find(&integrations).Error; err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}
	context.JSON(200, integrations)
}

func GetAllProducts(dbOperations ops.DatabaseOperations, context *gin.Context) {
	var products []tables.Product
	name := context.Query("name")
	query := dbOperations.Connection()
	if name != "" {
		query = query.Where("full_name ILIKE ?", "%"+name+"%")
	}
	if err := query.Find(&products).Error; err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}
	context.JSON(200, products)
}
