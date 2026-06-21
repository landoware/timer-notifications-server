package main

import (
	"testing"
)

func TestWikiThumbnailService_EmptyWikiTitle(t *testing.T) {
	svc := newWikiThumbnailService()
	crop := Crop{Name: "Test", Value: "test", WikiTitle: ""}

	url, err := svc.ThumbnailURL(crop)
	if err != nil {
		t.Errorf("ThumbnailURL() error = %v", err)
	}
	if url != "" {
		t.Errorf("ThumbnailURL() = %q, want empty string", url)
	}
}

func TestWikiThumbnailService_CacheHit(t *testing.T) {
	svc := newWikiThumbnailService()

	svc.mu.Lock()
	svc.cache["Ranarr"] = "https://example.com/ranarr.png"
	svc.mu.Unlock()

	crop := Crop{Name: "Ranarr", Value: "ranarr", WikiTitle: "Ranarr"}
	url, err := svc.ThumbnailURL(crop)
	if err != nil {
		t.Errorf("ThumbnailURL() error = %v", err)
	}
	if url != "https://example.com/ranarr.png" {
		t.Errorf("ThumbnailURL() = %q, want %q", url, "https://example.com/ranarr.png")
	}
}
