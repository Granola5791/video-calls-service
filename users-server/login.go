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

	userInfo, err := GetUserInfoFromDB(input.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": GetStringFromConfig("errors.internal_server_error")})
		return
	}
	jwtKey := []byte(os.Getenv("JWT_SECRET"))
	token, err := GenerateLoginToken(userInfo.UserID, userInfo.Username, userInfo.Role, jwtKey, GetIntFromConfig("jwt.token_exp"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": GetStringFromConfig("errors.internal_server_error")})
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     GetStringFromConfig("jwt.token_cookie_name"),
		Value:    token,
		Path:     "/", // visible to all paths
		Domain:   GetStringFromConfig("jwt.domain"),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		MaxAge:   GetIntFromConfig("jwt.token_exp"),
	})

	c.JSON(http.StatusOK, gin.H{"id": userInfo.UserID, "username": userInfo.Username, "role": userInfo.Role})
}
