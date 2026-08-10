package i18n

import "testing"

func TestNormalize(t *testing.T) {
	cases := []struct{ in, want string }{
		{"en", "en"},
		{"am", "am"},
		{"", "en"},
		{"fr", "en"},
		{"AM", "en"}, // case-sensitive; unknown coerces to English
	}
	for _, c := range cases {
		if got := string(Normalize(c.in)); got != c.want {
			t.Errorf("Normalize(%q)=%q want %q", c.in, got, c.want)
		}
	}
}

func TestIsValid(t *testing.T) {
	if !IsValid("en") || !IsValid("am") {
		t.Fatal("en and am should be valid")
	}
	if IsValid("fr") || IsValid("") {
		t.Fatal("unknown/empty should be invalid")
	}
}

func TestCatalogGet(t *testing.T) {
	c := New()
	// English message for a known key.
	en := c.Get(LocaleEnglish, "claim_finance_approved")
	if en == "" || en == "claim_finance_approved" {
		t.Fatalf("expected English translation, got %q", en)
	}
	// Amharic must differ from English (real translation present).
	am := c.Get(LocaleAmharic, "claim_finance_approved")
	if am == en {
		t.Fatalf("expected Amharic translation to differ from English; en=%q am=%q", en, am)
	}
	// Unknown key falls back to the key itself.
	if c.Get(LocaleEnglish, "no_such_key") != "no_such_key" {
		t.Fatal("unknown key should fall back to the key")
	}
}
