package client

// Board represents a flashbay board as returned by the coordinator API.
type Board struct {
	ID               string  `json:"id"`
	PodID            *string `json:"pod_id,omitempty"`
	BoardType        string  `json:"board_type"`
	Device           string  `json:"device,omitempty"`
	Label            string  `json:"label"`
	State            string  `json:"state"`
	Disabled         bool    `json:"disabled"`
	SledVersion      string  `json:"sled_version"`
	CurrentSessionID *string `json:"current_session_id,omitempty"`
	CreatedAt        string  `json:"created_at"`
}

// Session represents a flashbay session as returned by the coordinator API.
type Session struct {
	ID          string  `json:"id"`
	UserID      string  `json:"user_id"`
	BoardID     *string `json:"board_id,omitempty"`
	BoardType   string  `json:"board_type"`
	State       string  `json:"state"`
	EndReason   *string `json:"end_reason,omitempty"`
	CreditsUsed int     `json:"credits_used"`
	StartedAt   *string `json:"started_at,omitempty"`
	EndedAt     *string `json:"ended_at,omitempty"`
	CreatedAt   string  `json:"created_at"`
}

// User represents the authenticated user profile.
type User struct {
	ID             string `json:"id"`
	Email          string `json:"email"`
	BalanceCredits int    `json:"balance_credits"`
	IsActive       bool   `json:"is_active"`
	IsAdmin        bool   `json:"is_admin"`
	CreatedAt      string `json:"created_at"`
	LastActiveAt   string `json:"last_active_at"`
}
