package endpoint_tests

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/emarifer/gocms/tests/system_tests/helpers"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"

	_ "github.com/go-sql-driver/mysql"
)

func TestPostExists(t *testing.T) {

	// This is gonna be the in-memory mysql
	appSettings := helpers.GetAppSettings(3)

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

	e, err := helpers.StartApp(ctx, appSettings, dbConn)
	require.Nil(t, err)

	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/post/1", nil)
	require.Nil(t, err)
	e.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)
}
