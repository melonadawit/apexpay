# Swarm Multi-Agent — Planner/Critic/Executor + Confirmation >100k + Audit

Per FULL_PLATFORM_SPEC + swarm/domain.go + service.go

## Architecture

```
Goal: "Create link 100 ETB for coffee and run payroll July with bonus"
  → Planner RulesPlanner keyword (if "link" → create_payment_link, if "payroll" → calculate_payroll_draft, if "payout" → create_payout, if "tpv" → get_tpv) + LLM planner OpenAI function calling prod
  → Critic checks policy/hallucination/amount thresholds >100k + no ledger invent blocked + totalAmount >100k => confirmation_required true confirmation_data total_amount steps
  → If confirmation_required true → status needs_confirmation outstanding modal breakdown + biometric → return early awaiting confirmation UI shows outstanding modal
  → Else Executors call Go services via ToolExecutorImpl function map O(1) paymentLinkCreator id New pl_01H url https://checkout... amount currency, payoutCreator payout_id NewPayout pending_approval, payrollCalculator payroll_run_id NewPayrollRun period_month 7 period_year 2026 total_net 150k pending_approval, tpvGetter tpv_today 125430, complianceAsker answer 5000 ETB per ONPS/10/2025 citation
  → Sequentially executing steps pending→executing→succeeded/failed + toolCalls audit latency + agent_runs swarm_session_id input_text intent tool_calls output_text model rules_v1 + final_output Created link + payroll run
  → Confirm endpoint if confirmed true → reset flag + re-execute all pending steps → completed final_output summary with links
```

## Domain

- SwarmSession ID merchant_id user_id goal plan []PlanStep status planning/executing/needs_confirmation/completed/failed/cancelled confirmation_required bool confirmation_data map final_output created_at updated_at
- PlanStep step tool create_payment_link/create_payout/calculate_payroll_draft/ask_compliance/get_tpv description args map status pending/executing/succeeded/failed/needs_confirmation result map
- AgentRun ID merchant_id thread_id swarm_session_id input_text intent tool_calls []ToolCall output_text model rules_v1/gpt-4 created_at
- ToolCall tool args result latency_ms status succeeded/failed
- ToolDefinition name description argsSchema JSON schema threshold amount requiring confirmation roleAllowed owner/admin/developer/finance
- DefaultRegistry 7 tools: create_payment_link threshold 100000 role owner/admin/developer, list_payments 0, create_refund 50000, create_payout 50000, calculate_payroll_draft 100000, ask_compliance 0, get_tpv 0

## Service — Optimal State Machine

- NewService repo executor planner registry map O(1) lookup
- DefaultRegistry
- Run merchantID userID goal → Planning steps via planner.Plan goal → CreateSession id NewSwarmSession goal plan status planning → Critic totalAmount loop steps args amount float64 totalAmount+=amt if amt>def.Threshold && Threshold>0 confirmationRequired true + payroll/payout total>100k confirmationRequired true → If confirmationRequired status needs_confirmation confirmation_data total_amount steps UpdateSession return early awaiting confirmation UI outstanding modal breakdown + biometric → Else status executing loop Plan i step status executing UpdateSession start now executor.Execute tool args merchantID latencyMS ToolCall tool args latency status succeeded/failed result + CreateAgentRun id New arun merchantID swarmSessionID inputText goal intent tool toolCalls outputText Step i+1 tool succeeded model rules_v1 + step status succeeded result + Final status completed finalOutput buildFinalOutput plan steps result payment_link_url payroll_run_id + UpdateSession
- Confirm sessionID confirmed bool → GetSession status needs_confirmation check else validation error not needing confirmation → If !confirmed status cancelled finalOutput Cancelled by user UpdateSession return → Else confirmationRequired false status executing UpdateSession Continue execution re-execute all pending steps status pending → result + status succeeded + finalOutput + UpdateSession completed
- buildFinalOutput plan steps result payment_link_url payroll_run_id + Completed steps successfully
- RulesPlanner Plan goal lower keyword matching if contains link/payment+create → create_payment_link amount 500 currency ETB description goal status pending, if payroll → calculate_payroll_draft period_month 7 period_year 2026 status pending, if payout → create_payout amount 10000 currency ETB status pending, if tpv/today → get_tpv, if len steps 0 → list_payments, Re-index steps sequentially step i+1

## Executor — ToolExecutorImpl Function Map O(1) Optimal

- ToolExecutorImpl struct paymentLinkCreator func ctx merchantID amount float64 currency desc → map id pl_01H url https://checkout... amount currency, payoutCreator payout_id amount status pending_approval, payrollCalculator payroll_run_id period_month 7 period_year total_net 150000 status pending_approval, tpvGetter tpv_today 125430 currency ETB count 42, complianceAsker answer 5000 ETB per ONPS/10/2025 citation
- Execute tool args merchantID start now switch tool create_payment_link amt float64/int handling currency default ETB desc → paymentLinkCreator, create_payout amt → payoutCreator, calculate_payroll_draft → payrollCalculator, get_tpv → tpvGetter, ask_compliance query → complianceAsker, list_payments → payments count 5, default unknown tool error, latency start

## Handler

- Routes POST /run goal → 201 sess needs_confirmation outstanding modal breakdown + biometric + plan steps tool call cards like Vercel AI SDK, POST /{id}/confirm confirmed bool → completed final_output Created link + payroll draft, GET /{id} → session plan status final_output

## UI Trace View Outstanding — Day 6 Merchant Guide

- Merchant command center chat right panel glassmorphic backdrop-blur-xl bg-white/70 border white/50 shadow-glass + swarm run shows stepper timeline with tool call cards like Vercel AI SDK — each card icon + description + args preview + status spinner check + result link, confirmation modal outstanding with breakdown total_amount steps + biometric option, final output summary with links
- Safety: no direct ledger_entries insert from LLM — must via domain service; ledger invent blocked by critic
- Metrics: swarm_run_total{status}, swarm_tool_call_duration_seconds{tool}
- Audit: agent_runs + swarm_sessions visible in admin exam + exam console reconstruct <60s per SAD A1

Next: [API Reference 21 Paths](/docs/api-reference)
