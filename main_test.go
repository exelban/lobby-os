package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/pkgz/rest"
)

func TestDashboardConfig(t *testing.T) {
	a := &app{linksPath: filepath.Join(t.TempDir(), "links.json")}
	name := `Quotes " and </script><script>alert(1)</script>`
	if err := a.writeLinks([]Link{{ID: "existing", Name: name, URL: "https://example.com"}}); err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	a.router().ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("status: %d", w.Code)
	}
	_, config, ok := strings.Cut(w.Body.String(), "const data = ")
	if !ok {
		t.Fatal("missing dashboard config")
	}
	var got struct {
		Version string `json:"version"`
		Links   []Link `json:"links"`
	}
	if err := json.NewDecoder(strings.NewReader(config)).Decode(&got); err != nil {
		t.Fatal(err)
	}
	if got.Version != version || len(got.Links) != 1 || got.Links[0].Name != name {
		t.Fatalf("unexpected config: %+v", got)
	}
	if strings.Contains(w.Body.String(), name) {
		t.Fatal("link name was not escaped for the script context")
	}
}

func TestLegacyLinks(t *testing.T) {
	a := &app{linksPath: filepath.Join(t.TempDir(), "links.json")}
	legacy := `[
		{"id":"old","name":"Missing fields","url":"https://example.com"},
		{"name":"Null fields","url":"https://example.com","color":null,"icon":null,"preset":null,"group":null,"groupName":null},
		{"name":"Empty fields","url":"https://example.com","color":"","icon":"","preset":"","group":"","groupName":""},
		{"name":"Named group","url":"https://example.com","group":"group-1","groupName":"Services","color":"#ffffff","icon":"/icon.png","preset":"grafana"}
	]`
	if err := os.WriteFile(a.linksPath, []byte(legacy), 0600); err != nil {
		t.Fatal(err)
	}
	links := a.readLinks()
	if len(links) != 4 {
		t.Fatalf("loaded %d links", len(links))
	}
	if err := a.writeLinks(links); err != nil {
		t.Fatal(err)
	}
	got := a.readLinks()
	if len(got) != 4 || got[0].ID != "old" {
		t.Fatalf("existing links changed: %+v", got)
	}
	ids := map[string]bool{}
	for i, link := range got {
		if link.ID == "" || ids[link.ID] {
			t.Fatalf("missing or duplicate ID: %q", link.ID)
		}
		ids[link.ID] = true
		if link != links[i] {
			t.Fatalf("link changed after saving: %+v", link)
		}
	}
	if got[1].Group != "" || got[2].Group != "" || got[3].GroupName != "Services" {
		t.Fatalf("legacy group information changed: %+v", got)
	}
}

func TestDashboardCaching(t *testing.T) {
	assets, err := fs.Glob(htmlFS, "web/dist/assets/*.js")
	if err != nil || len(assets) == 0 {
		t.Fatalf("missing built assets: %v", err)
	}
	a := &app{linksPath: filepath.Join(t.TempDir(), "links.json")}
	router := a.router()
	for _, tc := range []struct {
		path   string
		status int
		cache  string
	}{
		{"/", http.StatusOK, "no-cache"},
		{"/manifest.json", http.StatusOK, "no-cache"},
		{strings.TrimPrefix(assets[0], "web/dist"), http.StatusOK, "public, max-age=31536000, immutable"},
		{"/assets/missing.js", http.StatusNotFound, ""},
	} {
		t.Run(tc.path, func(t *testing.T) {
			w := httptest.NewRecorder()
			router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, tc.path, nil))
			if w.Code != tc.status || w.Header().Get("Cache-Control") != tc.cache {
				t.Fatalf("status %d, cache %q", w.Code, w.Header().Get("Cache-Control"))
			}
		})
	}
}

func TestRunReturnsListenError(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	a := &app{srv: rest.NewServer(listener.Addr().(*net.TCPAddr).Port)}
	a.srv.Address = "127.0.0.1"
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := a.run(ctx); err == nil {
		t.Fatal("expected a listen error for the occupied port")
	}
}

