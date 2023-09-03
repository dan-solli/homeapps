package views

import (
	"embed"
	"html/template"
)

var (
	//go:embed "templates/*"
	trackerTemplates embed.FS
)

func NewTemplates() (*template.Template, error) {
	return template.ParseFS(trackerTemplates, "templates/*/*.gohtml", "templates/*.gohtml")
}
