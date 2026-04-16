package login

import (
	"net/http"
	"os"

	"github.com/Granola5791/video-calls-service/internal/auth"
	"github.com/Granola5791/video-calls-service/internal/config"
	"github.com/Granola5791/video-calls-service/internal/db"
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
		c.JSON(http.StatusBadRequest, gin.H{"error": config.GetStringFromConfig("errors.invalid_input")})
		return
	}

	userExists, err := db.UserExistsInDB(input.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.GetStringFromConfig("errors.internal_server_error")})
		return
	}
	if !userExists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": config.GetStringFromConfig("errors.invalid_input")})
		return
	}
	hashedPassword, salt, err := db.GetUserAuthFromDB(input.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.GetStringFromConfig("errors.internal_server_error")})
		return
	}

	if !auth.VerifyPassword(hashedPassword, input.Password, salt) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": config.GetStringFromConfig("errors.invalid_input")})
		return
	}

	userInfo, err := db.GetUserInfoFromDB(input.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.GetStringFromConfig("errors.internal_server_error")})
		return
	}
	jwtKey := []byte(os.Getenv("JWT_SECRET"))
	token, err := auth.GenerateLoginToken(userInfo.UserID, userInfo.Username, userInfo.Role, jwtKey, config.GetIntFromConfig("jwt.token_exp"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.GetStringFromConfig("errors.internal_server_error")})
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     config.GetStringFromConfig("jwt.token_cookie_name"),
		Value:    token,
		Path:     "/", // visible to all paths
		Domain:   config.GetStringFromConfig("jwt.domain"),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteNoneMode,
		MaxAge:   config.GetIntFromConfig("jwt.token_exp"),
	})

	c.JSON(http.StatusOK, gin.H{"id": userInfo.UserID, "username": userInfo.Username, "role": userInfo.Role})
}
