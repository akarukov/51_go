package handlers

import (
	"github.com/akarukov/51_go.git/internal/config"
	"github.com/akarukov/51_go.git/internal/service"
	"github.com/go-chi/chi/v5"
	"io"
	"log"
	"net/http"
)

type ShortenerServiceInterface interface {
	GetShortenedURL(req *service.GetShortenedURLRequest) (*service.GetShortenedURLResponse, error)
	SetShortenedURL(req *service.SetShortenedURLRequest) (*service.SetShortenedURLResponse, error)
}

func Serve(cfg *config.Config, shortener ShortenerServiceInterface) error {
	h := newHandlers(cfg.ReferenceAddr, shortener)
	router := newRouter(h)

	srv := http.Server{
		Addr:    cfg.ServerAddr,
		Handler: router,
	}

	return srv.ListenAndServe()
}

func newRouter(h *handlers) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/{shortUrl}", h.GetShortenedURL)
	r.Post("/", h.SetShortenedURL)

	return r
}

type handlers struct {
	ShortenerService ShortenerServiceInterface
	ReferenceAddr    string
}

func newHandlers(refAddr string, shortenedService ShortenerServiceInterface) *handlers {
	return &handlers{
		ShortenerService: shortenedService,
		ReferenceAddr:    refAddr,
	}
}

func (h *handlers) GetShortenedURL(w http.ResponseWriter, r *http.Request) {
	str := r.PathValue("shortUrl")
	resp, err := h.ShortenerService.GetShortenedURL(&service.GetShortenedURLRequest{
		ShortURL: str,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if resp != nil {
		http.Redirect(w, r, resp.URL, http.StatusTemporaryRedirect)
	} else {
		log.Printf("unknown error")
		http.Error(w, "", http.StatusInternalServerError)
	}
}

func (h *handlers) SetShortenedURL(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	str := string(body)
	resp, err := h.ShortenerService.SetShortenedURL(&service.SetShortenedURLRequest{
		URL: str,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if resp != nil {
		w.WriteHeader(http.StatusCreated)
		result := h.ReferenceAddr + "/" + resp.ShortURL
		_, err = w.Write([]byte(result))
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
		}
	} else {
		log.Printf("can't generate unique shortURL")
		http.Error(w, "", http.StatusInternalServerError)
	}
}
