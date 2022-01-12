package handler

import (
	"goginapi/helper"
	"goginapi/payment"
	"goginapi/transaction"
	"goginapi/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type transactionHandler struct {
	service        transaction.Service
	paymentService payment.Service
}

func NewTransactionHanler(service transaction.Service, paymentService payment.Service) *transactionHandler {
	return &transactionHandler{service, paymentService}
}

func (h *transactionHandler) CreateTransaction(c *gin.Context) {
	var input transaction.CreateTransactionInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Failed to get create campaign", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	input.User = currentUser

	newTransaction, err := h.service.CreateTransaction(input)
	if err != nil {
		response := helper.APIResponse("Failed to get create campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := helper.APIResponse("Create of campaigns", http.StatusOK, "success", newTransaction)
	c.JSON(http.StatusOK, response)
}

func (h *transactionHandler) GetNotification(c *gin.Context) {
	var input transaction.TransactionNotificationInput

	err := c.ShouldBindJSON(&input)
	if err != nil {

		response := helper.APIResponse("Failed", http.StatusUnprocessableEntity, "error", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	err = h.service.ProssesPayment(input)
	if err != nil {

		response := helper.APIResponse("Failed", http.StatusUnprocessableEntity, "error", nil)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	c.JSON(http.StatusOK, input)
}
