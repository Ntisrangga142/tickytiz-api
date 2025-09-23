package handlers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/Ntisrangga142/API_tickytiz/internals/models"
	"github.com/Ntisrangga142/API_tickytiz/internals/repositories"
	"github.com/Ntisrangga142/API_tickytiz/internals/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type OrderHandler struct {
	Repo *repositories.OrderRepo
	Rdb  *redis.Client
}

func NewOrderHandler(repo *repositories.OrderRepo, rdb *redis.Client) *OrderHandler {
	return &OrderHandler{Repo: repo, Rdb: rdb}
}

// CreateOrder godoc
// @Summary Create a new order
// @Description Create a new order with seats and associate it with the logged-in user
// @Tags Orders
// @Accept json
// @Produce json
// @Param request body models.OrderRequest true "Order request body"
// @Success 200 {object} models.ResponseOrders
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /order [post]
func (h *OrderHandler) CreateOrder(ctx *gin.Context) {
	userID, err := utils.GetUserIDFromJWT(ctx)
	if err != nil {
		utils.HandleError(ctx, http.StatusUnauthorized, "Unauthorized", err.Error())
		return
	}

	var req models.OrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "Bad Request", err.Error())
		return
	}

	res, err := h.Repo.CreateOrder(ctx.Request.Context(), req, userID)
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	listSeat, err := h.Repo.CreateOrderDetails(ctx.Request.Context(), req, res.ID)
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	res.Seat = listSeat

	redisKey := fmt.Sprintf("Ntisrangga142-UserHistory-%d", userID)
	if err := utils.InvalidateCache(ctx.Request.Context(), h.Rdb, redisKey); err != nil {
		log.Printf("Failed to invalidate chace : %s\n", err.Error())
	}

	ctx.JSON(http.StatusOK, models.Response[models.OrderResponse]{
		Success: true,
		Message: "Success Order",
		Data:    *res,
	})
}
