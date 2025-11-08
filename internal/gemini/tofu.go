package gemini

import (
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	ErrCertificateChanged = errors.New("certificate has changed since first use")
	ErrCertificateExpired = errors.New("certificate has expired")
)

// CertificateInfo stores information about a trusted certificate
type CertificateInfo struct {
	Fingerprint string    `json:"fingerprint"`
	FirstSeen   time.Time `json:"first_seen"`
	LastSeen    time.Time `json:"last_seen"`
	Subject     string    `json:"subject"`
	NotBefore   time.Time `json:"not_before"`
	NotAfter    time.Time `json:"not_after"`
}

// TOFUStore manages trusted certificates using Trust On First Use
type TOFUStore struct {
	mu           sync.RWMutex
	certs        map[string]*CertificateInfo // hostname -> cert info
	storePath    string
	OnNewCert    func(host string, cert *x509.Certificate) bool     // Callback for new certs
	OnCertChange func(host string, old, new *x509.Certificate) bool // Callback for changed certs
}

// NewTOFUStore creates a new TOFU certificate store
func NewTOFUStore(storePath string) (*TOFUStore, error) {
	store := &TOFUStore{
		certs:     make(map[string]*CertificateInfo),
		storePath: storePath,
	}

	// Try to load existing certificates
	if err := store.Load(); err != nil {
		// If file doesn't exist, that's okay, we'll create it on first save
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to load TOFU store: %w", err)
		}
	}

	return store, nil
}

// Verify verifies a certificate using TOFU
func (t *TOFUStore) Verify(host string, cert *x509.Certificate) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	// Check if certificate is expired
	now := time.Now()
	if now.Before(cert.NotBefore) || now.After(cert.NotAfter) {
		return ErrCertificateExpired
	}

	// Calculate fingerprint
	fingerprint := calculateFingerprint(cert)

	// Check if we've seen this host before
	stored, exists := t.certs[host]

	if !exists {
		// First time seeing this host
		// If callback is set, ask user for confirmation
		if t.OnNewCert != nil && !t.OnNewCert(host, cert) {
			return errors.New("certificate rejected by user")
		}

		// Trust on first use
		t.certs[host] = &CertificateInfo{
			Fingerprint: fingerprint,
			FirstSeen:   now,
			LastSeen:    now,
			Subject:     cert.Subject.String(),
			NotBefore:   cert.NotBefore,
			NotAfter:    cert.NotAfter,
		}

		// Save the updated store
		_ = t.save() // Ignore save errors for now

		return nil
	}

	// We've seen this host before, check if certificate matches
	if stored.Fingerprint != fingerprint {
		// Certificate has changed!
		// If callback is set, ask user for confirmation
		// Note: We pass nil for the old certificate because we only store
		// certificate metadata (fingerprint, dates), not the full certificate.
		// Callers can access stored.Fingerprint, stored.Subject, etc. for old cert info.
		if t.OnCertChange != nil && !t.OnCertChange(host, nil, cert) {
			return ErrCertificateChanged
		}

		// User accepted the change, update the certificate
		t.certs[host] = &CertificateInfo{
			Fingerprint: fingerprint,
			FirstSeen:   stored.FirstSeen, // Keep original first seen date
			LastSeen:    now,
			Subject:     cert.Subject.String(),
			NotBefore:   cert.NotBefore,
			NotAfter:    cert.NotAfter,
		}

		_ = t.save() // Ignore save errors for now

		return nil
	}

	// Certificate matches, update last seen
	stored.LastSeen = now
	_ = t.save() // Ignore save errors for now

	return nil
}

// GetCertInfo returns certificate information for a host
func (t *TOFUStore) GetCertInfo(host string) (*CertificateInfo, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	info, exists := t.certs[host]
	return info, exists
}

// RemoveCert removes a certificate from the store
func (t *TOFUStore) RemoveCert(host string) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	delete(t.certs, host)
	return t.save()
}

// ListHosts returns all hosts with stored certificates
func (t *TOFUStore) ListHosts() []string {
	t.mu.RLock()
	defer t.mu.RUnlock()

	hosts := make([]string, 0, len(t.certs))
	for host := range t.certs {
		hosts = append(hosts, host)
	}
	return hosts
}

// Load loads certificates from disk
func (t *TOFUStore) Load() error {
	data, err := os.ReadFile(t.storePath)
	if err != nil {
		return err
	}

	t.mu.Lock()
	defer t.mu.Unlock()

	return json.Unmarshal(data, &t.certs)
}

// Save saves certificates to disk
func (t *TOFUStore) Save() error {
	t.mu.RLock()
	defer t.mu.RUnlock()

	return t.save()
}

// save is the internal save function (must be called with lock held)
func (t *TOFUStore) save() error {
	// Ensure directory exists
	dir := filepath.Dir(t.storePath)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	data, err := json.MarshalIndent(t.certs, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal certificates: %w", err)
	}

	if err := os.WriteFile(t.storePath, data, 0600); err != nil {
		return fmt.Errorf("failed to write certificates: %w", err)
	}

	return nil
}

// calculateFingerprint calculates the SHA-256 fingerprint of a certificate
func calculateFingerprint(cert *x509.Certificate) string {
	hash := sha256.Sum256(cert.Raw)
	return hex.EncodeToString(hash[:])
}

// FormatFingerprint formats a fingerprint for display (e.g., "3A:F2:89:...")
func FormatFingerprint(fingerprint string) string {
	result := ""
	for i := 0; i < len(fingerprint); i += 2 {
		if i > 0 {
			result += ":"
		}
		result += fingerprint[i : i+2]
	}
	return result
}
