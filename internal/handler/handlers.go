package handlers

import (
	"github.com/akarukov/51_go.git/internal/config"
	"github.com/akarukov/51_go.git/internal/service"
	"io"
	"log"
	"net/http"
)

type ShortenerServiceInterface interface {
	GetShortenedURL(req *service.GetShortenedURLRequest) (*service.GetShortenedURLResponse, error)
	SetShortenedURL(req *service.SetShortenedURLRequest) (*service.SetShortenedURLResponse, error)
}

func Serve(cfg *config.Config, shortener ShortenerServiceInterface) error {
	h := newHandlers(cfg.ServerAddr, shortener)
	router := newRouter(h)

	srv := http.Server{
		Addr:    cfg.ServerAddr,
		Handler: router,
	}

	return srv.ListenAndServe()
}

func newRouter(h *handlers) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /{shortUrl}", h.GetShortenedURL)
	mux.HandleFunc("POST /", h.SetShortenedURL)

	return mux
}

type handlers struct {
	ShortenerService ShortenerServiceInterface
	ServerAddr       string
}

func newHandlers(serverAddr string, shortenedService ShortenerServiceInterface) *handlers {
	return &handlers{
		ShortenerService: shortenedService,
		ServerAddr:       serverAddr,
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
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
		str := "http://" + h.ServerAddr + "/" + resp.ShortURL
		_, err = w.Write([]byte(str))
		if err != nil {
			http.Error(w, "", http.StatusInternalServerError)
		}
	} else {
		log.Printf("can't generate unique shortURL")
		http.Error(w, "", http.StatusInternalServerError)
	}
}
