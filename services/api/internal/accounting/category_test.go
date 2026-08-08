package accounting

import "testing"

func TestCategoryForCode(t *testing.T) {
	cases := map[string]string{
		"asset:clearing:mock":        "asset",
		"liability:merchant_payable": "liability",
		"equity:owner":               "equity",
		"revenue:fees":               "revenue",
		"expense:salary":             "expense",
		// Numeric-first heuristic
		"1000": "asset",
		"2000": "liability",
		"3000": "equity",
		"4000": "revenue",
		"5000": "expense",
	}
	for code, want := range cases {
		if got := categoryForCode(code); got != want {
			t.Errorf("categoryForCode(%q) = %q, want %q", code, got, want)
		}
	}
}
