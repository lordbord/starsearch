package storage

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"starsearch/internal/types"
)

// SessionManager manages browser session persistence
type SessionManager struct {
	sessionPath string
}

// NewSessionManager creates a new session manager
func NewSessionManager(sessionPath string) *SessionManager {
	return &SessionManager{
		sessionPath: sessionPath,
	}
}

// Save saves the current session state
func (s *SessionManager) Save(tabs []types.Tab, activeIndex int) error {
	sessionTabs := make([]types.SessionTab, 0, len(tabs))
	for _, tab := range tabs {
		sessionTabs = append(sessionTabs, types.SessionTab{
			URL:    tab.URL,
			Title:  tab.Title,
			Scroll: tab.Scroll,
		})
	}

	session := types.Session{
		Tabs:        sessionTabs,
		ActiveIndex: activeIndex,
		Timestamp:   time.Now().Unix(),
	}

	// Ensure directory exists
	dir := filepath.Dir(s.sessionPath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.sessionPath, data, 0600)
}

// Load loads a saved session
func (s *SessionManager) Load() (*types.Session, error) {
	data, err := os.ReadFile(s.sessionPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // No session file exists
		}
		return nil, err
	}

	var session types.Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}

	return &session, nil
}

// Clear clears the saved session
func (s *SessionManager) Clear() error {
	if _, err := os.Stat(s.sessionPath); os.IsNotExist(err) {
		return nil // File doesn't exist, nothing to clear
	}
	return os.Remove(s.sessionPath)
}

