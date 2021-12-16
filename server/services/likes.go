package services

import (
	"ren0503/goapi/config"
	"ren0503/goapi/models"
)

func LikesCount(article *models.Article) uint {
	database := config.ConnectDB()
	var count uint

	database.Model(&models.Like{}).Where(models.Like{
		ArticleId: article.ID,
	}).Count(&count)

	return count
}
