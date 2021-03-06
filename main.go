package main

import (
	"goginapi/auth"
	"goginapi/campaign"
	"goginapi/handler"
	"goginapi/helper"
	"goginapi/payment"
	"goginapi/transaction"
	"goginapi/user"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
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
	campaignRepository := campaign.NewRepository(db)
	transactionRepository := transaction.NewRepository(db)

	userService := user.NewService(userRepository)
	authService := auth.NewService()
	campaignService := campaign.NewService(campaignRepository)
	paymentService := payment.NewService()
	transactionService := transaction.NewService(transactionRepository, paymentService, campaignRepository)

	userHandler := handler.NewUserHandler(userService, authService)
	campaignHandler := handler.NewCampaignHanler(campaignService)
	transactionHandler := handler.NewTransactionHanler(transactionService, paymentService)

	router := gin.Default()
	router.Static("images", "./images")
	api := router.Group("/api/v1")

	api.POST("users", userHandler.RegisterUser)
	api.POST("sessions", userHandler.Login)
	api.POST("email-checkers", userHandler.CheckEmailAvailibbility)
	api.POST("avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)

	api.GET("campaings", campaignHandler.GetCampaigns)
	api.GET("campaing/:id", campaignHandler.GetCampaign)
	api.POST("campaing", authMiddleware(authService, userService), campaignHandler.CreateCampaign)
	api.PUT("campaing/:id", authMiddleware(authService, userService), campaignHandler.UpdateCampaign)
	api.POST("campaing-images", authMiddleware(authService, userService), campaignHandler.UploadImage)
	api.POST("create-campaign", authMiddleware(authService, userService), transactionHandler.CreateTransaction)

	router.Run()
}

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userID := int(claim["user_id"].(float64))

		user, err := userService.GetUserByID(userID)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)
	}
}
