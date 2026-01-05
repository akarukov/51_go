package service

import (
	"errors"
	"fmt"
	"github.com/akarukov/51_go.git/internal/repository"
	"math/rand"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const shortUrlSize = 8
const maxGenerationTries = 5

type Repository interface {
	GetShortenedUrl(req *repository.GetShortenedUrlRequest) (*repository.GetShortenedUrlResponse, error)
	SetShortenedUrl(req *repository.SetShortenedUrlRequest) error
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

type GetShortenedUrlRequest struct {
	ShortUrl string
}
type GetShortenedUrlResponse struct {
	Url string
}

type SetShortenedUrlRequest struct {
	Url string
}
type SetShortenedUrlResponse struct {
	ShortUrl string
}

var (
	ErrGetShortenedUrlInvalidRequest = errors.New("invalid get shortenedUrl request")
	ErrRepoFailed                    = errors.New("repo failed")
)

func (f *ShortenerService) GetShortenedUrl(req *GetShortenedUrlRequest) (*GetShortenedUrlResponse, error) {
	repositoryResp, err := f.store.GetShortenedUrl(&repository.GetShortenedUrlRequest{
		ShortUrl: req.ShortUrl,
	})

	if repositoryResp != nil {
		return &GetShortenedUrlResponse{
			Url: repositoryResp.Url,
		}, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to fetch shortened url result from store: %w", err)
	}

	return nil, ErrRepoFailed

}

func (f *ShortenerService) SetShortenedUrl(req *SetShortenedUrlRequest) (*SetShortenedUrlResponse, error) {
	var err error
	var newShortUrl string
	for i := 0; i < maxGenerationTries; i++ {
		newShortUrl = f.generateNewShortUrl(shortUrlSize)
		err = f.store.SetShortenedUrl(&repository.SetShortenedUrlRequest{
			Url:      req.Url,
			ShortUrl: newShortUrl,
		})

		if err == nil {
			break
		}
	}

	if err != nil {
		return nil, fmt.Errorf("failed to set shortened url result from store: %w", err)
	}

	return &SetShortenedUrlResponse{
		ShortUrl: newShortUrl,
	}, nil
}

func (f *ShortenerService) generateNewShortUrl(length int) string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	shortKey := make([]byte, length)
	for i := range shortKey {
		shortKey[i] = charset[r.Intn(len(charset))]
	}
	return string(shortKey)
}
