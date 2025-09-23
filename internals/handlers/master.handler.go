package handlers

import (
	"net/http"

	"github.com/Ntisrangga142/API_tickytiz/internals/models"
	"github.com/Ntisrangga142/API_tickytiz/internals/repositories"

	"github.com/gin-gonic/gin"
)

type MasterHandler struct {
	Repo *repositories.MasterRepo
}

func NewMasterHandler(repo *repositories.MasterRepo) *MasterHandler {
	return &MasterHandler{Repo: repo}
}

func (h *MasterHandler) GetDirectors(c *gin.Context) {
	directors, err := h.Repo.GetDirectors(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.Response[[]models.MasterDirector]{
		Success: true,
		Message: "Success Load Director",
		Data:    directors,
	})
}

func (h *MasterHandler) GetActors(c *gin.Context) {
	actors, err := h.Repo.GetActors(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.Response[[]models.MasterActor]{
		Success: true,
		Message: "Success Load Actors",
		Data:    actors,
	})
}

func (h *MasterHandler) GetLocations(c *gin.Context) {
	locations, err := h.Repo.GetLocations(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.Response[[]models.MasterLocation]{
		Success: true,
		Message: "Success Load Location",
		Data:    locations,
	})
}

func (h *MasterHandler) GetTimes(c *gin.Context) {
	times, err := h.Repo.GetTimes(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.Response[[]models.MasterTime]{
		Success: true,
		Message: "Success Load Times",
		Data:    times,
	})
}

func (h *MasterHandler) GetCinemas(c *gin.Context) {
	cinemas, err := h.Repo.GetCinemas(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.Response[[]models.MasterCinema]{
		Success: true,
		Message: "Success Load Cinemas",
		Data:    cinemas,
	})
}
