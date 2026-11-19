package Infrastructure

import "golang.org/x/crypto/bcrypt"

// PasswordService hashes and verifies passwords with bcrypt.
type PasswordService struct{}

// NewPasswordService creates a PasswordService.
func NewPasswordService() *PasswordService {
	return &PasswordService{}
}

// Hash returns the bcrypt hash of a plaintext password.
func (s *PasswordService) Hash(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// Compare reports whether password matches the bcrypt hash.
func (s *PasswordService) Compare(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
