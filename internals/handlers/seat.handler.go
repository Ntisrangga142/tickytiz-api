package handlers

import (
	"net/http"
	"strconv"

	models "github.com/Ntisrangga142/API_tickytiz/internals/models"
	"github.com/Ntisrangga142/API_tickytiz/internals/repositories"
	"github.com/Ntisrangga142/API_tickytiz/internals/utils"
	"github.com/gin-gonic/gin"
)

type SeatHandler struct {
	Repo *repositories.SeatRepository
}

func NewSeatHandler(repo *repositories.SeatRepository) *SeatHandler {
	return &SeatHandler{Repo: repo}
}

// GetSoldSeats godoc
// @Summary Get sold seats
// @Description Retrieve all sold seats for a specific schedule
// @Tags Schedules
// @Accept json
// @Produce json
// @Param id path int true "Schedule ID"
// @Success 200 {object} models.ResponseSeats
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /schedule/seat/{id} [get]
func (h *SeatHandler) GetSoldSeats(ctx *gin.Context) {
	scheduleIDStr := ctx.Param("id")
	scheduleID, err := strconv.Atoi(scheduleIDStr)
	if err != nil || scheduleID < 1 {
		utils.HandleError(ctx, http.StatusBadRequest, "Bad Request", "invalid schedule id")
		return
	}

	seats, err := h.Repo.GetSoldSeats(ctx.Request.Context(), scheduleID)
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	ctx.JSON(http.StatusOK, models.Response[models.SeatResponse]{
		Success: true,
		Message: "Success Load Seats Sold",
		Data: models.SeatResponse{
			ScheduleID: scheduleID,
			Seat:       seats.ID,
		},
	})
}
