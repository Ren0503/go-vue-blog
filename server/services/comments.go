package services

import (
	"ren0503/goapi/config"
	"ren0503/goapi/models"
)

func DeleteComment(condition interface{}) error {
	database := config.GetDB()
	err := database.Where(condition).Delete(models.Comment{}).Error

	return err
}

func FetchCommentById(id int, includes ...bool) models.Comment {
	includeUser := false
	if len(includes) > 0 {
		includeUser = includes[0]
	}

	includeArticle := false
	if len(includes) > 1 {
		includeArticle = includes[1]
	}

	database := config.GetDB()
	var comment models.Comment

	if includeArticle && includeUser {
		database.Preload("Article").Preload("User").Find(&comment, id)
	} else if includeUser {
		database.Preload("User").Find(&comment, id)
	} else if includeArticle {
		database.Preload("Article").Find(&comment, id)
	} else {
		database.Find(&comment, id)
	}

	return comment
}
