package main

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(userID int, role string, jwtKey []byte) (string, error) {
	expTimeSec := GetIntFromConfig("jwt.token_exp")
	claims := jwt.MapClaims{
		"userID": userID,
		"role":   role,
		"exp":    time.Now().Add(time.Second * time.Duration(expTimeSec)).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ParseToken(tokenString string, jwtKey []byte) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return jwtKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, err
	}

	return token, nil
}

func RequireAuthentication(c *gin.Context) {

	// get cookie
	tokenString, err := c.Cookie(GetStringFromConfig("server.auth_cookie_name"))
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// validate token
	jwtKey := []byte(os.Getenv("JWT_SECRET"))
	token, err := ParseToken(tokenString, jwtKey)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// validate token claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	// check if token is expired
	if claims["exp"].(float64) < float64(time.Now().Unix()) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Set("userID", int(claims["userID"].(float64)))
	c.Set("role", claims["role"].(string))

	c.Next()
}

// can only be called after RequireAuthentication
func RequireAdmin(c *gin.Context) {
	role, _ := c.Get("role")
	roleString, ok := role.(string)
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	if !strings.Contains(roleString, "admin") {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	c.Next()
}
