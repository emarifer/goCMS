package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/emarifer/gocms/database"
	"github.com/emarifer/gocms/internal/app/api"
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
			address := fmt.Sprintf(":%d", s.WebserverPort)
			cache := api.MakeCache(4, time.Minute*10)
			go func() {
				// e.Logger.Fatal(a.Start(e, address))
				a.Start(e, address, &cache)
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			// return e.Close()
			return nil
		},
	})
}

/* REFERENCES:
https://www.youtube.com/playlist?list=PLZ51_5WcvDvDBhgBymjAGEk0SR1qmOh4q
https://github.com/matheusgomes28/urchin

https://github.com/golang-standards/project-layout
https://www.hostinger.es/tutoriales/headless-cms
https://www.bacancytechnology.com/blog/golang-web-frameworks
https://github.com/gin-gonic/gin

Gin vs. Echo:
https://www.codeavail.com/blog/gin-vs-echo/

The client-side-templates Extension for HTMX:
https://htmx.org/extensions/client-side-templates/

Error handler middleware:
https://blog.ruangdeveloper.com/membuat-global-error-handler-go-gin/

Model binding and validation:
https://gin-gonic.com/docs/examples/binding-and-validation/
https://gin-gonic.com/docs/examples/bind-uri/
https://dev.to/ankitmalikg/api-validation-in-gin-ensuring-data-integrity-in-your-api-2p40#:~:text=To%20perform%20input%20validation%20in,query%20parameters%2C%20and%20request%20bodies.
https://www.google.com/search?q=validation+function+gin&oq=validation+function+gin&aqs=chrome..69i57j33i160j33i671.11287j0j7&sourceid=chrome&ie=UTF-8

Go: how to check if a string contains multiple substrings?:
https://stackoverflow.com/questions/47131996/go-how-to-check-if-a-string-contains-multiple-substrings#47134680

Debugging Go with VSCode and Air:
https://www.arhea.net/posts/2023-08-25-golang-debugging-with-air-and-vscode/

What is the purpose of .PHONY in a Makefile?:
https://stackoverflow.com/questions/2145590/what-is-the-purpose-of-phony-in-a-makefile

A-H.TEMPL:
https://templ.guide/syntax-and-usage/template-composition/
https://templ.guide/syntax-and-usage/rendering-raw-html/

INJECTION OF A LOGGER WITH FX:
https://uber-go.github.io/fx/container.html#providing-values

USING LOG/SLOG:
https://www.youtube.com/watch?v=bDpB6k-Q_GY
https://github.com/disturb16/go_examples/tree/main/btrlogs
https://betterstack.com/community/guides/logging/logging-in-go/
https://go.dev/blog/slog

SHARDEDMAP:
https://pkg.go.dev/github.com/zutto/shardedmap?utm_source=godoc
https://github.com/zutto/shardedmap

BENCHMARK WITH/WITHOUT CACHE:
https://github.com/fcsonline/drill

COMMAND-LINE FLAGS:
https://gobyexample.com/command-line-flags

DECODING AND ENCODING OF TOML FILES:
https://github.com/BurntSushi/toml
https://godocs.io/github.com/BurntSushi/toml

MISCELLANEOUS:
https://github.com/a-h/templ/tree/1f30f822a6edfdbfbab9e6851b1ff61e0ab01d4f/examples/integration-gin

https://github.com/stackus/todos

https://toml.io/en/

https://github.com/pelletier/go-toml

[Git: See my last commit]
https://stackoverflow.com/questions/2231546/git-see-my-last-commit
*/

/* CHECKS FUNCTIONS:

// settings test:
func(s *settings.Settings) {
	fmt.Println(s.DB.Name)
},

// Database operation check function:
func(db *sqlx.DB) {
	_, err := db.Query("SELECT * FROM users")
	if err != nil {
		panic(err)
	}
},
*/
