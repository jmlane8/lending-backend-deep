// Package openapi provides HTTP handlers for serving the OpenAPI specification
// and an interactive Redoc UI — with no external framework dependency.
//
// The spec is read from the embedded FS (resources/gen) at startup time so
// any missing file causes a panic-on-start rather than a runtime error.
package openapi

import (
	"embed"
	"fmt"
	"io/fs"
	"lending-backend-deep/internal/rest/constant"
	"net/http"
)

// redocPage is a minimal HTML page that loads Redoc from the unpkg CDN
// and points it at the /openapi.json endpoint served by SpecHandler.
const redocPage = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8"/>
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>Loan Submission API — API Docs</title>
  <style>body{margin:0;padding:0}</style>
</head>
<body>
  <redoc spec-url="/openapi.json"></redoc>
  <script src="https://unpkg.com/redoc@2/bundles/redoc.standalone.js"></script>
</body>
</html>`

// SpecHandler returns an http.HandlerFunc that serves the pre-bundled OpenAPI
// JSON specification embedded at fsPath inside specFS.
//
// The file is read once at call time (startup); a missing file panics early.
//
// Usage:
//
//	router.Get("/openapi.json", openapi.SpecHandler(gen.OpenAPISpec, "openapi/openapi.json"))
func SpecHandler(specFS embed.FS, fsPath string) http.HandlerFunc {
	data, err := fs.ReadFile(specFS, fsPath)
	if err != nil {
		panic(fmt.Sprintf("openapi.SpecHandler: cannot read embedded spec %q: %v", fsPath, err))
	}
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", constant.MimeApplicationJSON)
		w.Header().Set("Cache-Control", "no-cache")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(data)
	}
}

// UIHandler returns an http.HandlerFunc that serves a Redoc HTML page.
// The page fetches the spec from /openapi.json via the browser, so no server-side
// URL configuration is needed.
//
// Usage:
//
//	router.Get("/redoc", openapi.UIHandler())
func UIHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, redocPage)
	}
}
