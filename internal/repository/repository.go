package repository

import (
	"errors"
	"fmt"
	"sync"
)

type Repository interface {
	GetShortenedUrl(req *GetShortenedUrlRequest) (*GetShortenedUrlResponse, error)
	SetShortenedUrlRequest(req *SetShortenedUrlRequest) error
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

type GetShortenedUrlRequest struct {
	ShortUrl string
}

type GetShortenedUrlResponse struct {
	Url string
}

var (
	ErrGetShortenedUrlNotFound = errors.New("url not found")
	ErrSetShortenedUrlExists   = errors.New("shortUrl is already exists")
)

func newErrGetShortenedUrlNotFound(shortUrl string) error {
	return fmt.Errorf("%w for shortUrl = %s", ErrGetShortenedUrlNotFound, shortUrl)
}

func newErrSetShortenedUrlExists(shortUrl string) error {
	return fmt.Errorf("%w for shortUrl = %s", ErrSetShortenedUrlExists, shortUrl)
}

func (s *Store) GetShortenedUrl(req *GetShortenedUrlRequest) (*GetShortenedUrlResponse, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	res, ok := s.s[req.ShortUrl]
	if !ok {
		return nil, newErrGetShortenedUrlNotFound(req.ShortUrl)
	}
	return &GetShortenedUrlResponse{
		Url: res,
	}, nil
}

type SetShortenedUrlRequest struct {
	Url      string
	ShortUrl string
}

func (s *Store) SetShortenedUrl(req *SetShortenedUrlRequest) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	originUrl, ok := s.s[req.ShortUrl]
	if ok && originUrl != req.Url {
		return newErrSetShortenedUrlExists(req.ShortUrl)
	} else {
		s.s[req.ShortUrl] = req.Url
		return nil
	}
}
