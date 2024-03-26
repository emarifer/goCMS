package api

import (
	"bytes"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/emarifer/gocms/internal/app/api/dto"
	"github.com/emarifer/gocms/views"
	"github.com/gin-gonic/gin"
)

// homeHandler will act as the handler for
// the home page
func (a *API) homeHandler(c *gin.Context) ([]byte, *customError) {
	ctx := c.Request.Context()

	posts, err := a.serv.RecoverPosts(ctx)
	if err != nil {
		if re.MatchString(err.Error()) {
			err := NewCustomError(
				http.StatusInternalServerError,
				"An unexpected condition was encountered. Please wait a few moments to reload the page.",
			)

			return nil, err
		} else {
			err := NewCustomError(
				http.StatusInternalServerError,
				fmt.Sprintf(
					"An unexpected condition was encountered: %s", err.Error(),
				),
			)

			return nil, err
		}
	}

	cmp := views.MakePage("| Home", "", views.Home(posts))
	html_buffer := bytes.NewBuffer(nil)
	err = cmp.Render(c.Request.Context(), html_buffer)
	if err != nil {
		err := NewCustomError(
			http.StatusInternalServerError,
			"An unexpected condition was encountered. The requested page could not be rendered.",
		)

		return nil, err
	}

	return html_buffer.Bytes(), nil
}

// postHandler will act as a controller
// the post details display page
func (a *API) postHandler(c *gin.Context) ([]byte, *customError) {
	tz := c.GetHeader("X-TimeZone")
	ctx := c.Request.Context()
	postBinding := &dto.PostBinding{}

	// localhost:8080/post/{id}
	if err := c.ShouldBindUri(postBinding); err != nil {
		err := NewCustomError(
			http.StatusBadRequest,
			"Invalid URL.",
		)

		return nil, err
	}

	// Get the post with the ID

	postId, err := strconv.Atoi(postBinding.Id)
	if err != nil {
		if strings.Contains(err.Error(), "strconv.Atoi") {
			err := NewCustomError(
				http.StatusBadRequest,
				"Invalid URL.",
			)

			return nil, err
		}
	}

	post, err := a.serv.RecoverPost(ctx, postId)
	if err != nil {
		if re.MatchString(err.Error()) {
			err := NewCustomError(
				http.StatusInternalServerError,
				"An unexpected condition was encountered. Please wait a few moments to reload the page.",
			)

			return nil, err
		}

		if strings.Contains(err.Error(), "no rows in result set") {
			err := NewCustomError(
				http.StatusNotFound,
				"The requested resource could not be found but may be available again in the future.",
			)

			return nil, err
		}
	}

	// Markdown to HTML the post content
	post.Content = string(a.mdToHTML([]byte(post.Content)))

	cmp := views.MakePage(
		fmt.Sprintf("| %s", post.Title),
		"",
		views.Post(*post, tz),
	)
	html_buffer := bytes.NewBuffer(nil)
	err = cmp.Render(c.Request.Context(), html_buffer)
	if err != nil {
		err := NewCustomError(
			http.StatusInternalServerError,
			"An unexpected condition was encountered. The requested page could not be rendered.",
		)

		return nil, err
	}

	return html_buffer.Bytes(), nil
}
