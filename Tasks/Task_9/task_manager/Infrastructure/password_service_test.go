package Infrastructure

import (
	"strings"
	"testing"
)

func TestPasswordService_HashAndCompare(t *testing.T) {
	svc := NewPasswordService()

	hash, err := svc.Hash("s3cret-pass")
	if err != nil {
		t.Fatalf("Hash returned error: %v", err)
	}
	if hash == "" || hash == "s3cret-pass" {
		t.Fatalf("hash must be non-empty and different from the plain password, got %q", hash)
	}
	if !strings.HasPrefix(hash, "$2") {
		t.Fatalf("expected a bcrypt hash, got %q", hash)
	}
	if !svc.Compare(hash, "s3cret-pass") {
		t.Fatal("the correct password should match its hash")
	}
}

func TestPasswordService_WrongPassword(t *testing.T) {
	svc := NewPasswordService()
	hash, err := svc.Hash("s3cret-pass")
	if err != nil {
		t.Fatalf("Hash returned error: %v", err)
	}

	tests := []struct {
		name  string
		input string
	}{
		{"completely different", "another-pass"},
		{"case changed", "S3CRET-PASS"},
		{"empty", ""},
		{"prefix only", "s3cret"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if svc.Compare(hash, tc.input) {
				t.Fatalf("password %q must not match the hash", tc.input)
			}
		})
	}
}

func TestPasswordService_InvalidHash(t *testing.T) {
	svc := NewPasswordService()
	if svc.Compare("not-a-bcrypt-hash", "anything") {
		t.Fatal("an invalid hash must never match")
	}
}

func TestPasswordService_SaltMakesHashesDifferent(t *testing.T) {
	svc := NewPasswordService()
	h1, err1 := svc.Hash("same-password")
	h2, err2 := svc.Hash("same-password")
	if err1 != nil || err2 != nil {
		t.Fatalf("Hash errors: %v, %v", err1, err2)
	}
	if h1 == h2 {
		t.Fatal("two hashes of the same password should differ because of the random salt")
	}
	if !svc.Compare(h1, "same-password") || !svc.Compare(h2, "same-password") {
		t.Fatal("both hashes should still match the password")
	}
}
