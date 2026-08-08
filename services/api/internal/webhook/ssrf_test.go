package webhook

import "testing"

func TestIsPrivateURL_BlocksInternalTargets(t *testing.T) {
	blocked := []string{
		"http://127.0.0.1:8080/hook",
		"http://localhost/webhook",
		"http://10.0.0.1/x",
		"http://172.16.5.5/x",
		"http://192.168.1.1/x",
		"http://169.254.169.254/latest/meta-data", // cloud metadata
		"http://100.64.0.1/x",                     // CGNAT
		"ftp://example.com/x",                     // non-http scheme
		"http:///no-host",                         // missing host
	}
	for _, u := range blocked {
		if !isPrivateURL(u) {
			t.Errorf("expected %q to be blocked (SSRF), but it was allowed", u)
		}
	}
}

func TestIsPrivateURL_AllowsPublicTargets(t *testing.T) {
	allowed := []string{
		"https://webhook.site/uuid",
		"https://example.com/callback",
		"http://8.8.8.8/hook",
		"https://api.stripe.com/v1",
	}
	for _, u := range allowed {
		if isPrivateURL(u) {
			t.Errorf("expected %q to be allowed, but it was blocked", u)
		}
	}
}
