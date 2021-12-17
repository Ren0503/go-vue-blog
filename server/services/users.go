package services

import (
	"errors"
	"ren0503/goapi/config"
	"ren0503/goapi/models"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

func SetPassword(user *models.User, password string) error {
	if len(password) == 0 {
		return errors.New("Password shout not be empty!")
	}

	bytePassword := []byte(password)

	passwordHash, _ := bcrypt.GenerateFromPassword(bytePassword, bcrypt.DefaultCost)
	user.Password = string(passwordHash)

	return nil
}

func FindOneUser(condition interface{}) (models.User, error) {
	database := config.GetDB()
	var user models.User
	err := database.Where(condition).Preload("Roles").First(&user).Error

	return user, err
}

func CreateOne(data interface{}) error {
	database := config.GetDB()
	err := database.Create(data).Error

	return err
}

func SaveOne(data interface{}) error {
	database := config.GetDB()
	err := database.Save(data).Error

	return err
}

func UpdateUser(user *models.User, data interface{}) error {
	database := config.GetDB()
	err := database.Model(user).Update(data).Error

	return err
}

func Following(user *models.User, v models.User) error {
	database := config.GetDB()
	var follow models.Subscription
	err := database.FirstOrCreate(&follow, &models.Subscription{
		FollowingId: v.ID,
		FollowerId:  user.ID,
	}).Error

	return err
}

func IsFollowing(follower *models.User, following models.User) bool {
	database := config.GetDB()
	var follow models.Subscription
	database.Where(models.Subscription{
		FollowingId: following.ID,
		FollowerId:  follower.ID,
	}).First(&follow)

	return follow.ID != 0
}

func UnFollowing(follower *models.User, following models.User) error {
	database := config.GetDB()
	err := database.Where(models.Subscription{
		FollowingId: following.ID,
		FollowerId:  follower.ID,
	}).Delete(models.Subscription{}).Error

	return err
}

func GetFollowings(user *models.User) []models.User {
	database := config.GetDB()
	tx := database.Begin()
	var follows []models.Subscription
	var followings []models.User

	tx.Where(models.Subscription{
		FollowerId: user.ID,
	}).Find(&follows)

	for _, follow := range follows {
		var userModel models.User
		tx.Model(&follow).Related(&userModel, "Following")
		followings = append(followings, userModel)
	}
	tx.Commit()

	return followings
}

func GetFollowingIds(user *models.User) []uint {
	database := config.GetDB()
	tx := database.Begin()
	var users []models.Subscription
	var followingIds []uint

	tx.Where(models.Subscription{
		FollowerId: user.ID,
	}).Find(&users).Pluck("following_id", &followingIds)
	tx.Commit()

	return followingIds
}

func GenerateJwtToken(user *models.User) string {
	jwt_token := jwt.New(jwt.GetSigningMethod("HS512"))

	var roles []string
	for _, role := range user.Roles {
		roles = append(roles, role.Name)
	}

	jwt_token.Claims = jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"roles":    roles,
		"exp":      time.Now().Add(time.Hour * 24 * 90).Unix(),
	}
	// Sign and get the complete encoded token as a string
	token, _ := jwt_token.SignedString([]byte("JWT_SUPER_SECRET"))
	return token
}
