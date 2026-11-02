package middleware

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"task_manager_api/models"
)

// TokenLifetime is how long an issued token stays valid.
const TokenLifetime = 24 * time.Hour

// Claims are the custom JWT claims: user ID, username, role, plus the
// standard registered claims (expiry, issued-at).
type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

var jwtSecret []byte

// Init sets the signing secret. Call once from main before serving requests.
func Init(secret string) {
	jwtSecret = []byte(secret)
}

// GenerateToken creates a signed HS256 token for the user.
func GenerateToken(user models.User) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   user.ID.Hex(),
		Username: user.Username,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(TokenLifetime)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret)
}

// parseToken validates the signature, the algorithm and the expiry.
// Restricting the method to HS256 blocks the "alg: none" and algorithm-swap attacks.
func parseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) { return jwtSecret, nil },
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

// AuthRequired rejects requests without a valid "Authorization: Bearer <token>"
// header with 401, and stores the caller's identity in the Gin context.
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		parts := strings.Fields(c.GetHeader("Authorization"))
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"error": "missing or malformed Authorization header"})
			return
		}

		claims, err := parseToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized,
				gin.H{"error": "invalid or expired token"})
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Next()
	}
}

// AdminOnly must run after AuthRequired. It answers 403 when the
// authenticated user is not an admin.
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		if role, _ := c.Get("role"); role != models.RoleAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden,
				gin.H{"error": "admin access required"})
			return
		}
		c.Next()
	}
}
