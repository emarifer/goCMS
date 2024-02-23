package main

import (
	"context"
	"fmt"

	"github.com/emarifer/gocms/database"
	"github.com/emarifer/gocms/internal/admin_app/api"
	"github.com/emarifer/gocms/internal/repository"
	"github.com/emarifer/gocms/internal/service"
	"github.com/emarifer/gocms/settings"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func main() {
	gocms := fx.New(
		fx.Provide(
			context.Background,
			settings.New,
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
	s *settings.Settings,
	e *gin.Engine,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			address := fmt.Sprintf(":%s", s.AdminPort)
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
