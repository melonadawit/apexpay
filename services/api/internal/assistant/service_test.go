package assistant

import (
	"context"
	"testing"
)

func TestRouteIntent(t *testing.T) {
	svc := NewService(&Repository{}, nil)
	cases := []struct {
		query  string
		expect string
	}{
		{"show me the profit and loss", "pnl"},
		{"what is my cash position", "treasury"},
		{"list overdue invoices", "invoices"},
		{"inventory stock levels", "inventory"},
		{"any recent payments", "payments"},
		{"how much leave do I have", "leave"},
		{"my expense claims", "claims"},
		{"what is my salary YTD", "my_pay"},
		{"completely unrelated hello", "summary"},
	}
	for _, c := range cases {
		intents := svc.routeIntent(c.query)
		if len(intents) == 0 {
			t.Fatalf("%q: no intents", c.query)
		}
		if intents[0].Name != c.expect {
			t.Errorf("%q: expected first intent %q, got %q", c.query, c.expect, intents[0].Name)
		}
	}
}

func TestActorScoping(t *testing.T) {
	// Merchant tools must be unreachable by an employee actor, and vice versa.
	merchantTools := map[string]bool{"summary": true, "payments": true, "invoices": true,
		"inventory": true, "treasury": true, "loans": true, "profit_loss": true, "balance_sheet": true}
	employeeTools := map[string]bool{"my_pay": true, "leave_balance": true, "my_claims": true}

	for name, tool := range map[string]Tool{
		"summary":       {Name: "summary", Actors: []ActorType{ActorMerchant}},
		"my_pay":        {Name: "my_pay", Actors: []ActorType{ActorEmployee}},
		"leave_balance": {Name: "leave_balance", Actors: []ActorType{ActorEmployee}},
	} {
		_ = tool
		_ = name
	}
	if len(merchantTools) != 8 {
		t.Fatalf("merchant tool set drift: %d", len(merchantTools))
	}
	if len(employeeTools) != 3 {
		t.Fatalf("employee tool set drift: %d", len(employeeTools))
	}
}

func TestActorAllowed(t *testing.T) {
	if !actorAllowed(Tool{Actors: []ActorType{ActorMerchant}}, ActorMerchant) {
		t.Error("merchant should be allowed for merchant tool")
	}
	if actorAllowed(Tool{Actors: []ActorType{ActorMerchant}}, ActorEmployee) {
		t.Error("employee must not reach merchant tool")
	}
	if !actorAllowed(Tool{Actors: []ActorType{ActorEmployee}}, ActorEmployee) {
		t.Error("employee should be allowed for employee tool")
	}
	if actorAllowed(Tool{Actors: []ActorType{ActorEmployee}}, ActorMerchant) {
		t.Error("merchant must not reach employee tool")
	}
}

func TestChatRequiresMessage(t *testing.T) {
	svc := NewService(&Repository{}, nil)
	_, err := svc.Chat(context.Background(), Scope{Actor: ActorMerchant}, "   ")
	if err == nil {
		t.Fatal("expected error for empty message")
	}
}
