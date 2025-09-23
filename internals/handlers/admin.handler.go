package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/Ntisrangga142/API_tickytiz/internals/models"
	"github.com/Ntisrangga142/API_tickytiz/internals/repositories"
	"github.com/Ntisrangga142/API_tickytiz/internals/utils"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type AdminHandler struct {
	Repo *repositories.AdminRepo
	Rdb  *redis.Client
}

func NewAdminHandler(repo *repositories.AdminRepo, rdb *redis.Client) *AdminHandler {
	return &AdminHandler{Repo: repo, Rdb: rdb}
}

// GetAllMovie godoc
// @Summary Get movie by ID
// @Description Retrieve detailed information of a movie by its ID
// @Tags Admin
// @Accept json
// @Produce json
// @Param id query int true "Movie ID"
// @Success 200 {object} models.AdminMovieResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 404 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security ApiKeyAuth
// @Router /admin/ [get]
func (h *AdminHandler) GetAllMovie(ctx *gin.Context) {

	var cachedData []models.AdminMovie
	if err := utils.CacheHit(ctx.Request.Context(), h.Rdb, "Ntisrangga142-AllMovies", &cachedData); err == nil {
		ctx.JSON(http.StatusOK, models.Response[[]models.AdminMovie]{
			Success: true,
			Message: "Success Load All Movie (from cache)",
			Data:    cachedData,
		})
		return
	}

	movie, err := h.Repo.GetAllMovie(ctx.Request.Context())
	if err != nil {
		utils.HandleError(ctx, http.StatusNotFound, "Not Found", "failed load all movies ")
		return
	}

	if err := utils.RenewCache(ctx.Request.Context(), h.Rdb, "Ntisrangga142-AllMovies", movie, 10); err != nil {
		log.Println("Failed to set redis cache:", err)
	}

	ctx.JSON(http.StatusOK, models.Response[[]models.AdminMovie]{
		Success: true,
		Message: "Success Load All Movie",
		Data:    movie,
	})
}

// UpdateMovie godoc
// @Summary Update movie
// @Description Update movie details, including poster and backdrop images
// @Tags Admin
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "Movie ID"
// @Param title formData string false "Movie title"
// @Param poster formData file false "Poster image"
// @Param backdrop formData file false "Backdrop image"
// @Param release_date formData string false "Release date in YYYY-MM-DD format"
// @Param duration formData int false "Duration in minutes"
// @Param synopsis formData string false "Movie synopsis"
// @Param rating formData number false "Movie rating"
// @Param id_director formData int false "Director ID"
// @Success 200 {object} models.AdminMovieUpdateResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security ApiKeyAuth
// @Router /admin/{id} [put]
func (h *AdminHandler) UpdateMovie(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id < 0 {
		utils.HandleError(ctx, http.StatusBadRequest, "Bad Request", err.Error())
		return
	}

	var req models.AdminUpdate
	if err := ctx.ShouldBind(&req); err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "Bad Request", "failed binding data")
		return
	}

	if file, err := ctx.FormFile("poster"); err == nil {
		filename := fmt.Sprintf("poster_%d_%s", id, file.Filename)
		path := filepath.Join("public/movie/poster", filename)
		if err := ctx.SaveUploadedFile(file, path); err != nil {
			utils.HandleError(ctx, http.StatusInternalServerError, "Internal Server Error", "failed save pster movie")
			return
		}
		req.Poster = &filename
	}

	if file, err := ctx.FormFile("backdrop"); err == nil {
		filename := fmt.Sprintf("backdrop_%d_%s", id, file.Filename)
		path := filepath.Join("public/movie/backdrop", filename)
		if err := ctx.SaveUploadedFile(file, path); err != nil {
			utils.HandleError(ctx, http.StatusInternalServerError, "Internal Server Error", "failed sva backdrop movie")
			return
		}
		req.Backdrop = &filename
	}

	if err := h.Repo.UpdateMovie(ctx, req, id); err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Internal Server Error", "failed update movie")
		return
	}

	if err := utils.InvalidateCache(ctx, h.Rdb, "Ntisrangga142-AllMovies"); err != nil {
		log.Println("Failed invalidate cache:", err)
	}
	if err := utils.InvalidateCache(ctx, h.Rdb, "Ntisrangga142-UpcomingMovies"); err != nil {
		log.Println("Failed invalidate cache:", err)
	}
	if err := utils.InvalidateCache(ctx, h.Rdb, "Ntisrangga142-PopularMovies"); err != nil {
		log.Println("Failed invalidate cache:", err)
	}
	if err := utils.InvalidateCache(ctx, h.Rdb, "Ntisrangga142-FilterMovies"); err != nil {
		log.Println("Failed invalidate cache:", err)
	}

	ctx.JSON(http.StatusOK, models.Response[models.AdminUpdate]{
		Success: true,
		Message: "Success Update Movie",
		Data:    req,
	})
}

