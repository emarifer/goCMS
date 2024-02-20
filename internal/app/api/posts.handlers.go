package api

import (
	"fmt"
	"html/template"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/emarifer/gocms/internal/app/api/dto"
	"github.com/emarifer/gocms/timezone_conversion"
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

	c.HTML(http.StatusOK, "home", gin.H{
		"title": "| Home",
		"posts": posts,
		"year":  time.Now().Year(),
	})
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
	md := string(a.mdToHTML([]byte(post.Content)))

	c.HTML(http.StatusOK, "post", gin.H{
		"title":       fmt.Sprintf("| %s", post.Title),
		"postTitle":   post.Title,
		"postContent": template.HTML(md),
		"createdAt":   timezone_conversion.ConvertDateTime(tz, post.CreatedAt),
		"year":        time.Now().Year(),
	})
}
