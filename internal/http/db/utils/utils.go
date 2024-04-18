package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type DatabaseOperations interface {
	Create(value interface{}) error
	First(out interface{}, where ...interface{}) error
}

func CreateResource(dbOperations DatabaseOperations, context *gin.Context, resource interface{}) {
	if err := context.ShouldBindJSON(resource); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := dbOperations.Create(resource); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, resource)
}

func GetResource(dbOperations DatabaseOperations, context *gin.Context, resource interface{}, idParam string, resourceName string) {
	resourceID := context.Param(idParam)

	if err := dbOperations.First(resource, "id = ?", resourceID); err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": resourceName + " not found"})
		return
	}

	context.JSON(http.StatusOK, resource)
}
