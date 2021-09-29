package handler

import (
	"goginapi/helper"
	"goginapi/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	var input user.RegisterUserInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	newUser, err := h.userService.RegisterUser(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	formatter := user.FormatUser(newUser, "token")

	response := helper.APIResponse("Account has registered", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}
