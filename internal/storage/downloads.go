package storage

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"starsearch/internal/types"
)

// Downloads manages active and completed downloads
type Downloads struct {
	downloads  map[string]*types.Download
	storePath  string
	mutex      sync.RWMutex
	maxConcurrent int
}

// NewDownloads creates a new downloads manager
func NewDownloads(storePath string, maxConcurrent int) *Downloads {
	if maxConcurrent <= 0 {
		maxConcurrent = 3
	}

	d := &Downloads{
		downloads:    make(map[string]*types.Download),
		storePath:    storePath,
		maxConcurrent: maxConcurrent,
	}

	// Load existing downloads
	_ = d.Load()

	return d
}

// Add adds a new download
func (d *Downloads) Add(url, filename string, size int64) (*types.Download, error) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	// Check if we've reached max concurrent downloads
	activeCount := 0
	for _, download := range d.downloads {
		if download.Status == types.Downloading {
			activeCount++
		}
	}

	if activeCount >= d.maxConcurrent {
		return nil, fmt.Errorf("maximum concurrent downloads (%d) reached", d.maxConcurrent)
	}

	// Generate unique ID
	id, err := generateID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate download ID: %w", err)
	}

	// Create download
	download := &types.Download{
		ID:        id,
		URL:       url,
		Filename:  filename,
		Size:      size,
		Downloaded: 0,
		Status:    types.DownloadPending,
		StartTime: time.Now().Unix(),
	}

	d.downloads[id] = download

	// Auto-save
	_ = d.Save()

	return download, nil
}

// Get gets a download by ID
func (d *Downloads) Get(id string) *types.Download {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	if download, ok := d.downloads[id]; ok {
		return download
	}
	return nil
}

// GetAll returns all downloads
func (d *Downloads) GetAll() []types.Download {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	downloads := make([]types.Download, 0, len(d.downloads))
	for _, download := range d.downloads {
		downloads = append(downloads, *download)
	}
	return downloads
}

// GetActive returns active downloads
func (d *Downloads) GetActive() []types.Download {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	active := make([]types.Download, 0)
	for _, download := range d.downloads {
		if download.Status == types.Downloading || download.Status == types.DownloadPending {
			active = append(active, *download)
		}
	}
	return active
}

// UpdateProgress updates download progress
func (d *Downloads) UpdateProgress(id string, downloaded int64) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if download, ok := d.downloads[id]; ok {
		download.Downloaded = downloaded
		if download.Downloaded >= download.Size {
			download.Status = types.DownloadCompleted
			download.FinishTime = time.Now().Unix()
		}
		_ = d.Save()
	}
}

// SetStatus sets download status
func (d *Downloads) SetStatus(id string, status types.DownloadStatus, errorMsg string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	if download, ok := d.downloads[id]; ok {
		download.Status = status
		if errorMsg != "" {
			download.Error = errorMsg
		}
		if status == types.DownloadCompleted || status == types.DownloadFailed || status == types.DownloadCancelled {
			download.FinishTime = time.Now().Unix()
		}
		_ = d.Save()
	}
}

// Remove removes a download
func (d *Downloads) Remove(id string) {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	delete(d.downloads, id)
	_ = d.Save()
}

// Clear removes all completed downloads
func (d *Downloads) Clear() error {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	for id, download := range d.downloads {
		if download.Status == types.DownloadCompleted || download.Status == types.DownloadFailed {
			delete(d.downloads, id)
		}
	}

	return d.Save()
}

// Load loads downloads from disk
func (d *Downloads) Load() error {
	data, err := os.ReadFile(d.storePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // No existing downloads file
		}
		return err
	}

	var downloads map[string]*types.Download
	if err := json.Unmarshal(data, &downloads); err != nil {
		return err
	}

	// Protect assignment with mutex to prevent race conditions
	d.mutex.Lock()
	d.downloads = downloads
	d.mutex.Unlock()
	return nil
}

// Save saves downloads to disk
func (d *Downloads) Save() error {
	// Ensure directory exists
	dir := filepath.Dir(d.storePath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(d.downloads, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(d.storePath, data, 0600)
}

// generateID generates a unique download ID
func generateID() (string, error) {
	bytes := make([]byte, 8)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GetProgressPercentage returns download progress as percentage
func (d *Downloads) GetProgressPercentage(id string) float64 {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	download := d.downloads[id]
	if download == nil || download.Size == 0 {
		return 0
	}

	percentage := float64(download.Downloaded) / float64(download.Size) * 100
	if percentage > 100 {
		percentage = 100
	}
	return percentage
}

// GetSpeed calculates download speed (bytes per second)
func (d *Downloads) GetSpeed(id string) float64 {
	d.mutex.RLock()
	defer d.mutex.RUnlock()

	download := d.downloads[id]
	if download == nil || download.Status != types.Downloading {
		return 0
	}

	elapsed := time.Now().Unix() - download.StartTime
	if elapsed <= 0 {
		return 0
	}

	return float64(download.Downloaded) / float64(elapsed)
}