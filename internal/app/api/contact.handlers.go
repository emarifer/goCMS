package api

import (
	"net/http"
	"net/mail"
	"time"

	"github.com/gin-gonic/gin"
)

// This function will act as a controller
// for the Contact page and the POST request
// made by the contact form
func (a *API) contactHandler(c *gin.Context) {
	if c.Request.Method == "POST" {
		name := c.Request.FormValue("name") // name := c.PostForm("name")
		email := c.Request.FormValue("email")
		message := c.Request.FormValue("message")

		// Check email
		_, err := mail.ParseAddress(email)
		if err != nil {
			c.HTML(http.StatusOK, "contact-failure.html", gin.H{
				"email": email,
				"error": "invalid email",
			})

			return
		}

		// Make sure name and message is reasonable
		if len(name) > 200 {
			c.HTML(http.StatusOK, "contact-failure.html", gin.H{
				"email": email,
				"error": "enter a name of less than 200 characters",
			})

			return
		}

		if len(message) > 1000 {
			c.HTML(http.StatusOK, "contact-failure.html", gin.H{
				"email": email,
				"error": "message too big",
			})

			return
		}

		c.HTML(http.StatusOK, "contact-success.html", gin.H{
			"email": email,
			"name":  name,
		})

		return
	}

	c.HTML(http.StatusOK, "contact", gin.H{
		"title": " | Contact",
		"year":  time.Now().Year(),
	})
}
