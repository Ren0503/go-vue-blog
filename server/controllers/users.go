package controllers

import (
	"net/http"
	"ren0503/goapi/dtos"
	"ren0503/goapi/models"
	"ren0503/goapi/services"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUserRoutes(router *gin.RouterGroup) {
	router.POST("/", RegisterUser)
	router.POST("/login", LoginUser)
}

func RegisterUser(c *gin.Context) {
	var json dtos.RegisterRequestDto

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, dtos.CreateBadRequestErrorDto(err))

		return
	}

	password, _ := bcrypt.GenerateFromPassword([]byte(json.Password), bcrypt.DefaultCost)
	if err := services.CreateOne(&models.User{
		Username:  json.Username,
		Password:  string(password),
		FirstName: json.FirstName,
		LastName:  json.LastName,
		Email:     json.Email,
	}); err != nil {
		c.JSON(http.StatusUnprocessableEntity, dtos.CreateDetailedErrorDto("db_error", err))
		return
	}

	c.JSON(http.StatusCreated, dtos.CreateSuccessWithMessageDto("User created successfully"))
}

func LoginUser(c *gin.Context) {
	var json dtos.Login

	if err := c.ShouldBindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, dtos.CreateBadRequestErrorDto(err))
		return
	}

	user, err := services.FindOneUser(&models.User{Username: json.Username})

	if err != nil {
		c.JSON(http.StatusForbidden, dtos.CreateDetailedErrorDto("login", err))
		return
	}

	if user.IsValidPassword(json.Password) != nil {
		c.JSON(http.StatusForbidden, dtos.CreateErrorDtoWithMessage("invalid credentials"))
		return
	}
	c.JSON(http.StatusOK, dtos.CreateLoginSuccessful(&user))
}
