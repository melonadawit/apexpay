package team

// Team (multi-user roles) domain types.

// Member is a merchant team member with their role and permissions.
type Member struct {
	UserID      string   `json:"user_id"`
	Email       string   `json:"email"`
	Name        string   `json:"name"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
	CreatedAt   string   `json:"created_at"`
}

// InviteRequest is the payload to add a member.
type InviteRequest struct {
	Email       string   `json:"email"`
	Name        string   `json:"name"`
	Role        string   `json:"role"`
	Permissions []string `json:"permissions"`
}

// Approval is a generic approval-request record (maker-checker capable).
type Approval struct {
	ID                string         `json:"id"`
	ResourceType      string         `json:"resource_type"`
	ResourceID        string         `json:"resource_id"`
	Action            string         `json:"action"`
	Summary           string         `json:"summary"`
	Amount            string         `json:"amount"`
	Currency          string         `json:"currency"`
	Status            string         `json:"status"`
	RequiredApprovals int            `json:"required_approvals"`
	Approvals         []ApprovalVote `json:"approvals"`
	CreatedAt         string         `json:"created_at"`
}

type ApprovalVote struct {
	UserID    string `json:"user_id"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	Decision  string `json:"decision"`
	DecidedAt string `json:"decided_at"`
}
