package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/Ntisrangga142/API_tickytiz/internals/models"
	"github.com/Ntisrangga142/API_tickytiz/internals/repositories"
	"github.com/Ntisrangga142/API_tickytiz/internals/utils"
	pkg "github.com/Ntisrangga142/API_tickytiz/pkg"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type UserHandler struct {
	repo *repositories.UserRepository
	Rdb  *redis.Client
}

func NewUserHandler(Repo *repositories.UserRepository, rdb *redis.Client) *UserHandler {
	return &UserHandler{repo: Repo, Rdb: rdb}
}

// GetProfile godoc
// @Summary Get user profile
// @Description Mengambil data profil user yang sedang login
// @Tags Users
// @Produce json
// @Success 200 {object} models.ResponseUserProfile
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security ApiKeyAuth
// @Router /user [get]
func (h *UserHandler) GetProfile(ctx *gin.Context) {
	userID, err := utils.GetUserIDFromJWT(ctx)
	if err != nil {
		utils.HandleError(ctx, http.StatusUnauthorized, "Unauthorized", err.Error())
		return
	}

	var cachedData models.UserProfile
	redisKey := fmt.Sprintf("Ntisrangga142-UserProfiles-%d", userID)
	if err := utils.CacheHit(ctx.Request.Context(), h.Rdb, redisKey, &cachedData); err == nil {
		ctx.JSON(http.StatusOK, models.Response[models.UserProfile]{
			Success: true,
			Message: "Success Load User Profile (from cache)",
			Data:    cachedData,
		})
		return
	}

	user, err := h.repo.GetProfileByID(ctx.Request.Context(), userID)
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	if err := utils.RenewCache(ctx.Request.Context(), h.Rdb, redisKey, user, 10); err != nil {
		log.Println("Failed to set redis cache:", err)
	}

	ctx.JSON(http.StatusOK, models.Response[models.UserProfile]{
		Success: true,
		Message: "Success Load User Profile",
		Data:    *user,
	})
}

// GetHistory godoc
// @Summary Get order history
// @Description Mengambil riwayat order user yang sedang login
// @Tags Users
// @Produce json
// @Success 200 {object} models.ResponseOrderHistory
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security ApiKeyAuth
// @Router /user/history [get]
func (h *UserHandler) GetHistory(ctx *gin.Context) {
	userID, err := utils.GetUserIDFromJWT(ctx)
	if err != nil {
		utils.HandleError(ctx, http.StatusUnauthorized, "Unauthorized", err.Error())
		return
	}

	var cachedData models.OrderHistoryResponse
	redisKey := fmt.Sprintf("Ntisrangga142-UserHistory-%d", userID)
	if err := utils.CacheHit(ctx.Request.Context(), h.Rdb, redisKey, &cachedData); err == nil {
		ctx.JSON(http.StatusOK, models.Response[models.OrderHistoryResponse]{
			Success: true,
			Message: "Success Load History (from cache)",
			Data:    cachedData,
		})
		return
	}

	history, err := h.repo.GetHistoryByUserID(ctx, userID)
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	if err := utils.RenewCache(ctx.Request.Context(), h.Rdb, redisKey, history, 10); err != nil {
		log.Println("Failed to set redis cache:", err)
	}

	ctx.JSON(http.StatusOK, models.Response[models.OrderHistoryResponse]{
		Success: true,
		Message: "Success Load History",
		Data:    history,
	})
}

// UpdateProfile godoc
// @Summary Update user profile
// @Description Mengupdate profil user termasuk upload profile image
// @Tags Users
// @Accept multipart/form-data
// @Produce json
// @Param profileimg formData file false "Profile Image"
// @Param firstname formData string false "First Name"
// @Param lastname formData string false "Last Name"
// @Param phone formData string false "Phone"
// @Success 200 {object} models.ResponseUpdateProfile
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Security ApiKeyAuth
// @Router /user/ [patch]
func (h *UserHandler) UpdateProfile(ctx *gin.Context) {
	userID, err := utils.GetUserIDFromJWT(ctx)
	if err != nil {
		utils.HandleError(ctx, http.StatusUnauthorized, "Unauthorized", err.Error())
		return
	}

	var req models.UpdateProfile
	if err := ctx.ShouldBind(&req); err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "Bad Request", err.Error())
		return
	}

	// 🔹 Ambil file gambar kalau ada
	file, err := ctx.FormFile("profileimg")
	if err == nil {
		filename := fmt.Sprintf("profile_%d.png", userID)
		path := filepath.Join("public/profiles/", filename)

		if err := ctx.SaveUploadedFile(file, path); err != nil {
			utils.HandleError(ctx, http.StatusInternalServerError, "Internal Server Error", err.Error())
			return
		}
		req.ProfileImg = &filename
	}

	if err := h.repo.UpdateUser(ctx, userID, req); err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	redisKey := fmt.Sprintf("Ntisrangga142-UserProfiles-%d", userID)
	if err := utils.InvalidateCache(ctx.Request.Context(), h.Rdb, redisKey); err != nil {
		log.Println("Failed to invalidate redis cache:", err)
	}

	ctx.JSON(http.StatusOK, models.Response[models.UpdateProfile]{
		Success: true,
		Message: "Success Update Profile",
		Data:    req,
	})
}

func (h *UserHandler) GetVirtualAccountHandler(c *gin.Context) {
	userID, err := utils.GetUserIDFromJWT(c)
	if err != nil {
		utils.HandleError(c, http.StatusUnauthorized, "Unauthorized", err.Error())
		return
	}

	users, err := h.repo.GetVirtualAccountByID(context.Background(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": users})
}

func (h *UserHandler) ChangePassword(ctx *gin.Context) {
	var req models.ChangePasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Ambil UserID dari JWT
	userID, err := utils.GetUserIDFromJWT(ctx)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Ambil password lama dari DB
	dbPassword, err := h.repo.GetPasswordByID(ctx, userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user password"})
		return
	}

	claims := pkg.NewHashConfig()
	claims.UseRecommended()

	// Cek password lama
	valid, err := claims.CompareHashAndPassword(req.OldPassword, dbPassword)
	if err != nil || !valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "old password is incorrect"})
		return
	}

	// Hash password baru
	newHashedPwd, err := claims.GenerateHash(req.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash new password"})
		return
	}

	// Update password ke DB
	if err := h.repo.UpdatePassword(ctx, userID, newHashedPwd); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update password"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "password changed successfully"})
}
