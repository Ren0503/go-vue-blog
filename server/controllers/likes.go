package controllers

import (
	"fmt"
	"net/http"
	"ren0503/goapi/config"
	"ren0503/goapi/dtos"
	"ren0503/goapi/middleware"
	"ren0503/goapi/models"
	"ren0503/goapi/services"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RegisterLikeRoutes(router *gin.RouterGroup) {
	router.Use(middleware.EnforceAuthenticatedMiddleware())
	{
		router.POST("/articles/:slug/likes", LikeArticle)
		router.DELETE("/articles/:slug/likes", DislikeArticle)
		router.GET("/likes", MyLikes)
	}
}

func MyLikes(c *gin.Context) {
	pageSizeStr := c.Query("page_size")
	pageStr := c.Query("page")
	currentUserId := c.MustGet("currentUserId").(uint)
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		pageSize = 5
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}
	var result []models.Like
	var ids []uint
	database := config.GetDB()

	database.Select("article_id").Where("likes.user_id = ?", currentUserId).Find(&result).Pluck("article_id", &ids)

	var articles []models.Article
	var likedArticles = 0
	database.Table("articles").Where("id in (?)", ids).Count(&likedArticles)
	database.Where("id in (?)", ids).
		Order("created_at desc").Offset((page - 1) * pageSize).
		Limit(pageSize).Preload("User").Preload("Tags").
		Preload("Categories").Find(&articles)

	for index := range articles {
		articles[index].CommentsCount = database.Model(&articles[index]).Association("Comments").Count()
	}

	c.JSON(http.StatusOK, dtos.CreatedArticlePagedResponse(c.Request, articles, page, pageSize, likedArticles))
}

func LikeArticle(c *gin.Context) {
	slug := c.Param("slug")

	database := config.GetDB()
	var article models.Article
	err := database.Model(&models.Article{}).Where("slug = ?", slug).Select([]string{"id", "title"}).First(&article).Error
	if err != nil {
		c.JSON(http.StatusNotFound, dtos.CreateDetailedErrorDto("articles", err))
		return
	}

	user := c.MustGet("currentUser").(models.User)

	if !services.IsLikedBy(&article, user) {
		like := models.Like{
			ArticleId: article.ID,
			UserId:    user.ID,
		}

		if err := services.CreateOne(&like); err != nil {
			c.JSON(http.StatusUnprocessableEntity, dtos.CreateDetailedErrorDto("database", err))
			return
		}

		c.JSON(http.StatusOK, dtos.GetSuccessDto(fmt.Sprintf("You liked the article \"%v\" successfully", article.Title)))
	} else {
		c.JSON(http.StatusForbidden, dtos.GetErrorDto("You have already liked this article"))
	}
}

func DislikeArticle(c *gin.Context) {
	slug := c.Param("slug")
	database := config.GetDB()
	var result struct {
		Id string
	}

	err := database.Table("articles").Select("id").Where("slug = ?", slug).Scan(&result).Error
	if err != nil {
		c.JSON(http.StatusNotFound, dtos.CreateDetailedErrorDto("articles", err))
		return
	}

	user := c.MustGet("currentUser").(models.User)
	var like models.Like
	err = database.Model(models.Like{}).Where("user_id = ? AND article_id = ?", user.ID, result.Id).First(&like).Error
	if err == nil {
		database.Delete(&like)
		c.JSON(http.StatusOK, dtos.GetSuccessDto("Article disliked successfully"))
	} else {
		c.JSON(http.StatusForbidden, dtos.GetSuccessDto("You were not liking this article, so you can not perform this operation"))
	}
}
