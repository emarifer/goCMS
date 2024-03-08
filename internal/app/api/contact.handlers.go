package api

import (
	"bytes"
	"net/http"
	"net/mail"

	"github.com/emarifer/gocms/views"
	"github.com/gin-gonic/gin"
)

// contactHandler will act as a controller
// for the Contact page
func (a *API) contactHandler(c *gin.Context) ([]byte, *customError) {

	cmp := views.MakePage("| Contact", "", views.Contact())
	html_buffer := bytes.NewBuffer(nil)
	err := cmp.Render(c.Request.Context(), html_buffer)
	if err != nil {
		err := NewCustomError(
			http.StatusInternalServerError,
			"An unexpected condition was encountered. The requested page could not be rendered.",
		)

		return nil, err
	}

	return html_buffer.Bytes(), nil
}

// contactFormHandle will act as a controller
// for the POST request made by the contact form
func (a *API) contactFormHandler(c *gin.Context) {
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
}
