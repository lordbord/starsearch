package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"

	"starsearch/internal/types"
)

// Bookmarks manages saved bookmarks
type Bookmarks struct {
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
	// Check if bookmark already exists
	for i, bm := range b.bookmarks {
		if bm.URL == url {
			// Update existing bookmark
			b.bookmarks[i].Title = title
			b.bookmarks[i].Tags = tags
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

	return b.Save()
}

// Remove removes a bookmark by URL
func (b *Bookmarks) Remove(url string) error {
	for i, bm := range b.bookmarks {
		if bm.URL == url {
			// Remove bookmark
			b.bookmarks = append(b.bookmarks[:i], b.bookmarks[i+1:]...)
			return b.Save()
		}
	}

	return nil // URL not found, nothing to remove
}

// Get gets a bookmark by URL
func (b *Bookmarks) Get(url string) *types.Bookmark {
	for _, bm := range b.bookmarks {
		if bm.URL == url {
			return &bm
		}
	}
	return nil
}

// GetAll returns all bookmarks
func (b *Bookmarks) GetAll() []types.Bookmark {
	return b.bookmarks
}

// GetByTag returns bookmarks with a specific tag
func (b *Bookmarks) GetByTag(tag string) []types.Bookmark {
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
	for _, bm := range b.bookmarks {
		if bm.URL == url {
			return true
		}
	}
	return false
}

// Clear clears all bookmarks
func (b *Bookmarks) Clear() error {
	b.bookmarks = make([]types.Bookmark, 0)
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

	b.bookmarks = bookmarks
	return nil
}

// Save saves bookmarks to disk
func (b *Bookmarks) Save() error {
	// Ensure directory exists
	dir := filepath.Dir(b.storePath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(b.bookmarks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(b.storePath, data, 0600)
}