// DeleteMovie godoc
// @Summary Delete movie
// @Description Soft-delete a movie by ID (sets delete_at timestamp)
// @Tags Admin
// @Accept json
// @Produce json
// @Param id path int true "Movie ID"
// @Success 200 {object} models.AdminMovieDeleteResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security ApiKeyAuth
// @Router /admin/{id} [delete]
func (h *AdminHandler) DeleteMovie(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id < 0 {
		utils.HandleError(ctx, http.StatusBadRequest, "Bad Request", "invalid movie id")
		return
	}

	res, err := h.Repo.DeleteMovie(ctx.Request.Context(), id)
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Internal Server Error", "failed delete movie")
		return
	}

	if err := utils.InvalidateCache(ctx, h.Rdb, "Ntisrangga142-AllMovies"); err != nil {
		log.Println("Failed invalidate cache:", err)
	}
	if err := utils.InvalidateCache(ctx, h.Rdb, "Ntisrangga142-UpcomingMovies"); err != nil {
		log.Println("Failed invalidate cache:", err)
	}
	if err := utils.InvalidateCache(ctx, h.Rdb, "Ntisrangga142-PopularMovies"); err != nil {
		log.Println("Failed invalidate cache:", err)
	}
	if err := utils.InvalidateCache(ctx, h.Rdb, "Ntisrangga142-FilterMovies"); err != nil {
		log.Println("Failed invalidate cache:", err)
	}

	ctx.JSON(http.StatusOK, models.Response[models.AdminDelete]{
		Success: true,
		Message: "Success Delete Movie",
		Data:    *res,
	})
}

