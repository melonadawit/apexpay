package twofa

import (
	"testing"
	"time"

	"github.com/pquerna/otp/totp"
)

func TestGenerateSecretAndVerify(t *testing.T) {
	p := New()
	secret, url, err := p.GenerateSecret("user@example.et")
	if err != nil {
		t.Fatal(err)
	}
	if secret == "" || url == "" {
		t.Fatal("secret and otpauth URL should be non-empty")
	}
	// Generate a valid code using the same secret to simulate an authenticator app.
	code, err := totp.GenerateCode(secret, time.Now().UTC())
	if err != nil {
		t.Fatal(err)
	}
	if !p.VerifyCode(secret, code) {
		t.Fatal("a correct TOTP code should verify")
	}
	// Wrong code must fail.
	if p.VerifyCode(secret, "000000") {
		t.Fatal("wrong code must not verify")
	}
}

func TestVerifyEmptyRejected(t *testing.T) {
	p := New()
	if p.VerifyCode("SOMEsecret", "") {
		t.Fatal("empty code must not verify")
	}
	if p.VerifyCode("", "123456") {
		t.Fatal("empty secret must not verify")
	}
}

func TestRandomSecret(t *testing.T) {
	s := RandomSecret()
	if len(s) < 16 {
		t.Fatalf("random secret too short: %q", s)
	}
}
