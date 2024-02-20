package api

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type customError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (ce *customError) Error() string {
	return fmt.Sprintf("Error %d: %s", ce.Code, ce.Message)
}

func NewCustomError(code int, message string) *customError {
	return &customError{
		Code:    code,
		Message: message,
	}
}

func (a *API) globalErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Catch errors that appear in the middleware or handler
		err := c.Errors.Last()
		if err != nil {
			// Handle errors here
			switch e := err.Err.(type) {
			case *customError:
				// Handle custom errors
				c.HTML(e.Code, "error.page", gin.H{
					"title": fmt.Sprintf("| Error %d", e.Code),
					"error": e.Message,
					"year":  time.Now().Year(),
				})
			default:
				// Handle other errors
				c.JSON(500, gin.H{
					"error": err.Error(),
				})
			}

			// Stop context execution
			c.Abort()
		}
	}
}

/* REFERENCES:
https://blog.ruangdeveloper.com/membuat-global-error-handler-go-gin/
*/
