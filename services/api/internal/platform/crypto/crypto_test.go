package crypto

import (
	"bytes"
	"testing"
)

func TestDeriveKey_DomainSeparationAndStability(t *testing.T) {
	master := "super-secret-master-key"

	k1 := DeriveKey(master, "webhook-secret")
	k2 := DeriveKey(master, "connector-config")
	same := DeriveKey(master, "webhook-secret")

	// Each purpose must get a distinct key.
	if bytes.Equal(k1, k2) {
		t.Error("different purposes must derive different keys")
	}
	// Same purpose + same master must be stable (idempotent).
	if !bytes.Equal(k1, same) {
		t.Error("same purpose+master must derive the same key")
	}
	// Keys must be 32 bytes for AES-256.
	if len(k1) != 32 {
		t.Errorf("expected 32-byte AES-256 key, got %d", len(k1))
	}
}

func TestDeriveKey_ShortMasterNoPanic(t *testing.T) {
	// The original code used ConnectorEncKey[:16] which panics on short input.
	// DeriveKey must tolerate any master length.
	for _, master := range []string{"", "a", "short", "0123456789abcdef"} {
		if k := DeriveKey(master, "webhook-secret"); len(k) != 32 {
			t.Errorf("master %q: expected 32-byte key, got %d", master, len(k))
		}
	}
}
