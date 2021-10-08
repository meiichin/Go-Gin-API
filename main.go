package main

import (
	"goginapi/auth"
	"goginapi/handler"
	"goginapi/user"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	// // refer https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dsn := "root:@tcp(127.0.0.1:3306)/golangapi?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	userRepository := user.NewRepository(db)
	userService := user.NewService(userRepository)
	authService := auth.NewService()

	userHandler := handler.NewUserHandler(userService, authService)

	router := gin.Default()
	api := router.Group("/api/v1")

	api.POST("users", userHandler.RegisterUser)
	api.POST("sessions", userHandler.Login)
	api.POST("email-checkers", userHandler.CheckEmailAvailibbility)
	api.POST("avatars", userHandler.UploadAvatar)

	router.Run()

}
