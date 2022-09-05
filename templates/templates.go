package templates

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func getTemplateGlob() (*template.Template, error) {
	var tg *template.Template
	var err error

	tg, err = template.ParseGlob("templates/*.html")

	if err != nil {
		tg, err = template.ParseGlob("../templates/*.html")
	}

	return tg, err
}

var Renderer *Template = &Template{templates: template.Must(getTemplateGlob())}
