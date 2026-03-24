package main

import (
	"context"
	"crypto/rand"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/pkgz/rest"
	"github.com/pkgz/service"
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
}

type Link struct {
	ID string `json:"id"`

	Name   string  `json:"name"`
	URL    string  `json:"url"`
	Color  *string `json:"color,omitempty"`
	Icon   *string `json:"icon,omitempty"`
	Preset *string `json:"preset,omitempty"`
	Group  *string `json:"group,omitempty"`
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

	app := &app{
		srv:       rest.NewServer(args.Port),
		linksPath: fmt.Sprintf("%s/links.json", args.DataPath),
	}

	if err = app.run(ctx); err != nil {
		log.Printf("[ERROR] start application failed: %+v", err)
		os.Exit(4)
	}

	log.Print("[INFO] application terminated")
}

func (a *app) run(ctx context.Context) error {
	log.Printf("[INFO] application started")

	go func() {
		if err := a.srv.Run(a.router()); err != nil {
			log.Printf("[ERROR] rest server: %v", err)
		}
	}()

	<-ctx.Done()

	if err := a.srv.Shutdown(); err != nil {
		log.Printf("[ERROR] shutdown rest server: %v", err)
	}

	return nil
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
				Links:   a.getLinks(),
			})
		})
		r.Post("/", func(w http.ResponseWriter, r *http.Request) {
			var req []Link
			if err := rest.ReadBody(r, &req); err != nil {
				rest.ErrorResponse(w, r, http.StatusBadRequest, nil, err.Error())
				return
			}
			if err := a.saveLinks(req); err != nil {
				log.Printf("[ERROR] failed to save links: %v", err)
				rest.ErrorResponse(w, r, http.StatusInternalServerError, nil, "failed to save links")
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

			faviconURL, err := extractFaviconURL(uri)
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

	tpl, err := template.ParseFS(htmlFS, "web/dist/index.html")
	if err != nil {
		log.Fatalf("[FATAL] failed to parse template: %v", err)
	}

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		b, err := json.Marshal(data{
			Version: version,
			Links:   a.getLinks(),
		})
		if err != nil {
			log.Printf("[ERROR] failed to marshal links: %v", err)
			rest.ErrorResponse(w, r, http.StatusInternalServerError, nil, "failed to marshal links")
			return
		}
		w.Header().Set("Content-Type", "text/html")
		if err := tpl.Execute(w, string(b)); err != nil {
			log.Printf("[ERROR] failed to execute template: %v", err)
		}
	})
	router.Handle("/*", http.FileServer(http.FS(staticFs)))

	return router
}

func (a *app) getLinks() []Link {
	list := []Link{}

	if _, err := os.Stat(a.linksPath); errors.Is(err, os.ErrNotExist) {
		return list
	}

	f, err := os.OpenFile(a.linksPath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("[ERROR] failed to open file: %v", err)
		return list
	}
	defer f.Close()

	if err = json.NewDecoder(f).Decode(&list); err != nil {
		log.Printf("[ERROR] failed to decode links: %v", err)
		return list
	}

	return list
}
func (a *app) saveLinks(links []Link) error {
	for i := range links {
		if links[i].ID == "" {
			links[i].ID = generateShortID()
		}
	}

	f, err := os.OpenFile(a.linksPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()

	b, err := json.Marshal(links)
	if err != nil {
		return fmt.Errorf("failed to marshal links: %w", err)
	}

	if _, err = f.Write(b); err != nil {
		return fmt.Errorf("failed to write links: %w", err)
	}

	return nil
}

func generateShortID() string {
	timestamp := time.Now().UnixNano()
	randomBytes := make([]byte, 4)
	if _, err := rand.Read(randomBytes); err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x%x", timestamp, randomBytes)
}

func extractFaviconURL(siteURL string) (string, error) {
	if !strings.HasPrefix(siteURL, "http://") && !strings.HasPrefix(siteURL, "https://") {
		siteURL = "https://" + siteURL
	}

	resp, err := faviconClient.Get(siteURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", err
	}

	candidates := []string{}
	if favicon, ok := doc.Find("link[rel='icon'], link[rel='shortcut icon'], link[rel='apple-touch-icon']").Attr("href"); ok && favicon != "" {
		candidates = append(candidates, favicon)
	}
	candidates = append(candidates, "/favicon.ico", "/apple-touch-icon.png", "/favicon.png")

	for _, raw := range candidates {
		full := raw
		if !strings.HasPrefix(raw, "http://") && !strings.HasPrefix(raw, "https://") {
			base, err := url.Parse(siteURL)
			if err != nil {
				continue
			}
			rel, err := url.Parse(raw)
			if err != nil {
				continue
			}
			full = base.ResolveReference(rel).String()
		}
		head, err := faviconClient.Head(full)
		if err == nil && head.StatusCode == http.StatusOK {
			head.Body.Close()
			return full, nil
		}
		if head != nil && head.Body != nil {
			head.Body.Close()
		}
	}

	base, err := url.Parse(siteURL)
	if err != nil {
		return "", err
	}
	return base.ResolveReference(&url.URL{Path: "/favicon.ico"}).String(), nil
}
