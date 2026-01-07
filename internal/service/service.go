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

type RepositoryInterface interface {
	GetShortenedURL(req *repository.GetShortenedURLRequest) (*repository.GetShortenedURLResponse, error)
	SetShortenedURL(req *repository.SetShortenedURLRequest) error
}

type ShortenerService struct {
	store RepositoryInterface
	addr  string
}

func NewShortenerService(store RepositoryInterface) *ShortenerService {
	return &ShortenerService{
		store: store,
	}
}

type GetShortenedURLRequest struct {
	ShortURL string
}
type GetShortenedURLResponse struct {
	URL string
}

type SetShortenedURLRequest struct {
	URL string
}
type SetShortenedURLResponse struct {
	ShortURL string
}

var (
	errImpossibleCase = errors.New("unknown error, impossible case")
	errFailedToFetch  = errors.New("failed to fetch shortened url result from store")
	errFailedToStore  = errors.New("failed to store url")
)

func (f *ShortenerService) GetShortenedURL(req *GetShortenedURLRequest) (*GetShortenedURLResponse, error) {
	repositoryResp, err := f.store.GetShortenedURL(&repository.GetShortenedURLRequest{
		ShortURL: req.ShortURL,
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", errFailedToFetch, err)
	}

	if repositoryResp != nil {
		return &GetShortenedURLResponse{
			URL: repositoryResp.URL,
		}, nil
	}

	return nil, errImpossibleCase

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
		return nil, fmt.Errorf("%s: %w", errFailedToStore, err)
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
