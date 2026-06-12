package utils

import "time"

const timeLayout = "2006-01-02 15:04:05"

// IsExpired checks if a given expiry string has passed
func IsExpired(expiresAt *string) bool {
	if expiresAt == nil || *expiresAt == "" {
		return false // no expiry = permanent
	}

	expTime, err := time.Parse(timeLayout, *expiresAt)
	if err != nil {
		return false
	}

	return time.Now().UTC().After(expTime)
}

// FormatExpiry converts a human input like "24h", "7d" into a datetime string
func FormatExpiry(duration string) (string, bool) {
	var d time.Duration

	switch duration {
	case "10s":
		d = 10 * time.Second
	case "1h":
		d = 1 * time.Hour
	case "24h":
		d = 24 * time.Hour
	case "7d":
		d = 7 * 24 * time.Hour
	case "30d":
		d = 30 * 24 * time.Hour
	default:
		return "", false // unknown format
	}

	expiry := time.Now().UTC().Add(d).Format(timeLayout)
	return expiry, true
}
