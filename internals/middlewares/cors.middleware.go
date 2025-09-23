package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORSMiddleware handle CORS for frontend development
func CORSMiddleware() gin.HandlerFunc {
	// Daftar origin yang diizinkan
	whitelist := []string{
		"http://127.0.0.1:5500",
		"http://127.0.0.1:5173",
		"http://127.0.0.1:3000",
		"http://127.0.0.1:8080",
		"http://localhost:5500",
		"http://localhost:5173",
		"http://localhost:3000",
		"http://localhost:8080",
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")

		// Jika origin ada dan termasuk whitelist
		allowed := false
		for _, o := range whitelist {
			if o == origin {
				allowed = true
				break
			}
		}

		if origin != "" && !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "CORS origin not allowed",
			})
			return
		}

		// Set CORS headers
		if origin != "" && allowed {
			c.Header("Access-Control-Allow-Origin", origin)
		} else {
			c.Header("Access-Control-Allow-Origin", "*")
		}

		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Authorization, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Expose-Headers", "Authorization, Content-Type")

		// Preflight request handling
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
