package api

import (
	// "github.com/a-h/templ"
	"github.com/a-h/templ"
	"github.com/emarifer/gocms/internal/app/service"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

type API struct {
	serv service.Service
}

func New(serv service.Service) *API {

	return &API{
		serv: serv,
	}
}

func (a *API) Start(e *gin.Engine, address string) error {
	e.Use(gzip.Gzip(gzip.DefaultCompression)) // gzip compression middleware
	e.Use(a.globalErrorHandler())             // Error handler middleware
	e.MaxMultipartMemory = 1                  // 8 MiB max. request

	e.Static("/assets", "./assets")
	e.LoadHTMLGlob("views/**/*")

	a.registerRoutes(e)

	return e.Run(address)
}

func (a *API) registerRoutes(e *gin.Engine) {
	e.GET("/", a.homeHandler)
	e.GET("/post/:id", a.postHandler)

	e.GET("/contact", a.contactHandler)
	e.POST("/contact", a.contactHandler)
}

// Markdown to HTML conversion
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

// This function will render the templ component into
// a gin context's Response Writer
func (a *API) RenderView( // TODO: to lowerCamelCase
	c *gin.Context, status int, cmp templ.Component,
) error {
	c.Status(status)

	return cmp.Render(c.Request.Context(), c.Writer)
}
