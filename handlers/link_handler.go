package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shrutip04/linkvault/database"
	"github.com/shrutip04/linkvault/models"
	"github.com/shrutip04/linkvault/utils"
	qrcode "github.com/skip2/go-qrcode"
	"golang.org/x/crypto/bcrypt"
)

// POST /shorten (protected)
func ShortenURL(c *gin.Context) {
	userID := c.GetInt("user_id")

	var input struct {
		Original  string `json:"original"    binding:"required"`
		Alias     string `json:"alias"`
		ExpiresIn string `json:"expires_in"`
		Password  string `json:"password"`
		Category  string `json:"category"`
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

	// Handle expiry
	var expiresAt *string
	if input.ExpiresIn != "" {
		expiry, valid := utils.FormatExpiry(input.ExpiresIn)
		if !valid {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expires_in. Use: 1h, 24h, 7d, 30d"})
			return
		}
		expiresAt = &expiry
	}

	// Handle password
	var hashedPassword *string
	if input.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
			return
		}
		h := string(hashed)
		hashedPassword = &h
	}

	// Default category
	category := input.Category
	if category == "" {
		category = "General"
	}

	_, err = database.DB.Exec(
		`INSERT INTO links (user_id, original, short, expires_at, password, category)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		userID, input.Original, shortCode, expiresAt, hashedPassword, category,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save link"})
		return
	}

	response := gin.H{
		"original":     input.Original,
		"short_url":    "http://localhost:8080/" + shortCode,
		"category":     category,
		"is_protected": input.Password != "",
	}
	if expiresAt != nil {
		response["expires_at"] = *expiresAt
	}

	c.JSON(http.StatusOK, response)
}

// GET /:code → redirect with expiry + password check
func RedirectURL(c *gin.Context) {
	code := c.Param("code")

	var link models.Link
	var passwordHash sql.NullString

	err := database.DB.QueryRow(
		`SELECT id, original, expires_at, password FROM links WHERE short = ?`, code,
	).Scan(&link.ID, &link.Original, &link.ExpiresAt, &passwordHash)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
		return
	}

	if utils.IsExpired(link.ExpiresAt) {
		c.JSON(http.StatusGone, gin.H{"error": "This link has expired"})
		return
	}

	// If password protected, require it in the request header
	if passwordHash.Valid && passwordHash.String != "" {
		inputPassword := c.GetHeader("X-Link-Password")
		if inputPassword == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "This link is password protected. Provide password in X-Link-Password header",
			})
			return
		}
		err := bcrypt.CompareHashAndPassword([]byte(passwordHash.String), []byte(inputPassword))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password"})
			return
		}
	}

	now := time.Now().UTC().Format("2006-01-02 15:04:05")
	database.DB.Exec(
		"UPDATE links SET clicks = clicks + 1, last_accessed = ? WHERE short = ?",
		now, code,
	)

	c.Redirect(http.StatusMovedPermanently, link.Original)
}

// GET /links (protected) — supports ?category=Work filter
func GetAllLinks(c *gin.Context) {
	userID := c.GetInt("user_id")
	categoryFilter := c.Query("category") // e.g. /links?category=Work

	var rows interface {
		Next() bool
		Scan(...interface{}) error
		Close() error
	}
	var err error

	if categoryFilter != "" {
		rows, err = database.DB.Query(
			`SELECT id, original, short, clicks, created_at, last_accessed, expires_at, password, category
			 FROM links WHERE user_id = ? AND category = ? ORDER BY created_at DESC`,
			userID, categoryFilter,
		)
	} else {
		rows, err = database.DB.Query(
			`SELECT id, original, short, clicks, created_at, last_accessed, expires_at, password, category
			 FROM links WHERE user_id = ? ORDER BY created_at DESC`,
			userID,
		)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch links"})
		return
	}
	defer rows.Close()

	var links []models.Link
	for rows.Next() {
		var l models.Link
		var passwordHash sql.NullString
		rows.Scan(
			&l.ID, &l.Original, &l.Short, &l.Clicks,
			&l.CreatedAt, &l.LastAccessed, &l.ExpiresAt,
			&passwordHash, &l.Category,
		)

		l.IsProtected = passwordHash.Valid && passwordHash.String != ""

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

// GET /links/stats (protected)
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

	// Get categories breakdown
	catRows, _ := database.DB.Query(
		"SELECT category, COUNT(*) FROM links WHERE user_id = ? GROUP BY category", userID,
	)
	defer catRows.Close()

	categories := map[string]int{}
	for catRows.Next() {
		var cat string
		var count int
		catRows.Scan(&cat, &count)
		categories[cat] = count
	}

	c.JSON(http.StatusOK, gin.H{
		"username":      username,
		"total_links":   totalLinks,
		"total_clicks":  totalClicks,
		"active_links":  activeLinks,
		"expired_links": expiredLinks,
		"categories":    categories,
		"most_popular": gin.H{
			"short":  popularShort,
			"clicks": popularClicks,
		},
	})
}

// DELETE /links/:id (protected)
func DeleteLink(c *gin.Context) {
	userID := c.GetInt("user_id")
	linkID := c.Param("id")

	// Make sure the link belongs to this user
	var ownerID int
	err := database.DB.QueryRow(
		"SELECT user_id FROM links WHERE id = ?", linkID,
	).Scan(&ownerID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
		return
	}

	if ownerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to delete this link"})
		return
	}

	database.DB.Exec("DELETE FROM links WHERE id = ?", linkID)
	c.JSON(http.StatusOK, gin.H{"message": "Link deleted successfully"})
}

// PUT /links/:id (protected)
func EditLink(c *gin.Context) {
	userID := c.GetInt("user_id")
	linkID := c.Param("id")

	var input struct {
		Original  string `json:"original"`
		Alias     string `json:"alias"`
		ExpiresIn string `json:"expires_in"`
		Password  string `json:"password"`
		Category  string `json:"category"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Verify ownership
	var ownerID int
	err := database.DB.QueryRow(
		"SELECT user_id FROM links WHERE id = ?", linkID,
	).Scan(&ownerID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Link not found"})
		return
	}

	if ownerID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You don't have permission to edit this link"})
		return
	}

	// Handle optional expiry update
	var expiresAt *string
	if input.ExpiresIn != "" {
		expiry, valid := utils.FormatExpiry(input.ExpiresIn)
		if !valid {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid expires_in. Use: 1h, 24h, 7d, 30d"})
			return
		}
		expiresAt = &expiry
	}

	// Handle optional password update
	var hashedPassword *string
	if input.Password != "" {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		h := string(hashed)
		hashedPassword = &h
	}

	database.DB.Exec(
		`UPDATE links SET
			original   = COALESCE(NULLIF(?, ''), original),
			short      = COALESCE(NULLIF(?, ''), short),
			expires_at = COALESCE(?, expires_at),
			password   = COALESCE(?, password),
			category   = COALESCE(NULLIF(?, ''), category)
		WHERE id = ?`,
		input.Original, input.Alias, expiresAt, hashedPassword, input.Category, linkID,
	)

	c.JSON(http.StatusOK, gin.H{"message": "Link updated successfully"})
}

// GET /qr/:code → public
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
