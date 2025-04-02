package crypto

import (
	"testing"
)

func TestEncryptAndDecrypt(t *testing.T) {
	password := "myStrongPassword"
	privateKey := "mySuperSecretPrivateKey"

	// Encrypt the private key
	encrypted, err := Encrypt(privateKey, password)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	// Decrypt it back
	decrypted, err := Decrypt(encrypted, password)
	if err != nil {
		t.Fatalf("Decryption failed: %v", err)
	}

	// Check if original matches decrypted
	if decrypted != privateKey {
		t.Errorf("Expected %s but got %s", privateKey, decrypted)
	}
}

func TestDecryptWithWrongPassword(t *testing.T) {
	password := "correct-password"
	wrongPassword := "wrong-password"
	privateKey := "sensitive-data"

	// Encrypt with correct password
	encrypted, err := Encrypt(privateKey, password)
	if err != nil {
		t.Fatalf("Encryption failed: %v", err)
	}

	// Try decrypting with the wrong password
	_, err = Decrypt(encrypted, wrongPassword)
	if err == nil {
		t.Errorf("Expected decryption to fail with wrong password, but it succeeded")
	}
}
