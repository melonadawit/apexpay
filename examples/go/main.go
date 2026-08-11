// Package main is a payments-only ApexPay example in Go, mirroring the SDKs in
// sdk/. It shows the classic gateway flow a merchant uses for e-commerce / a
// custom backend: initialize -> redirect to checkout_url -> verify on return ->
// verify the webhook HMAC signature.
//
// Run:  go run . --key sk_test_... --amount 100.00
package main

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

const baseURL = "http://localhost:8080"

type Payment struct {
	ID          string `json:"id"`
	TxRef       string `json:"tx_ref"`
	Amount      string `json:"amount"`
	Currency    string `json:"currency"`
	Status      string `json:"status"`
	CheckoutURL string `json:"checkout_url"`
}

func call(method, path, apiKey, idem string, body any) (*Payment, error) {
	var rdr io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		rdr = bytes.NewReader(b)
	}
	req, err := http.NewRequest(method, baseURL+"/v1/"+path, rdr)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	if idem != "" {
		req.Header.Set("Idempotency-Key", idem)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	raw, _ := io.ReadAll(resp.Body)

	var env struct {
		Data Payment `json:"data"`
	}
	if err := json.Unmarshal(raw, &env); err != nil {
		return nil, fmt.Errorf("bad response %d: %s", resp.StatusCode, raw)
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("ApexPay error %d: %s", resp.StatusCode, raw)
	}
	return &env.Data, nil
}

// verifyWebhookSignature mirrors the SDKs: HMAC-SHA256(secret, rawBody).
func verifyWebhookSignature(secret, rawBody, signature string) bool {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(rawBody))
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}

func main() {
	apiKey := flag.String("key", "", "ApexPay secret key (sk_test_... / sk_live_...)")
	amount := flag.String("amount", "100.00", "amount to charge in ETB")
	flag.Parse()
	if *apiKey == "" {
		fmt.Fprintln(os.Stderr, "usage: go run . --key sk_test_... [--amount 100.00]")
		os.Exit(2)
	}

	// 1. Initialize a payment
	p, err := call("POST", "transactions/initialize", *apiKey, "go-example-idem", map[string]any{
		"tx_ref":         "go-example-001",
		"amount":         *amount,
		"currency":       "ETB",
		"method":         "telebirr",
		"callback_url":   "https://store.example/api/apexpay-webhook",
		"customer_email": "buyer@example.com",
	})
	if err != nil {
		fmt.Println("initialize error:", err)
		os.Exit(1)
	}
	fmt.Printf("INITIALIZE ok: %s status=%s checkout_url=%s\n", p.ID, p.Status, p.CheckoutURL)

	// 2. (simulate customer paying, then) verify on return
	verified, err := call("GET", "transactions/verify/go-example-001", *apiKey, "", nil)
	if err != nil {
		fmt.Println("verify error:", err)
		os.Exit(1)
	}
	fmt.Printf("VERIFY ok: status=%s\n", verified.Status)

	// 3. Webhook signature verification
	sig := verifyWebhookSignature("whsec_secret", `{"event_type":"payment.succeeded"}`, "<signature>")
	fmt.Printf("Webhook signature helper ready (sample check=%v)\n", sig)

	// 4. Create a shareable payment link
	link, err := call("POST", "payment_links", *apiKey, "", map[string]any{
		"amount": "50.00", "currency": "ETB", "description": "Go example",
	})
	if err != nil {
		fmt.Println("link error:", err)
		os.Exit(1)
	}
	fmt.Printf("PAYMENT LINK ok: %s\n", link.CheckoutURL)
}
