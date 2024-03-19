package api

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/emarifer/gocms/internal/service"
	"github.com/emarifer/gocms/settings"
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	lua "github.com/yuin/gopher-lua"
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

// Find all the occurences of {{ and }} (including whitespace)
var shortcodesFound = regexp.MustCompile(`{{[\s\w.-]+(:[\s\w.-]+)+}}`)

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

func (a *API) loadShortcodeHandlers() (map[string]*lua.LState, error) {
	shortcodeHandlers := map[string]*lua.LState{}

	for _, shortcode := range a.settings.Shortcodes {
		// Read the LState (Lua state)
		state := lua.NewState()
		err := state.DoFile(shortcode.Plugin)
		// TODO : check that the function HandleShortcode(args)
		//        exists and returns the correct type
		if err != nil {
			return shortcodeHandlers, fmt.Errorf(
				"could not load shortcode %s: %v", shortcode.Name, err,
			)
		}
		shortcodeHandlers[shortcode.Name] = state
	}

	return shortcodeHandlers, nil
}

// partitionString will partition the strings by
// removing the given ranges
func (a *API) partitionString(text string, indexes [][]int) []string {
	partitions := []string{}
	start := 0

	for _, window := range indexes {
		partitions = append(partitions, text[start:window[0]])
		start = window[1]
	}

	partitions = append(partitions, text[start:len(text)-1])

	return partitions
}

func (a *API) shortcodeToMarkdown(
	shortcode string, shortcodeHandlers map[string]*lua.LState,
) (string, error) {
	keyValue := strings.Split(shortcode, ":")
	key := keyValue[0]
	values := keyValue[1:]

	/* if key == "img" {
		if len(values) == 1 {
			imageSrc := fmt.Sprintf("/media/%s", values[0])

			return fmt.Sprintf("![image](%s)", imageSrc), nil
		} else if len(values) == 2 {
			imageSrc := fmt.Sprintf("/media/%s", values[0])
			altText := values[1]

			return fmt.Sprintf("![%s](%s)", altText, imageSrc), nil
		} else {
			return "", fmt.Errorf("invalid shortcode: %s", shortcode)
		}
	} */

	if handler, found := shortcodeHandlers[key]; found {
		// Need to quote all values for a valid lua syntax
		quotedValues := []string{}
		for _, value := range values {
			quotedValues = append(quotedValues, fmt.Sprintf("%q", value))
		}

		err := handler.DoString(fmt.Sprintf(
			`result = HandleShortcode({%s})`, strings.Join(quotedValues, ","),
		))
		if err != nil {
			return "", fmt.Errorf("error running %s shortcode: %v", key, err)
		}

		value := handler.GetGlobal("result")
		if retType := value.Type().String(); retType != "string" {
			return "", fmt.Errorf(
				"error running %s shortcode: invalid return type %s", key, retType,
			)
		} else if retType == "" {
			return "", fmt.Errorf(
				"error running %s shortcode: returned empty string", key,
			)
		}

		return value.String(), nil
	}

	return "", fmt.Errorf("unsupported shortcode: %s", key)
}

func (a *API) transformContent(
	content string, shortcodeHandlers map[string]*lua.LState,
) (string, error) {
	// see note below
	shortcodes := shortcodesFound.FindAllStringIndex(content, -1)
	// content without shortcodes
	partitions := a.partitionString(content, shortcodes)

	builder := strings.Builder{}
	// i := 0
	for i, shortcode := range shortcodes {
		builder.WriteString(partitions[i])

		// +2 is added, or -2 is regressed, due to the double curly brackets
		md, err := a.shortcodeToMarkdown(
			content[shortcode[0]+2:shortcode[1]-2], shortcodeHandlers,
		)
		if err != nil {
			return "", fmt.Errorf("could not transform post: %v", err)
		}

		builder.WriteString(md)
	}

	// Guaranteed to have +1 than the number of
	// shortcodes by algorithm
	if len(shortcodes) > 0 {
		builder.WriteString(partitions[len(partitions)-1])
	} else {
		builder.WriteString(content)
	}

	return builder.String(), nil
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

/* FINDING THE INDEX OF ALL OCCURRENCES OF A REGULAR EXPRESSION IN A STRING. SEE:
https://www.tutorialspoint.com/finding-index-of-the-regular-expression-present-in-string-of-golang
*/

/* NO PARAMETERIZED METHODS. SEE:
https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#No-parameterized-methods

JSON Escape / Unescape Online. SEE:
https://www.freeformatter.com/json-escape.html#before-output

What HTTP status code should I return for POST when no resource is created?. SEE:
https://stackoverflow.com/questions/55685576/what-http-status-code-should-i-return-for-post-when-no-resource-is-created
*/
