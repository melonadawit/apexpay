package notify

// Notification preferences domain.

type Preference struct {
	EventType string `json:"event_type"`
	Email     bool   `json:"email"`
	SMS       bool   `json:"sms"`
	Push      bool   `json:"push"`
	InApp     bool   `json:"in_app"`
}

// The set of event types a merchant user can configure.
var EventTypes = []string{
	"bulk_payouts_approval", "pending_payout", "payout_failed",
	"payroll_run_pending_approval", "payroll_run_completed",
	"tax_payment_due", "compliance_alert", "bank_file_generated",
	"loan_emi_due", "leave_request_pending", "claim_pending",
	"escrow_held", "escrow_released", "current_account_opened",
	"corporate_card_transaction", "forex_rate_alert", "accounting_sync_failed",
}
