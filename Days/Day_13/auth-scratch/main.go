package main

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Claims holds our custom data plus the standard registered claims (exp, iat, sub).
type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// hashPassword hashes a plaintext password with bcrypt (salt is added automatically).
func hashPassword(password string) (string, error) {
	if len(password) > 72 {
		return "", errors.New("password too long: bcrypt only uses the first 72 bytes")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// checkPassword returns true only if the password matches the stored hash.
func checkPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// createToken builds and signs an HS256 JWT that expires after ttl.
func createToken(secret []byte, userID, username, role string, ttl time.Duration) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

// parseToken validates the signature and expiry and returns the claims.
func parseToken(secret []byte, tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		// Only accept HMAC signing, otherwise an attacker could pick another algorithm.
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return secret, nil
	}, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

// tamperRole keeps the old signature but rewrites the payload so role becomes admin.
func tamperRole(tokenString string) (string, error) {
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return "", errors.New("token must have 3 parts")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return "", err
	}
	changed := strings.Replace(string(payload), `"role":"user"`, `"role":"admin"`, 1)
	if changed == string(payload) {
		return "", errors.New("role not found in payload")
	}
	parts[1] = base64.RawURLEncoding.EncodeToString([]byte(changed))
	return strings.Join(parts, "."), nil
}

func main() {
	secretStr := os.Getenv("JWT_SECRET")
	if secretStr == "" {
		fmt.Println("JWT_SECRET is not set. In PowerShell run:")
		fmt.Println(`  $env:JWT_SECRET="some-long-random-secret-for-testing"`)
		os.Exit(1)
	}
	secret := []byte(secretStr)

	// ---------- 1. Password hashing ----------
	fmt.Println("=== 1. bcrypt password hashing ===")
	password := "MyStrongPassword123"
	hash1, err := hashPassword(password)
	if err != nil {
		fmt.Println("hash error:", err)
		return
	}
	hash2, _ := hashPassword(password)
	fmt.Println("hash 1:", hash1)
	fmt.Println("hash 2:", hash2)
	fmt.Println("same password, different hashes (salt):", hash1 != hash2)
	fmt.Println("correct password matches:", checkPassword(hash1, password))
	fmt.Println("wrong password matches:  ", checkPassword(hash1, "WrongPassword"))

	// ---------- 2. Create a JWT ----------
	fmt.Println("\n=== 2. Create a JWT ===")
	tokenStr, err := createToken(secret, "42", "michael", "user", 15*time.Minute)
	if err != nil {
		fmt.Println("create error:", err)
		return
	}
	fmt.Println("token:", tokenStr)
	parts := strings.Split(tokenStr, ".")
	fmt.Println("number of parts (header.payload.signature):", len(parts))

	// ---------- 3. Parse a valid token ----------
	fmt.Println("\n=== 3. Parse a valid token ===")
	claims, err := parseToken(secret, tokenStr)
	if err != nil {
		fmt.Println("unexpected error:", err)
		return
	}
	fmt.Printf("user_id=%s username=%s role=%s expires=%s\n",
		claims.UserID, claims.Username, claims.Role, claims.ExpiresAt.Time.Format(time.RFC3339))

	// ---------- 4. Expired token ----------
	fmt.Println("\n=== 4. Expired token ===")
	expired, _ := createToken(secret, "42", "michael", "user", -1*time.Minute)
	_, err = parseToken(secret, expired)
	fmt.Println("rejected:", err != nil)
	fmt.Println("is expired error:", errors.Is(err, jwt.ErrTokenExpired))
	fmt.Println("error:", err)

	// ---------- 5. Tampered token ----------
	fmt.Println("\n=== 5. Tampered token (role changed from user to admin) ===")
	tampered, err := tamperRole(tokenStr)
	if err != nil {
		fmt.Println("tamper error:", err)
		return
	}
	_, err = parseToken(secret, tampered)
	fmt.Println("rejected:", err != nil)
	fmt.Println("is signature error:", errors.Is(err, jwt.ErrTokenSignatureInvalid))
	fmt.Println("error:", err)

	// ---------- 6. Wrong secret ----------
	fmt.Println("\n=== 6. Token checked with a different secret ===")
	_, err = parseToken([]byte("a-completely-different-secret"), tokenStr)
	fmt.Println("rejected:", err != nil)
	fmt.Println("error:", err)
}
