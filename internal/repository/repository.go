package repository

import (
	"errors"
	"fmt"
	"sync"
)

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
	errGetShortenedURLNotFound = errors.New("url not found")
	errSetShortenedURLExists   = errors.New("shortUrl is already exists")
)

func (s *Store) GetShortenedURL(req *GetShortenedURLRequest) (*GetShortenedURLResponse, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	res, ok := s.s[req.ShortURL]
	if !ok {
		return nil, fmt.Errorf("%w for shortUrl = %s", errGetShortenedURLNotFound, req.ShortURL)
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
		return fmt.Errorf("%w for shortUrl = %s", errSetShortenedURLExists, req.ShortURL)
	} else {
		s.s[req.ShortURL] = req.URL
		return nil
	}
}
