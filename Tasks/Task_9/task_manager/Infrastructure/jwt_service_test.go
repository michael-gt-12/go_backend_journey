package Infrastructure

import (
	"strings"
	"testing"
	"time"

	"task_manager/Domain"
)

const testSecret = "unit-test-secret-at-least-32-bytes-long!!"

func testUser() Domain.User {
	return Domain.User{ID: "user-1", Username: "michael", Role: "admin"}
}

func TestJWTService_ValidToken(t *testing.T) {
	svc := NewJWTService(testSecret, time.Hour)

	token, err := svc.Generate(testUser())
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}
	if strings.Count(token, ".") != 2 {
		t.Fatalf("a JWT has 3 dot-separated parts, got %q", token)
	}

	claims, err := svc.Validate(token)
	if err != nil {
		t.Fatalf("a valid token should be accepted, got: %v", err)
	}
	if claims.UserID != "user-1" || claims.Username != "michael" || claims.Role != "admin" {
		t.Fatalf("unexpected claims: %+v", claims)
	}
	if claims.ExpiresAt == nil || !claims.ExpiresAt.After(time.Now()) {
		t.Fatal("expiry should be set in the future")
	}
}

func TestJWTService_ExpiredToken(t *testing.T) {
	svc := NewJWTService(testSecret, -time.Minute) // already expired when issued

	token, err := svc.Generate(testUser())
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}
	if _, err := svc.Validate(token); err == nil {
		t.Fatal("an expired token must be rejected")
	}
}

func TestJWTService_TamperedToken(t *testing.T) {
	svc := NewJWTService(testSecret, time.Hour)
	token, err := svc.Generate(testUser())
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}
	parts := strings.Split(token, ".")

	// Change one character in the payload so it no longer matches the signature.
	payload := []byte(parts[1])
	if payload[0] == 'A' {
		payload[0] = 'B'
	} else {
		payload[0] = 'A'
	}
	modifiedPayload := parts[0] + "." + string(payload) + "." + parts[2]

	// Break the signature itself.
	modifiedSignature := parts[0] + "." + parts[1] + "." + parts[2][:len(parts[2])-2] + "xx"

	tests := []struct {
		name  string
		token string
	}{
		{"modified payload", modifiedPayload},
		{"modified signature", modifiedSignature},
		{"signature removed", parts[0] + "." + parts[1] + "."},
		{"garbage", "not.a.jwt"},
		{"empty", ""},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := svc.Validate(tc.token); err == nil {
				t.Fatalf("%s must be rejected", tc.name)
			}
		})
	}
}

func TestJWTService_WrongSecret(t *testing.T) {
	issuer := NewJWTService(testSecret, time.Hour)
	other := NewJWTService("a-different-secret-that-is-also-long!!", time.Hour)

	token, err := issuer.Generate(testUser())
	if err != nil {
		t.Fatalf("Generate returned error: %v", err)
	}
	if _, err := other.Validate(token); err == nil {
		t.Fatal("a token signed with another secret must be rejected")
	}
}

func TestJWTService_RejectsNoneAlgorithm(t *testing.T) {
	svc := NewJWTService(testSecret, time.Hour)
	// header {"alg":"none","typ":"JWT"}, payload {"user_id":"user-1","role":"admin"}, empty signature
	unsigned := "eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0.eyJ1c2VyX2lkIjoidXNlci0xIiwicm9sZSI6ImFkbWluIn0."
	if _, err := svc.Validate(unsigned); err == nil {
		t.Fatal("a token with alg=none must be rejected")
	}
}
