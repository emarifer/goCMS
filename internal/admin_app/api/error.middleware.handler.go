package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type customError struct {
	Code      int    `json:"code"`
	ErrMsg    string `json:"errMsg"`
	CustomMsg string `json:"customMsg"`
}

func (ce *customError) Error() string {
	return fmt.Sprintf("Error %d: %s. %s", ce.Code, ce.ErrMsg, ce.CustomMsg)
}

func NewCustomError(code int, errMsg, customMsg string) *customError {
	return &customError{
		Code:      code,
		ErrMsg:    errMsg,
		CustomMsg: customMsg,
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
				c.JSON(e.Code, gin.H{
					"errMsg":    e.ErrMsg,
					"customMsg": e.CustomMsg,
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
