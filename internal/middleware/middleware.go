package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Logger is a middleware function that logs the request method, path, and time
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Starting time
		startTime := time.Now()

		// Process request
		c.Next()

		// End time
		endTime := time.Now()

		// Execution time
		latency := endTime.Sub(startTime)

		// Request information
		method := c.Request.Method
		path := c.Request.URL.Path
		statusCode := c.Writer.Status()

		// Log request details
		log.Printf("[%s] %s %s %d %s", method, path, latency, statusCode, c.ClientIP())
	}
}

// CORS is a middleware function that adds CORS headers to the response
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// ErrorHandler is a middleware function that handles errors
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			// Log the errors
			for _, err := range c.Errors {
				log.Printf("Error: %s", err.Error())
			}

			// Return the first error to the client
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": c.Errors[0].Error(),
			})
		}
	}
}