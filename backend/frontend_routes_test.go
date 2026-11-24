package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestEmbeddedFrontendServesSPAAndAssets(t *testing.T) {
	t.Chdir(t.TempDir())
	s := &Server{staticFS: StaticFiles}
	mux := http.NewServeMux()
	s.registerWebRoutes(mux)

	var index string
	for _, route := range []string{"/", "/login", "/devices", "/messages", "/messages/", "/send", "/config"} {
		w := frontendRequest(mux, http.MethodGet, route)
		if w.Code != http.StatusOK || !strings.HasPrefix(w.Header().Get("Content-Type"), "text/html") ||
			!strings.Contains(w.Body.String(), `id="app"`) {
			t.Fatalf("GET %s = %d, %q; want Vue index", route, w.Code, w.Header().Get("Content-Type"))
		}
		if index == "" {
			index = w.Body.String()
		} else if w.Body.String() != index {
			t.Fatalf("route %s returned a different app shell", route)
		}
	}
	assets := regexp.MustCompile(`(?:src|href)="(/assets/[^"]+)"`).FindAllStringSubmatch(index, -1)
	if len(assets) < 2 {
		t.Fatalf("expected script and CSS assets in the Vue index, found %v", assets)
	}
	for _, pair := range assets {
		path := pair[1]
		w := frontendRequest(mux, http.MethodGet, path)
		want := "text/css"
		if strings.HasSuffix(path, ".js") {
			want = "javascript"
		}
		if w.Code != http.StatusOK || !strings.Contains(w.Header().Get("Content-Type"), want) || w.Body.Len() == 0 ||
			strings.HasPrefix(strings.ToLower(w.Body.String()), "<!doctype html") {
			t.Fatalf("GET %s = %d, %q; want nonempty %s, never SPA HTML", path, w.Code, w.Header().Get("Content-Type"), want)
		}
		head := frontendRequest(mux, http.MethodHead, path)
		if head.Code != http.StatusOK || head.Body.Len() != 0 || head.Header().Get("Content-Type") != w.Header().Get("Content-Type") {
			t.Fatalf("HEAD %s = %d, %q with %d body bytes", path, head.Code, head.Header().Get("Content-Type"), head.Body.Len())
		}
	}
	for _, path := range []string{"/app-logo.png", "/static/app-logo.png", "/art/van-gogh-sunflowers-bg.jpg"} {
		w := frontendRequest(mux, http.MethodGet, path)
		if w.Code != http.StatusOK || !strings.HasPrefix(w.Header().Get("Content-Type"), "image/") {
			t.Fatalf("GET %s = %d, %q; want embedded image", path, w.Code, w.Header().Get("Content-Type"))
		}
	}
	for _, path := range []string{"/assets/not-real.js", "/art/not-real.png", "/api/not-real", "/ws/not-real", "/not-a-route"} {
		if w := frontendRequest(mux, http.MethodGet, path); w.Code != http.StatusNotFound {
			t.Fatalf("GET %s = %d; want 404", path, w.Code)
		}
	}
	if w := frontendRequest(mux, http.MethodGet, "/app"); w.Code != http.StatusTemporaryRedirect || w.Header().Get("Location") != "/" {
		t.Fatalf("legacy /app = %d, %q; want redirect to /", w.Code, w.Header().Get("Location"))
	}
}

func TestFrontendDistIsOnlyUsedForDevelopment(t *testing.T) {
	root := t.TempDir()
	backend := filepath.Join(root, "backend")
	dist := filepath.Join(root, "frontend", "dist")
	for _, path := range []string{backend, filepath.Join(dist, "assets")} {
		if err := os.MkdirAll(path, 0755); err != nil {
			t.Fatal(err)
		}
	}
	for name, body := range map[string]string{"index.html": "<title>disk build</title>", "assets/dev.js": "const dev = 1;"} {
		if err := os.WriteFile(filepath.Join(dist, name), []byte(body), 0644); err != nil {
			t.Fatal(err)
		}
	}
	t.Chdir(backend)
	original := frontendEmbeddedOnly
	t.Cleanup(func() { frontendEmbeddedOnly = original })

	s := &Server{staticFS: StaticFiles}
	mux := http.NewServeMux()
	s.registerWebRoutes(mux)
	frontendEmbeddedOnly = "false"
	if w := frontendRequest(mux, http.MethodGet, "/login"); w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "disk build") {
		t.Fatalf("developer mode failed to serve disk index: %d, %q", w.Code, w.Body.String())
	}
	if w := frontendRequest(mux, http.MethodGet, "/assets/dev.js"); w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "const dev") {
		t.Fatalf("developer mode failed to serve disk asset: %d, %q", w.Code, w.Body.String())
	}

	frontendEmbeddedOnly = "true"
	if w := frontendRequest(mux, http.MethodGet, "/login"); w.Code != http.StatusOK || !strings.Contains(w.Body.String(), `id="app"`) {
		t.Fatalf("standalone mode did not serve its own index: %d", w.Code)
	}
	if w := frontendRequest(mux, http.MethodGet, "/assets/dev.js"); w.Code != http.StatusNotFound {
		t.Fatalf("standalone mode unexpectedly served external asset: %d", w.Code)
	}
}

func frontendRequest(mux *http.ServeMux, method, path string) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	mux.ServeHTTP(w, httptest.NewRequest(method, path, nil))
	return w
}
