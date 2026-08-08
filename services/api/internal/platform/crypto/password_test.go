package crypto

import (
	"strings"
	"testing"
)

func TestHashAndVerifyPassword(t *testing.T) {
	hash, err := HashPassword("correct horse battery staple")
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	if !VerifyPassword("correct horse battery staple", hash) {
		t.Fatal("correct password should verify")
	}
	if VerifyPassword("wrong password", hash) {
		t.Fatal("wrong password must not verify")
	}
}

func TestPasswordHashesAreSaltedAndFormatCorrect(t *testing.T) {
	h1, _ := HashPassword("same-password")
	h2, _ := HashPassword("same-password")
	if h1 == h2 {
		t.Fatal("argon2id must use a random salt so identical passwords hash differently")
	}
	// PHC format: $argon2id$v=19$m=...,t=...,p=...$salt$hash
	if len(h1) < 40 || !strings.HasPrefix(h1, "$argon2id$v=19$") {
		t.Fatalf("unexpected hash format: %q", h1)
	}
}

func TestVerifyPasswordRejectsMalformed(t *testing.T) {
	if VerifyPassword("x", "not-a-valid-hash") {
		t.Fatal("malformed hash must not verify")
	}
	if VerifyPassword("x", "$argon2id$v=19$m=bad$AA$AA") {
		t.Fatal("malformed params must not verify")
	}
}
