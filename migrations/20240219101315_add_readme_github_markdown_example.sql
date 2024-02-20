-- +goose Up
-- +goose StatementBegin
INSERT INTO posts(title, content, excerpt) VALUES(
    'Second post with Markdown',
    
    '<div align="center">
  
<h1 align="center">Shortify</h1>

<img src="https://raw.githubusercontent.com/emarifer/url-shortener-echo-templ-htmx/main/doc/banner.png" width="100%">

<hr />

<p style="margin-bottom: 8px;">

A Full Stack Url Shortener App using Golang\'s Echo framework + </>HTMX & Templ.

</p>

</div>

<hr />

## 🤔 What Stack have we used?

In the implementation of this application we have used the following technologies:

- ✅ **Bootstrapping**: [Go programming language (v1.21)](https://go.dev/)
- ✅ **Backend Framework**: [Echo v4.11.4 ](https://echo.labstack.com/)
- ✅ **Auth & Session middleware**: [Echo Contrib Session](https://github.com/labstack/echo-contrib/)
- ✅ **Dependency Injection**: [Fx - dependency injection system for Go](https://uber-go.github.io/fx/)
- ✅ **Database**: [PostgreSQL](https://www.postgresql.org/)
- ✅ **Styling**: [TailwindCSS + DaisyUI](https://tailwindcss.com)
- ✅ **Frontend interactivity**: [</>Htmx + _Hyperscript](https://htmx.org/)
- ✅ **Templating Language**: [</>Templ - build HTML with Go](https://templ.guide/)
- ✅ **Popup Boxes (Alerts)**: [Sweetalert2 - responsive & customizable replacement for JavaScript\'s popup boxes](https://sweetalert2.github.io/)
  
<br />


>***The use of [</>htmx](https://htmx.org/) allows behavior similar to that of a SPA, without page reloads when switching from one route to another or when making requests (via AJAX) to the backend. Likewise, the [_hyperscript](https://hyperscript.org/) library allows you to add some dynamic features to the frontend in a very easy way.***
  
<br />
  

>***In this application, instead of using the [html/template](https://pkg.go.dev/html/template) package (native Golang templates), we use the [a-h/templ](https://github.com/a-h/templ) library. This amazing library implements a templating language (very similar to JSX) that compiles to Go code. Templ will allow us to write code almost identical to Go (with expressions, control flow, if/else, for loops, etc.) and have autocompletion thanks to strong typing. This means that errors appear at compile time and any calls to these templates (which are compiled as Go functions) from the handlers side will always require the correct data, minimizing errors and thus increasing the security and speed of our coding.***

---

## 🖼️ Screenshots:

<div align="center">

<h6>Home & Login Pages with success alert:</h6>

<img src="https://raw.githubusercontent.com/emarifer/url-shortener-echo-templ-htmx/main/doc/screenshot-01.png" width="35%" align="top">&nbsp;&nbsp;&nbsp;&nbsp;<img src="https://raw.githubusercontent.com/emarifer/url-shortener-echo-templ-htmx/main/doc/screenshot-02-v.2.png" width="35%">

<h6>Short link creator Page and Dashboard Page with alerts:</h6>

<img src="https://raw.githubusercontent.com/emarifer/url-shortener-echo-templ-htmx/main/doc/screenshot-03.png" width="35%" align="top">&nbsp;&nbsp;&nbsp;&nbsp;<img src="https://raw.githubusercontent.com/emarifer/url-shortener-echo-templ-htmx/main/doc/screenshot-04-v.2.png" width="35%">

<h6>Dashboard Page with alert and Short Link Update Modal:</h6>

<img src="https://raw.githubusercontent.com/emarifer/url-shortener-echo-templ-htmx/main/doc/screenshot-05.png" width="35%" align="top">&nbsp;&nbsp;&nbsp;&nbsp;<img src="https://raw.githubusercontent.com/emarifer/url-shortener-echo-templ-htmx/main/doc/screenshot-06.png" width="35%">

<h6>Centralized HTTP error handling:</h6>

<img src="https://raw.githubusercontent.com/emarifer/url-shortener-echo-templ-htmx/main/doc/screenshot-07.png" width="40%" align="top">

</div>

---

## 📦 Project structure

```
- assets
  |- css
  |- img
- database
- encryption
- internal
  |- api
    |- dto
  |- entity
  |- model
  |- repository
  |- service
- postgresdb_init
- settings
- timezone_conversion
- views
  |- auth_views
  |- components
  |- errors_pages
  |- layout
  |- links_views
```

<br />

<img src="https://raw.githubusercontent.com/emarifer/url-shortener-echo-templ-htmx/main/doc/flow.svg" width="100%">

<br />


>***In this application, we have tried to apply a [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html) pattern. The architecture follows a typical "onion model" where each layer doesn\'t know about the layer above it, and each layer is responsible for a specific thing. Although the application is extremely simple, we use this pattern to illustrate its use in more complex applications. Layering an application in this way can simplify code structure, since the responsibility of each type is clear. To ensure that each part of the application is initialized with its dependencies, each struct defines a constructor (the New function in this example). Related to the latter, we have used a dependency injector ([Fx](https://uber-go.github.io/fx/) from uber-go), which helps to remove the global state in the application and add new components and have them instantly accessible across the application.***

<br />

>***When applying this approach in a real-life application, as with most things, taking the layering approach to an extreme level can have a negative effect. Ask yourself if what you are doing actually helps make the code understandable or simply spreads the application logic across many files and makes the overall structure difficult to see.***


## 👨‍🚀 Getting Started

Besides the obvious prerequisite of having Go! on your machine, you must have Air installed for hot reloading when editing code.

Since we use the PostgreSQL database from a Docker container, it is necessary to have the latter also installed and execute this command in the project folder:

```
$ docker compose up -d
```

These other commands will also be useful to manage the database from its container:

```
$ docker container start shortify-db # start container
$ docker container stop shortify-db # stop container
$ docker exec -it shortify-db psql -U admin -W shortify_db # (pass: admin) access the database
```

Download the necessary dependencies:

```
$ go mod tidy
```

Start the app in development mode:

```
$ air # Ctrl + C to stop the application
```

Build for production:

```
$ go build -ldflags="-s -w" -o ./bin/main . # ./bin/main to run the application / Ctrl + C to stop the application
```


>***In order to have autocompletion and syntax highlighting in VS Code for the Teml templating language, you will have to install the [templ-vscode](https://marketplace.visualstudio.com/items?itemName=a-h.templ) extension (for vim/nvim install this [plugin](https://github.com/joerdav/templ.vim)). To generate the Go code corresponding to these templates you will have to download this [executable binary](https://github.com/a-h/templ/releases/tag/v0.2.476) from Github and place it in the PATH of your system. The command:***

```
$ templ generate --watch
```

>***This will allow us to monitor changes to the .templ files and compile them as we save them. Review the documentation on Templ [installation](https://templ.guide/quick-start/installation) and [support](https://templ.guide/commands-and-tools/ide-support/) for your IDE.***

---

### Happy coding 😀!!',
    'Example post of markdown rendering with its own syntax from the GitHub README.md file.'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM posts ORDER BY id DESC LIMIT 1;
-- +goose StatementEnd
