package admin_endpoint_tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/emarifer/gocms/internal/admin_app/api/dto"
	"github.com/emarifer/gocms/internal/model"
	"github.com/emarifer/gocms/tests/system_tests/helpers"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"

	_ "github.com/go-sql-driver/mysql"
)

type AddPostResponse struct {
	Id int `json:"id"`
}

func TestAdminPostExists(t *testing.T) {

	// This is gonna be the in-memory mysql
	appSettings := helpers.GetAppSettings(4)

	go helpers.RunDatabaseServer(appSettings)

	ctx := context.Background()
	dbConn, err := helpers.WaitForDb(ctx, appSettings)
	require.Nil(t, err)

	// make migrations
	goose.SetBaseFS(helpers.EmbedMigrations)

	if err := goose.SetDialect("mysql"); err != nil {
		require.Nil(t, err)
	}

	if err := goose.Up(dbConn.DB, "migrations"); err != nil {
		require.Nil(t, err)
	}

	// send the post in
	addPostRequest := dto.AddPostRequest{
		Title:   "Test Post Title",
		Excerpt: "test post excerpt",
		Content: "test post content",
	}

	e, err := helpers.StartAdminApp(ctx, appSettings, dbConn)
	require.Nil(t, err)

	w := httptest.NewRecorder()
	requestBytes, err := json.Marshal(addPostRequest)
	require.Nil(t, err)
	req, err := http.NewRequest(
		"POST", "/api/v1/post", bytes.NewBuffer(requestBytes),
	)
	require.Nil(t, err)
	e.ServeHTTP(w, req)

	require.Equal(t, http.StatusCreated, w.Code)

	// make sure it's an expected response
	addPostResponse := &AddPostResponse{}
	err = json.Unmarshal(w.Body.Bytes(), addPostResponse)
	require.Nil(t, err)

	// Get the post
	w = httptest.NewRecorder()
	req, err = http.NewRequest(
		"GET", fmt.Sprintf("/api/v1/post/%d", addPostResponse.Id), nil,
	)
	require.Nil(t, err)
	e.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	// make sure it's the expected content
	getPostResponse := &model.Post{}
	err = json.Unmarshal(w.Body.Bytes(), getPostResponse)
	require.Nil(t, err)

	require.Equal(t, getPostResponse.ID, addPostResponse.Id)
	require.Equal(t, getPostResponse.Title, addPostRequest.Title)
	require.Equal(t, getPostResponse.Excerpt, addPostRequest.Excerpt)
	require.Equal(t, getPostResponse.Content, addPostRequest.Content)
}
