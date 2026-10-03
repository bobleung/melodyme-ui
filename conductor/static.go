package conductor

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"io/fs"
	"net/http"
	"strings"
)

//go:generate go tool templ generate

// Prefix is where an app mounts Static. The page asks for its files there, from the app's own
// address, so a Content-Security-Policy of 'self' covers them.
const Prefix = "/static/melodyme-ui/"

//go:embed static
var static embed.FS

var versions = map[string]string{}

func init() {
	fs.WalkDir(static, "static", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		b, _ := static.ReadFile(path)
		sum := sha256.Sum256(b)
		versions[strings.TrimPrefix(path, "static/")] = hex.EncodeToString(sum[:])[:10]
		return nil
	})
}

// asset is a file's address, versioned by its content so it can be cached for ever.
func asset(name string) string {
	v, ok := versions[name]
	if !ok {
		panic("conductor: no static file " + name)
	}
	return Prefix + name + "?v=" + v
}

// Static serves the page's styles, script and fonts. Mount it at Prefix:
//
//	mux.Handle("GET "+conductor.Prefix+"{path...}", conductor.Static())
func Static() http.Handler {
	sub, _ := fs.Sub(static, "static")
	files := http.StripPrefix(Prefix, http.FileServerFS(sub))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v, ok := versions[strings.TrimPrefix(r.URL.Path, Prefix)]
		if !ok {
			http.NotFound(w, r)
			return
		}
		if r.URL.Query().Get("v") == v {
			w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
		} else {
			w.Header().Set("Cache-Control", "no-cache")
		}
		files.ServeHTTP(w, r)
	})
}
