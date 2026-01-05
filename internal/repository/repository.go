package repository

import (
	"errors"
	"fmt"
	"sync"
)

type Repository interface {
	GetShortenedUrl(req *GetShortenedURLRequest) (*GetShortenedURLResponse, error)
	SetShortenedUrlRequest(req *SetShortenedURLRequest) error
}

type Store struct {
	mux *sync.Mutex
	s   map[string]string
}

func NewStore() *Store {
	return &Store{
		mux: &sync.Mutex{},
		s:   make(map[string]string),
	}
}

type GetShortenedURLRequest struct {
	ShortURL string
}

type GetShortenedURLResponse struct {
	URL string
}

var (
	ErrGetShortenedURLNotFound = errors.New("url not found")
	ErrSetShortenedURLExists   = errors.New("shortUrl is already exists")
)

func newErrGetShortenedURLNotFound(shortURL string) error {
	return fmt.Errorf("%w for shortUrl = %s", ErrGetShortenedURLNotFound, shortURL)
}

func newErrSetShortenedURLExists(shortURL string) error {
	return fmt.Errorf("%w for shortUrl = %s", ErrSetShortenedURLExists, shortURL)
}

func (s *Store) GetShortenedURL(req *GetShortenedURLRequest) (*GetShortenedURLResponse, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	res, ok := s.s[req.ShortURL]
	if !ok {
		return nil, newErrGetShortenedURLNotFound(req.ShortURL)
	}
	return &GetShortenedURLResponse{
		URL: res,
	}, nil
}

type SetShortenedURLRequest struct {
	URL      string
	ShortURL string
}

func (s *Store) SetShortenedURL(req *SetShortenedURLRequest) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	originURL, ok := s.s[req.ShortURL]
	if ok && originURL != req.URL {
		return newErrSetShortenedURLExists(req.ShortURL)
	} else {
		s.s[req.ShortURL] = req.URL
		return nil
	}
}
