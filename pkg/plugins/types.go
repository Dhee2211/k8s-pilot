package plugins

// Issue represents a problem detected by a plugin
type Issue struct {
	ID          string                 `json:"id,omitempty"`
	Severity    string                 `json:"severity"` // "critical", "high", "medium", "low"
	Title       string                 `json:"title"`
	Description string                 `json:"description"`
	Resource    string                 `json:"resource"`
	Namespace   string                 `json:"namespace,omitempty"`
	Type        string                 `json:"type,omitempty"`
	Details     map[string]interface{} `json:"details,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	Remediation []RemediationStep      `json:"remediation,omitempty"`
}

// RemediationStep represents a step to fix an issue
type RemediationStep struct {
	Description string `json:"description"`
	Command     string `json:"command,omitempty"`
	Action      string `json:"action,omitempty"`
	Safe        bool   `json:"safe"`
	Confidence  string `json:"confidence,omitempty"` // "high", "medium", "low"
}

// Plugin interface that all plugins must implement
type Plugin interface {
	Name() string
	Version() string
	Description() string
	Analyze(resource interface{}) ([]Issue, error)
}
