package api

import (
	"regexp"

	"github.com/emarifer/gocms/internal/service"
	"github.com/emarifer/gocms/settings"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type API struct {
	serv          service.Service
	dataValidator *validator.Validate
	settings      *settings.AppSettings
}

type ErrorResponse struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value,omitempty"`
}

func New(
	serv service.Service,
	dataValidator *validator.Validate,
	settings *settings.AppSettings,
) *API {

	return &API{
		serv:          serv,
		dataValidator: dataValidator,
		settings:      settings,
	}
}

var re = regexp.MustCompile(`Table|refused`)

func (a *API) Start(e *gin.Engine, address string) error {
	e.Use(gzip.Gzip(gzip.DefaultCompression)) // gzip compression middleware
	e.Use(a.globalErrorHandler())             // Error handler middleware
	e.MaxMultipartMemory = 1                  // 8 MiB max. request

	a.registerRoutes(e)

	return e.Run(address)
}

func (a *API) registerRoutes(e *gin.Engine) {
	v1 := e.Group("api/v1")
	v1.POST("/post", a.addPostHandler)
	v1.GET("/post", a.getPostsHandler)
	v1.GET("/post/:id", a.postHandler)
	v1.PUT("/post", a.updatePostHandler)
	v1.DELETE("/post/:id", a.deletePostHandler)

	v1.GET("/image/:uuid", a.getImageHandler)
	v1.GET("/image", a.getAllImagesHandler)
	v1.POST("image", a.addImageHandler)
	v1.DELETE("/image/:uuid", a.deleteImageHandler)

	/* e.GET("/contact", a.contactHandler)
	e.POST("/contact", a.contactHandler) */
}

func (a *API) validateStruct(payload any) []ErrorResponse { // SEE NOTE BELOW
	var errors []ErrorResponse

	err := a.dataValidator.Struct(payload)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ErrorResponse
			element.Field = err.StructNamespace()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, element)
		}
	}

	return errors
}

/*
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
func (a *API) renderView(
	c *gin.Context, status int, cmp templ.Component,
) error {
	c.Status(status)

	return cmp.Render(c.Request.Context(), c.Writer)
}
*/

/* NO PARAMETERIZED METHODS. SEE:
https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#No-parameterized-methods

JSON Escape / Unescape Online. SEE:
https://www.freeformatter.com/json-escape.html#before-output

What HTTP status code should I return for POST when no resource is created?. SEE:
https://stackoverflow.com/questions/55685576/what-http-status-code-should-i-return-for-post-when-no-resource-is-created
*/
