package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shrutip04/linkvault/database"
	"github.com/shrutip04/linkvault/models"
	"github.com/shrutip04/linkvault/utils"
)

// POST /shorten
func ShortenURL(c *gin.Context) {
	var input struct {
		Original string `json:"original" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Original URL is required"})
		return
	}

	// Generate a unique short code
	shortCode := utils.GenerateShortCode(6)

	// Save to database
	_, err := database.DB.Exec(
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

// GET /:code  → redirect
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

	// Increment click count
	database.DB.Exec("UPDATE links SET clicks = clicks + 1 WHERE short = ?", code)

	c.Redirect(http.StatusMovedPermanently, link.Original)
}

// GET /links  → list all links
func GetAllLinks(c *gin.Context) {
	rows, err := database.DB.Query(
		"SELECT id, original, short, clicks, created_at FROM links",
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch links"})
		return
	}
	defer rows.Close()

	var links []models.Link
	for rows.Next() {
		var l models.Link
		rows.Scan(&l.ID, &l.Original, &l.Short, &l.Clicks, &l.CreatedAt)
		links = append(links, l)
	}

	c.JSON(http.StatusOK, links)
}