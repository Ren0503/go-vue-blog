package services

import (
	"ren0503/goapi/config"
	"ren0503/goapi/models"
)

func FetchAllTags() ([]models.Tag, error) {
	database := config.ConnectDB()
	var tags []models.Tag
	err := database.Find(&tags).Error

	return tags, err
}
