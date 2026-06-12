package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shrutip04/linkvault/database"
	"github.com/shrutip04/linkvault/models"
	"github.com/shrutip04/linkvault/utils"
)

// POST /shorten
func ShortenURL(c *gin.Context) {
	var input struct {
		Original string `json:"original" binding:"required"`
		Alias    string `json:"alias"` // optional
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Original URL is required"})
		return
	}

	// Use alias if provided, otherwise generate a random code
	shortCode := input.Alias
	if shortCode == "" {
		shortCode = utils.GenerateShortCode(6)
	}

	// Check if alias/code already exists
	var existing string
	err := database.DB.QueryRow(
		"SELECT short FROM links WHERE short = ?", shortCode,
	).Scan(&existing)

	if err == nil {
		// No error means a row was found — alias is taken
		c.JSON(http.StatusConflict, gin.H{"error": "This alias is already taken, please choose another"})
		return
	}

	// Save to database
	_, err = database.DB.Exec(
		"INSERT INTO links (original, short) VALUES (?, ?)",
		input.Original, shortCode,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save link"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"original":  input.Original,
		"short_url": "http://localhost:8080/" + shortCode,
	})
}

// GET /:code → redirect
func RedirectURL(c *gin.Context) {
	code := c.Param("code")

	var link models.Link
	err := database.DB.QueryRow(
		"SELECT id, original FROM links WHERE short = ?", code,
	).Scan(&link.ID, &link.Original)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
		return
	}

	// Update click count and last accessed time
	now := time.Now().UTC().Format("2006-01-02 15:04:05")
	database.DB.Exec(
		"UPDATE links SET clicks = clicks + 1, last_accessed = ? WHERE short = ?",
		now, code,
	)

	c.Redirect(http.StatusMovedPermanently, link.Original)
}

// GET /links → list all links with analytics
func GetAllLinks(c *gin.Context) {
	rows, err := database.DB.Query(
		"SELECT id, original, short, clicks, created_at, last_accessed FROM links ORDER BY created_at DESC",
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch links"})
		return
	}
	defer rows.Close()

	var links []models.Link
	for rows.Next() {
		var l models.Link
		rows.Scan(&l.ID, &l.Original, &l.Short, &l.Clicks, &l.CreatedAt, &l.LastAccessed)
		links = append(links, l)
	}

	if links == nil {
		links = []models.Link{} // return empty array instead of null
	}

	c.JSON(http.StatusOK, links)
}

// GET /links/stats → overall analytics summary
func GetStats(c *gin.Context) {
	var totalLinks, totalClicks, activeLinks int

	database.DB.QueryRow("SELECT COUNT(*) FROM links").Scan(&totalLinks)
	database.DB.QueryRow("SELECT COALESCE(SUM(clicks), 0) FROM links").Scan(&totalClicks)

	// For now, all links are "active" (expiry comes in Phase 3)
	activeLinks = totalLinks

	// Most popular link
	var popularShort string
	var popularClicks int
	database.DB.QueryRow(
		"SELECT short, clicks FROM links ORDER BY clicks DESC LIMIT 1",
	).Scan(&popularShort, &popularClicks)

	c.JSON(http.StatusOK, gin.H{
		"total_links":  totalLinks,
		"total_clicks": totalClicks,
		"active_links": activeLinks,
		"most_popular": gin.H{
			"short":  popularShort,
			"clicks": popularClicks,
		},
	})
}
