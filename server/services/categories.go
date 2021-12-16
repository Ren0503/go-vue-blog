package services

import (
	"ren0503/goapi/config"
	"ren0503/goapi/models"
)

func FetchAllCategories() ([]models.Category, error) {
	database := config.ConnectDB()

	var categories []models.Category
	err := database.Find(&categories).Error

	return categories, err
}