// CreateMovie godoc
// @Summary Insert movie
// @Description Insert movie with poster, backdrop, genres, actors, schedule
// @Tags Admin
// @Accept multipart/form-data
// @Produce json
// @Param title formData string true "Movie title"
// @Param release_date formData string true "Release date yyyy-mm-dd"
// @Param duration formData int true "Duration in minutes"
// @Param director_name formData string true "Director name"
// @Param poster formData file true "Poster image"
// @Param backdrop formData file true "Backdrop image"
// @Param genres formData string true "Comma separated genres"
// @Param cast_name formData string true "Comma separated actor names"
// @Param location formData string true "Comma separated locations"
// @Param date_time formData string true "Comma separated datetimes (yyyy-mm-ddTHH:MM)"
// @Success 200 {object} models.AdminMovie
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /admin/movies [post]
func (h *AdminHandler) CreateMovie(c *gin.Context) {
	title := c.PostForm("title")
	releaseStr := strings.TrimSpace(c.PostForm("release"))
	duration, _ := strconv.Atoi(c.PostForm("duration"))
	directorID, _ := strconv.Atoi(c.PostForm("director"))
	synopsis := c.PostForm("synopsis")

	if title == "" || releaseStr == "" || duration == 0 || directorID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All required fields must be filled"})
		return
	}

	releaseDate, err := time.Parse("2006-01-02", releaseStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid release date format"})
		return
	}

	// Parse genres & actors
	var genreIDs, actorIDs []int
	for _, g := range strings.Split(c.PostForm("added_genres"), ",") {
		if g != "" {
			id, _ := strconv.Atoi(g)
			genreIDs = append(genreIDs, id)
		}
	}
	for _, a := range strings.Split(c.PostForm("added_actors"), ",") {
		if a != "" {
			id, _ := strconv.Atoi(a)
			actorIDs = append(actorIDs, id)
		}
	}

	// Parse schedules JSON
	var combos []models.ScheduleComboAdminInsert
	schedulesJSON := c.PostForm("schedules")
	if schedulesJSON != "" {
		if err := json.Unmarshal([]byte(schedulesJSON), &combos); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid schedules format"})
			return
		}
	}

	movie := models.MovieAdminInsert{
		Title:       title,
		ReleaseDate: releaseDate,
		Duration:    duration,
		Synopsis:    synopsis,
		IdDirector:  directorID,
	}

	// Simpan movie + relasi
	movieID, err := h.Repo.CreateMovieWithRelations(context.Background(), &movie, genreIDs, actorIDs, combos)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// ===== Handle Upload File =====
	uploadDir := filepath.Join("public", "movie")
	_ = os.MkdirAll(filepath.Join(uploadDir, "poster"), os.ModePerm)
	_ = os.MkdirAll(filepath.Join(uploadDir, "backdrop"), os.ModePerm)

	posterFile, _ := c.FormFile("poster")
	backdropFile, _ := c.FormFile("backdrop")

	var posterFilename, backdropFilename string

	if posterFile != nil {
		ext := strings.ToLower(filepath.Ext(posterFile.Filename))
		posterFilename = fmt.Sprintf("poster_%d%s", movieID, ext)
		posterPath := filepath.Join(uploadDir, "poster", posterFilename)
		_ = c.SaveUploadedFile(posterFile, posterPath)
	}

	if backdropFile != nil {
		ext := strings.ToLower(filepath.Ext(backdropFile.Filename))
		backdropFilename = fmt.Sprintf("backdrop_%d%s", movieID, ext)
		backdropPath := filepath.Join(uploadDir, "backdrop", backdropFilename)
		_ = c.SaveUploadedFile(backdropFile, backdropPath)
	}

	_ = h.Repo.UpdateMovieFiles(context.Background(), movieID, posterFilename, backdropFilename)

	if err := utils.InvalidateCache(c, h.Rdb, "Ntisrangga142-AllMovies"); err != nil {
		log.Println("Failed invalidate cache:", err)
	}
	if err := utils.InvalidateCache(c, h.Rdb, "Ntisrangga142-UpcomingMovies"); err != nil {
		log.Println("Failed invalidate cache:", err)
	}
	if err := utils.InvalidateCache(c, h.Rdb, "Ntisrangga142-PopularMovies"); err != nil {
		log.Println("Failed invalidate cache:", err)
	}
	if err := utils.InvalidateCache(c, h.Rdb, "Ntisrangga142-FilterMovies"); err != nil {
		log.Println("Failed invalidate cache:", err)
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Movie created successfully",
		"movie_id": movieID,
		"poster":   posterFilename,
		"backdrop": backdropFilename,
	})
}

func (h *AdminHandler) GetMovieDetailUpdate(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie id"})
		return
	}

	movie, err := h.Repo.GetMovieDetail(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "movie not found"})
		return
	}

	c.JSON(http.StatusOK, movie)
}

