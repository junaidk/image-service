package token

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"time"
)

type Manager struct {
	hmacSecret string
}

func New(secret string) *Manager {
	return &Manager{hmacSecret: secret}
}

func (m Manager) Create(expirationTime time.Duration) (string, error) {
	expiry := time.Now().Add(expirationTime).Unix()
	token := m.generateHMACToken(expiry)

	return token, nil
}

func (m Manager) generateHMACToken(expiry int64) string {
	// Convert expiry timestamp to byte slice
	expiryBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(expiryBytes, uint64(expiry))

	message := expiryBytes

	// Generate the HMAC signature
	h := hmac.New(sha256.New, []byte(m.hmacSecret))
	h.Write(message)
	signature := h.Sum(nil)

	// Combine the expiry time and signature into the final token
	token := append(expiryBytes, signature...)
	return base64.URLEncoding.EncodeToString(token)
}

func (m Manager) Validate(tokenString string) bool {
	token, err := base64.URLEncoding.DecodeString(tokenString)
	if err != nil || len(token) < 8+sha256.Size {
		return false
	}

	// Extract expiry timestamp from the token
	expiry := int64(binary.BigEndian.Uint64(token[:8]))

	if time.Now().Unix() > expiry {
		return false
	}

	expectedToken := m.generateHMACToken(expiry)

	return tokenString == expectedToken
}
