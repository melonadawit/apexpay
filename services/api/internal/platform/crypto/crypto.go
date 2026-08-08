package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"regexp"
	"strings"

	"golang.org/x/crypto/argon2"
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

// DeriveKey returns a 32-byte AES-256 key derived from a master secret and a purpose label.
//
// Purpose-scoped derivation is a defensive standard: it gives every subsystem (webhooks,
// connector configs, Fayda records, file refs, ...) its own key so a leak in one place does
// not compromise another, and it never panics on a short master secret (unlike raw [:16]
// slicing). SHA-256(prefixed) produces a fixed 32-byte output regardless of master length.
//
// The master secret itself must still be strong and stored in a secret manager in production.
func DeriveKey(master, purpose string) []byte {
	h := sha256.New()
	h.Write([]byte("apexpay:v1:" + purpose + ":"))
	h.Write([]byte(master))
	return h.Sum(nil)
}

// Mask helpers for outstanding UI.
func MaskFINLast4(finLast4 string) string { return "****-****-" + finLast4 }
func MaskAccount(acct string) string {
	if len(acct) <= 4 {
		return "****"
	}
	return "****" + acct[len(acct)-4:]
}

// ---------------------------------------------------------------------------
// Password hashing (argon2id). Argon2id is the memory-hard, GPU-resistant choice
// for credential hashing. The stored string embeds all parameters:
//   $argon2id$v=19$m=<MiB>,t=<time>,p=<threads>$<salt b64>$<hash b64>
// This lets parameters evolve over time without breaking old hashes.
// ---------------------------------------------------------------------------

const (
	argonTime    = 1
	argonMemory  = 64 * 1024 // 64 MiB
	argonThreads = 4
	argonKeyLen  = 32
)

// HashPassword returns an argon2id PHC-formatted hash of the password.
func HashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	key := argon2.IDKey([]byte(password), salt, argonTime, argonMemory, argonThreads, argonKeyLen)
	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, argonMemory, argonTime, argonThreads,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(key)), nil
}

// VerifyPassword checks a plaintext password against an argon2id PHC hash.
// It returns false (and no error) on any format/decoding failure rather than
// leaking which component was invalid.
func VerifyPassword(password, encoded string) bool {
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[1] != "argon2id" {
		return false
	}
	// parts[2] = v=19, parts[3] = m=...,t=...,p=...
	var memory uint32
	var timeIters uint32
	var threads uint8
	if _, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &timeIters, &threads); err != nil {
		return false
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false
	}
	expected, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false
	}
	actual := argon2.IDKey([]byte(password), salt, timeIters, memory, threads, uint32(len(expected)))
	return subtle.ConstantTimeCompare(actual, expected) == 1
}
