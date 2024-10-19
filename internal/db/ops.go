package db

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-orm/gorm"
	"github.com/google/uuid"
)

type DatabaseOperations interface {
	Connection() *gorm.DB
	Create(value interface{}) error
	First(out interface{}, where ...interface{}) error
}

func CreateResource(dbOps DatabaseOperations, context *gin.Context, resource interface{}) {
	if err := context.ShouldBindJSON(resource); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := dbOps.Create(resource); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	context.JSON(http.StatusOK, resource)
}

func GetResource(dbOps DatabaseOperations, context *gin.Context, resource interface{}, idParam string, resourceName string) {
	resourceID := context.Param(idParam)

	if err := dbOps.First(resource, "id = ?", resourceID); err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": resourceName + " not found"})
		return
	}

	context.JSON(http.StatusOK, resource)
}

func GenerateUniqueID() string {
	return uuid.New().String()
}
