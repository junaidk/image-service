package token

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/binary"
	"time"
)

const (
	// Secret key for generating the HMAC signature
	hmacSecret = "my_secret_key"
)

func Create(expirationTime time.Duration, secretToken string) (string, error) {

	expiry := time.Now().Add(expirationTime).Unix()

	// Generate HMAC token that includes the expiry timestamp
	token := generateHMACToken(secretToken, expiry)

	return token, nil
}

func generateHMACToken(secretToken string, expiry int64) string {
	// Convert expiry timestamp to byte slice
	expiryBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(expiryBytes, uint64(expiry))

	// Create the message to be signed, including the expiry time
	message := append([]byte(secretToken), expiryBytes...)

	// Generate the HMAC signature
	h := hmac.New(sha256.New, []byte(hmacSecret))
	h.Write(message)
	signature := h.Sum(nil)

	// Combine the expiry time and signature into the final token
	token := append(expiryBytes, signature...)
	return base64.URLEncoding.EncodeToString(token)
}

func Validate(tokenString string) bool {
	token, err := base64.URLEncoding.DecodeString(tokenString)
	if err != nil || len(token) < 8+sha256.Size {
		return false
	}

	// Extract expiry timestamp from the token
	expiry := int64(binary.BigEndian.Uint64(token[:8]))

	if time.Now().Unix() > expiry {
		return false
	}

	expectedToken := generateHMACToken("", expiry)

	return tokenString == expectedToken
}
