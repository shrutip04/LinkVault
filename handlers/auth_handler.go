package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"github.com/shrutip04/linkvault/database"
	"github.com/shrutip04/linkvault/utils"
)

// POST /register
func Register(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email"    binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username, email and password are required"})
		return
	}

	// Hash the password before storing
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	_, err = database.DB.Exec(
		"INSERT INTO users (username, email, password) VALUES (?, ?, ?)",
		input.Username, input.Email, string(hashed),
	)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Username or email already exists"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Account created successfully"})
}

// POST /login
func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email"    binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required"})
		return
	}

	// Fetch user from DB
	var id int
	var username, hashedPassword string
	err := database.DB.QueryRow(
		"SELECT id, username, password FROM users WHERE email = ?", input.Email,
	).Scan(&id, &username, &hashedPassword)

	if err != nil {
		// Don't reveal whether email exists or not
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	// Compare password with hash
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, err := utils.GenerateToken(id, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":    token,
		"username": username,
		"message":  "Login successful",
	})
}