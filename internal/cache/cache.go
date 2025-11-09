package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"

	"starsearch/internal/types"
)

// CacheEntry represents a cached page
type CacheEntry struct {
	URL       string
	Response  *types.Response
	Timestamp int64
	TTL       int64 // Time to live in seconds
}

// Cache manages page caching
type Cache struct {
	entries    map[string]*CacheEntry
	mutex      sync.RWMutex
	maxSize    int64 // Maximum cache size in bytes
	currentSize int64 // Current cache size in bytes
	defaultTTL int64 // Default TTL in seconds
}

// NewCache creates a new cache
func NewCache(maxSizeMB int, defaultTTLSeconds int64) *Cache {
	return &Cache{
		entries:    make(map[string]*CacheEntry),
		maxSize:    int64(maxSizeMB) * 1024 * 1024,
		currentSize: 0,
		defaultTTL: defaultTTLSeconds,
	}
}

// Get retrieves a cached entry if it exists and is still valid
func (c *Cache) Get(url string) (*types.Response, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	key := c.key(url)
	entry, exists := c.entries[key]
	if !exists {
		return nil, false
	}

	// Check if entry has expired
	now := time.Now().Unix()
	if entry.Timestamp+entry.TTL < now {
		// Entry expired, but don't delete here (lazy deletion)
		return nil, false
	}

	return entry.Response, true
}

// Set stores a response in the cache
func (c *Cache) Set(url string, resp *types.Response, ttl int64) {
	if resp == nil {
		return
	}

	// Only cache text/gemini and text/plain responses
	if resp.Meta != "text/gemini" && resp.Meta != "text/plain" {
		return
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	key := c.key(url)
	entrySize := int64(len(resp.Body))

	// Remove old entry if exists
	if oldEntry, exists := c.entries[key]; exists {
		c.currentSize -= int64(len(oldEntry.Response.Body))
		delete(c.entries, key)
	}

	// Check if we need to evict entries to make room
	if c.currentSize+entrySize > c.maxSize {
		c.evictOldest()
	}

	// If still too large, don't cache
	if c.currentSize+entrySize > c.maxSize {
		return
	}

	if ttl <= 0 {
		ttl = c.defaultTTL
	}

	entry := &CacheEntry{
		URL:       url,
		Response:  resp,
		Timestamp: time.Now().Unix(),
		TTL:       ttl,
	}

	c.entries[key] = entry
	c.currentSize += entrySize
}

// Clear removes all cached entries
func (c *Cache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.entries = make(map[string]*CacheEntry)
	c.currentSize = 0
}

// Invalidate removes a specific URL from cache
func (c *Cache) Invalidate(url string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	key := c.key(url)
	if entry, exists := c.entries[key]; exists {
		c.currentSize -= int64(len(entry.Response.Body))
		delete(c.entries, key)
	}
}

// evictOldest removes the oldest entries until we have enough space
func (c *Cache) evictOldest() {
	// Simple eviction: remove entries older than half the TTL
	now := time.Now().Unix()
	for key, entry := range c.entries {
		if entry.Timestamp+entry.TTL/2 < now {
			c.currentSize -= int64(len(entry.Response.Body))
			delete(c.entries, key)
		}
	}
}

// key generates a cache key from URL
func (c *Cache) key(url string) string {
	hash := sha256.Sum256([]byte(url))
	return hex.EncodeToString(hash[:])
}

// GetSize returns the current cache size in bytes
func (c *Cache) GetSize() int64 {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.currentSize
}

// GetEntryCount returns the number of cached entries
func (c *Cache) GetEntryCount() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return len(c.entries)
}

