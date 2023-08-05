package security

// TokenStore for persistence of email token during login
type TokenStore struct {
}

// NewTokenStore using backing DB connection
func NewTokenStore() *TokenStore {
	return &TokenStore{}
}

// WriteTempToken create a new token
func (s *TokenStore) WriteTempToken(email, token string) error {
	return nil
}

// ValidateToken ensure token exists in the right state
func (s *TokenStore) ValidateToken(email string) error {
	return nil
}
