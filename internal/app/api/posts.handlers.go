package api

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/emarifer/gocms/internal/app/api/dto"
	"github.com/emarifer/gocms/views"
	"github.com/gin-gonic/gin"
)

// This function will act as the handler for
// the home page
func (a *API) homeHandler(c *gin.Context) {
	ctx := c.Request.Context()
	re := regexp.MustCompile(`Table|refused`)

	posts, err := a.serv.RecoverPosts(ctx)
	if err != nil {
		if re.MatchString(err.Error()) {
			err := NewCustomError(
				http.StatusInternalServerError,
				"An unexpected condition was encountered. Please wait a few moments to reload the page.",
			)
			c.Error(err)
		}

		return
	}

	a.renderView(c, http.StatusOK, views.MakePage(
		"| Home",
		"",
		views.Home(posts),
	))
}

// This function will act as a controller
// the post details display page
func (a *API) postHandler(c *gin.Context) {
	tz := c.GetHeader("X-TimeZone")
	ctx := c.Request.Context()
	postBinding := &dto.PostBinding{}
	re := regexp.MustCompile(`Table|refused`)

	// localhost:8080/post/{id}
	if err := c.ShouldBindUri(postBinding); err != nil {
		// TODO redo this error to serve error page
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})

		return
	}

	// Get the post with the ID

	postId, err := strconv.Atoi(postBinding.Id)
	if err != nil {
		if strings.Contains(err.Error(), "strconv.Atoi") {
			err := NewCustomError(
				http.StatusBadRequest,
				"Invalid URL.",
			)
			c.Error(err)
		}

		return
	}

	post, err := a.serv.RecoverPost(ctx, postId)
	if err != nil {
		if re.MatchString(err.Error()) {
			err := NewCustomError(
				http.StatusInternalServerError,
				"An unexpected condition was encountered. Please wait a few moments to reload the page.",
			)
			c.Error(err)
		}

		if strings.Contains(err.Error(), "no rows in result set") {
			err := NewCustomError(
				http.StatusNotFound,
				"The requested resource could not be found but may be available again in the future.",
			)
			c.Error(err)
		}

		return
	}

	// Markdown to HTML the post content
	post.Content = string(a.mdToHTML([]byte(post.Content)))

	a.renderView(c, http.StatusOK, views.MakePage(
		fmt.Sprintf("| %s", post.Title),
		"",
		views.Post(*post, tz),
	))
}
