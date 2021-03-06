package controllers

import (
	"fmt"
	"net/http"
	"ren0503/goapi/config"
	"ren0503/goapi/dtos"
	"ren0503/goapi/middleware"
	"ren0503/goapi/models"
	"ren0503/goapi/services"
	"ren0503/goapi/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RegisterUserSubscriptionRoutes(router *gin.RouterGroup) {
	router.Use(middleware.EnforceAuthenticatedMiddleware())
	{
		router.GET("/followers", ListFollowers)
		router.GET("/by_id/:id/followers", ListFollowers)

		router.GET("/following", ListFollowing)
		router.GET("/by_id/:id/following", ListFollowing)

		router.GET("/user_subscriptions", ListSubscriptions)

		router.POST("/by_id/:id/followers", SubscribeToUser)

		router.DELETE("/by_id/:id/followers", UnsubscribeFromUser)
	}
}

func ListFollowers(c *gin.Context) {
	userId, page, pageSize, offset, err := getInfo(c)
	if err != nil {
		return
	}

	database := config.GetDB()
	var followerIds []uint
	var followersCount int
	var followers []models.User

	database.Table("subscriptions").Where("following_id = ?", userId).Count(&followersCount)

	database.Table("subscriptions").Select("follower_id").Where("following_id = ?", userId).Offset(offset).Limit(pageSize).Pluck("follower_id", &followerIds)

	database.Table("users").Select("id,username").Where("id in (?)", followerIds).Find(&followers)
	c.JSON(http.StatusOK, dtos.CreateUserSubscriptionPageResponse(c.Request, followers, "followers", page, pageSize, followersCount))
}

func ListFollowing(c *gin.Context) {
	userId, page, pageSize, offset, err := getInfo(c)
	if err != nil {
		return
	}
	database := config.GetDB()
	var followingIds []uint
	var followingCount int
	var following []models.User

	// TODO it also counts where delete_at is not null (soft deleted)
	database.Table("subscriptions").Where("follower_id = ?", userId).Count(&followingCount)

	database.Table("subscriptions").Select("following_id").Where("follower_id = ?", userId).Offset(offset).Limit(pageSize).Pluck("following_id", &followingIds)

	database.Table("users").Select("id,username").Where("id in (?)", followingIds).Find(&following)
	c.JSON(http.StatusOK, dtos.CreateUserSubscriptionPageResponse(c.Request, following, "following", page, pageSize, followingCount))
}

func ListSubscriptions(c *gin.Context) {
	userId, page, pageSize, offset, err := getInfo(c)
	if err != nil {
		return
	}
	database := config.GetDB()

	var userSubscriptionsCount int
	var following []models.User
	var followers []models.User
	var userRelations []models.User
	var followingIds []uint
	var followerIds []uint

	var userSubscriptions []models.Subscription

	database.Model(&models.Subscription{}).Where(&models.Subscription{FollowingId: userId}).Or(&models.Subscription{FollowerId: userId}).Count(&userSubscriptionsCount)

	database.Table("subscriptions").Where("following_id = ?", userId).Or("follower_id = ?", userId).Offset(offset).Limit(pageSize).Scan(&userSubscriptions)

	for index := range userSubscriptions {
		if userSubscriptions[index].FollowerId == userId {
			followingIds = append(followingIds, userSubscriptions[index].FollowingId)
		} else {
			followerIds = append(followerIds, userSubscriptions[index].FollowerId)
		}
	}
	database.Table("users").Select("id,username").Where("id in (?)", append(followingIds, followerIds...)).Find(&userRelations)

	for index := range userRelations {
		if utils.Contains(followingIds, userRelations[index].ID) {
			following = append(following, userRelations[index])
		} else {
			followers = append(followers, userRelations[index])
		}
	}
	c.JSON(http.StatusOK, dtos.CreateFollowersFollowingPageResponse(c.Request, following, followers, page, pageSize, userSubscriptionsCount))
}

func SubscribeToUser(c *gin.Context) {
	userId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	database := config.GetDB()
	var followingUser models.User
	err = database.Model(&models.User{}).Where("id = ?", userId).Select([]string{"id", "username"}).Preload("Roles").First(&followingUser).Error
	if err != nil {
		c.JSON(http.StatusNotFound, dtos.CreateDetailedErrorDto("articles", err))
		return
	}

	currentUser := c.MustGet("currentUser").(models.User)

	if currentUser.ID == followingUser.ID {
		c.AbortWithStatusJSON(http.StatusForbidden, dtos.CreateErrorDtoWithMessage("You can not follow yourself"))
	}

	if !followingUser.IsAuthor() && !followingUser.IsAdmin() {
		c.JSON(http.StatusNotFound, dtos.GetErrorDto("You can not subscribe to a non-author nor admin user"))
		return
	}

	var userSubscription models.Subscription
	err = database.Model(models.Subscription{}).Where("following_id = ?", userId, currentUser.ID).First(&userSubscription).Error

	if userSubscription.ID == 0 {
		subscription := models.Subscription{
			FollowingId: followingUser.ID,
			FollowerId:  currentUser.ID,
		}
		if err := services.CreateOne(&subscription); err != nil {
			c.JSON(http.StatusUnprocessableEntity, dtos.CreateDetailedErrorDto("database", err))
			return
		}

		c.JSON(http.StatusOK, dtos.GetSuccessDto(fmt.Sprintf("You subscribed to \"%v\" successfully", followingUser.Username)))

	} else {
		c.JSON(http.StatusForbidden, dtos.GetSuccessDto("Permission denied, You are already subscribed to this user"))
	}
}

func UnsubscribeFromUser(c *gin.Context) {
	userId, err := strconv.ParseUint(c.Param("id"), 10, 32)
	database := config.GetDB()
	var followingUser models.User
	err = database.Model(&models.User{}).Where("id = ?", userId).Select([]string{"id", "username"}).First(&followingUser).Error
	if err != nil {
		c.JSON(http.StatusNotFound, dtos.CreateDetailedErrorDto("articles", err))
		return
	}

	currentUser := c.MustGet("currentUser").(models.User)

	var userSubscription models.Subscription
	err = database.Model(models.Subscription{}).Where("following_id = ? AND follower_id = ?", userId, currentUser.ID).First(&userSubscription).Error

	if userSubscription.ID == 0 {
		c.JSON(http.StatusForbidden, dtos.GetSuccessDto("Permission denied, You were not subscribed to this user"))
	} else {
		database.Delete(&userSubscription)
		c.JSON(http.StatusOK, dtos.GetSuccessDto(fmt.Sprintf("You just unsubscribed from %v", currentUser.Username)))
	}
}

func getInfo(c *gin.Context) (uint, int, int, int, error) {
	page_size_str := c.Query("page_size")
	page_str := c.Query("page")

	pageSize, err := strconv.Atoi(page_size_str)
	if err != nil {
		pageSize = 5
	}

	page, err := strconv.Atoi(page_str)
	if err != nil {
		page = 1
	}

	err = nil
	offset := (page - 1) * pageSize
	username := c.Param("username")
	userIdStr := c.Param("id")
	database := config.GetDB()
	var userId uint
	var user models.User

	if username != "" {
		database.Model(&models.User{}).Where("username = ?", username).Pluck("id", &userId)
	} else if userIdStr != "" {
		ui, err := strconv.Atoi(userIdStr)
		userId = uint(ui)
		if err != nil {
			c.JSON(http.StatusOK, dtos.GetErrorDto("user id must be an integer"))
		}
	} else {
		user = c.MustGet("currentUser").(models.User)
		userId = user.ID
	}

	return userId, page, pageSize, offset, err
}
