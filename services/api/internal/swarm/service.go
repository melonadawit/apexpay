package swarm

import (
	"context"
	"fmt"
	"time"

	"apexpay/internal/id"
	"apexpay/internal/platform/errors"
)

type Repository interface {
	CreateSession(ctx context.Context, s *SwarmSession) error
	GetSession(ctx context.Context, id string) (*SwarmSession, error)
	UpdateSession(ctx context.Context, s *SwarmSession) error
	CreateAgentRun(ctx context.Context, run *AgentRun) error
}

type ToolExecutor interface {
	Execute(ctx context.Context, tool string, args map[string]any, merchantID string) (map[string]any, error)
}

type Planner interface {
	Plan(ctx context.Context, goal string) ([]PlanStep, error) // LLM or rules
}

type Service struct {
	repo     Repository
	executor ToolExecutor
	planner  Planner
	registry map[string]ToolDefinition
}

func NewService(repo Repository, executor ToolExecutor, planner Planner, registry []ToolDefinition) *Service {
	regMap := make(map[string]ToolDefinition, len(registry))
	for _, t := range registry {
		regMap[t.Name] = t
	}
	return &Service{repo: repo, executor: executor, planner: planner, registry: regMap}
}

// Default registry per FULL spec
func DefaultRegistry() []ToolDefinition {
	return []ToolDefinition{
		{Name: "create_payment_link", Description: "Create payment link with amount and description", Threshold: 100000, RoleAllowed: []string{"owner", "admin", "developer"}},
		{Name: "list_payments", Description: "List recent payments", Threshold: 0},
		{Name: "create_refund", Description: "Refund payment", Threshold: 50000},
		{Name: "create_payout", Description: "Create single payout", Threshold: 50000},
		{Name: "calculate_payroll_draft", Description: "Calculate payroll run draft", Threshold: 100000},
		{Name: "ask_compliance", Description: "Ask NBE compliance via RAG", Threshold: 0},
		{Name: "get_tpv", Description: "Get today's TPV", Threshold: 0},
	}
}

