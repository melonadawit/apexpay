package assistant

import "time"

// ActorType identifies who is talking to the assistant. Scoping is enforced per actor:
// a merchant actor may query whole-org finance/inventory/accounting; an employee actor
// may only query their own payslip, leave and claims. Read-only by construction.
type ActorType string

const (
	ActorMerchant ActorType = "merchant"
	ActorEmployee ActorType = "employee"
)

// Thread is a persistent chat conversation. Storing threads lets us surface history to the
// user and gives an append-only audit trail of what the assistant was asked and answered.
type Thread struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	MerchantID string    `json:"merchant_id"`
	Actor      ActorType `json:"actor"`
	Title      string    `json:"title"`
	CreatedAt  time.Time `json:"created_at"`
}

// Message is one user turn (role "user") or assistant turn (role "assistant") in a thread.
type Message struct {
	ID        string    `json:"id"`
	ThreadID  string    `json:"thread_id"`
	Role      string    `json:"role"` // user | assistant
	Content   string    `json:"content"`
	Intent    string    `json:"intent,omitempty"`
	ToolsUsed []string  `json:"tools_used,omitempty"`
	Data      string    `json:"data,omitempty"` // compact JSON of tool results (audit, not rendered)
	CreatedAt time.Time `json:"created_at"`
}

// Reply is the structured assistant response to a single user turn.
type Reply struct {
	ThreadID  string         `json:"thread_id"`
	Actor     ActorType      `json:"actor"`
	Answer    string         `json:"answer"`
	Intent    string         `json:"intent"`
	ToolsUsed []string       `json:"tools_used"`
	Data      map[string]any `json:"data,omitempty"`
}

// ToolResult is what a single read-only tool returns: a human-readable line plus structured data.
type ToolResult struct {
	Line string
	Data map[string]any
}
