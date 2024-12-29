package static

import (
	"embed"
	"net/http"
)

//go:embed  css images js
var staticfs embed.FS // for things exposed via /static

//go:embed partials templates data
var embedfs embed.FS // for things NOT exposed via /static

// HttpFS returns either an embedded file system or a reference to the local ./static folder.
// It is used to serve on /static.
// Sample: mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(static.HttpFS(useLocal))))
func HttpFS(useLocal bool) http.FileSystem {
	if useLocal {
		return http.Dir("./static")
	} else {
		return http.FS(staticfs)
	}
}

// StaticFS returns the embedded file system for /static
func StaticFS() embed.FS {
	return staticfs
}

// EmbeddedFS returns the file system for non-public embedded resources
func EmbeddedFS() embed.FS {
	return embedfs
}
