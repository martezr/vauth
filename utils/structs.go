package utils

// VMRecords
type VMRecords struct {
	Records []VMRecord `json:"vms"`
	Total   int        `json:"total"`
}

// VMRecord
type VMRecord struct {
	Name          string `json:"name"`
	LatestEventId string `json:"latest_event_id"`
	Role          string `json:"role"`
}

// Config is the
type Config struct {
	UIPort               string `yaml:"ui_port" mapstructure:"UI_PORT"`
	DataDir              string `yaml:"data_dir" mapstructure:"DATA_DIR"`
	VsphereServer        string `yaml:"vsphere_server" mapstructure:"VSPHERE_SERVER"`
	VsphereTLSSkipVerify bool   `yaml:"vsphere_tls_skip_verify,optional" mapstructure:"VSPHERE_TLS_SKIP_VERIFY"`
	VsphereUsername      string `yaml:"vsphere_username" mapstructure:"VSPHERE_USERNAME"`
	VspherePassword      string `yaml:"vsphere_password" mapstructure:"VSPHERE_PASSWORD"`
	VaultAddress         string `yaml:"vault_address" mapstructure:"VAULT_ADDRESS"`
	VaultToken           string `yaml:"vault_token" mapstructure:"VAULT_TOKEN"`
	VaultAppRoleMount    string `yaml:"vault_approle_mount" mapstructure:"VAULT_APPROLE_MOUNT"`
	VaultTLSSkipVerify   bool   `yaml:"vault_tls_skip_verify,optional" mapstructure:"VAULT_TLS_SKIP_VERIFY"`
	VaultWrapResponse    bool   `yaml:"vault_wrap_response,optional" mapstructure:"VAULT_WRAP_RESPONSE"`
}

type HealthStatus struct {
	Version       string `json:"version"`
	VaultStatus   string `json:"vault_status"`
	VsphereStatus string `json:"vsphere_status"`
}
