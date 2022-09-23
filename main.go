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

const db_user = "DB_USER"
const db_pass = "DB_PASS"
const db_name = "DB_NAME"

func main() {
	//get env vars
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	var user = os.Getenv(db_user)
	var pass = os.Getenv(db_pass)
	var name = os.Getenv(db_name)

	//connect to database
	dsn := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s", user, pass, name)
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
		var error = false
		ctx.BindJSON(&signupForm)

		var newUser models.User = models.User{Username: signupForm.Username, Email: signupForm.Email, Passwd: signupForm.Password}
		db.Create(&newUser).Error(); dbc.Error != nil {
			// Create failed, do something e.g. return, panic etc.
			ctx.JSON(http.StatusOK, gin.H{
				"message": "error",
			})
			return null
			//TODO: handle errors 
		}
		
		if(error){
			
		}else{
			ctx.JSON(http.StatusOK, gin.H{
				"message": "success",
			})
		}
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
