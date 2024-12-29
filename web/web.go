package web

import (
	"html/template"
	"io"
	"path/filepath"
	"pochtmx/static"
)

// Local indicates whether to use the local file system
// or the embedded file system. The default is false
// which will use the embedded templates.
var Local bool

// Execute parses and executes a template.  It takes a variable number
// of templates which probably includes a base, content and any other
// shared templates.
func Execute(w io.Writer, data any, names ...string) error {
	tmpl, err := parseTemplates(names...)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, data)
}

func parseTemplates(names ...string) (*template.Template, error) {
	// local needs "static" prepended to the path of each template
	if Local {
		newNames := make([]string, 0, len(names))
		for _, str := range names {

			newNames = append(newNames, filepath.ToSlash(filepath.Join("static", str)))
		}
		return template.ParseFiles(newNames...)
	}
	return template.ParseFS(static.EmbeddedFS(), names...)
}
