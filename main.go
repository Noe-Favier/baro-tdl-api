package main

import (
	"baro-todo-list/forms"
	"baro-todo-list/models"

	"log"
	"os"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"fmt"
)

const env_db_user = "DB_USER"
const env_db_pass = "DB_PASS"
const env_db_name = "DB_NAME"
const env_db_url  = "DB_URL"
const env_db_port = "DB_PORT"

func main() {
	//get env vars
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	var db_user = os.Getenv(env_db_user)
	var db_pass = os.Getenv(env_db_pass)
	var db_url  = os.Getenv(env_db_url)
	var db_port = os.Getenv(env_db_port)
	var db_name = os.Getenv(env_db_name)

	//connect to database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", db_user, db_pass, db_url, db_port, db_name)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("Failed to connect database")
	}

	db.Logger.LogMode(logger.Error)

	db.AutoMigrate(models.User{})
	db.AutoMigrate(models.Category{})
	db.AutoMigrate(models.Element{})

	//routing
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.GET("/user", func(c *gin.Context) {
		users := []models.User{}
		db.Find(&users)
		c.JSON(http.StatusOK, &users)
	})

	r.POST("/user", func(ctx *gin.Context) {
		var signupForm forms.FormCreateUser
		ctx.BindJSON(&signupForm)

		var newUser models.User = models.User{Username: signupForm.Username, Email: signupForm.Email, Passwd: signupForm.Password}
		db.Create(&newUser)
		ctx.JSON(http.StatusOK, gin.H{
			"message": "succes",
		})
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
