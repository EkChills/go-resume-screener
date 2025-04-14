package auth

import (
    "fmt"
    "os"
    "time"
    "github.com/golang-jwt/jwt/v5"
)

type JwtService struct {
    AccessTokenSecret  string
    RefreshTokenSecret string
    AccessTokenTTL     time.Duration
    RefreshTokenTTL    time.Duration
}

// NewJwtService creates a new JWT service with default or custom settings
func NewJwtService() *JwtService {
    return &JwtService{
        AccessTokenSecret:  os.Getenv("JWT_ACCESS_SECRET"),
        RefreshTokenSecret: os.Getenv("JWT_REFRESH_SECRET"),
        AccessTokenTTL:     time.Hour * 1,     // 1 hour
        RefreshTokenTTL:    time.Hour * 24 * 7, // 7 days
    }
}

// TokenPair contains access and refresh tokens
type TokenPair struct {
    AccessToken  string `json:"access_token"`
    RefreshToken string `json:"refresh_token"`
}

// GenerateTokenPair creates both access and refresh tokens
func (j *JwtService) GenerateTokenPair(email string, userId uint) (*TokenPair, error) {
    // Generate access token
    accessToken, err := j.generateAccessToken(email, userId)
    if err != nil {
        return nil, fmt.Errorf("failed to generate access token: %w", err)
    }

    // Generate refresh token (now including email)
    refreshToken, err := j.generateRefreshToken(email, userId)
    if err != nil {
        return nil, fmt.Errorf("failed to generate refresh token: %w", err)
    }

    return &TokenPair{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
    }, nil
}

// generateAccessToken creates a new access token
func (j *JwtService) generateAccessToken(email string, userId uint) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "email":  email,
        "userId": userId,
        "exp":    time.Now().Add(j.AccessTokenTTL).Unix(),
        "type":   "access",
    })
    return token.SignedString([]byte(j.AccessTokenSecret))
}

// generateRefreshToken creates a new refresh token - now including email
func (j *JwtService) generateRefreshToken(email string, userId uint) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "email":  email,
        "userId": userId,
        "exp":    time.Now().Add(j.RefreshTokenTTL).Unix(),
        "type":   "refresh",
    })
    return token.SignedString([]byte(j.RefreshTokenSecret))
}

// VerifyAccessToken validates an access token and extracts claims
func (j *JwtService) VerifyAccessToken(tokenString string) (string, uint, error) {
    parsedToken, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
        return []byte(j.AccessTokenSecret), nil
    })
    if err != nil {
        return "", 0, err
    }

    claims, ok := parsedToken.Claims.(jwt.MapClaims)
    if !ok || !parsedToken.Valid {
        return "", 0, fmt.Errorf("invalid token")
    }

    // Verify token type
    tokenType, ok := claims["type"].(string)
    if !ok || tokenType != "access" {
        return "", 0, fmt.Errorf("invalid token type")
    }

    email := claims["email"].(string)
    userId := claims["userId"].(float64)
    return email, uint(userId), nil
}

// VerifyRefreshToken validates a refresh token and extracts user ID and email
func (j *JwtService) VerifyRefreshToken(tokenString string) (string, uint, error) {
    parsedToken, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
        return []byte(j.RefreshTokenSecret), nil
    })
    if err != nil {
        return "", 0, err
    }

    claims, ok := parsedToken.Claims.(jwt.MapClaims)
    if !ok || !parsedToken.Valid {
        return "", 0, fmt.Errorf("invalid token")
    }

    // Verify token type
    tokenType, ok := claims["type"].(string)
    if !ok || tokenType != "refresh" {
        return "", 0, fmt.Errorf("invalid token type")
    }

    email := claims["email"].(string)
    userId := claims["userId"].(float64)
    return email, uint(userId), nil
}

// RefreshTokens generates a new token pair using a valid refresh token
// No need to pass email separately since it's already in the refresh token
func (j *JwtService) RefreshTokens(refreshToken string) (*TokenPair, error) {
    // Verify the refresh token and extract data
    email, userId, err := j.VerifyRefreshToken(refreshToken)
    if err != nil {
        return nil, err
    }

    // Generate new token pair
    return j.GenerateTokenPair(email, userId)
}