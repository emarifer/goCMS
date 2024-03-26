<div align="center">
  
<h1 align="center">goCMS</h1>

<img src="doc/cms-logo.png" width="55%">

<hr />

<p style="margin-bottom: 8px;">

goCMS is a headless CMS (Content Management System) written in Golang using Gin framework + </>Htmx & A-H Templ, designed to be fast, efficient, and easily extensible. It allows you to create a website or blog, with any template you like, in only a few commands.

</p>
  
![GitHub License](https://img.shields.io/github/license/emarifer/url-shortener-echo-templ-htmx) ![Static Badge](https://img.shields.io/badge/Go-%3E=1.22-blue)

</div>

<hr />

## Features 🚀

- [x] **Headless Architecture:** Adding pages, posts, or forms should all
  be done with easy requests to the API.
- [x] **Golang-Powered:** Leverage the performance and safety of one of the
  best languages in the market for backend development.
- [x] **SQL Database Integration:** Store your posts and pages in SQL databases for reliable and scalable data storage.
- [x] **Centralized HTTP error handling:** The user receives feedback about the Http errors that their actions may cause, through the use of middleware that centralizes the Http errors that occur.
- [x] **Caching HTML responses from endpoints:** Own implementation of an in-memory cache that stores HTML responses for 10 minutes in a map with mutex lock R/W access.
- [x] **Live Reload:** through the use of `air`.
- [x] **Possibility for the user to add their own plugins written in `Lua`:** this feature allows you to customize the admin application at runtime.
- [ ] **Post**: We can add, update, and delete posts. Posts can be served
  through a unique URL.
- [ ] **Pages**: TODO.
- [ ] **Menus**: TODO
  
<br />

>[!IMPORTANT]
>***The Go language uses [html/template](https://pkg.go.dev/html/template) package to render HTML. In this application we have used the [a-h/templ](https://github.com/a-h/templ) library instead. The main difference is that templ uses a generation step to compile the files .templ into Go code (as functions). This means that the templates are type-safe and can be checked at compile time. This amazing library implements a templating language (very similar to JSX) which allows you to write code almost identical to Go (with expressions, control flow, if/else, for loops, etc.) and have autocompletion. For all these reasons, calling these templates from the controllers side will always require the correct data, minimizing errors and thus increasing the security and speed of our coding.***

## Installation

Ensure you have Golang installed on your system before proceeding with the installation.

```bash
go get -u github.com/emarifer/gocms
```

### Example - Running the App (user application) manually

First, make sure you have the necessary executable binaries to run and work with the application.

```bash
make install-tools
```

After that, with the `MariaDB` database engine running, with the idea of populating the database with some sample data, make sure to run the migrations with the previously installed `Goose` tool. We recommend creating a database called `cms_db` and running the following command:

```bash
GOOSE_DRIVER="mysql" GOOSE_DBSTRING="root:root@/cms_db" goose up
```
Replace the database connection string with the appropriate string
dependending on where your database is.

After you've replaced the default template files with your prefered
template, simply build and start the app with the following commands.

```bash
go mod tidy && go build -ldflags="-s -w" -v -o ./tmp/gocms ./cmd/gocms && ./tmp/gocms
```

Alternatively, the `air` command will allow us to start the user application (also creating the admin application executable), having, however, with said command the possibility of hot reloading after any change in the user/admin applications code.

This will start `goCMS` on `http://localhost:8080`. If we have used the `air` command we can start the admin application with the `make run` (on `http://localhost:8081`) command. You can customize the configuration by providing the necessary environment variables.

```bash
# e.g.

DATABASE_PORT=3306 ./tmp/gocms-admin
```

For more information, see the [configuration settings](#configuration).

### Example - Running with Docker Compose (user & admin applications)

In this case the only requirement is to have `Docker` installed and running.

To create the image and the `Docker` containers and start the application you only need to run the following command in the project folder:

```bash
make run-containers
```
The above will create an Ubuntu:jammy image and, within that OS, will install Golang, Goose, A-H.Templ and Air. Next, from said image and the mariadb:jammy image, you will create and start two containers: one containing the `goCMS` app, serving on port `8080`, and another one serving the `mariadb` database internally. This will also run the migrations automatically to setup the database!

To stop and eliminate both containers we will execute the following in another terminal:

```bash
docker compose down # to stop and remove containers (run in another terminal)
```

If we do not plan to delete the containers with the idea of continuing to reuse them, we will simply press `Ctrl+C` in the same terminal. This will stop the containers without deleting them. The next time we want to start the application we will run `make run-containers` again.

As long as we have created/run the aforementioned containers, the management application executable file will have been created. To start it (within `Docker`), simply run the following commands:

```bash
docker exec -it docker-gocms-1 sh # to enter the `docker-gocms-1` container

cd gocms && make start-admin-container # to enter the project folder (inside the container) and start the admin application
```

>[!NOTE]
>***the above serves the application in `http://localhost:8081`.***

## Architecture

Currently, the architecture of `goCMS` is still in its early days.
The plan is to have two main applications: the public facing application
to serve the content through a website, and the admin application that
can be hidden, where users can modify the settings, add posts, pages, etc.

## Configuration

The runtime configuration can be done through a [toml](https://toml.io/en/) configuration file or by setting the mandatory environment variables (*fallback*). This approach was chosen because configuration via toml supports advanced features (i.e. *relationships*, *arrays*, etc.). The `.dev.env`-file used only for the `goose up` command, they are not needed for `Docker` files.

### `.toml` configuration

The application can be started by providing the `config` flag which has to be set to a toml configuration file. The file has to contain the following mandatory values:

```toml
webserver_port = 8080 # port to run the webserver on
admin_webserver_port = 8081 # port to run the webserver (admin) on
database_host = "localhost" # database host (address to the MariaDB database)
database_port = 3306 # database port
database_user = "root" # database user
database_password = "my-secret-pw" # database password
database_name = "cms_db" # name of the database that is created through `Docker`
image_dir = "./media" # directory to use for storing uploaded images

# optional: directives containing the name and path of user-supplied `Lua` plugins
# e.g.

[[shortcodes]]
name = "img"
# must have function "HandleShortcode(arguments []string) -> string"
plugin = "plugins/image_shortcode.lua"
```

>[!NOTE]
>***The above configuration values are used to start the local development database, in addition to the user/admin application ports, media storage folder, or optionally, admin plugins directives.***

### Environment variables configuration (fallback)

If chosen, by setting the following environment variables the application can be started without providing a toml configuration file (although a file of this type is necessary to establish the directives of user plugins written in `Lua`). 

- `WEBSERVER_PORT` port the application should run on
- `ADMIN_WEBSERVER_PORT` the same as the previous one but for the admin app
- `DATABASE_HOST` should contain the database addres, e.g. `localhost`
- `DATABASE_PORT` should be the connection port to the db, e.g. `3306`
- `DATABASE_USER` is the database username.
- `DATABASE_PASSWORD` needs to contain the database password for the given user.
- `DATABASE_NAME` sets the name of the database `goCMS` will use.
- `IMAGE_DIRECTORY` directory images should be stored to if uploaded to `goCMS`

To the above (as we have already mentioned), we would have to add an environment variable (`CONFIG_FILE_PATH`) that contains the path to a `.toml` file that contains the directives for the plugins (`Lua` scripts) that the user wants add. This file would have the form:

```toml
# e.g.

[[shortcodes]]
name = "img"
# must have function "HandleShortcode(arguments []string) -> string"
plugin = "plugins/image_shortcode.lua"
```

## Development

To facilitate the development process, `Docker` is highly recommended. This way you can use `docker/mariadb.yml` to configure a predefined MariaDB database server. The file `mariadb.yml` creates the database `cms_db`.

```bash
$ make start-devdb
```

To populate the aforementioned db with some sample data you can use this command:

```bash
$ make run-migrations
```

## License

`goCMS` is released under the MIT License. See LICENSE for
details. Feel free to fork, modify, and use it in your projects!

---

## Happy coding 😀!!