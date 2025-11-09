package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"starsearch/internal/types"
)

// Bookmarks manages saved bookmarks
type Bookmarks struct {
	mu        sync.RWMutex
	bookmarks []types.Bookmark
	storePath string
}

// NewBookmarks creates a new bookmarks manager
func NewBookmarks(storePath string) *Bookmarks {
	b := &Bookmarks{
		bookmarks: make([]types.Bookmark, 0),
		storePath: storePath,
	}

	// Try to load existing bookmarks
	_ = b.Load() // Ignore errors

	return b
}

// Add adds a new bookmark
func (b *Bookmarks) Add(url, title string, tags []string) error {
	b.mu.Lock()

	// Check if bookmark already exists
	for i, bm := range b.bookmarks {
		if bm.URL == url {
			// Update existing bookmark
			b.bookmarks[i].Title = title
			b.bookmarks[i].Tags = tags
			b.mu.Unlock()
			return b.Save()
		}
	}

	// Add new bookmark
	bookmark := types.Bookmark{
		URL:   url,
		Title: title,
		Tags:  tags,
	}

	b.bookmarks = append(b.bookmarks, bookmark)

	// Sort bookmarks by title
	sort.Slice(b.bookmarks, func(i, j int) bool {
		return b.bookmarks[i].Title < b.bookmarks[j].Title
	})

	b.mu.Unlock()
	return b.Save()
}

// Remove removes a bookmark by URL
func (b *Bookmarks) Remove(url string) error {
	b.mu.Lock()

	for i, bm := range b.bookmarks {
		if bm.URL == url {
			// Remove bookmark
			b.bookmarks = append(b.bookmarks[:i], b.bookmarks[i+1:]...)
			b.mu.Unlock()
			return b.Save()
		}
	}

	b.mu.Unlock()
	return nil // URL not found, nothing to remove
}

// Get gets a bookmark by URL
func (b *Bookmarks) Get(url string) *types.Bookmark {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, bm := range b.bookmarks {
		if bm.URL == url {
			// Return a copy to prevent external modification
			bmCopy := bm
			return &bmCopy
		}
	}
	return nil
}

// GetAll returns all bookmarks
func (b *Bookmarks) GetAll() []types.Bookmark {
	b.mu.RLock()
	defer b.mu.RUnlock()

	// Return a copy to prevent external modification
	bookmarks := make([]types.Bookmark, len(b.bookmarks))
	copy(bookmarks, b.bookmarks)
	return bookmarks
}

// GetByTag returns bookmarks with a specific tag
func (b *Bookmarks) GetByTag(tag string) []types.Bookmark {
	b.mu.RLock()
	defer b.mu.RUnlock()

	result := make([]types.Bookmark, 0)
	for _, bm := range b.bookmarks {
		for _, t := range bm.Tags {
			if t == tag {
				result = append(result, bm)
				break
			}
		}
	}
	return result
}

// HasBookmark checks if a URL is bookmarked
func (b *Bookmarks) HasBookmark(url string) bool {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, bm := range b.bookmarks {
		if bm.URL == url {
			return true
		}
	}
	return false
}

// Clear clears all bookmarks
func (b *Bookmarks) Clear() error {
	b.mu.Lock()
	b.bookmarks = make([]types.Bookmark, 0)
	b.mu.Unlock()
	return b.Save()
}

// Load loads bookmarks from disk
func (b *Bookmarks) Load() error {
	data, err := os.ReadFile(b.storePath)
	if err != nil {
		return err
	}

	var bookmarks []types.Bookmark
	if err := json.Unmarshal(data, &bookmarks); err != nil {
		return err
	}

	b.mu.Lock()
	b.bookmarks = bookmarks
	b.mu.Unlock()
	return nil
}

// Save saves bookmarks to disk
func (b *Bookmarks) Save() error {
	b.mu.RLock()
	bookmarks := make([]types.Bookmark, len(b.bookmarks))
	copy(bookmarks, b.bookmarks)
	b.mu.RUnlock()

	// Ensure directory exists
	dir := filepath.Dir(b.storePath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(bookmarks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(b.storePath, data, 0600)
}
