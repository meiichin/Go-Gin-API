package handler

import (
	"goginapi/campaign"
	"goginapi/helper"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type campainHandler struct {
	service campaign.Service
}

func NewCampaignHanler(service campaign.Service) *campainHandler {
	return &campainHandler{service}
}

func (h *campainHandler) GetCampaigns(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Query("user_id"))

	campaigns, err := h.service.GetCampaigns(userID)
	if err != nil {
		response := helper.APIResponse("Error to get campaign", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := campaign.FormatCampaigns(campaigns)

	response := helper.APIResponse("List of campaigns", http.StatusOK, "success", formatter)
	c.JSON(http.StatusOK, response)
}
