<div align="center">
  
<h1 align="center">goCMS</h1>

<img src="doc/cms-logo.png" width="55%">

<hr />

<p style="margin-bottom: 8px;">

goCMS is a headless CMS (Content Management System) written in Golang using Gin framework + </>Htmx & A-H Templ, designed to be fast, efficient, and easily extensible. It allows you to create a website or blog, with any template you like, in only a few commands.

</p>
  
![GitHub License](https://img.shields.io/github/license/emarifer/url-shortener-echo-templ-htmx) ![Static Badge](https://img.shields.io/badge/Go-%3E=1.18-blue)

</div>

<hr />

## Features 🚀

- [x] **Headless Architecture:** Adding pages, posts, or forms should all
  be done with easy requests to the API.
- [x] **Golang-Powered:** Leverage the performance and safety of one of the
  best languages in the market for backend development.
- [x] **SQL Database Integration:** Store your posts and pages in SQL databases for reliable and scalable data storage.
- [x] **Centralized HTTP error handling:** The user receives feedback about the Http errors that their actions may cause, through the use of middleware that centralizes the Http errors that occur.
- [ ] **Post**: We can add, update, and delete posts. Posts can be served
  through a unique URL.
- [ ] **Pages**: TODO.
- [ ] **Menus**: TODO
- [ ] **Live Reload** through the use of `air`.

## Installation

Ensure you have Golang installed on your system before proceeding with the installation.

```bash
go get -u github.com/emarifer/gocms
```

## Example - Running the App

After you've replaced the default template files with your prefered
template, simply build and start the app with the following commands.

```bash
go build
./gocms
```

This will start `goCMS` on `http://localhost:8080`. You can customize
the configuration by providing the necessary environment variables (settings/settings.yaml).

## Architecture

Currently, the architecture of `goCMS` is still in its early days.
The plan is to have two main applications: the public facing application
to serve the content through a website, and the admin application that
can be hidden, where users can modify the settings, add posts, pages, etc.

## License

`goCMS` is released under the MIT License. See LICENSE for
details. Feel free to fork, modify, and use it in your projects!

---

## Happy coding 😀!!