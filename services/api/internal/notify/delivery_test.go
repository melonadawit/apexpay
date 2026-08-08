package notify

import (
	"context"
	"strings"
	"testing"
)

func TestEmailValidation(t *testing.T) {
	if err := validateEmail("a@b.com"); err != nil {
		t.Fatalf("valid email rejected: %v", err)
	}
	if err := validateEmail("not-an-email"); err == nil {
		t.Fatal("invalid email accepted")
	}
	if err := validateEmail("a@b"); err == nil {
		t.Fatal("missing TLD accepted")
	}
}

func TestDelivererFallsBackToConsole(t *testing.T) {
	var got []string
	d := NewEmailSMSDeliverer(nil) // no SMTP configured -> console fallback
	d.SetConsole(func(ch, to, subj, body string) {
		got = append(got, ch+":"+to+":"+subj)
	})
	if err := d.SendEmail(context.Background(), "x@y.com", "Hello", "Body"); err != nil {
		t.Fatal(err)
	}
	if err := d.SendSMS(context.Background(), "0911", "Hi"); err != nil {
		t.Fatal(err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 console deliveries, got %d", len(got))
	}
	if !strings.Contains(got[0], "email") || !strings.Contains(got[1], "sms") {
		t.Fatalf("unexpected console logs: %v", got)
	}
}
