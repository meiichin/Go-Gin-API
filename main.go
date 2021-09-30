package main

import (
	"fmt"
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

	userInput := user.LoginInput{
		Email:    "coco@gmail.com",
		Password: "password",
	}
	userByEmail, err := userService.Login(userInput)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(userByEmail)

	return

	userHandler := handler.NewUserHandler(userService)

	router := gin.Default()
	api := router.Group("/api/v1")

	api.POST("users", userHandler.RegisterUser)

	router.Run()

}
