package risk

import "fmt"

// Engine evaluates an EvaluationContext against a set of rules and returns an Evaluation.
// It is pure (no I/O) so it is fully unit-testable.
type Engine struct{}

func NewEngine() *Engine { return &Engine{} }

// Evaluate runs every enabled rule against the context.
func (e *Engine) Evaluate(ctx EvaluationContext, rules []Rule) Evaluation {
	eval := Evaluation{Findings: []Finding{}}
	for _, rule := range rules {
		if !rule.Enabled {
			continue
		}
		finding := e.evalRule(ctx, rule)
		eval.Findings = append(eval.Findings, finding)
	}
	return eval
}

// evalRule applies a single rule to the context.
func (e *Engine) evalRule(ctx EvaluationContext, rule Rule) Finding {
	base := Finding{
		RuleID: rule.ID, RuleName: rule.Name, RuleType: rule.RuleType,
		Severity: rule.Severity, Action: rule.Action,
	}
	switch rule.RuleType {
	case "velocity_amount":
		if ctx.AmountInWindow.GreaterThan(rule.Parameters.AmountLimit) {
			base.Matched = true
			base.Reason = fmt.Sprintf("amount %s in %dmin exceeds limit %s",
				ctx.AmountInWindow.String(), rule.Parameters.WindowMinutes, rule.Parameters.AmountLimit.String())
		}
	case "velocity_count":
		if ctx.CountInWindow >= rule.Parameters.CountLimit {
			base.Matched = true
			base.Reason = fmt.Sprintf("count %d in %dmin exceeds limit %d",
				ctx.CountInWindow, rule.Parameters.WindowMinutes, rule.Parameters.CountLimit)
		}
	case "threshold_amount":
		if ctx.Amount.GreaterThan(rule.Parameters.AmountLimit) {
			base.Matched = true
			base.Reason = fmt.Sprintf("single amount %s exceeds threshold %s", ctx.Amount.String(), rule.Parameters.AmountLimit.String())
		}
	case "high_ticket":
		if ctx.Amount.GreaterThan(rule.Parameters.AmountLimit) {
			base.Matched = true
			base.Reason = fmt.Sprintf("high-ticket amount %s", ctx.Amount.String())
		}
	case "high_failure_rate":
		if ctx.FailureRate.GreaterThan(rule.Parameters.RatioLimit) {
			base.Matched = true
			base.Reason = fmt.Sprintf("failure rate %s exceeds limit %s", ctx.FailureRate.StringFixed(2), rule.Parameters.RatioLimit.StringFixed(2))
		}
	case "new_device":
		if ctx.DeviceID != "" && rule.Parameters.WindowMinutes > 0 {
			// The repository reports count of prior distinct devices; if the entity is a
			// fresh device and we've seen <1, flag. Simplified: a new device is flagged when
			// ctx.CountInWindow == 0 (no prior history for this device).
			base.Matched = ctx.CountInWindow == 0
			if base.Matched {
				base.Reason = "payment from a previously unseen device"
			}
		}
	case "new_ip":
		if ctx.IP != "" && ctx.CountInWindow == 0 {
			base.Matched = true
			base.Reason = "payment from a previously unseen IP"
		}
	case "manual":
		// Manual rules never auto-match; they are for analyst-created review cases.
	default:
	}
	return base
}
