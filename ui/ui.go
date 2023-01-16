package ui

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/MrWong99/adventofcode/db"
	"github.com/MrWong99/adventofcode/log"
	"github.com/MrWong99/adventofcode/templatefunc"
)

//go:embed *
var htmlDir embed.FS

var files map[string][]byte

var templates map[string]*template.Template

func init() {
	files = make(map[string][]byte)
	templates = make(map[string]*template.Template)
	readDir(".")
	log.Log.Debug("Initialized embedded fs", "filecount", len(files), "templatecount", len(templates))
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
			if path == "sitelayout.html" {
				continue
			}
			content, err := htmlDir.ReadFile(path)
			if err != nil {
				log.Log.Warn("Could not read fs file", "path", path, "error", err)
			} else {
				if filepath.Ext(path) == ".gohtml" {
					tpl, err := siteTpl().ParseFS(htmlDir, "sitelayout.gohtml", path)
					if err == nil {
						templates[strings.ReplaceAll(path, ".gohtml", ".html")] = tpl
						continue
					} else {
						log.Log.Debug("Error while loading template", "path", path, "error", err)
					}
				}
				files[path] = content
			}
		}
	}
}

func siteTpl() (tpl *template.Template) {
	tpl = template.New("sitelayout.gohtml").Funcs(template.FuncMap{
		"AllPlugins":   db.AllPlugins,
		"AddPlugin":    templatefunc.AddPlugin,
		"DeletePlugin": templatefunc.DeletePlugin,
		"Calculate":    templatefunc.Calculate,
	})
	return
}

var ui = &Ui{}

type Ui struct{}

func Instance() *Ui {
	return ui
}

func (*Ui) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	input := make(url.Values)
	if r.Body != nil {
		if r.Method == "POST" {
			r.ParseForm()
			input = r.Form
			log.Log.Debug("User input form data on site", "url", r.URL, "input", input)
		}
		r.Body.Close()
	}

	toServe := r.URL.EscapedPath()[1:]
	if toServe == "" {
		toServe = "index.html"
	}
	switch filepath.Ext(toServe) {
	case ".html":
		w.Header().Set("Content-Type", "text/html")
	case ".css":
		w.Header().Set("Content-Type", "text/css")
	case ".ico":
		w.Header().Set("Content-Type", "image/x-icon")
	default:
		serveErrorPage(&ErrorDisplay{Error: fmt.Errorf("file not found %s", toServe), StatusCode: http.StatusNotFound}, w)
		return
	}
	tpl, ok := templates[toServe]
	if ok {
		w.Header().Add("Cache-Control", "no-store")
		err := serveAsTemplate(tpl, w, input)
		if err == nil {
			return
		}
		serveErrorPage(&ErrorDisplay{Error: err, StatusCode: http.StatusInternalServerError}, w)
		return
	}
	file, ok := files[toServe]
	if !ok {
		log.Log.Debug("Unknown file requested", "url", toServe)
		serveErrorPage(&ErrorDisplay{Error: fmt.Errorf("file not found %s", toServe), StatusCode: http.StatusNotFound}, w)
		return
	}
	w.Write(file)
}

func serveAsTemplate(tpl *template.Template, w http.ResponseWriter, input url.Values) error {
	err := tpl.Execute(w, input)
	if err == nil {
		return nil
	}
	log.Log.Warn("Could not parse template file", "error", err)
	return err
}

type ErrorDisplay struct {
	Error      error
	StatusCode int
}

func serveErrorPage(input *ErrorDisplay, w http.ResponseWriter) {
	w.Header().Add("Cache-Control", "no-store")
	w.WriteHeader(input.StatusCode)
	tpl := templates["error.html"]
	err := tpl.Execute(w, input)
	if err == nil {
		return
	}
	log.Log.Error("Could not serve error page", "input", input, "error", err)
}
