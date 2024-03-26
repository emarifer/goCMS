package api

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/emarifer/gocms/views"
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
				// SEE NOTE BELOW (this is hacky):
				if os.Getenv("GO_ENV") == "testing" {
					a.logger.Error(
						"error message",
						slog.String("could not render HTML", err.Error()),
					)
					c.JSON(e.Code, gin.H{
						"error": "could not render HTML",
						"msg":   e.Error(),
					})
				} else {
					a.renderView(c, e.Code, views.MakePage(
						fmt.Sprintf("| Error %d", e.Code),
						e.Message,
						views.ErrorPage(e.Message),
					))
				}
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

/* HOW DO I KNOW I'M RUNNING WITHIN "GO TEST". SEE:
https://stackoverflow.com/questions/14249217/how-do-i-know-im-running-within-go-test#59444829
*/
