package risk

import (
	"testing"

	"github.com/shopspring/decimal"
)

func rule(id, rtype string, params RuleParams) Rule {
	return Rule{ID: id, Name: id, RuleType: rtype, Parameters: params, Action: "flag", Severity: "medium", Enabled: true}
}

func TestEngine_ThresholdAmount(t *testing.T) {
	e := NewEngine()
	rules := []Rule{rule("r1", "threshold_amount", RuleParams{AmountLimit: decimal.NewFromInt(1000)})}
	ctx := EvaluationContext{Amount: decimal.NewFromInt(5000)}
	eval := e.Evaluate(ctx, rules)
	if !eval.HasBlock() && len(eval.Matched()) != 1 {
		t.Fatalf("expected 1 match, got %d", len(eval.Matched()))
	}
	if !eval.Matched()[0].Matched {
		t.Fatal("threshold rule should match")
	}
}

func TestEngine_VelocityAmount(t *testing.T) {
	e := NewEngine()
	rules := []Rule{rule("r2", "velocity_amount", RuleParams{WindowMinutes: 60, AmountLimit: decimal.NewFromInt(100000)})}
	ctx := EvaluationContext{AmountInWindow: decimal.NewFromInt(250000)}
	eval := e.Evaluate(ctx, rules)
	if len(eval.Matched()) != 1 {
		t.Fatalf("velocity amount should match, got %d", len(eval.Matched()))
	}
	// Under the limit -> no match.
	ctx2 := EvaluationContext{AmountInWindow: decimal.NewFromInt(1000)}
	if len(e.Evaluate(ctx2, rules).Matched()) != 0 {
		t.Fatal("under-limit should not match")
	}
}

func TestEngine_VelocityCount(t *testing.T) {
	e := NewEngine()
	rules := []Rule{rule("r3", "velocity_count", RuleParams{WindowMinutes: 60, CountLimit: 10})}
	if len(e.Evaluate(EvaluationContext{CountInWindow: 15}, rules).Matched()) != 1 {
		t.Fatal("count over limit should match")
	}
	if len(e.Evaluate(EvaluationContext{CountInWindow: 5}, rules).Matched()) != 0 {
		t.Fatal("count under limit should not match")
	}
}

func TestEngine_HighFailureRate(t *testing.T) {
	e := NewEngine()
	rules := []Rule{rule("r4", "high_failure_rate", RuleParams{WindowMinutes: 60, RatioLimit: decimal.NewFromFloat(0.5)})}
	if len(e.Evaluate(EvaluationContext{FailureRate: decimal.NewFromFloat(0.8)}, rules).Matched()) != 1 {
		t.Fatal("high failure rate should match")
	}
	if len(e.Evaluate(EvaluationContext{FailureRate: decimal.NewFromFloat(0.1)}, rules).Matched()) != 0 {
		t.Fatal("low failure rate should not match")
	}
}

func TestEngine_DisabledRuleIgnored(t *testing.T) {
	e := NewEngine()
	r := rule("r5", "threshold_amount", RuleParams{AmountLimit: decimal.NewFromInt(1)})
	r.Enabled = false
	eval := e.Evaluate(EvaluationContext{Amount: decimal.NewFromInt(10000)}, []Rule{r})
	if len(eval.Matched()) != 0 {
		t.Fatal("disabled rule must be ignored")
	}
}

func TestEngine_BlockAction(t *testing.T) {
	e := NewEngine()
	r := rule("r6", "threshold_amount", RuleParams{AmountLimit: decimal.NewFromInt(1000)})
	r.Action = "block"
	eval := e.Evaluate(EvaluationContext{Amount: decimal.NewFromInt(50000)}, []Rule{r})
	if !eval.HasBlock() {
		t.Fatal("block action should make HasBlock true")
	}
}
