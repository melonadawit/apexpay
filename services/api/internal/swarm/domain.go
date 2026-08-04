package swarm

import "time"

type SessionStatus string
const (
	StatusPlanning          SessionStatus = "planning"
	StatusExecuting         SessionStatus = "executing"
	StatusNeedsConfirmation SessionStatus = "needs_confirmation"
	StatusCompleted         SessionStatus = "completed"
	StatusFailed            SessionStatus = "failed"
	StatusCancelled         SessionStatus = "cancelled"
)

type SwarmSession struct {
	ID                   string
	MerchantID           string
	UserID               *string
	Goal                 string
	Plan                 []PlanStep
	Status               SessionStatus
	ConfirmationRequired bool
	ConfirmationData     map[string]any
	FinalOutput          string
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type PlanStep struct {
	Step        int              `json:"step"`
	Tool        string           `json:"tool"` // create_payment_link, create_payout, calculate_payroll, etc
	Description string           `json:"description"`
	Args        map[string]any   `json:"args"`
	Status      string           `json:"status"` // pending, executing, succeeded, failed, needs_confirmation
	Result      map[string]any   `json:"result,omitempty"`
}

type AgentRun struct {
	ID             string
	MerchantID     string
	ThreadID       *string
	SwarmSessionID *string
	InputText      string
	Intent         string
	ToolCalls      []ToolCall
	OutputText     string
	Model          string // rules_v1, gpt-4, etc
	CreatedAt      time.Time
}

type ToolCall struct {
	Tool      string         `json:"tool"`
	Args      map[string]any `json:"args"`
	Result    map[string]any `json:"result"`
	LatencyMS int            `json:"latency_ms"`
	Status    string         `json:"status"`
}

// Tool registry - optimal map O(1) lookup
type ToolDefinition struct {
	Name        string
	Description string
	ArgsSchema  map[string]any // JSON schema
	Threshold   float64        // amount threshold requiring confirmation
	RoleAllowed []string       // owner, admin, etc
}