func TestExtractFaviconURL(t *testing.T) {
	for _, tc := range []struct {
		name      string
		href      string
		available string
		want      string
	}{
		{"relative after redirect", "icons/site.png", "/app/icons/site.png", "/app/icons/site.png"},
		{"root relative", "/site.png", "/site.png", "/site.png"},
		{"absolute", "{base}/site.png", "/site.png", "/site.png"},
		{"protocol relative", "//{host}/site.png", "/site.png", "/site.png"},
		{"invalid candidate", "%zz", "/apple-touch-icon.png", "/apple-touch-icon.png"},
		{"missing candidate", "/missing.png", "/favicon.png", "/favicon.png"},
		{"default fallback", "", "", "/favicon.ico"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodHead {
					if r.URL.Path != tc.available {
						http.NotFound(w, r)
					}
					return
				}
				if r.URL.Path == "/" {
					http.Redirect(w, r, "/app/", http.StatusFound)
					return
				}
				href := strings.NewReplacer("{base}", "http://"+r.Host, "{host}", r.Host).Replace(tc.href)
				fmt.Fprintf(w, `<html><head><link rel="icon" href="%s"></head></html>`, href)
			}))
			defer server.Close()

			got, err := extractFaviconURL(context.Background(), server.URL)
			if err != nil {
				t.Fatal(err)
			}
			if want := server.URL + tc.want; got != want {
				t.Fatalf("got %q, want %q", got, want)
			}
		})
	}
}

func TestFaviconHead(t *testing.T) {
	for _, tc := range []struct {
		name string
		page string
		want string
	}{
		{"rel tokens", `<head><link rel="alternate ICON" href=" /site.png "></head>`, "/site.png"},
		{"head only", `<head></head><body><link rel="icon" href="/site.png"></body>`, "/favicon.ico"},
		{"implicit head end", `<body><link rel="icon" href="/site.png"></body>`, "/favicon.ico"},
		{"read limit", "<head>" + strings.Repeat(" ", 1<<20) + `<link rel="icon" href="/site.png"></head>`, "/favicon.ico"},
		{"unsupported scheme", `<head><link rel="icon" href="data:image/png;base64,AAAA"></head>`, "/favicon.ico"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodGet {
					fmt.Fprint(w, tc.page)
				}
			}))
			defer server.Close()

			got, err := extractFaviconURL(context.Background(), server.URL)
			if err != nil || got != server.URL+tc.want {
				t.Fatalf("got %q, error %v", got, err)
			}
		})
	}
}

func TestFaviconErrorPageFallback(t *testing.T) {
	for _, status := range []int{http.StatusUnauthorized, http.StatusForbidden, http.StatusNotFound} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == http.MethodHead && r.URL.Path == "/favicon.ico" {
					return
				}
				w.WriteHeader(status)
				w.(http.Flusher).Flush()
				<-r.Context().Done()
			}))
			defer server.Close()

			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			got, err := extractFaviconURL(ctx, server.URL+"/protected")
			if err != nil || got != server.URL+"/favicon.ico" {
				t.Fatalf("got %q, error %v", got, err)
			}
		})
	}
}

func TestFaviconStreaming(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodHead {
			return
		}
		fmt.Fprint(w, `<head><link rel="icon" href="/site.png">`)
		w.(http.Flusher).Flush()
		<-r.Context().Done()
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	got, err := extractFaviconURL(ctx, server.URL)
	if err != nil || got != server.URL+"/site.png" {
		t.Fatalf("got %q, error %v", got, err)
	}
}

func TestFaviconCancellation(t *testing.T) {
	for _, method := range []string{http.MethodGet, http.MethodHead} {
		t.Run(method, func(t *testing.T) {
			started := make(chan struct{})
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == method {
					close(started)
					<-r.Context().Done()
					return
				}
				fmt.Fprint(w, `<head><link rel="icon" href="/site.png"></head>`)
			}))
			defer server.Close()

			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			done := make(chan error, 1)
			go func() {
				_, err := extractFaviconURL(ctx, server.URL)
				done <- err
			}()
			select {
			case <-started:
				cancel()
			case <-ctx.Done():
				t.Fatal("favicon request did not start")
			}
			if err := <-done; !errors.Is(err, context.Canceled) {
				t.Fatalf("expected cancellation, got %v", err)
			}
		})
	}
}
