package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pkgz/rest"
	"github.com/pkgz/service"
	"golang.org/x/net/html"
)

var (
	//go:embed web/dist
	htmlFS  embed.FS
	version = "debug"
)

type args struct {
	DataPath string `long:"data-path" env:"DATA_PATH" default:"/srv/data" description:"path to data files"`
	service.ARGS
}

type app struct {
	srv       *rest.Server
	linksPath string

	storageMu sync.Mutex
}

type Link struct {
	ID string `json:"id"`

	Name      string `json:"name"`
	URL       string `json:"url"`
	Color     string `json:"color,omitempty"`
	Icon      string `json:"icon,omitempty"`
	Preset    string `json:"preset,omitempty"`
	Group     string `json:"group,omitempty"`
	GroupName string `json:"groupName,omitempty"`
}

var faviconClient = &http.Client{Timeout: 10 * time.Second}

func main() {
	fmt.Println(version)

	var args args
	ctx, cancel, err := service.Init(&args)
	if err != nil {
		log.Printf("[ERROR] %v", err)
		os.Exit(1)
	}
	defer cancel()

	srv := &app{
		srv:       rest.NewServer(args.Port),
		linksPath: filepath.Join(args.DataPath, "links.json"),
	}

	if err = srv.run(ctx); err != nil {
		log.Printf("[ERROR] application failed: %v", err)
		os.Exit(4)
	}

	log.Print("[INFO] application terminated")
}

func (a *app) run(ctx context.Context) error {
	log.Printf("[INFO] application started")

	router := a.router()
	done := make(chan error, 1)
	go func() { done <- a.srv.Run(router) }()

	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return a.srv.Shutdown()
	}
}

func (a *app) router() chi.Router {
	router := chi.NewRouter()

	router.Use(middleware.Heartbeat("/ping"))
	router.Use(middleware.Recoverer)
	router.Use(rest.Logger)
	router.NotFound(rest.NotFound)

	type data struct {
		Version string `json:"version"`
		Links   []Link `json:"links"`
	}

	router.Route("/api", func(r chi.Router) {
		r.Use(middleware.AllowContentType("application/x-www-form-urlencoded", "application/json"))

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			rest.JsonResponse(w, data{
				Version: version,
				Links:   a.readLinks(),
			})
		})
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			var req []Link
			if err := rest.ReadBody(r, &req); err != nil {
				rest.ErrorResponse(w, r, http.StatusBadRequest, nil, err.Error())
				return
			}
			if err := a.writeLinks(req); err != nil {
				log.Printf("[ERROR] failed to write links: %v", err)
				rest.ErrorResponse(w, r, http.StatusInternalServerError, nil, "failed to write links")
				return
			}
			rest.OkResponse(w)
		})

		r.Get("/favicon", func(w http.ResponseWriter, r *http.Request) {
			uri := r.URL.Query().Get("url")
			if uri == "" {
				rest.ErrorResponse(w, r, http.StatusBadRequest, nil, "missing url parameter")
				return
			}

			faviconURL, err := extractFaviconURL(r.Context(), uri)
			if err != nil {
				log.Printf("[ERROR] failed to extract favicon: %v", err)
				rest.ErrorResponse(w, r, http.StatusInternalServerError, nil, "failed to extract favicon")
				return
			}

			rest.TextResponse(w, faviconURL)
		})
	})

	staticFs, err := fs.Sub(htmlFS, "web/dist")
	if err != nil {
		log.Fatalf("[FATAL] failed to create sub filesystem: %v", err)
	}

	// Include the quotes so the placeholder is valid JavaScript in the standalone demo.
	tpl, err := template.New("index.html").Delims(`"{{`, `}}"`).ParseFS(htmlFS, "web/dist/index.html")
	if err != nil {
		log.Fatalf("[FATAL] failed to parse template: %v", err)
	}

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Cache-Control", "no-cache")
		if err := tpl.Execute(w, data{
			Version: version,
			Links:   a.readLinks(),
		}); err != nil {
			log.Printf("[ERROR] failed to execute template: %v", err)
		}
	})
	static := http.FileServer(http.FS(staticFs))
	router.Handle("/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-cache")
		if strings.HasPrefix(r.URL.Path, "/assets/") {
			if info, err := fs.Stat(staticFs, strings.TrimPrefix(r.URL.Path, "/")); err == nil && info.Mode().IsRegular() {
				w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
			}
		}
		static.ServeHTTP(w, r)
	}))

	return router
}

