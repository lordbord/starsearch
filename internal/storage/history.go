package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"starsearch/internal/types"
)

// History manages browsing history with back/forward navigation
type History struct {
	mu           sync.RWMutex
	entries      []types.HistoryEntry
	currentIndex int // Current position in history
	maxSize      int
	storePath    string
}

// NewHistory creates a new history manager
func NewHistory(storePath string, maxSize int) *History {
	if maxSize <= 0 {
		maxSize = 1000 // Default max size
	}

	h := &History{
		entries:      make([]types.HistoryEntry, 0),
		currentIndex: -1,
		maxSize:      maxSize,
		storePath:    storePath,
	}

	// Try to load existing history
	_ = h.Load() // Ignore errors

	return h
}

// Add adds a new entry to the history
func (h *History) Add(url, title string) {
	h.mu.Lock()

	// If we're not at the end of history, remove everything after current position
	if h.currentIndex < len(h.entries)-1 {
		h.entries = h.entries[:h.currentIndex+1]
	}

	// Add new entry
	entry := types.HistoryEntry{
		URL:       url,
		Title:     title,
		Timestamp: time.Now().Unix(),
	}

	h.entries = append(h.entries, entry)
	h.currentIndex = len(h.entries) - 1

	// Trim if exceeded max size
	if len(h.entries) > h.maxSize {
		// Remove oldest entries
		excess := len(h.entries) - h.maxSize
		h.entries = h.entries[excess:]
		h.currentIndex -= excess
		if h.currentIndex < 0 {
			h.currentIndex = 0
		}
	}

	h.mu.Unlock()

	// Auto-save (release lock before saving to avoid deadlock)
	_ = h.Save()
}

// Back moves back in history and returns the URL, or empty string if can't go back
func (h *History) Back() string {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.currentIndex <= 0 {
		return ""
	}

	h.currentIndex--
	return h.entries[h.currentIndex].URL
}

// Forward moves forward in history and returns the URL, or empty string if can't go forward
func (h *History) Forward() string {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.currentIndex >= len(h.entries)-1 {
		return ""
	}

	h.currentIndex++
	return h.entries[h.currentIndex].URL
}

// CanGoBack returns true if we can go back in history
func (h *History) CanGoBack() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.currentIndex > 0
}

// CanGoForward returns true if we can go forward in history
func (h *History) CanGoForward() bool {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.currentIndex < len(h.entries)-1
}

// Current returns the current history entry, or nil if none
func (h *History) Current() *types.HistoryEntry {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if h.currentIndex >= 0 && h.currentIndex < len(h.entries) {
		return &h.entries[h.currentIndex]
	}
	return nil
}

// GetAll returns all history entries
func (h *History) GetAll() []types.HistoryEntry {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// Return a copy to prevent external modification
	entries := make([]types.HistoryEntry, len(h.entries))
	copy(entries, h.entries)
	return entries
}

// Clear clears all history
func (h *History) Clear() error {
	h.mu.Lock()
	h.entries = make([]types.HistoryEntry, 0)
	h.currentIndex = -1
	h.mu.Unlock()
	return h.Save()
}

// Load loads history from disk
func (h *History) Load() error {
	data, err := os.ReadFile(h.storePath)
	if err != nil {
		return err
	}

	var entries []types.HistoryEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return err
	}

	h.mu.Lock()
	h.entries = entries
	// Set current index to end
	if len(h.entries) > 0 {
		h.currentIndex = len(h.entries) - 1
	} else {
		h.currentIndex = -1
	}
	h.mu.Unlock()

	return nil
}

// Save saves history to disk
func (h *History) Save() error {
	h.mu.RLock()
	entries := make([]types.HistoryEntry, len(h.entries))
	copy(entries, h.entries)
	h.mu.RUnlock()

	// Ensure directory exists
	dir := filepath.Dir(h.storePath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(h.storePath, data, 0600)
}
