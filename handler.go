package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Handler handles HTTP requests
type Handler struct {
	store *URLStore
}

// NewHandler creates a new HTTP handler
func NewHandler(store *URLStore) *Handler {
	return &Handler{store: store}
}

// ShortenRequest represents the request body for shortening a URL
type ShortenRequest struct {
	URL        string `json:"url"`
	CustomCode string `json:"custom_code,omitempty"`
}

// ShortenResponse represents the response for a shortened URL
type ShortenResponse struct {
	ShortCode   string `json:"short_code"`
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// HandleShorten handles POST requests to create short URLs
func (h *Handler) HandleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Limit request body to 1MB to avoid abuse
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)

	var req ShortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate and normalize URL
	if !ValidateURL(req.URL) {
		h.respondError(w, "Invalid URL format. URL must start with http:// or https://", http.StatusBadRequest)
		return
	}

	normalized, err := NormalizeURL(req.URL)
	if err != nil {
		h.respondError(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	// Check if URL already exists (using normalized form)
	if existingCode, exists := h.store.GetByOriginalURL(normalized); exists {
		h.respondSuccess(w, existingCode, normalized, r)
		return
	}

	// Generate or use custom short code
	var shortCode string
	if req.CustomCode != "" {
		// Validate custom code
		if !isValidShortCode(req.CustomCode) {
			h.respondError(w, "Invalid custom code. Use only alphanumeric characters", http.StatusBadRequest)
			return
		}
		if h.store.Exists(req.CustomCode) {
			h.respondError(w, "Custom code already exists", http.StatusConflict)
			return
		}
		shortCode = req.CustomCode
	} else {
		// Generate short code with collision handling
		maxAttempts := 10
		for i := 0; i < maxAttempts; i++ {
			shortCode = GenerateShortCode(req.URL, 6)
			if !h.store.Exists(shortCode) {
				break
			}
			if i == maxAttempts-1 {
				h.respondError(w, "Failed to generate unique short code", http.StatusInternalServerError)
				return
			}
		}
	}

	// Save the mapping (store normalized URL)
	if err := h.store.Save(shortCode, normalized); err != nil {
		h.respondError(w, err.Error(), http.StatusInternalServerError)
		return
	}

	h.respondSuccess(w, shortCode, normalized, r)
}

// HandleRedirect handles GET requests to redirect short URLs
func (h *Handler) HandleRedirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract short code from path
	shortCode := strings.TrimPrefix(r.URL.Path, "/")

	// Skip API endpoints and empty paths
	// if shortCode == "" || strings.HasPrefix(shortCode, "api/") || shortCode == "shorten" {
	// 	http.NotFound(w, r)
	// 	return
	// }
	if shortCode == "" || strings.HasPrefix(shortCode, "api/") || shortCode == "shorten" {
		// If root path, serve the modular UI. Otherwise return 404 for api/ or shorten path collisions.
		if shortCode == "" {
			ServeUI(w)
			return
		}
		http.NotFound(w, r)
		return
	}

	// Get original URL
	mapping, err := h.store.Get(shortCode)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// Increment click counter
	h.store.IncrementClicks(shortCode)

	// Redirect to original URL
	http.Redirect(w, r, mapping.OriginalURL, http.StatusMovedPermanently)
}

// HandleListURLs handles GET requests to list all URLs
func (h *Handler) HandleListURLs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.respondError(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	mappings := h.store.GetAll()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"count": len(mappings),
		"urls":  mappings,
	})
}

// respondSuccess sends a successful response
func (h *Handler) respondSuccess(w http.ResponseWriter, shortCode, originalURL string, r *http.Request) {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	host := r.Host
	if host == "" {
		host = "localhost:8080"
	}

	response := ShortenResponse{
		ShortCode:   shortCode,
		ShortURL:    fmt.Sprintf("%s://%s/%s", scheme, host, shortCode),
		OriginalURL: originalURL,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// respondError sends an error response
func (h *Handler) respondError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(ErrorResponse{Error: message})
}

// isValidShortCode validates a custom short code
func isValidShortCode(code string) bool {
	if len(code) < 3 || len(code) > 20 {
		return false
	}

	for _, char := range code {
		if !((char >= 'a' && char <= 'z') ||
			(char >= 'A' && char <= 'Z') ||
			(char >= '0' && char <= '9') ||
			char == '-' || char == '_') {
			return false
		}
	}

	return true
}
