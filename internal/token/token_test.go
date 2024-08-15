package token

import (
	"testing"
	"time"
)

func TestManager_CreateAndValidate(t *testing.T) {
	manager := Manager{hmacSecret: "mysecret"}

	// Test 1: Valid Token
	expirationTime := time.Minute * 5
	token, err := manager.Create(expirationTime)
	if err != nil {
		t.Fatalf("unexpected error in Create: %v", err)
	}

	if !manager.Validate(token) {
		t.Error("expected token to be valid")
	}

	// Test 2: Expired Token
	expirationTime = time.Second * 1
	token, err = manager.Create(expirationTime)
	if err != nil {
		t.Fatalf("unexpected error in Create: %v", err)
	}

	time.Sleep(2 * time.Second)
	if manager.Validate(token) {
		t.Error("expected token to be invalid (expired)")
	}

	// Test 3: Invalid Token (tampered)
	token, err = manager.Create(expirationTime)
	if err != nil {
		t.Fatalf("unexpected error in Create: %v", err)
	}

	tamperedToken := token + "tampered"
	if manager.Validate(tamperedToken) {
		t.Error("expected tampered token to be invalid")
	}

	// Test 4: Invalid Token (wrong length)
	if manager.Validate("short") {
		t.Error("expected short token to be invalid")
	}
}
