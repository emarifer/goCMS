-- +goose Up
-- +goose StatementBegin
INSERT INTO posts(title, content, excerpt) VALUES(
    'Gofiber Templ Htmx',
    '# Go/Fiber🧬+<span style="color:yellow"></></span>Templ to-do list app with user login and HTMx-powered frontend (Demo)

A full-stack application using Golang´s Fiber framework with session-based authentication. Once we are authenticated we can enter a view from which we can manage a list of tasks (list, update and delete). Requests to the backend are controlled by [</>htmx](https://htmx.org/) ([hypermedia](https://hypermedia.systems/) only).

### Explanation

The application allows us to perform a complete CRUD on the database, in this case SQLite3.

The DB stores both a table with the users and another table for each user´s to-do. Both tables are related using a foreign key.

>***In this application, instead of using the html/template package (gofiber/template/html, specifically), we use the [a-h/templ](https://github.com/a-h/templ) library. This amazing library implements a templating language (very similar to JSX) that compiles to Go code. Templ will allow us to write code almost identical to Go (with expressions, control flow, if/else, for loops, etc.) and have autocompletion thanks to strong typing. This means that errors appear at compile time and any calls to these templates (which are compiled as Go functions) from the handlers side will always require the correct data, minimizing errors and thus increasing the security and speed of our coding.***

The use of </>htmx allows behavior similar to that of a SPA, without page reloads when switching from one route to another or when making requests (via AJAX) to the backend.

On the other hand, the styling of the views is achieved through Tailwind CSS and DaisyUI that are obtained from their respective CDNs.

Finally, minimal use of [_hyperscript](https://hyperscript.org/) is made to achieve the action of closing the alerts when they are displayed.

>***This application is identical to that of a previous [repository](https://github.com/emarifer/gofiber-htmx-todolist) of mine, which is developed in GoFiber-template/html instead of the [Templ](https://templ.guide/) templating language, as in this case.***

---

## Screenshots:

###### Todo List Page with success alert:

<img src="/media/27e248cb-0785-40bb-80a6-5f65fafe9916.png" width="40%">

<br>

###### Another App:

<img src="/media/70d17e81-e825-44fd-89e5-b3a0697e706a.png" width="40%">

<br>

###### Task update page:

<img src="/media/fcca46f4-3645-4780-9303-ef117a572092.png" width="40%">

---

## Setup:

Besides the obvious prerequisite of having Go! on your machine, you must have Air installed for hot reloading when editing code.

Start the app in development mode:

```
$ air # Ctrl + C to stop the application
```

Build for production:

```
$ go build -ldflags="-s -w" -o ./bin/main . # ./bin/main to run the application
```

>***In order to have autocompletion and syntax highlighting in VS Code for the Teml templating language, you will have to install the [templ-vscode](https://marketplace.visualstudio.com/items?itemName=a-h.templ) extension (for vim/nvim install this [plugin](https://github.com/joerdav/templ.vim)). To generate the Go code corresponding to these templates you will have to download this [executable binary](https://github.com/a-h/templ/releases/tag/v0.2.476) from Github and place it in the PATH of your system. The command:***

```
$ templ generate --watch
```

>***Will allow us to watch the .templ files and compile them as we save. Review the documentation on Templ [installation](https://templ.guide/quick-start/installation) and [support](https://templ.guide/commands-and-tools/ide-support/) for your IDE.***

---

### Happy coding 😀!!',
    'A full-stack application using Golang´s Fiber framework.'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM posts ORDER BY id DESC LIMIT 1;
-- +goose StatementEnd
