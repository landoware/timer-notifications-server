package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type wikiThumbnailService struct {
	client *http.Client
	mu     sync.RWMutex
	cache  map[string]string
}

func newWikiThumbnailService() *wikiThumbnailService {
	return &wikiThumbnailService{
		client: &http.Client{Timeout: 10 * time.Second},
		cache:  make(map[string]string),
	}
}

func (s *wikiThumbnailService) ThumbnailURL(crop Crop) (string, error) {
	if crop.WikiTitle == "" {
		return "", nil
	}

	s.mu.RLock()
	thumbnailURL, ok := s.cache[crop.WikiTitle]
	s.mu.RUnlock()
	if ok {
		return thumbnailURL, nil
	}

	thumbnailURL, err := s.fetchThumbnailURL(crop.WikiTitle)
	if err != nil {
		return "", err
	}

	s.mu.Lock()
	s.cache[crop.WikiTitle] = thumbnailURL
	s.mu.Unlock()

	return thumbnailURL, nil
}

func (s *wikiThumbnailService) fetchThumbnailURL(title string) (string, error) {
	apiURL := fmt.Sprintf(
		"https://oldschool.runescape.wiki/api.php?action=query&format=json&prop=pageimages&piprop=thumbnail&pithumbsize=128&titles=%s",
		url.QueryEscape(title),
	)

	resp, err := s.client.Get(apiURL)
	if err != nil {
		return "", fmt.Errorf("fetch wiki thumbnail: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("fetch wiki thumbnail: unexpected status %d", resp.StatusCode)
	}

	var payload struct {
		Query struct {
			Pages map[string]struct {
				Thumbnail struct {
					Source string `json:"source"`
				} `json:"thumbnail"`
			} `json:"pages"`
		} `json:"query"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return "", fmt.Errorf("decode wiki thumbnail: %w", err)
	}

	for _, page := range payload.Query.Pages {
		return page.Thumbnail.Source, nil
	}

	return "", nil
}
