package handlers

import (
	"hypha/api/internal/db"

	"github.com/gin-gonic/gin"
)

func CreateProduct(dbOps db.DatabaseOperations, context *gin.Context) {
	var newProduct db.Product
	newProduct.ID = db.GenerateUniqueID()
	db.CreateResource(dbOps, context, &newProduct)
}

func GetProduct(dbOps db.DatabaseOperations, context *gin.Context) {
	var existingProduct db.Product
	db.GetResource(dbOps, context, &existingProduct, "id", "Product")
}

func GetProductIntegrations(dbOps db.DatabaseOperations, context *gin.Context) {
	var integrations []db.Integration
	productID := context.Param("id")
	if err := dbOps.Connection().
		Where("product_id1 = ? OR product_id2 = ?", productID, productID).
		Preload("Product1").
		Preload("Product2").
		Find(&integrations).Error; err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}
	context.JSON(200, integrations)
}

func GetAllProducts(dbOps db.DatabaseOperations, context *gin.Context) {
	var products []db.Product
	name := context.Query("name")
	query := dbOps.Connection()
	if name != "" {
		query = query.Where("full_name ILIKE ? OR short_name ILIKE ?", "%"+name+"%", "%"+name+"%")
	}
	if err := query.Find(&products).Error; err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}
	context.JSON(200, products)
}
