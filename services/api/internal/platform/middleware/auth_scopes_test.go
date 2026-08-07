package middleware

import "testing"

func TestRoleFromScopesRequiresExplicitOperationalScope(t *testing.T) {
	if got := roleFromScopes(`["payments:write"]`); got != "" { t.Fatalf("merchant payment scope must not become admin role: %q", got) }
	if got := roleFromScopes(`["role:ops","payments:read"]`); got != "ops" { t.Fatalf("expected ops role, got %q", got) }
	if got := roleFromScopes(`not-json`); got != "" { t.Fatalf("invalid scopes must fail closed, got %q", got) }
}
