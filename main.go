package main

import (
	"baro-todo-list/forms"
	"baro-todo-list/models"
	"strconv"
	"time"

	"log"
	"os"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	cors "github.com/itsjamie/gin-cors"
	"github.com/joho/godotenv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"golang.org/x/crypto/bcrypt"

	"fmt"
)

const env_db_user = "DB_USER"
const env_db_pass = "DB_PASS"
const env_db_name = "DB_NAME"
const env_db_url = "DB_URL"
const env_db_port = "DB_PORT"

const env_bcrypt_cost = "BCRYPT_COST"

const env_app_secret = "APP_SECRET"

func main() {
	//get env vars
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	var db_user = os.Getenv(env_db_user)
	var db_pass = os.Getenv(env_db_pass)
	var db_url = os.Getenv(env_db_url)
	var db_port = os.Getenv(env_db_port)
	var db_name = os.Getenv(env_db_name)

	var bcrypt_cost, bcrypt_cost_error = strconv.ParseInt(os.Getenv(env_bcrypt_cost), 10, 64)

	var app_secret = os.Getenv(env_app_secret)

	if bcrypt_cost_error != nil {
		log.Fatalf("Error parsing BCrypt Cost")
	}

	var ROLE_USER = "USER"
	//var ROLE_ADMIN = "ADMIN"

	//connect to database
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", db_user, db_pass, db_url, db_port, db_name)
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

	r.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     false,
		ValidateHeaders: false,
	}))

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

		hashed_password, bcrypt_err := bcrypt.GenerateFromPassword([]byte(signupForm.Password), int(bcrypt_cost))

		if bcrypt_err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": "error (bcrypt)",
			})
		} else {
			//password successfully hashed
			var newUser models.User = models.User{Username: signupForm.Username, Email: signupForm.Email, Passwd: string(hashed_password), Roles: ROLE_USER}
			result := db.Create(&newUser)

			if result.Error != nil {
				ctx.JSON(http.StatusUnprocessableEntity, gin.H{
					"message":  "error (sql)",
					"errorMsg": result.Error.Error(),
				})
			} else {
				ctx.JSON(http.StatusOK, gin.H{
					"message": "succes",
				})
			}
		}
	})

	r.POST("/login", func(ctx *gin.Context) {
		/* TOKEN GENERATOR */
		var createToken = func(user models.User) (string, error) {
			var err error
			//Creating Access Token
			os.Setenv("ACCESS_SECRET", "jdnfksdmfksd") //this should be in an env file
			atClaims := jwt.MapClaims{}
			atClaims["authorized"] = true
			atClaims["user"] = user
			atClaims["exp"] = time.Now().Add(time.Hour * 48).Unix() //token lasts for 48 hours TODO: add this to .env
			at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
			token, err := at.SignedString([]byte(app_secret))
			if err != nil {
				return "", err
			}
			return token, nil
		}
		/* /// */

		var loggedUser models.User
		var success bool = false

		//Get value from request
		var loginForm forms.FormLoginUser
		ctx.BindJSON(&loginForm)

		//Get users
		users := []models.User{}
		db.Find(&users)

		//Foreach User : is password valid ?
		for _, user := range users {
			println("user check : " + loginForm.Login + ">>" + user.Username + "||" + user.Email)
			println("pwd check" + loginForm.Password + "][" + user.Passwd + ">>")
			if bcrypt.CompareHashAndPassword([]byte(user.Passwd), []byte(loginForm.Password)) == nil && (loginForm.Login == user.Username || loginForm.Login == user.Email) {
				//if password is ok and the right username|email has been supplied :
				println("OK !")
				loggedUser = user
				success = true
				break
			}
		}
		if !success {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "error (logins invalid)",
			})
			return
		}

		token, err := createToken(loggedUser)
		if err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, err.Error())
			return
		}
		ctx.JSON(http.StatusOK, token)
	})

	r.POST("/category", func(ctx *gin.Context) {
		var categoryForm forms.FormCreateCategory
		ctx.BindJSON(&categoryForm)
		//
		var creator models.User
		db.First(&creator, "username = ?", categoryForm.CreatedByUsername)
		if len(creator.Passwd) <= 0 {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"message":  "error (request)",
				"errorMsg": "Username unknown",
			})
			return
		}
		//
		var code string = uuid.New().String()
		var newCategory models.Category = models.Category{Label: categoryForm.Label, CreatedByUsername: categoryForm.CreatedByUsername, Users: []models.User{creator}, Code: code}
		result := db.Create(&newCategory)

		if result.Error != nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"message":  "error (sql)",
				"errorMsg": result.Error.Error(),
			})
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "succes",
			})
		}
	})

	/** REPLACE USERS BY INPUTED USERS */
	r.POST("/category/link/user", func(ctx *gin.Context) {
		var categoryLinkForm forms.FormLinkCategoryToUser
		ctx.BindJSON(&categoryLinkForm)
		//
		var category models.Category
		db.First(&category, "code = ?", categoryLinkForm.CategoryCode)
		//
		newUsers := []models.User{}
		db.Preload("Users").Find(&category) //PreLoad Relations
		db.Where("username IN ?", categoryLinkForm.Usernames).Find(&newUsers)

		var creator []models.User
		db.First(&creator, "username = ?", category.CreatedByUsername)

		var finalUsers []models.User = append(creator, newUsers...)

		result := db.Debug().Model(&category).Association("Users").Replace(&finalUsers)

		if result != nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"message":  "error (sql)",
				"errorMsg": result.Error(),
			})
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "succes",
			})
		}
	})

	r.POST("/element", func(ctx *gin.Context) {
		var elementForm forms.FormCreateElement
		ctx.BindJSON(&elementForm)
		//
		var creator models.User
		db.First(&creator, "username = ?", elementForm.CreatedByUsername)
		if len(creator.Passwd) <= 0 {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"message":  "error (request)",
				"errorMsg": "Username unknown",
			})
			return
		}
		//
		var category models.Category
		db.First(&category, "code = ?", elementForm.CategoryCode)
		//
		var code string = uuid.New().String()
		var newElement models.Element = models.Element{Label: elementForm.Label, CreatedByUsername: elementForm.CreatedByUsername, Code: code, CategoryID: category.ID, Checked: false}
		result := db.Create(&newElement)

		if result.Error != nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"message":  "error (sql)",
				"errorMsg": result.Error.Error(),
			})
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "succes",
			})
		}
	})

	r.POST("/element/check", func(ctx *gin.Context) {
		var elementForm forms.FormCheckElement
		ctx.BindJSON(&elementForm)
		//
		var element models.Element
		db.First(&element, "code = ?", elementForm.Code)
		println(element.Code + " // " + strconv.FormatUint(uint64(element.ID), 10))
		//
		//db.Debug().Model(&element).Updates(models.Element{Checked: element.Checked})
		result := db.Debug().Model(&element).Update("checked", elementForm.Checked)

		if result.Error != nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"message":  "error (sql)",
				"errorMsg": result.Error.Error(),
			})
		} else {
			ctx.JSON(http.StatusOK, gin.H{
				"message": "succes",
			})
		}
	})

	r.GET("category/:categoryCode/elements", func(ctx *gin.Context) {
		var code string = ctx.Param("categoryCode")
		var category models.Category
		//
		db.First(&category, "code = ?", code)
		result := db.Preload("Elements").Find(&category)

		if result.Error != nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"message":  "error (sql)",
				"errorMsg": result.Error.Error(),
			})
		} else {
			ctx.JSON(http.StatusOK, &category.Elements)
		}
	})

	r.GET("user/:username/categories", func(ctx *gin.Context) {
		var username string = ctx.Param("username")
		var user models.User
		//
		db.First(&user, "username = ?", username)
		result := db.Preload("Categories").Find(&user)

		if result.Error != nil {
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"message":  "error (sql)",
				"errorMsg": result.Error.Error(),
			})
		} else {
			ctx.JSON(http.StatusOK, &user.Categories)
		}
	})

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
