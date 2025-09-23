package handlers

import (
	"net/http"

	"github.com/Ntisrangga142/API_tickytiz/internals/models"
	"github.com/Ntisrangga142/API_tickytiz/internals/repositories"
	"github.com/Ntisrangga142/API_tickytiz/internals/utils"
	"github.com/gin-gonic/gin"
)

type GenreHandler struct {
	Repo *repositories.GenreRepo
}

func NewGenreHandler(repo *repositories.GenreRepo) *GenreHandler {
	return &GenreHandler{Repo: repo}
}

// GetGenres godoc
// @Summary Get list of genres
// @Description Get all movie genres
// @Tags Genres
// @Accept json
// @Produce json
// @Success 200 {object} models.GenreResponse
// @Failure 500 {object} models.ErrorResponse "Internal Server Error"
// @Router /genres [get]
func (h *GenreHandler) GetGenres(ctx *gin.Context) {
	genres, err := h.Repo.GetGenres(ctx.Request.Context())
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	ctx.JSON(http.StatusOK, models.Response[[]models.Genre]{
		Success: true,
		Message: "Success Load Genres",
		Data:    genres,
	})
}
