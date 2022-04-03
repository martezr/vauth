package utils

// LogEntry defines the audit log format
type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Operation string `json:"operation"`
}
