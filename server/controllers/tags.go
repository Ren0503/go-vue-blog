package controllers

import (
	"net/http"
	"ren0503/goapi/dtos"
	"ren0503/goapi/services"

	"github.com/gin-gonic/gin"
)

func RegisterTagRoutes(router *gin.RouterGroup) {
	router.GET("", TagList)
}

func TagList(c *gin.Context) {
	tags, err := services.FetchAllTags()
	if err != nil {
		c.JSON(http.StatusNotFound, dtos.CreateDetailedErrorDto("db_error", err))
		return
	}

	c.JSON(http.StatusOK, dtos.CreateTagListDto(tags))
}
