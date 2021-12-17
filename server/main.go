package main

import (
	"fmt"
	"os"
	"ren0503/goapi/config"
	"ren0503/goapi/controllers"
	"ren0503/goapi/middleware"
	"ren0503/goapi/models"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

func drop(db *gorm.DB) {
	db.DropTableIfExists(&models.Like{},
		models.ArticleTag{}.TableName(), models.ArticleCategory{}.TableName(),
		&models.Tag{}, &models.Category{},
		&models.Comment{}, &models.Article{},
		&models.Subscription{},
		&models.UserRole{}, &models.Role{}, &models.User{})
}
func migrate(db *gorm.DB) {
	db.AutoMigrate(&models.Subscription{})
	db.AutoMigrate(&models.User{})
	db.AutoMigrate(&models.Article{})
	db.AutoMigrate(&models.Tag{})
	db.AutoMigrate(&models.Like{})
	db.AutoMigrate(&models.Comment{})
	db.AutoMigrate(&models.Category{})
	db.AutoMigrate(&models.Role{})
}

func addDbConstraints(database *gorm.DB) {
	dialect := database.Dialect().GetName() // mysql, sqlite3
	if dialect != "sqlite3" {
		database.Model(&models.Comment{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
		database.Model(&models.Comment{}).AddForeignKey("article_id", "articles(id)", "CASCADE", "CASCADE")

		database.Model(&models.Article{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
		database.Model(&models.Article{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")

		database.Model(&models.Comment{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")

		database.Model(&models.Like{}).AddForeignKey("article_id", "articles(id)", "CASCADE", "CASCADE")
		database.Model(&models.Like{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")

		database.Model(&models.UserRole{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
		database.Model(&models.UserRole{}).AddForeignKey("role_id", "roles(id)", "CASCADE", "CASCADE")

		database.Model(&models.ArticleTag{}).AddForeignKey("article_id", "articles(id)", "CASCADE", "CASCADE")
		database.Model(models.ArticleTag{}.TableName()).AddForeignKey("tag_id", "tags(id)", "CASCADE", "CASCADE")

		database.Table(models.ArticleCategory{}.TableName()).AddForeignKey("article_id", "articles(id)", "CASCADE", "CASCADE")
		database.Model(&models.ArticleCategory{}).AddForeignKey("category_id", "categories(id)", "CASCADE", "CASCADE")

		database.Model(&models.Subscription{}).AddForeignKey("following_id", "users(id)", "CASCADE", "CASCADE")
		database.Table("subscriptions").AddForeignKey("follower_id", "users(id)", "CASCADE", "CASCADE")
	} else if dialect == "sqlite3" {
		database.Table("comments").AddIndex("comments__idx_article_id", "article_id")
		database.Table("comments").AddIndex("comments__idx_user_id", "user_id")

		database.Model(&models.Comment{}).AddIndex("comments__idx_created_at", "created_at")
	}

	database.Model(&models.UserRole{}).AddIndex("user_roles__idx_user_id", "user_id")
	database.Table("articles_tags").AddIndex("articles_tags__idx_article_id", "article_id")
}

func create(database *gorm.DB) {
	drop(database)
	migrate(database)
	addDbConstraints(database)
}

func main() {
	e := godotenv.Load()
	if e != nil {
		fmt.Print(e)
	}

	database := config.ConnectDB()
	defer database.Close()
	args := os.Args
	if len(args) > 1 {
		first := args[1]
		second := ""
		if len(args) > 2 {
			second = args[2]
		}

		if first == "create" {
			create(database)
		} else if first == "migrate" {
			migrate(database)
		}

		if first == "migrate" {
			migrate(database)
		}

		if first != "" && second == "" {
			os.Exit(0)
		}
	}
	migrate(database)
	router := gin.Default() // gin with the Logger and Recovery Middlewares attached
	router.Use(middleware.Benchmark())
	router.Use(cors.Default())
	router.Use(middleware.UserLoaderMiddleware())
	apiRouteGroup := router.Group("/api")

	controllers.RegisterUserRoutes(apiRouteGroup.Group("/users"))
	controllers.RegisterArticleRoutes(apiRouteGroup.Group("/articles"))
	controllers.RegisterCommentRoutes(apiRouteGroup.Group("/"))
	controllers.RegisterLikeRoutes(apiRouteGroup.Group("/"))
	controllers.RegisterTagRoutes(apiRouteGroup.Group("/tags"))
	controllers.RegisterCategoryRoutes(apiRouteGroup.Group("/categories"))
	controllers.RegisterUserSubscriptionRoutes(apiRouteGroup.Group("/users"))

	router.Run(":8080")
}
