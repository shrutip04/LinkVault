package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shrutip04/linkvault/database"
	"github.com/shrutip04/linkvault/models"
	"github.com/shrutip04/linkvault/utils"
	qrcode "github.com/skip2/go-qrcode"
)

// POST /shorten  (protected)
func ShortenURL(c *gin.Context) {
	userID := c.GetInt("user_id")

	var input struct {
		Original  string `json:"original"   binding:"required"`
		Alias     string `json:"alias"`
		ExpiresIn string `json:"expires_in"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Original URL is required"})
		return
	}

	shortCode := input.Alias
	if shortCode == "" {
		shortCode = utils.GenerateShortCode(6)
	}

	var existing string
	err := database.DB.QueryRow(
		"SELECT short FROM links WHERE short = ?", shortCode,
	).Scan(&existing)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "This alias is already taken"})
		return
	}

	var expiresAt *string
	if input.ExpiresIn != "" {
		expiry, valid := utils.FormatExpiry(input.ExpiresIn)
		if !valid {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expires_in. Use: 1h, 24h, 7d, 30d"})
			return
		}
		expiresAt = &expiry
	}

	_, err = database.DB.Exec(
		"INSERT INTO links (user_id, original, short, expires_at) VALUES (?, ?, ?, ?)",
		userID, input.Original, shortCode, expiresAt,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save link"})
		return
	}

	response := gin.H{
		"original":  input.Original,
		"short_url": "http://localhost:8080/" + shortCode,
	}
	if expiresAt != nil {
		response["expires_at"] = *expiresAt
	}

	c.JSON(http.StatusOK, response)
}

// GET /:code → public redirect
func RedirectURL(c *gin.Context) {
	code := c.Param("code")

	var link models.Link
	err := database.DB.QueryRow(
		"SELECT id, original, expires_at FROM links WHERE short = ?", code,
	).Scan(&link.ID, &link.Original, &link.ExpiresAt)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
		return
	}

	if utils.IsExpired(link.ExpiresAt) {
		c.JSON(http.StatusGone, gin.H{"error": "This link has expired"})
		return
	}

	now := time.Now().UTC().Format("2006-01-02 15:04:05")
	database.DB.Exec(
		"UPDATE links SET clicks = clicks + 1, last_accessed = ? WHERE short = ?",
		now, code,
	)

	c.Redirect(http.StatusMovedPermanently, link.Original)
}

// GET /links → user's own links only (protected)
func GetAllLinks(c *gin.Context) {
	userID := c.GetInt("user_id")

	rows, err := database.DB.Query(
		`SELECT id, original, short, clicks, created_at, last_accessed, expires_at
		 FROM links WHERE user_id = ? ORDER BY created_at DESC`, userID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch links"})
		return
	}
	defer rows.Close()

	var links []models.Link
	for rows.Next() {
		var l models.Link
		rows.Scan(&l.ID, &l.Original, &l.Short, &l.Clicks, &l.CreatedAt, &l.LastAccessed, &l.ExpiresAt)

		if utils.IsExpired(l.ExpiresAt) {
			l.Status = "expired"
		} else {
			l.Status = "active"
		}

		links = append(links, l)
	}

	if links == nil {
		links = []models.Link{}
	}

	c.JSON(http.StatusOK, links)
}

// GET /links/stats → user's personal dashboard stats (protected)
func GetStats(c *gin.Context) {
	userID := c.GetInt("user_id")
	username, _ := c.Get("username")

	var totalLinks, totalClicks int
	database.DB.QueryRow(
		"SELECT COUNT(*) FROM links WHERE user_id = ?", userID,
	).Scan(&totalLinks)
	database.DB.QueryRow(
		"SELECT COALESCE(SUM(clicks), 0) FROM links WHERE user_id = ?", userID,
	).Scan(&totalClicks)

	var popularShort string
	var popularClicks int
	database.DB.QueryRow(
		"SELECT short, clicks FROM links WHERE user_id = ? ORDER BY clicks DESC LIMIT 1", userID,
	).Scan(&popularShort, &popularClicks)

	rows, _ := database.DB.Query(
		"SELECT expires_at FROM links WHERE user_id = ?", userID,
	)
	defer rows.Close()

	activeLinks, expiredLinks := 0, 0
	for rows.Next() {
		var exp *string
		rows.Scan(&exp)
		if utils.IsExpired(exp) {
			expiredLinks++
		} else {
			activeLinks++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"username":      username,
		"total_links":   totalLinks,
		"total_clicks":  totalClicks,
		"active_links":  activeLinks,
		"expired_links": expiredLinks,
		"most_popular": gin.H{
			"short":  popularShort,
			"clicks": popularClicks,
		},
	})
}

// GET /qr/:code → public QR code
func GenerateQR(c *gin.Context) {
	code := c.Param("code")

	var original string
	err := database.DB.QueryRow(
		"SELECT original FROM links WHERE short = ?", code,
	).Scan(&original)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
		return
	}

	shortURL := "http://localhost:8080/" + code
	png, err := qrcode.Encode(shortURL, qrcode.Medium, 256)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate QR code"})
		return
	}

	c.Data(http.StatusOK, "image/png", png)
}
