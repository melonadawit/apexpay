package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"regexp"
)

// Never log plain FIN or full account numbers - only last4 + hash.

var (
	finRe = regexp.MustCompile(`^\d{12}$`)       // Fayda FIN 12-digit
	fanRe = regexp.MustCompile(`^[A-Z0-9]{16}$`) // FAN alias 16 chars example
	tinRe = regexp.MustCompile(`^\d{10}$`)       // ET TIN example
)

// ValidateFaydaFIN checks 12-digit format + optional check digit (Luhn placeholder, real check is NIDP side).
func ValidateFaydaFIN(fin string) bool {
	return finRe.MatchString(fin)
}

func ValidateFAN(fan string) bool {
	return fanRe.MatchString(fan)
}

func ValidateTIN(tin string) bool {
	return tinRe.MatchString(tin)
}

// HashFIN returns sha256(salt+fin) hex.
func HashFIN(salt, fin string) string {
	h := sha256.New()
	h.Write([]byte(salt + fin))
	return hex.EncodeToString(h.Sum(nil))
}

// Last4 returns last 4 chars.
func Last4(s string) string {
	if len(s) <= 4 {
		return s
	}
	return s[len(s)-4:]
}

// Encrypt AES-GCM for connector configs, file refs if needed.
func Encrypt(key []byte, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func Decrypt(key []byte, ciphertext []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}
	nonce, ct := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ct, nil)
}

// Mask helpers for outstanding UI.
func MaskFINLast4(finLast4 string) string { return "****-****-" + finLast4 }
func MaskAccount(acct string) string {
	if len(acct) <= 4 {
		return "****"
	}
	return "****" + acct[len(acct)-4:]
}
