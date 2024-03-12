package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/emarifer/gocms/database"
	"github.com/emarifer/gocms/internal/admin_app/api"
	"github.com/emarifer/gocms/internal/repository"
	"github.com/emarifer/gocms/internal/service"
	"github.com/emarifer/gocms/settings"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/fx"
)

func main() {
	gocms := fx.New(
		fx.Provide(
			context.Background,
			func() *os.File { return os.Stdout },
			func() *slog.JSONHandler {
				return slog.NewJSONHandler(os.Stdout, nil)
			},
			func(h *slog.JSONHandler) *slog.Logger { return slog.New(h) },
			func(l *slog.Logger) *string {
				config_toml := flag.String("config", "", "path to the config to be used")
				flag.Parse()
				l.Info("reading config file", "from path", config_toml)

				return config_toml
			},
			settings.New,
			validator.New,
			database.NewMariaDBConnection,
			repository.New,
			service.New,
			gin.Default,
			api.New,
		),

		fx.Invoke(
			setLifeCycle,
		),
	)

	gocms.Run()
}

func setLifeCycle(
	lc fx.Lifecycle,
	a *api.API,
	s *settings.AppSettings,
	e *gin.Engine,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			address := fmt.Sprintf(":%d", s.AdminWebserverPort)
			go func() {
				// e.Logger.Fatal(a.Start(e, address))
				a.Start(e, address)
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			// return e.Close()
			return nil
		},
	})
}
