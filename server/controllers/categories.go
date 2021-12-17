package controllers

import (
	"net/http"
	"ren0503/goapi/dtos"
	"ren0503/goapi/services"

	"github.com/gin-gonic/gin"
)

func RegisterCategoryRoutes(router *gin.RouterGroup) {
	router.GET("", ListCategories)
}
func ListCategories(c *gin.Context) {
	categories, err := services.FetchAllCategories()
	if err != nil {
		c.JSON(http.StatusNotFound, dtos.CreateDetailedErrorDto("db_error", err))
		return
	}

	c.JSON(http.StatusOK, dtos.CreateCategoryListDto(categories))
}
