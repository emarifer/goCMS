package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/emarifer/gocms/internal/admin_app/api/dto"
	"github.com/emarifer/gocms/internal/entity"
	"github.com/emarifer/gocms/internal/model"
	"github.com/gin-gonic/gin"
)

func (a *API) getPostsHandler(c *gin.Context) {
	ctx := c.Request.Context()

	posts, err := a.serv.RecoverPosts(ctx)
	if err != nil {
		if re.MatchString(err.Error()) {
			customError := NewCustomError(
				http.StatusInternalServerError,
				err.Error(),
				"An unexpected condition was encountered.",
			)
			c.Error(customError)

			return
		}
	}

	pps := []model.PostSummary{}
	for _, post := range posts {
		ps := model.PostSummary{
			ID:      post.ID,
			Title:   post.Title,
			Excerpt: post.Excerpt,
		}
		pps = append(pps, ps)
	}

	c.JSON(http.StatusOK, gin.H{
		"post": pps,
	})
}

func (a *API) postHandler(c *gin.Context) {
	ctx := c.Request.Context()
	postBinding := &dto.PostBinding{}

	// localhost:8080/post/{id}
	if err := c.ShouldBindUri(postBinding); err != nil {
		customError := NewCustomError(
			http.StatusBadRequest,
			err.Error(),
			"could not get post id",
		)
		c.Error(customError)

		return
	}

	// Get the post with the ID

	postId, err := strconv.Atoi(postBinding.Id)
	if err != nil {
		customError := NewCustomError(
			http.StatusBadRequest,
			err.Error(),
			"invalid post id type",
		)
		c.Error(customError)

		return
	}

	post, err := a.serv.RecoverPost(ctx, postId)
	if err != nil {
		if re.MatchString(err.Error()) {
			customError := NewCustomError(
				http.StatusInternalServerError,
				err.Error(),
				"An unexpected condition was encountered.",
			)
			c.Error(customError)

			return
		}

		if strings.Contains(err.Error(), "no rows in result set") {
			customError := NewCustomError(
				http.StatusNotFound,
				err.Error(),
				"The requested resource could not be found but may be available again in the future.",
			)
			c.Error(customError)

			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         post.ID,
		"title":      post.Title,
		"excerpt":    post.Excerpt,
		"content":    post.Content,
		"created_at": post.CreatedAt,
	})
}

func (a *API) addPostHandler(c *gin.Context) {
	ctx := c.Request.Context()
	addPostBinding := &dto.AddPostRequest{}

	if c.Request.Body == nil {
		customError := NewCustomError(
			http.StatusBadRequest,
			"invalid request",
			"null request body",
		)
		c.Error(customError)

		return
	}

	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(addPostBinding)
	if err != nil {
		customError := NewCustomError(
			http.StatusBadRequest,
			err.Error(),
			"invalid request body",
		)
		c.Error(customError)

		return
	}

	// Validating request body data
	if errors := a.validateStruct(*addPostBinding); errors != nil {
		errMsgs := []string{}
		for _, err := range errors {
			errStr := fmt.Sprintf("%#v", err)
			errMsgs = append(errMsgs, errStr)
		}
		customError := NewCustomError(
			http.StatusBadRequest,
			strings.Join(errMsgs, "\n"),
			"errors occurred while validating request body",
		)
		c.Error(customError)

		return
	}

	// Getting shortcodeHandlers
	shortcodeHandlers, err := a.loadShortcodeHandlers()
	if err != nil {
		customError := NewCustomError(
			http.StatusInternalServerError,
			err.Error(),
			"the server encountered an unexpected condition that prevented add the post",
		)
		c.Error(customError)

		return
	}

	transformedContent, err := a.transformContent(
		addPostBinding.Content, shortcodeHandlers,
	)
	if err != nil {
		customError := NewCustomError(
			http.StatusInternalServerError,
			err.Error(),
			"the server encountered an unexpected condition that prevented add the post",
		)
		c.Error(customError)

		return
	}

	post := &entity.Post{
		Title:   addPostBinding.Title,
		Excerpt: addPostBinding.Excerpt,
		Content: transformedContent,
	}

	id, err := a.serv.CreatePost(ctx, post)
	if err != nil {
		customError := NewCustomError(
			http.StatusInternalServerError,
			err.Error(),
			"the server encountered an unexpected condition that prevented add the post",
		)
		c.Error(customError)

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id": id,
	})
}

func (a *API) updatePostHandler(c *gin.Context) {
	ctx := c.Request.Context()
	updatePostRequest := &dto.UpdatePostRequest{}
	decoder := json.NewDecoder(c.Request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(updatePostRequest); err != nil {
		customError := NewCustomError(
			http.StatusBadRequest,
			err.Error(),
			"invalid request body",
		)
		c.Error(customError)

		return
	}

	post := &entity.Post{
		ID:      updatePostRequest.ID,
		Title:   updatePostRequest.Title,
		Excerpt: updatePostRequest.Excerpt,
		Content: updatePostRequest.Content,
	}

	rowsAffected, err := a.serv.ChangePost(ctx, post)
	if err != nil {
		// fmt.Println(err)
		if re.MatchString(err.Error()) {
			customError := NewCustomError(
				http.StatusInternalServerError,
				err.Error(),
				"An unexpected condition was encountered",
			)
			c.Error(customError)

			return
		}
	}

	if rowsAffected == -1 {
		customError := NewCustomError(
			http.StatusNotFound,
			"no rows in result set",
			"could not change post",
		)
		c.Error(customError)

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rowsAffected": rowsAffected,
	})
}

func (a *API) deletePostHandler(c *gin.Context) {
	ctx := c.Request.Context()
	postBinding := &dto.PostBinding{}

	if err := c.ShouldBindUri(postBinding); err != nil {
		customError := NewCustomError(
			http.StatusBadRequest,
			err.Error(),
			"could not get post id",
		)
		c.Error(customError)

		return
	}

	// Remove the post with the ID

	postId, err := strconv.Atoi(postBinding.Id)
	if err != nil {
		customError := NewCustomError(
			http.StatusBadRequest,
			err.Error(),
			"invalid post id type",
		)
		c.Error(customError)

		return
	}

	rowsAffected, err := a.serv.RemovePost(ctx, postId)
	if err != nil {
		if re.MatchString(err.Error()) {
			customError := NewCustomError(
				http.StatusInternalServerError,
				err.Error(),
				"An unexpected condition was encountered",
			)
			c.Error(customError)

			return
		}
	}

	if rowsAffected == -1 {
		customError := NewCustomError(
			http.StatusNotFound,
			"no rows in result set",
			"could not delete post",
		)
		c.Error(customError)

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rowsAffected": rowsAffected,
	})
}
