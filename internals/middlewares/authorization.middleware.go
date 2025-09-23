package middlewares

import (
	"net/http"
	"slices"

	"github.com/Ntisrangga142/API_tickytiz/internals/utils"
	"github.com/Ntisrangga142/API_tickytiz/pkg"
	"github.com/gin-gonic/gin"
)

func Authorization(roles ...string) func(*gin.Context) {
	return func(ctx *gin.Context) {
		claims, isExist := ctx.Get("claims")
		if !isExist {
			utils.HandleMiddlewareError(ctx, http.StatusUnauthorized, "Unauthorized", "Unauthorized Access")
			return
		}
		user, ok := claims.(pkg.Claims)
		if !ok {
			utils.HandleMiddlewareError(ctx, http.StatusInternalServerError, "Internal Server Error", "Cannot Cast into pkg.claims")
			return
		}
		if !slices.Contains(roles, user.Role) {
			utils.HandleMiddlewareError(ctx, http.StatusForbidden, "For Bidden", "You don't have access rights to this resource.")
			return
		}
		ctx.Next()
	}
}
