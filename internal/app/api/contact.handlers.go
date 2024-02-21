package api

import (
	"net/http"
	"net/mail"

	"github.com/emarifer/gocms/views"
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
			a.renderView(c, http.StatusOK, views.ContactFailure(
				email, "invalid email",
			))

			return
		}

		// Make sure name and message is reasonable
		if len(name) > 200 {
			a.renderView(c, http.StatusOK, views.ContactFailure(
				email, "enter a name of less than 200 characters",
			))

			return
		}

		if len(message) > 1000 {
			a.renderView(c, http.StatusOK, views.ContactFailure(
				email, "message too big",
			))

			return
		}

		a.renderView(c, http.StatusOK, views.ContactSuccess(email, name))

		return
	}

	a.renderView(c, http.StatusOK, views.MakePage(
		"| Contact",
		"",
		views.Contact(),
	))
}
