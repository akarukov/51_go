package handlers

import (
	"fmt"
	"github.com/akarukov/51_go.git/internal/repository"
	"github.com/akarukov/51_go.git/internal/service"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_handlersShortenedURL(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}
	//serverAddr := "localhost:8080"
	referenceAddr := "localhost:8081"
	store := repository.NewStore()
	shortenerService := service.NewShortenerService(store)
	handlers := newHandlers(referenceAddr, shortenerService)
	router := newRouter(handlers)

	var rurl string
	localURL := referenceAddr
	longURL := "https://practicum.yandex.ru/ "
	t.Run("Create shortUrl", func(t *testing.T) {
		writer := httptest.NewRecorder()
		router.ServeHTTP(writer, httptest.NewRequest(http.MethodPost, "/", strings.NewReader(longURL)))
		assert.Equal(t, 201, writer.Code)
		body := writer.Body.String()
		assert.NotEmpty(t, body)
		_, err := fmt.Sscanf(body, localURL+"/%s", &rurl)
		assert.NoError(t, err)
		assert.NotEmpty(t, rurl)

	})

	t.Run("Fetch shortUrl", func(t *testing.T) {
		writer := httptest.NewRecorder()
		router.ServeHTTP(writer, httptest.NewRequest(http.MethodGet, "/"+rurl, nil))
		assert.Equal(t, 307, writer.Code)
		assert.Equal(t, longURL, writer.Header().Get("Location"))
	})

	t.Run("Fetch unknown Url", func(t *testing.T) {
		writer := httptest.NewRecorder()
		router.ServeHTTP(writer, httptest.NewRequest(http.MethodGet, "/BlaBla", nil))
		assert.Equal(t, 400, writer.Code)
	})

}
