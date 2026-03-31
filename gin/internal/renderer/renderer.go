package renderer

import (
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin/render"
)

type HTMLRenderer struct {
	templates map[string]*template.Template
}

func NewHTMLRenderer(templateDir string) *HTMLRenderer {
	r := &HTMLRenderer{templates: make(map[string]*template.Template)}
	pages := []string{"login", "signup", "home", "dashboard", "error"}

	for _, page := range pages {
		r.templates[page] = template.Must(template.ParseFiles(
			filepath.Join(templateDir, "base.html"),
			filepath.Join(templateDir, page+".html"),
		))
	}

	return r
}

func (r *HTMLRenderer) Instance(name string, data interface{}) render.Render {
	return &HTMLInstance{
		tmpl: r.templates[name],
		data: data,
	}
}

type HTMLInstance struct {
	tmpl *template.Template
	data interface{}
}

func (h *HTMLInstance) Render(w http.ResponseWriter) error {
	return h.tmpl.ExecuteTemplate(w, "base", h.data)
}

func (h *HTMLInstance) WriteContentType(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
}
