package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/hyphaene/hexa/internal/jira"
)

const (
	// DefaultTTL is the default cache TTL in seconds (5 minutes)
	DefaultTTL = 300
	// CacheDirName is the cache directory name under home directory
	CacheDirName = ".hexa/cache"
)

// ReadCache reads cached sprint ticket data from filesystem
func ReadCache(sprintID int) (*CacheEntry, error) {
	cachePath, err := getCachePath(sprintID)
	if err != nil {
		return nil, fmt.Errorf("failed to get cache path: %w", err)
	}

	data, err := os.ReadFile(cachePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // Cache miss, not an error
		}
		return nil, fmt.Errorf("failed to read cache file: %w", err)
	}

	var entry CacheEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return nil, fmt.Errorf("corrupted cache file: %w", err)
	}

	return &entry, nil
}

// WriteCache writes sprint ticket data to filesystem cache
func WriteCache(sprintID int, tickets []jira.Ticket, total int) error {
	cachePath, err := getCachePath(sprintID)
	if err != nil {
		return fmt.Errorf("failed to get cache path: %w", err)
	}

	// Create cache directory if missing
	cacheDir := filepath.Dir(cachePath)
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %w", err)
	}

	entry := CacheEntry{
		SprintID:   sprintID,
		CachedAt:   time.Now(),
		TTLSeconds: DefaultTTL,
		Total:      total,
		Issues:     tickets,
	}

	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal cache entry: %w", err)
	}

	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write cache file: %w", err)
	}

	return nil
}

// ShouldRefresh determines if cache should be refreshed
func ShouldRefresh(entry *CacheEntry, noCache bool) bool {
	if noCache {
		return true // --no-cache flag forces refresh
	}
	if entry == nil {
		return true // Cache miss
	}
	return entry.IsExpired() // TTL expired
}

// getCachePath returns the full path to cache file for a sprint
func getCachePath(sprintID int) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	return filepath.Join(home, CacheDirName, fmt.Sprintf("sprint_%d.json", sprintID)), nil
}
