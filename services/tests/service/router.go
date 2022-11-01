package service

import (
	"net/http"
	"os"
	"regexp"

	"github.com/costal/go-misc-tools/httpapp"
)

var CorsManager = struct {
	allowedOrigins map[string]struct{}
}{
	allowedOrigins: make(map[string]struct{}),
}

func Router(app *HTTPApplication) http.Handler {
	CorsManager.allowedOrigins[os.Getenv("MAINAPP_URL")] = struct{}{}
	m := httpapp.HTTPMethods()
	app.HTTP.CorsOptions.AllowCredentials = true
	app.HTTP.CorsOptions.AllowOriginRequestFunc = app.allowedOriginValidator
	app.HTTP.AddStandardMiddleware(app.HTTP.Cors([]string{
		"authorization",
		"content-type",
	}))
	app.HTTP.
		URL("info", m.Get(info)).
		URL("/api/").Middleware(app.JWTGuard).
		URL("/api/tests",
			m.Get(app.showTests).
				Middleware(
					app.withPermissionScope("view"),
					app.withViewConstraints,
				),
			m.Post(app.createTest),
		)
	return app.HTTP.Router()
}

type route struct {
	pattern *regexp.Regexp
	handler http.Handler
}

type RegexpHandler struct {
	routes []*route
}

func (h *RegexpHandler) Handler(pattern *regexp.Regexp, handler http.Handler) {
	h.routes = append(h.routes, &route{pattern, handler})
}

func (h *RegexpHandler) HandleFunc(pattern *regexp.Regexp, handler func(http.ResponseWriter, *http.Request)) {
	h.routes = append(h.routes, &route{pattern, http.HandlerFunc(handler)})
}

func (h *RegexpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range h.routes {
		if route.pattern.MatchString(r.URL.Path) {
			route.handler.ServeHTTP(w, r)
			return
		}
	}
	// no pattern matched; send 404 response
	http.NotFound(w, r)
}
