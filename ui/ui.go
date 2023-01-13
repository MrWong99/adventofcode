package ui

import (
	"embed"
	"html/template"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/MrWong99/adventofcode/log"
)

//go:embed *
var htmlDir embed.FS

var files map[string][]byte

func init() {
	files = make(map[string][]byte)
	readDir(".")
	log.Log.Debug("Htmls", "htmlDir", htmlDir)
	log.Log.Debug("Initialized embedded fs", "filecount", len(files))
}

func readDir(dir string) {
	entries, err := htmlDir.ReadDir(dir)
	if err != nil {
		return
	}
	if !strings.HasSuffix(dir, "/") {
		dir = dir + "/"
	}
	if dir == "./" {
		dir = ""
	}
	for _, entry := range entries {
		if entry.IsDir() {
			readDir(dir + entry.Name())
		} else {
			path := dir + entry.Name()
			content, err := htmlDir.ReadFile(path)
			if err != nil {
				log.Log.Warn("Could not read fs file", "path", path, "error", err)
			} else {
				files[path] = content
			}
		}
	}
}

var ui = &Ui{}

type Ui struct{}

func Instance() *Ui {
	return ui
}

func (*Ui) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var input url.Values
	if r.Body != nil && r.Method == "POST" {
		r.ParseForm()
		input = r.Form
		log.Log.Info("User input form data on site", "url", r.URL, "input", input)
	} else {
		input = make(url.Values)
	}
	toServe := r.URL.EscapedPath()[1:]
	if toServe == "" {
		toServe = "index.html"
	}
	file, ok := files[toServe]
	if !ok {
		log.Log.Debug("Unknown file requested", "url", toServe)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	switch filepath.Ext(toServe) {
	case ".html":
		w.Header().Set("Content-Type", "text/html")
		templateParsed := serveAsTemplate(toServe, w, input, r)
		if templateParsed {
			return
		}
	case ".css":
		w.Header().Set("Content-Type", "text/css")
	case ".ico":
		w.Header().Set("Content-Type", "image/x-icon")
	default:
		w.WriteHeader(http.StatusForbidden)
		return
	}
	w.Write(file)
}

func serveAsTemplate(toServe string, w http.ResponseWriter, input url.Values, r *http.Request) bool {
	tpl, err := template.New("sitelayout.html").ParseFS(htmlDir, "sitelayout.html", toServe)
	if err == nil {
		err = tpl.Execute(w, input)
		if err == nil {
			return true
		}
		log.Log.Warn("Could not parse template file. Serving raw instead", "url", r.URL, "toServe", toServe, "error", err)
	} else {
		log.Log.Warn("Could not create template file. Serving raw instead", "url", r.URL, "toServe", toServe, "error", err)
	}
	return false
}
