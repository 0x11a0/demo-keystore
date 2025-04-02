package crypto

import (
	"crypto/sha256"
	"errors"
	"time"

	"github.com/fernet/fernet-go"
	"golang.org/x/crypto/pbkdf2"
)

const (
	salt       = "super-secret-salt"
	iterations = 390000
)

// DeriveKey generates a Fernet-compatible key from a password using PBKDF2
func DeriveKey(password string) (*fernet.Key, error) {
	keyBytes := pbkdf2.Key([]byte(password), []byte(salt), iterations, 32, sha256.New)
	key := &fernet.Key{}
	copy(key[:], keyBytes)
	return key, nil
}

// Encrypt encrypts the plaintext using Fernet
func Encrypt(plaintext, password string) (string, error) {
	key, err := DeriveKey(password)
	if err != nil {
		return "", err
	}
	token, err := fernet.EncryptAndSign([]byte(plaintext), key)
	if err != nil {
		return "", err
	}
	return string(token), nil
}

// Decrypt decrypts the token using Fernet
func Decrypt(token, password string) (string, error) {
	key, err := DeriveKey(password)
	if err != nil {
		return "", err
	}
	msg := fernet.VerifyAndDecrypt([]byte(token), time.Hour, []*fernet.Key{key})
	if msg == nil {
		return "", errors.New("decryption failed")
	}
	return string(msg), nil
}
