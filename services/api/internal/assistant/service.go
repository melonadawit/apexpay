package assistant

import (
	"context"
	"fmt"
	"strings"

	"apexpay/internal/i18n"
	"apexpay/internal/platform/errors"
)

// cat is the shared message catalog used to localize assistant framing and tool lines.
var cat = i18n.New()

// Scope carries everything a read-only tool needs and is fixed before any tool runs. All
// tools receive the same scope; each tool enforces its own actor allowlist so no tool can
// be reached by the wrong actor.
type Scope struct {
	MerchantID string
	UserID     string
	EmployeeID string // empty unless ActorEmployee
	Actor      ActorType
	Locale     i18n.Locale // resolved language for answers (en/am)
}

// Tool is a single read-only capability. Run must never mutate state.
type Tool struct {
	Name        string
	Description string
	Actors      []ActorType // which actors may invoke this tool
	Run         func(ctx context.Context, scope Scope) (ToolResult, error)
}

// Service routes a free-text query to the matching read-only tools and composes an answer.
type Service struct {
	repo  *Repository
	tools map[string]Tool
}

func NewService(repo *Repository, tools []Tool) *Service {
	m := make(map[string]Tool, len(tools))
	for _, t := range tools {
		m[t.Name] = t
	}
	return &Service{repo: repo, tools: m}
}

// Intent is a user-turn classification used to select tools.
type Intent struct {
	Name  string
	Tools []string
}

// Chat processes one user turn for the given actor scope, returns the reply, and persists
// both the user and assistant messages (append-only audit).
func (s *Service) Chat(ctx context.Context, scope Scope, text string) (*Reply, error) {
	text = strings.TrimSpace(text)
	if text == "" {
		return nil, errors.Validation("message required")
	}

	// 1. Route intent (deterministic rules; LLM routing is a later phase per ADR-006).
	intents := s.routeIntent(text)

	// 2. Execute the allowed tools for this actor, deduped and in order.
	seen := map[string]bool{}
	results := []ToolResult{}
	toolsUsed := []string{}
	for _, it := range intents {
		for _, tname := range it.Tools {
			if seen[tname] {
				continue
			}
			seen[tname] = true
			tool, ok := s.tools[tname]
			if !ok {
				continue
			}
			if !actorAllowed(tool, scope.Actor) {
				continue
			}
			res, err := tool.Run(ctx, scope)
			if err != nil {
				// Read-only tool errors must not crash the turn; report gracefully.
				results = append(results, ToolResult{Line: fmt.Sprintf("⚠ %s %s.", tname, cat.Get(scope.Locale, "assistant_unavailable"))})
				toolsUsed = append(toolsUsed, tname)
				continue
			}
			results = append(results, res)
			toolsUsed = append(toolsUsed, tname)
		}
	}

	// 3. Compose answer.
	answer, intentName := s.compose(text, results, intents, scope.Actor, scope.Locale)

	// 4. Persist thread + messages.
	thread := &Thread{
		UserID: scope.UserID, MerchantID: scope.MerchantID, Actor: scope.Actor,
		Title: shortTitle(text),
	}
	if err := s.repo.CreateThread(ctx, thread); err != nil {
		return nil, err
	}
	_ = s.repo.AppendMessage(ctx, &Message{ThreadID: thread.ID, Role: "user", Content: text})
	_ = s.repo.AppendMessage(ctx, &Message{
		ThreadID: thread.ID, Role: "assistant", Content: answer,
		Intent: intentName, ToolsUsed: toolsUsed,
	})

	data := map[string]any{}
	for _, r := range results {
		if r.Data != nil {
			data["result"] = r.Data
			break // first structured result is enough for P1
		}
	}

	return &Reply{
		ThreadID: thread.ID, Actor: scope.Actor, Answer: answer,
		Intent: intentName, ToolsUsed: toolsUsed, Data: data,
	}, nil
}

// routeIntent returns candidate intents in priority order for a free-text query.
func (s *Service) routeIntent(text string) []Intent {
	lower := strings.ToLower(text)
	var out []Intent

	// Employee self-service intents take priority so a merchant employee's "my claim"
	// is never misclassified as org expense analytics.
	if containsAny(lower, "claim", "expense report", "reimburse") {
		out = append(out, Intent{Name: "claims", Tools: []string{"my_claims"}})
	}
	if containsAny(lower, "payslip", "my salary", "my pay", "net pay", "ytd", "my earning") {
		out = append(out, Intent{Name: "my_pay", Tools: []string{"my_pay"}})
	}
	if containsAny(lower, "leave", "vacation", "holiday", "days remaining") {
		out = append(out, Intent{Name: "leave", Tools: []string{"leave_balance"}})
	}

	// Merchant/org analytics intents.
	if containsAny(lower, "profit", "pnl", "income statement", "revenue", "loss") {
		out = append(out, Intent{Name: "pnl", Tools: []string{"profit_loss"}})
	}
	if containsAny(lower, "balance sheet", "asset", "liability", "equity") {
		out = append(out, Intent{Name: "balance_sheet", Tools: []string{"balance_sheet"}})
	}
	if containsAny(lower, "cash", "treasury", "position", "forecast", "liquidity") {
		out = append(out, Intent{Name: "treasury", Tools: []string{"treasury"}})
	}
	if containsAny(lower, "inventory", "stock", "product", "low stock", "sku") {
		out = append(out, Intent{Name: "inventory", Tools: []string{"inventory"}})
	}
	if containsAny(lower, "invoice", "overdue", "aging", "receivable", "due") {
		out = append(out, Intent{Name: "invoices", Tools: []string{"invoices"}})
	}
	if containsAny(lower, "loan", "credit line", "borrowing") {
		out = append(out, Intent{Name: "loans", Tools: []string{"loans"}})
	}
	if containsAny(lower, "expense", "spend", "cost", "outflow") {
		out = append(out, Intent{Name: "expenses", Tools: []string{"expenses"}})
	}
	if containsAny(lower, "payment", "transaction", "tpv", "sale", "refund") {
		out = append(out, Intent{Name: "payments", Tools: []string{"payments"}})
	}

	// Fallback: summary covers the org.
	if len(out) == 0 {
		out = append(out, Intent{Name: "summary", Tools: []string{"summary"}})
	}
	return out
}

// compose builds the natural-language answer from tool results.
func (s *Service) compose(text string, results []ToolResult, intents []Intent, actor ActorType, locale i18n.Locale) (string, string) {
	name := "summary"
	if len(intents) > 0 {
		name = intents[0].Name
	}
	if len(results) == 0 {
		return cat.Get(locale, "assistant_no_results"), name
	}
	var b strings.Builder
	if actor == ActorEmployee {
		b.WriteString(cat.Get(locale, "assistant_found"))
	} else {
		b.WriteString(cat.Get(locale, "assistant_overview"))
	}
	for i, r := range results {
		if i > 0 {
			b.WriteString(" ")
		}
		b.WriteString(r.Line)
	}
	return b.String(), name
}

func actorAllowed(t Tool, a ActorType) bool {
	for _, x := range t.Actors {
		if x == a {
			return true
		}
	}
	return false
}

func containsAny(s string, subs ...string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}

func shortTitle(text string) string {
	t := strings.TrimSpace(text)
	if len(t) > 60 {
		return t[:60] + "…"
	}
	return t
}
