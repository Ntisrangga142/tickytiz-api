package handlers

import (
	"context"
	"net/http"
	"time"

	models "github.com/Ntisrangga142/API_tickytiz/internals/models"
	repo "github.com/Ntisrangga142/API_tickytiz/internals/repositories"
	"github.com/Ntisrangga142/API_tickytiz/internals/utils"
	pkg "github.com/Ntisrangga142/API_tickytiz/pkg"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	Repo *repo.Auth
}

func NewAuthHandler(repo *repo.Auth) *AuthHandler {
	return &AuthHandler{Repo: repo}
}

// Register godoc
// @Summary Register User
// @Description Register Akun Baru
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param request body models.RegisterRequest true "Register request"
// @Success 200 {object} models.RegisterDocs
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /auth/register [post]
func (h *AuthHandler) Register(ctx *gin.Context) {
	var req models.RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "Bad Request", "failed binding data")
		return
	}

	// Hash Password
	hashConfig := pkg.NewHashConfig()
	hashConfig.UseRecommended()
	hashedPassword, err := hashConfig.GenerateHash(req.Password)
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Internal Server Error", "failed hashed password")
		return
	}

	if err := h.Repo.Register(ctx, req.Email, hashedPassword); err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Internal Server Error", err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, models.Response[string]{
		Success: true,
		Message: "Register successful",
		Data:    "",
	})
}

// Login godoc
// @Summary Login user
// @Description Login dengan Email dan Password untuk mendapatkan JWT token
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param request body models.LoginRequest true "Login request"
// @Success 200 {object} models.LoginDocs
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /auth [post]
func (h *AuthHandler) Login(ctx *gin.Context) {
	var req models.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, "Bad Request", "failed binding data")
		return
	}

	// Cari akun
	userID, role, hashedPassword, err := h.Repo.Login(ctx.Request.Context(), req.Email)
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Internal Server Error", "user not found")
		return
	}
	if userID == 0 {
		utils.HandleError(ctx, http.StatusUnauthorized, "Unauthorized", "user not found")
		return
	}

	// Verifikasi password
	hashConfig := pkg.NewHashConfig()
	match, err := hashConfig.CompareHashAndPassword(req.Password, hashedPassword)
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Internal Server Error", "failed compare password")
		return
	}
	if !match {
		utils.HandleError(ctx, http.StatusUnauthorized, "Unauthorized", "invalid password")
		return
	}

	// Generate JWT
	claims := pkg.NewJWTClaims(userID, role)
	token, err := claims.GenToken()
	if err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Internal Server Error", "failed generate token")
		return
	}

	ctx.JSON(http.StatusOK, models.Response[models.LoginResponse]{
		Success: true,
		Message: "Login successful",
		Data:    models.LoginResponse{Token: token, Role: role},
	})
}

func (h *AuthHandler) Logout(ctx *gin.Context) {
	token, err := utils.GetToken(ctx)
	if err != nil {
		utils.HandleError(ctx, http.StatusUnauthorized, "Unauthorized", "failed get token")
		return
	}

	expiresAt, err := utils.GetExpiredFromJWT(ctx)
	if err != nil {
		utils.HandleError(ctx, http.StatusUnauthorized, "Unauthorized", "failed get expired time token")
		return
	}

	expiresIn := time.Until(expiresAt)
	if expiresIn <= 0 {
		utils.HandleError(ctx, http.StatusUnauthorized, "Unauthorized", "token already expired")
		return
	}

	if err = h.Repo.BlacklistToken(context.Background(), token, expiresIn); err != nil {
		utils.HandleError(ctx, http.StatusInternalServerError, "Internal Server Error", "failed to blacklist token")
		return
	}

	ctx.JSON(http.StatusOK, models.Response[string]{
		Success: true,
		Message: "Successfully logged out",
		Data:    "",
	})
}
