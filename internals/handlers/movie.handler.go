package handlers

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/Ntisrangga142/API_tickytiz/internals/models"
	"github.com/Ntisrangga142/API_tickytiz/internals/repositories"
	"github.com/Ntisrangga142/API_tickytiz/internals/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type MovieHandler struct {
	Repo *repositories.MovieRepo
	Rdb  *redis.Client
}

func NewMovieHandler(repo *repositories.MovieRepo, rdb *redis.Client) *MovieHandler {
	return &MovieHandler{Repo: repo, Rdb: rdb}
}

// UpcomingMovies godoc
// @Summary Get upcoming movies
// @Description Retrieve a list of upcoming movies (cached in Redis for 10 minutes)
// @Tags Movies
// @Accept json
// @Produce json
// @Success 200 {object} models.ResponseMovies
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security ApiKeyAuth
// @Router /movies/upcoming [get]
func (h *MovieHandler) UpcomingMovies(ctx *gin.Context) {
	page := 1
	if p := ctx.Query("page"); p != "" {
		if convertPage, err := strconv.Atoi(p); err == nil && convertPage > 0 {
			page = convertPage
		}
	}

	if page == 1 {
		var cachedData []models.Movie
		if err := utils.CacheHit(ctx.Request.Context(), h.Rdb, "Ntisrangga142-UpcomingMovies", &cachedData); err == nil {
			ctx.JSON(http.StatusOK, models.Response[[]models.Movie]{
				Success: true,
				Message: "Success Load Upcoming Movie (from cache)",
				Data:    cachedData,
			})
			return
		}
	}

	movies, err := h.Repo.GetUpcoming(ctx.Request.Context(), page)

	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Internal server Error", err.Error())
		return
	}

	if len(movies) == 0 {
		utils.HandleError(ctx, http.StatusNotFound, "Not Found", "no upcoming movies found")
		return
	}

	if page == 1 {
		if err := utils.RenewCache(ctx.Request.Context(), h.Rdb, "Ntisrangga142-UpcomingMovies", movies, 10); err != nil {
			log.Println("Failed to set redis cache:", err)
		}
	}

	ctx.JSON(http.StatusOK, models.Response[[]models.Movie]{
		Success: true,
		Message: "Success Load Upcoming Movies",
		Data:    movies,
	})
}

// PopularMovies godoc
// @Summary Get popular movies
// @Description Retrieve a list of popular movies (cached in Redis for 10 minutes)
// @Tags Movies
// @Accept json
// @Produce json
// @Success 200 {object} models.ResponseMovies
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security ApiKeyAuth
// @Router /movies/popular [get]
func (h *MovieHandler) PopularMovies(ctx *gin.Context) {
	var cachedData []models.Movie
	if err := utils.CacheHit(ctx.Request.Context(), h.Rdb, "Ntisrangga142-PopularMovies", &cachedData); err == nil {
		ctx.JSON(http.StatusOK, models.Response[[]models.Movie]{
			Success: true,
			Message: "Success Load Popular Movies (from cache)",
			Data:    cachedData,
		})
		return
	}

	movies, err := h.Repo.GetPopular(context.Background())
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	if len(movies) == 0 {
		utils.HandleError(ctx, http.StatusNotFound, "Not Found", "No upcoming movies found")
		return
	}

	if err := utils.RenewCache(ctx.Request.Context(), h.Rdb, "Ntisrangga142-PopularMovies", movies, 10); err != nil {
		log.Println("Failed to set redis cache:", err)
	}

	ctx.JSON(http.StatusOK, models.Response[[]models.Movie]{
		Success: true,
		Message: "Success Load Popular Movies",
		Data:    movies,
	})
}

// FilteredMovies godoc
// @Summary Get filtered movies
// @Description Retrieve movies filtered by title, genres, and page number (supports Redis caching)
// @Tags Movies
// @Accept json
// @Produce json
// @Param title query string false "Filter by movie title"
// @Param genres query string false "Filter by genres, comma-separated"
// @Param page query int false "Page number for pagination"
// @Success 200 {object} models.ResponseMovies
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security ApiKeyAuth
// @Router /movies/ [get]
func (h *MovieHandler) FilteredMovies(ctx *gin.Context) {
	title := ctx.Query("title")
	genresStr := ctx.Query("genres")
	page := 1
	if p := ctx.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}

	var genres []string
	if genresStr != "" {
		genres = strings.Split(genresStr, ",")
	}

	movies, totalItems, err := h.Repo.GetFilteredMovies(ctx.Request.Context(), title, genres, page)
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	if len(movies) == 0 {
		utils.HandleError(ctx, http.StatusNotFound, "Not Found", "No movies found")
		return
	}

	limit := 12
	totalPages := int(math.Ceil(float64(totalItems) / float64(limit)))

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Success Load Movies",
		"data":    movies,
		"pagination": gin.H{
			"page":       page,
			"totalPages": totalPages,
			"totalItems": totalItems,
		},
	})
}

// MovieDetail godoc
// @Summary Get movie details
// @Description Retrieve detailed information for a single movie by ID (cached in Redis for 10 minutes)
// @Tags Movies
// @Accept json
// @Produce json
// @Param id path int true "Movie ID"
// @Success 200 {object} models.ResponseMovieDetail
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security ApiKeyAuth
// @Router /movies/{id} [get]
func (h *MovieHandler) MovieDetail(ctx *gin.Context) {
	movieIDStr := ctx.Param("id")
	if movieIDStr == "" {
		utils.HandleError(ctx, http.StatusBadRequest, "Bad Request", "Movie ID is required")
		return
	}

	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil || movieID < 1 {
		utils.HandleError(ctx, http.StatusBadRequest, "Bad Request", "Invalid movie ID")
		return
	}

	var movie models.MovieDetail
	moviePtr, err := h.Repo.GetMovieDetail(ctx.Request.Context(), movieID)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	movie = *moviePtr

	// 🔹 Return response JSON
	ctx.JSON(http.StatusOK, models.Response[models.MovieDetail]{
		Success: true,
		Message: "Success Load Detail Movie",
		Data:    movie,
	})
}
