package controllers

import (
	"net/http"
	"strings"

	"github.com/ekchills/go-resume-screener/auth"
	"github.com/ekchills/go-resume-screener/database"
	"github.com/ekchills/go-resume-screener/services"
	"github.com/gin-gonic/gin"
)

type userDetails struct {
	Email string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
func RegisterController(ctx *gin.Context) {
	var userInput userDetails
	if err := ctx.ShouldBindJSON(&userInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userService := services.UserService{Db: database.DB}
	err := userService.Register(userInput.Email, userInput.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func LoginController(ctx *gin.Context) {
	var userInput userDetails
	if err := ctx.ShouldBindJSON(&userInput); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	userService := services.UserService{Db: database.DB}
	user, err := userService.Login(userInput.Email, userInput.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	jwtS := auth.NewJwtService()

	token, err := jwtS.GenerateTokenPair(user.Email, user.ID)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}


	ctx.JSON(http.StatusOK, gin.H{"message": "Login successful", "user": user, "tokens": struct{AccessToken string; RefreshToken string}{AccessToken: token.AccessToken, RefreshToken: token.RefreshToken}})
}

func RefreshToken (ctx *gin.Context) {
	type refTokenType struct {
		RefreshToken string `json:"refreshToken" binding:"required"`
	}
	var requestBody refTokenType 
	err := ctx.ShouldBindJSON(&requestBody)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	j := auth.NewJwtService()
	tokenPair, err := j.RefreshTokens(requestBody.RefreshToken)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Could'nt refresh token",
		})
		return
	}
	ctx.JSON(http.StatusOK, struct{AccessToken string; RefreshToken string}{AccessToken: tokenPair.AccessToken, RefreshToken: tokenPair.RefreshToken})
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Request.Header.Get("Authorization")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No authorization header provided"})
			c.Abort()
			return
		}
		jwtS := &auth.JwtService{}
		extractedToken := strings.Split(token, " ")[1]
		email, userId,  err := jwtS.VerifyAccessToken(extractedToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		c.Set("email", email)
		c.Set("userId", userId)
		c.Next()
	}
}