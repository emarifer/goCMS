package main

import (
	"context"
	"fmt"

	"github.com/emarifer/gocms/database"
	"github.com/emarifer/gocms/internal/app/api"
	"github.com/emarifer/gocms/internal/app/repository"
	"github.com/emarifer/gocms/internal/app/service"
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
			address := fmt.Sprintf(":%s", s.Port)
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
