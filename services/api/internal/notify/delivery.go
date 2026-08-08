package notify

import (
	"context"
	"errors"
	"fmt"
	"net/smtp"
	"strings"
)

// Sender is the interface for delivering a notification on a channel.
type Sender interface {
	// SendEmail delivers an email. Returns error if delivery fails.
	SendEmail(ctx context.Context, to, subject, body string) error
	// SendSMS delivers an SMS.
	SendSMS(ctx context.Context, to, body string) error
}

// SMTPConfig holds the settings for the real email transport (net/smtp).
type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

// EmailSMSDeliverer sends email via SMTP and SMS via a pluggable provider. If no SMTP or
// SMS provider is configured, it falls back to structured console logging (dev) and returns
// nil so the pipeline doesn't break — the framework is production-ready once credentials
// are supplied.
type EmailSMSDeliverer struct {
	smtpCfg *SMTPConfig
	// smsHTTP is an optional callback that POSTs to an SMS provider API (Twilio-compatible).
	smsHTTP func(to, body string) error
	// console is the fallback logger when no real provider is set.
	onConsole func(channel, to, subject, body string)
}

func NewEmailSMSDeliverer(smtpCfg *SMTPConfig) *EmailSMSDeliverer {
	return &EmailSMSDeliverer{smtpCfg: smtpCfg, onConsole: func(c, to, s, b string) {
		fmt.Printf("[notify:%s] to=%s subject=%s body=%s\n", c, to, s, b)
	}}
}

// SetSMSProvider wires a real SMS provider (e.g. Twilio-compatible HTTP).
func (d *EmailSMSDeliverer) SetSMSProvider(fn func(to, body string) error) { d.smsHTTP = fn }

// SetConsole sets a custom fallback logger.
func (d *EmailSMSDeliverer) SetConsole(fn func(channel, to, subject, body string)) { d.onConsole = fn }

// SendEmail sends via SMTP, or console-logs when no SMTP is configured.
func (d *EmailSMSDeliverer) SendEmail(ctx context.Context, to, subject, body string) error {
	if d.smtpCfg == nil || d.smtpCfg.Host == "" {
		d.onConsole("email", to, subject, body)
		return nil
	}
	addr := d.smtpCfg.Host + ":" + d.smtpCfg.Port
	msg := "From: " + d.smtpCfg.From + "\r\n" +
		"To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"MIME-Version: 1.0\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n" + body
	auth := smtp.PlainAuth("", d.smtpCfg.Username, d.smtpCfg.Password, d.smtpCfg.Host)
	return smtp.SendMail(addr, auth, d.smtpCfg.From, []string{to}, []byte(msg))
}

// SendSMS sends via the configured SMS provider, or console-logs when none is set.
func (d *EmailSMSDeliverer) SendSMS(ctx context.Context, to, body string) error {
	if d.smsHTTP == nil {
		d.onConsole("sms", to, "", body)
		return nil
	}
	return d.smsHTTP(to, body)
}

// validateEmail is a light email sanity check for recipients.
func validateEmail(email string) error {
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") {
		return errors.New("invalid email address")
	}
	return nil
}
