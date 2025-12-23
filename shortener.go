package main

import (
	"crypto/rand"
	"math/big"
	"net/url"
	"strings"
)

const base62Charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// GenerateShortCode generates a cryptographically secure short code
// using base62 characters.
func GenerateShortCode(src string, length int) string {
	if length <= 0 {
		length = 6
	}

	b := make([]byte, length)
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(base62Charset))))
		if err != nil {
			// fallback to 'a' on unlikely error
			b[i] = 'a'
			continue
		}
		b[i] = base62Charset[n.Int64()]
	}
	return string(b)
}

// ValidateURL performs robust URL validation using net/url.
func ValidateURL(u string) bool {
	if u == "" {
		return false
	}
	parsed, err := url.ParseRequestURI(u)
	if err != nil {
		return false
	}
	scheme := strings.ToLower(parsed.Scheme)
	return scheme == "http" || scheme == "https"
}

// NormalizeURL returns a canonical form of the URL suitable for
// deduplication. It lowercases the host, strips default ports and
// trims trailing slashes (except for root path).
func NormalizeURL(u string) (string, error) {
	parsed, err := url.Parse(u)
	if err != nil {
		return "", err
	}

	parsed.Host = strings.ToLower(parsed.Host)

	// Remove default ports
	if (parsed.Scheme == "http" && strings.HasSuffix(parsed.Host, ":80")) ||
		(parsed.Scheme == "https" && strings.HasSuffix(parsed.Host, ":443")) {
		hostParts := strings.Split(parsed.Host, ":")
		parsed.Host = hostParts[0]
	}

	if parsed.Path != "/" {
		parsed.Path = strings.TrimRight(parsed.Path, "/")
	}

	return parsed.String(), nil
}
