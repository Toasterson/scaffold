package template

import (
	"fmt"
	"html/template"
	"io"

	"github.com/labstack/echo"
)

var DefaultRenderer *SimpleRenderer

type SimpleRenderer struct {
	templates *template.Template
}

func (t *SimpleRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func InitTemplates(tplDir string) {
	DefaultRenderer = &SimpleRenderer{
		templates: template.Must(template.ParseGlob(fmt.Sprintf("%s/*.html", tplDir))),
	}
}
