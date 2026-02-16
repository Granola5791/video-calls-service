package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

type UserLoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func HandleLogin(c *gin.Context) {
	var input UserLoginInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": GetStringFromConfig("errors.invalid_input")})
		return
	}

	userExists, err := UserExistsInDB(input.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": GetStringFromConfig("errors.internal_server_error")})
		return
	}
	if !userExists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": GetStringFromConfig("errors.invalid_input")})
		return
	}
	hashedPassword, salt, err := GetUserAuthFromDB(input.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": GetStringFromConfig("errors.internal_server_error")})
		return
	}

	if !VerifyPassword(hashedPassword, input.Password, salt) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": GetStringFromConfig("errors.invalid_input")})
		return
	}

	userID, role, err := GetUserIDAndRoleFromDB(input.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": GetStringFromConfig("errors.internal_server_error")})
		return
	}
	jwtKey := []byte(os.Getenv("JWT_SECRET"))
	token, err := GenerateLoginToken(userID, role, jwtKey, GetIntFromConfig("jwt.token_exp"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": GetStringFromConfig("errors.internal_server_error")})
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     GetStringFromConfig("jwt.token_cookie_name"),
		Value:    token,
		Path:     "/", // visible to all paths
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   GetIntFromConfig("jwt.token_exp"),
	})

	c.JSON(http.StatusOK, gin.H{"id": userID, "role": role})
}