func (a *app) readLinks() []Link {
	links := []Link{}

	data, err := os.ReadFile(a.linksPath)
	if errors.Is(err, os.ErrNotExist) {
		return links
	}
	if err != nil {
		log.Printf("[ERROR] read links: %v", err)
		return links
	}
	if err := json.Unmarshal(data, &links); err != nil {
		log.Printf("[ERROR] decode links: %v", err)
		return []Link{}
	}

	return links
}
func (a *app) writeLinks(links []Link) error {
	a.storageMu.Lock()
	defer a.storageMu.Unlock()

	if links == nil {
		links = []Link{}
	}

	for i := range links {
		if links[i].ID == "" {
			links[i].ID = rand.Text()
		}
	}

	data, err := json.Marshal(links)
	if err != nil {
		return err
	}

	previous, err := os.ReadFile(a.linksPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	if err == nil && bytes.Equal(previous, data) {
		return nil
	}

	dir := filepath.Dir(a.linksPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	file, err := os.CreateTemp(dir, ".links-*")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	defer file.Close()

	if _, err := file.Write(data); err != nil {
		return err
	}
	if err := file.Sync(); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}

	return os.Rename(file.Name(), a.linksPath)
}

func extractFaviconURL(ctx context.Context, siteURL string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if !strings.HasPrefix(siteURL, "http://") && !strings.HasPrefix(siteURL, "https://") {
		siteURL = "https://" + siteURL
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, siteURL, nil)
	if err != nil {
		return "", err
	}
	resp, err := faviconClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	candidates := []string{}
	tokenizer := html.NewTokenizer(io.LimitReader(resp.Body, 1<<20))
readHead:
	for resp.StatusCode == http.StatusOK {
		switch tokenizer.Next() {
		case html.ErrorToken:
			if err := tokenizer.Err(); err != io.EOF {
				return "", err
			}
			break readHead
		case html.EndTagToken:
			if tokenizer.Token().Data == "head" {
				break readHead
			}
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			if token.Data == "body" {
				break readHead
			}
			if token.Data != "link" {
				continue
			}
			var rel, href string
			for _, attr := range token.Attr {
				switch attr.Key {
				case "rel":
					rel = attr.Val
				case "href":
					href = strings.TrimSpace(attr.Val)
				}
			}
			for value := range strings.FieldsSeq(strings.ToLower(rel)) {
				if href != "" && (value == "icon" || value == "apple-touch-icon") {
					candidates = append(candidates, href)
					break readHead
				}
			}
		}
	}
	resp.Body.Close()
	candidates = append(candidates, "/favicon.ico", "/apple-touch-icon.png", "/favicon.png")

	base := resp.Request.URL
	for _, raw := range candidates {
		rel, err := url.Parse(raw)
		if err != nil {
			continue
		}
		full := base.ResolveReference(rel)
		if full.Scheme != "http" && full.Scheme != "https" {
			continue
		}
		req, err := http.NewRequestWithContext(ctx, http.MethodHead, full.String(), nil)
		if err != nil {
			continue
		}
		head, err := faviconClient.Do(req)
		if head != nil && head.Body != nil {
			head.Body.Close()
		}
		if err == nil && head.StatusCode == http.StatusOK {
			return full.String(), nil
		}
		if err := ctx.Err(); err != nil {
			return "", err
		}
	}

	return base.ResolveReference(&url.URL{Path: "/favicon.ico"}).String(), nil
}
