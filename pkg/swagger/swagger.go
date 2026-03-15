package swagger

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

const uiTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>%s</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    SwaggerUIBundle({ url: "/swagger/openapi.yaml", dom_id: "#swagger-ui" });
  </script>
</body>
</html>`

// RegisterRoutes mounts Swagger UI on /swagger/ serving the provided OpenAPI spec.
// title is displayed in the browser tab. spec is the raw OpenAPI YAML bytes.
func RegisterRoutes(r chi.Router, title string, spec []byte) {
	html := []byte(fmt.Sprintf(uiTemplate, title))

	r.Get("/swagger/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yaml")
		w.Write(spec)
	})
	r.Get("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
	})
	r.Get("/swagger/*", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write(html)
	})
}
