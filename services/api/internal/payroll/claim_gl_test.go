package payroll

import "testing"

// TestNilStrApprov ensures an empty approver maps to NULL (FK-safe) while a real id passes.
func TestNilStrApprov(t *testing.T) {
	if got := nilStrApprov(""); got != nil {
		t.Fatalf("expected nil for empty approver, got %v", got)
	}
	if got := nilStrApprov("user_abc"); got != "user_abc" {
		t.Fatalf("expected user_abc, got %v", got)
	}
}
