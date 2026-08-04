package swarm

import (
	"context"
	"encoding/json"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgRepository struct{ pool *pgxpool.Pool }

func NewPgRepository(pool *pgxpool.Pool) *PgRepository { return &PgRepository{pool: pool} }

func (r *PgRepository) CreateSession(ctx context.Context, s *SwarmSession) error {
	planBytes, _ := json.Marshal(s.Plan)
	_, err := r.pool.Exec(ctx, `INSERT INTO swarm_sessions (id, merchant_id, user_id, goal, plan, status, confirmation_required, confirmation_data) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`,
		s.ID, s.MerchantID, s.UserID, s.Goal, planBytes, s.Status, s.ConfirmationRequired, s.ConfirmationData)
	return err
}

func (r *PgRepository) GetSession(ctx context.Context, id string) (*SwarmSession, error) {
	row := r.pool.QueryRow(ctx, `SELECT id, merchant_id, user_id, goal, plan, status, confirmation_required, confirmation_data, final_output FROM swarm_sessions WHERE id=$1`, id)
	var s SwarmSession
	var planBytes []byte
	err := row.Scan(&s.ID, &s.MerchantID, &s.UserID, &s.Goal, &planBytes, &s.Status, &s.ConfirmationRequired, &s.ConfirmationData, &s.FinalOutput)
	if err != nil { return nil, err }
	_ = json.Unmarshal(planBytes, &s.Plan)
	return &s, nil
}

func (r *PgRepository) UpdateSession(ctx context.Context, s *SwarmSession) error {
	planBytes, _ := json.Marshal(s.Plan)
	_, err := r.pool.Exec(ctx, `UPDATE swarm_sessions SET plan=$1, status=$2, confirmation_required=$3, confirmation_data=$4, final_output=$5, updated_at=now() WHERE id=$6`,
		planBytes, s.Status, s.ConfirmationRequired, s.ConfirmationData, s.FinalOutput, s.ID)
	return err
}

func (r *PgRepository) CreateAgentRun(ctx context.Context, run *AgentRun) error {
	toolCallsBytes, _ := json.Marshal(run.ToolCalls)
	_, err := r.pool.Exec(ctx, `INSERT INTO agent_runs (id, merchant_id, thread_id, swarm_session_id, input_text, intent, tool_calls, output_text, model) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`,
		run.ID, run.MerchantID, run.ThreadID, run.SwarmSessionID, run.InputText, run.Intent, toolCallsBytes, run.OutputText, run.Model)
	return err
}
