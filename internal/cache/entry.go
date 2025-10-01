package cache

import (
	"time"

	"github.com/hyphaene/hexa/internal/jira"
)

// CacheEntry represents cached sprint ticket data with expiry metadata
type CacheEntry struct {
	SprintID   int           `json:"sprintId"`   // Sprint ID this cache represents
	CachedAt   time.Time     `json:"cachedAt"`   // Timestamp when cache was created
	TTLSeconds int           `json:"ttlSeconds"` // Time-to-live in seconds (default: 300)
	Total      int           `json:"total"`      // Total ticket count in sprint
	Issues     []jira.Ticket `json:"issues"`     // Ticket array
}

// IsExpired checks if cache has exceeded its TTL
func (c *CacheEntry) IsExpired() bool {
	return time.Since(c.CachedAt) > time.Duration(c.TTLSeconds)*time.Second
}

// Age returns human-readable cache age
func (c *CacheEntry) Age() time.Duration {
	return time.Since(c.CachedAt)
}
