package api

import (
	"log/slog"
	"net/http"
	"regexp"

	"github.com/a-h/templ"
	"github.com/emarifer/gocms/internal/service"
	"github.com/emarifer/gocms/settings"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type API struct {
	serv     service.Service
	logger   *slog.Logger
	settings *settings.AppSettings
}

func New(
	serv service.Service, logger *slog.Logger, settings *settings.AppSettings,
) *API {

	return &API{
		serv:     serv,
		logger:   logger,
		settings: settings,
	}
}

type generator = func(*gin.Context) ([]byte, *customError)

var re = regexp.MustCompile(`Table|refused`)

func (a *API) Start(e *gin.Engine, address string, cache *Cache) error {
	e.Use(gzip.Gzip(gzip.DefaultCompression)) // gzip compression middleware
	e.Use(a.globalErrorHandler())             // Error handler middleware
	e.MaxMultipartMemory = 1                  // 8 MiB max. request

	e.Static("/assets", "./assets")
	e.Static("/media", a.settings.ImageDirectory)
	// e.LoadHTMLGlob("views/**/*") // Used for Go Html templates

	a.registerRoutes(e, cache)

	return e.Run(address)
}

func (a *API) registerRoutes(e *gin.Engine, cache *Cache) {
	// ↓ injected into the Start function ↓
	// logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	// cache := MakeCache(4, time.Minute*1)

	// All cache-able endpoints
	a.addCachableHandler(e, "GET", "/", a.homeHandler, cache)
	a.addCachableHandler(e, "GET", "/post/:id", a.postHandler, cache)
	a.addCachableHandler(e, "GET", "/contact", a.contactHandler, cache)

	// Do not cache as it needs to handle new form values
	e.POST("/contact-send", a.contactFormHandler)
}

// mdToHTML converts markdown to HTML
func (a *API) mdToHTML(md []byte) []byte {
	// create markdown parser with extensions
	extensions := parser.CommonExtensions | parser.AutoHeadingIDs | parser.NoEmptyLineBeforeBlock
	p := parser.NewWithExtensions(extensions)
	doc := p.Parse(md)

	// create HTML renderer with extensions
	htmlFlags := html.CommonFlags | html.HrefTargetBlank
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)

	return markdown.Render(doc, renderer)
}

// renderView will render the templ component into
// a gin context's Response Writer
func (a *API) renderView(
	c *gin.Context, status int, cmp templ.Component,
) error {
	c.Status(status)

	return cmp.Render(c.Request.Context(), c.Writer)
}

func (a *API) addCachableHandler(
	e *gin.Engine, method, endpoint string, gen generator, cache *Cache,
) {

	handler := func(c *gin.Context) {
		var errCache error
		// if endpoint is cached, get it from cache
		cachedEndpoint, errCache := (*cache).Get(c.Request.RequestURI)
		if errCache == nil {
			c.Data(
				http.StatusOK,
				"text/html; charset=utf-8",
				cachedEndpoint.Contents,
			)

			return
		} else {
			a.logger.Info(
				"cache info",
				slog.String("could not get page from cache", errCache.Error()),
			)
		}

		// If the endpoint data is not recovered, the handler (gen) is called
		html_buffer, err := gen(c)
		if err != nil {
			c.Error(err)

			return
		}

		// After handler call, add to cache
		errCache = (*cache).Store(c.Request.RequestURI, html_buffer)
		if errCache != nil {
			a.logger.Warn(
				"cache warning",
				slog.String("could not add page to cache", errCache.Error()),
			)
		}

		c.Data(http.StatusOK, "text/html; charset=utf-8", html_buffer)
	}

	// Hacky
	switch method {
	case "GET":
		e.GET(endpoint, handler)
	case "POST":
		e.POST(endpoint, handler)
	case "DELETE":
		e.DELETE(endpoint, handler)
	case "PUT":
		e.PUT(endpoint, handler)
	}
}
