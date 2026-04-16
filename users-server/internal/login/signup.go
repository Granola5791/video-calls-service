package login

import (
	"net/http"

	"github.com/Granola5791/video-calls-service/internal/auth"
	"github.com/Granola5791/video-calls-service/internal/config"
	"github.com/Granola5791/video-calls-service/internal/db"
	"github.com/gin-gonic/gin"
)

type UserSignupInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func IsValidPassword(password string) bool {
	return len(password) >= config.GetIntFromConfig("signup.min_password_length") &&
		len(password) <= config.GetIntFromConfig("signup.max_password_length")
}

func IsValidUsername(username string) bool {
	return len(username) >= config.GetIntFromConfig("signup.min_username_length") &&
		len(username) <= config.GetIntFromConfig("signup.max_username_length")
}

func SignupUser(username string, password string) error {
	hashedPassword, salt, err := auth.GenerateNewHashedPassword(password)
	if err != nil {
		return err
	}
	err = db.InsertUserToDB(username, hashedPassword, salt)
	if err != nil {
		return err
	}

	return nil
}

func HandleSignup(c *gin.Context) {
	var input UserSignupInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": config.GetStringFromConfig("errors.invalid_input")})
		return
	}

	if !IsValidUsername(input.Username) || !IsValidPassword(input.Password) {
		c.JSON(http.StatusBadRequest, gin.H{"error": config.GetStringFromConfig("errors.invalid_input")})
		return
	}

	userExists, err := db.UserExistsInDB(input.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.GetStringFromConfig("errors.internal_server_error")})
		return
	}
	if userExists {
		c.JSON(http.StatusConflict, gin.H{"error": config.GetStringFromConfig("errors.user_exists")})
		return
	}

	err = SignupUser(input.Username, input.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": config.GetStringFromConfig("errors.internal_server_error")})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": config.GetStringFromConfig("success.user_created")})
}
