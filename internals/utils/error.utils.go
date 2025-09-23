package utils

import (
	"log"

	"github.com/Ntisrangga142/API_tickytiz/internals/models"
	"github.com/gin-gonic/gin"
)

func HandleError(ctx *gin.Context, code int, status, err string) {
	log.Printf("%s\nCause: %s\n", status, err)
	ctx.JSON(code, models.NewErrorResponse(status, err, code))
}

func HandleMiddlewareError(ctx *gin.Context, code int, status, err string) {
	log.Printf("%s\nCause: %s\n", status, err)
	ctx.AbortWithStatusJSON(code, models.NewErrorResponse(status, err, code))
}
