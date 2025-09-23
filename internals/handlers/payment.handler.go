package handlers

import (
	"net/http"

	"github.com/Ntisrangga142/API_tickytiz/internals/models"
	"github.com/Ntisrangga142/API_tickytiz/internals/repositories"
	"github.com/gin-gonic/gin"
)

type PaymentMethodHandler struct {
	repo *repositories.PaymentMethodRepository
}

func NewPaymentMethodHandler(repo *repositories.PaymentMethodRepository) *PaymentMethodHandler {
	return &PaymentMethodHandler{repo: repo}
}

func (h *PaymentMethodHandler) GetAll(c *gin.Context) {
	methods, err := h.repo.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch payment methods"})
		return
	}

	c.JSON(http.StatusOK, models.Response[[]models.PaymentMethod]{
		Success: true,
		Message: "Success Load Schedules",
		Data:    methods,
	})
}
