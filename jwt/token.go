// Package jwt provides utilities for generating and validating JSON Web Tokens (JWT)
// for user authentication and authorization in a gaming application.
package jwt

import (
	"gaming/model"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

// SecretKey is used to sign and verify the JWT tokens.
var SecretKey = []byte("qwertuiouplkhgfdsazxcvbnm")

// Userdetails holds the information of the authenticated user.
var Userdetails model.User

// BlacklistedToken keeps track of revoked tokens.
var BlacklistedToken = make(map[string]bool)

// Claims represents the structure of the JWT claims.
type Claims struct {
	ID    uint   `json:"id"`    // User ID
	Email string `json:"email"` // User email
	Role  string `json:"role"`  // User role
	jwt.StandardClaims
}

// JwtToken generates a new JWT token for a user with the given ID, email, and role.
// It sends the signed token as a JSON response.
func JwtToken(c *gin.Context, id uint, email string, role string) {
	claims := Claims{
		ID:    id,
		Email: email,
		Role:  role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 2).Unix(), // Token expiration time
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(SecretKey)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Failed to sign token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Token": signedToken})
}

// AuthMiddleware is a middleware that checks if the JWT token is valid
func AuthMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenstring := c.GetHeader("Authorization")
		if tokenstring == "" {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Token not provided"})
			c.Abort()
			return
		}
		if BlacklistedToken[tokenstring] {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Token removed"})
			c.Abort()
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenstring, claims, func(token *jwt.Token) (interface{}, error) {
			return SecretKey, nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "Invalid token"})
			c.Abort()
			return
		}
		if claims.Role != requiredRole {
			c.JSON(http.StatusBadRequest, gin.H{"Error": "No permission"})
			c.Abort()
			return
		}

		c.Set("userid", claims.ID) // Store user ID in the context for further processing
		c.Next()                   // Proceed to the next handler
	}
}
