package service

import (
	"errors"
	"fmt"
	"github.com/akarukov/51_go.git/internal/repository"
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const shortURLSize = 8
const maxGenerationTries = 5

type Repository interface {
	GetShortenedURL(req *repository.GetShortenedURLRequest) (*repository.GetShortenedURLResponse, error)
	SetShortenedURL(req *repository.SetShortenedURLRequest) error
}

type ShortenerService struct {
	store Repository
	addr  string
}

func NewShortenerService(store Repository) *ShortenerService {
	return &ShortenerService{
		store: store,
	}
}

type GetShortenedURLRequest struct {
	ShortUrl string
}
type GetShortenedUrlResponse struct {
	URL string
}

type SetShortenedURLRequest struct {
	URL string
}
type SetShortenedURLResponse struct {
	ShortURL string
}

var (
	ErrGetShortenedURLInvalidRequest = errors.New("invalid get shortenedUrl request")
	ErrRepoFailed                    = errors.New("repo failed")
)

func (f *ShortenerService) GetShortenedURL(req *GetShortenedURLRequest) (*GetShortenedUrlResponse, error) {
	repositoryResp, err := f.store.GetShortenedURL(&repository.GetShortenedURLRequest{
		ShortURL: req.ShortUrl,
	})

	if repositoryResp != nil {
		return &GetShortenedUrlResponse{
			URL: repositoryResp.URL,
		}, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to fetch shortened url result from store: %w", err)
	}

	return nil, ErrRepoFailed

}

func (f *ShortenerService) SetShortenedURL(req *SetShortenedURLRequest) (*SetShortenedURLResponse, error) {
	var err error
	var newShortURL string
	for i := 0; i < maxGenerationTries; i++ {
		newShortURL = f.generateNewShortURL(shortURLSize)
		err = f.store.SetShortenedURL(&repository.SetShortenedURLRequest{
			URL:      req.URL,
			ShortURL: newShortURL,
		})

		if err == nil {
			break
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to set shortened url result from store: %w", err)
	}

	return &SetShortenedURLResponse{
		ShortURL: newShortURL,
	}, nil
}

func (f *ShortenerService) generateNewShortURL(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	shortKey := make([]byte, length)
	for i := range shortKey {
		shortKey[i] = charset[r.Intn(len(charset))]
	}
	return string(shortKey)
}
