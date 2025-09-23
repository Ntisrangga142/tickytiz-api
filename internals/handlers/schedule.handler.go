package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	models "github.com/Ntisrangga142/API_tickytiz/internals/models"
	"github.com/Ntisrangga142/API_tickytiz/internals/repositories"
	"github.com/Ntisrangga142/API_tickytiz/internals/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type ScheduleHandler struct {
	Repo *repositories.ScheduleRepo
	Rdb  *redis.Client
}

func NewScheduleHandler(repo *repositories.ScheduleRepo, rdb *redis.Client) *ScheduleHandler {
	return &ScheduleHandler{Repo: repo, Rdb: rdb}
}

// ScheduleMovie godoc
// @Summary Get movie schedules
// @Description Get all schedules for a specific movie by its ID
// @Tags Schedules
// @Accept json
// @Produce json
// @Param id path int true "Movie ID"
// @Success 200 {object} models.Response[models.ScheduleResponse]
// @Failure 400 {object} models.ErrorResponse "Bad Request"
// @Failure 401 {object} models.ErrorResponse "Unauthorized"
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Security ApiKeyAuth
// @Router /schedule/{id} [get]
func (h *ScheduleHandler) ScheduleMovie(ctx *gin.Context) {
	movieIDStr := ctx.Param("id")
	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil || movieID < 1 {
		utils.HandleError(ctx, http.StatusBadRequest, "Bad Request", "invalid movie id")
		return
	}

	// Ambil query params
	date := ctx.Query("date")
	location := ctx.Query("location")
	showTime := ctx.Query("time")

	// Pointer untuk optional filter
	var datePtr, locPtr, timePtr *string
	if date != "" {
		datePtr = &date
	}
	if location != "" {
		locPtr = &location
	}
	if showTime != "" {
		timePtr = &showTime
	}

	// Redis key hanya untuk data default (tanpa filter)
	redisKey := fmt.Sprintf("Ntisrangga142-Schedule-%d", movieID)

	if date == "" && location == "" && showTime == "" {
		var cachedData models.ScheduleResponse
		if err := utils.CacheHit(ctx.Request.Context(), h.Rdb, redisKey, &cachedData); err == nil && cachedData.MovieID != 0 {
			ctx.JSON(http.StatusOK, models.Response[models.ScheduleResponse]{
				Success: true,
				Message: "Success Load Schedules",
				Data:    cachedData,
			})
			return
		}
	}

	// Ambil dari DB
	schedules, err := h.Repo.Schedule(ctx.Request.Context(), movieID, datePtr, locPtr, timePtr)
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	respData := models.ScheduleResponse{
		MovieID:  movieID,
		Schedule: schedules,
	}

	// Simpan ke Redis jika tanpa filter
	if date == "" && location == "" && showTime == "" {
		if err := utils.RenewCache(ctx.Request.Context(), h.Rdb, redisKey, respData, 10); err != nil {
			log.Println("Failed to set redis cache:", err)
		}
	}

	ctx.JSON(http.StatusOK, models.Response[models.ScheduleResponse]{
		Success: true,
		Message: "Success Load Schedules",
		Data:    respData,
	})
}
