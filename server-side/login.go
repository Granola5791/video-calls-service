package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
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

	c.JSON(http.StatusOK, gin.H{"message": GetStringFromConfig("success.login_successful")})
}