// UpdateMovieAdmin handler
func (h *AdminHandler) UpdateMovieAdmin(c *gin.Context) {
	idParam := c.Param("id")
	movieID, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid movie id"})
		return
	}

	// Ambil field form
	title := c.PostForm("title")
	release := c.PostForm("release")
	durationStr := c.PostForm("duration")
	director := c.PostForm("director")
	synopsis := c.PostForm("synopsis")

	addedGenres := parseIDs(c.PostForm("added_genres"))
	removedGenres := parseIDs(c.PostForm("removed_genres"))
	addedActors := parseIDs(c.PostForm("added_actors"))
	removedActors := parseIDs(c.PostForm("removed_actors"))

	schedulesJSON := c.PostForm("schedules")

	duration, _ := strconv.Atoi(durationStr)
	releaseDate, _ := time.Parse("2006-01-02", release)

	// Prepare model
	var m models.MovieUpdateAdmin
	m.ID = movieID
	m.Title = title
	m.ReleaseDate = releaseDate
	m.Duration = duration
	m.DirectorID, _ = strconv.Atoi(director)
	m.Synopsis = synopsis

	// Parse schedules
	var schedules []models.ScheduleUpdate
	if err := json.Unmarshal([]byte(schedulesJSON), &schedules); err == nil {
		m.Schedules = schedules
	}

	// Handle poster & backdrop
	uploadDir := "./public/movie"
	os.MkdirAll(path.Join(uploadDir, "poster"), os.ModePerm)
	os.MkdirAll(path.Join(uploadDir, "backdrop"), os.ModePerm)

	// Poster
	if posterFile, posterHeader, err := c.Request.FormFile("poster"); err == nil && posterFile != nil {
		defer posterFile.Close()
		posterExt := path.Ext(posterHeader.Filename)
		posterPath := path.Join(uploadDir, "poster", fmt.Sprintf("poster_%d%s", movieID, posterExt))
		out, err := os.Create(posterPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save poster"})
			return
		}
		defer out.Close()
		if _, err := io.Copy(out, posterFile); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save poster"})
			return
		}
		m.PosterPath = fmt.Sprintf("poster_%d%s", movieID, posterExt)
	}

	// Backdrop
	if backdropFile, backdropHeader, err := c.Request.FormFile("backdrop"); err == nil && backdropFile != nil {
		defer backdropFile.Close()
		backdropExt := path.Ext(backdropHeader.Filename)
		backdropPath := path.Join(uploadDir, "backdrop", fmt.Sprintf("backdrop_%d%s", movieID, backdropExt))
		out, err := os.Create(backdropPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save backdrop"})
			return
		}
		defer out.Close()
		if _, err := io.Copy(out, backdropFile); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save backdrop"})
			return
		}
		m.BackdropPath = fmt.Sprintf("backdrop_%d%s", movieID, backdropExt)
	}

	// Update di repository
	ctx := context.Background()
	if err := h.Repo.UpdateMovieAdmin(ctx, movieID, m, addedGenres, removedGenres, addedActors, removedActors); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := utils.InvalidateCache(ctx, h.Rdb, "Ntisrangga142-AllMovies"); err != nil {
		log.Println("Failed invalidate cache:", err)
	}
	if err := utils.InvalidateCache(ctx, h.Rdb, "Ntisrangga142-UpcomingMovies"); err != nil {
		log.Println("Failed invalidate cache:", err)
	}
	if err := utils.InvalidateCache(ctx, h.Rdb, "Ntisrangga142-PopularMovies"); err != nil {
		log.Println("Failed invalidate cache:", err)
	}
	if err := utils.InvalidateCache(ctx, h.Rdb, "Ntisrangga142-FilterMovies"); err != nil {
		log.Println("Failed invalidate cache:", err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "movie updated"})
}

// Helper untuk parse list of IDs dari CSV string
func parseIDs(s string) []int {
	if s == "" {
		return []int{}
	}
	parts := strings.Split(s, ",")
	var res []int
	for _, p := range parts {
		if id, err := strconv.Atoi(p); err == nil {
			res = append(res, id)
		}
	}
	return res
}