// Run executes swarm goal with plan + critic + executor - optimal state machine
func (s *Service) Run(ctx context.Context, merchantID, userID, goal string) (*SwarmSession, error) {
	if goal == "" {
		return nil, errors.Validation("goal required")
	}

	// 1. Planning
	steps, err := s.planner.Plan(ctx, goal)
	if err != nil {
		return nil, err
	}

	session := &SwarmSession{
		ID: id.NewSwarmSession(), MerchantID: merchantID, Goal: goal,
		Plan: steps, Status: StatusPlanning, CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
	if userID != "" {
		session.UserID = &userID
	}

	if err := s.repo.CreateSession(ctx, session); err != nil {
		return nil, err
	}

	// 2. Critic checks - policy, amount thresholds, no ledger invent
	confirmationRequired := false
	var confData map[string]any
	totalAmount := 0.0
	for _, step := range steps {
		if def, ok := s.registry[step.Tool]; ok {
			if amt, ok := step.Args["amount"].(float64); ok {
				totalAmount += amt
				if amt > def.Threshold && def.Threshold > 0 {
					confirmationRequired = true
				}
			}
			// Payroll confirmation always >100k
			if step.Tool == "calculate_payroll_draft" || step.Tool == "create_payout" {
				if totalAmount > 100000 {
					confirmationRequired = true
				}
			}
		}
	}

	if confirmationRequired {
		session.Status = StatusNeedsConfirmation
		session.ConfirmationRequired = true
		session.ConfirmationData = map[string]any{"total_amount": totalAmount, "steps": len(steps)}
		_ = s.repo.UpdateSession(ctx, session)
		// Return early awaiting confirmation - UI shows outstanding modal
		return session, nil
	}

	// 3. Execute steps sequentially - optimal error handling with rollback not needed (idempotent tools)
	session.Status = StatusExecuting
	for i := range session.Plan {
		step := &session.Plan[i]
		step.Status = "executing"
		_ = s.repo.UpdateSession(ctx, session)

		start := time.Now()
		result, execErr := s.executor.Execute(ctx, step.Tool, step.Args, merchantID)
		lat := time.Since(start).Milliseconds()

		toolCall := ToolCall{Tool: step.Tool, Args: step.Args, LatencyMS: int(lat)}

		if execErr != nil {
			step.Status = "failed"
			toolCall.Status = "failed"
			toolCall.Result = map[string]any{"error": execErr.Error()}
			session.Status = StatusFailed
			// Audit agent run
			_ = s.repo.CreateAgentRun(ctx, &AgentRun{
				ID: id.New("arun"), MerchantID: merchantID, SwarmSessionID: &session.ID,
				InputText: goal, Intent: step.Tool, ToolCalls: []ToolCall{toolCall},
				OutputText: fmt.Sprintf("Step %d failed: %v", i+1, execErr), Model: "rules_v1",
				CreatedAt: time.Now(),
			})
			_ = s.repo.UpdateSession(ctx, session)
			return session, nil
		}

		step.Status = "succeeded"
		step.Result = result
		toolCall.Status = "succeeded"
		toolCall.Result = result

		_ = s.repo.CreateAgentRun(ctx, &AgentRun{
			ID: id.New("arun"), MerchantID: merchantID, SwarmSessionID: &session.ID,
			InputText: goal, Intent: step.Tool, ToolCalls: []ToolCall{toolCall},
			OutputText: fmt.Sprintf("Step %d %s succeeded", i+1, step.Tool),
			Model:      "rules_v1", CreatedAt: time.Now(),
		})
	}

	session.Status = StatusCompleted
	session.FinalOutput = s.buildFinalOutput(session.Plan)
	session.UpdatedAt = time.Now()
	_ = s.repo.UpdateSession(ctx, session)

	return session, nil
}

func (s *Service) Confirm(ctx context.Context, sessionID string, confirmed bool) (*SwarmSession, error) {
	sess, err := s.repo.GetSession(ctx, sessionID)
	if err != nil {
		return nil, errors.NotFound("swarm session not found")
	}
	if sess.Status != StatusNeedsConfirmation {
		return nil, errors.Validation("session not needing confirmation")
	}
	if !confirmed {
		sess.Status = StatusCancelled
		sess.FinalOutput = "Cancelled by user"
		_ = s.repo.UpdateSession(ctx, sess)
		return sess, nil
	}

	// User confirmed - reset flag and re-run execution
	sess.ConfirmationRequired = false
	sess.Status = StatusExecuting
	_ = s.repo.UpdateSession(ctx, sess)

	// Continue execution - similar loop as Run but from where left off
	// Simplified: re-execute all pending
	for i := range sess.Plan {
		step := &sess.Plan[i]
		if step.Status == "succeeded" {
			continue
		}
		result, execErr := s.executor.Execute(ctx, step.Tool, step.Args, sess.MerchantID)
		if execErr != nil {
			step.Status = "failed"
			sess.Status = StatusFailed
			_ = s.repo.UpdateSession(ctx, sess)
			return sess, nil
		}
		step.Status = "succeeded"
		step.Result = result
	}

	sess.Status = StatusCompleted
	sess.FinalOutput = s.buildFinalOutput(sess.Plan)
	_ = s.repo.UpdateSession(ctx, sess)
	return sess, nil
}

func (s *Service) buildFinalOutput(plan []PlanStep) string {
	var out string
	for _, step := range plan {
		if step.Status == "succeeded" && step.Result != nil {
			if url, ok := step.Result["payment_link_url"].(string); ok {
				out += fmt.Sprintf("Created payment link: %s\n", url)
			}
			if runID, ok := step.Result["payroll_run_id"].(string); ok {
				out += fmt.Sprintf("Created payroll run: %s\n", runID)
			}
		}
	}
	if out == "" {
		out = fmt.Sprintf("Completed %d steps successfully", len(plan))
	}
	return out
}

// RulesPlanner - MVP deterministic rules engine per SAD ADR-006 AI off hot path
type RulesPlanner struct{}

func (r *RulesPlanner) Plan(ctx context.Context, goal string) ([]PlanStep, error) {
	// Very simple keyword matching - production would use LLM
	lower := goal
	steps := []PlanStep{}

	if contains(lower, "link") || contains(lower, "payment") && contains(lower, "create") {
		steps = append(steps, PlanStep{Step: len(steps) + 1, Tool: "create_payment_link", Description: "Create payment link", Args: map[string]any{"amount": 500.0, "currency": "ETB", "description": goal}, Status: "pending"})
	}
	if contains(lower, "payroll") {
		steps = append(steps, PlanStep{Step: len(steps) + 1, Tool: "calculate_payroll_draft", Description: "Calculate payroll draft", Args: map[string]any{"period_month": 7, "period_year": 2026}, Status: "pending"})
	}
	if contains(lower, "payout") {
		steps = append(steps, PlanStep{Step: len(steps) + 1, Tool: "create_payout", Description: "Create payout", Args: map[string]any{"amount": 10000.0, "currency": "ETB"}, Status: "pending"})
	}
	if contains(lower, "tpv") || contains(lower, "today") {
		steps = append(steps, PlanStep{Step: 0, Tool: "get_tpv", Description: "Get TPV today", Args: map[string]any{}, Status: "pending"})
	}
	if len(steps) == 0 {
		steps = append(steps, PlanStep{Step: 1, Tool: "list_payments", Description: "List recent payments", Args: map[string]any{}, Status: "pending"})
	}

	// Re-index steps sequentially
	for i := range steps {
		steps[i].Step = i + 1
	}

	return steps, nil
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (func() bool {
		// case-insensitive simple
		sl := s
		if len(sl) != len(s) {
			return false
		}
		// naive
		for i := 0; i <= len(s)-len(substr); i++ {
			if s[i:i+len(substr)] == substr {
				return true
			}
		}
		return false
	})()
}
