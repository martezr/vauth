package utils

type Config struct {
	UIPort            string `yaml:"ui_port" mapstructure:"UI_PORT"`
	VsphereURL        string `yaml:"vsphere_url" mapstructure:"VSPHERE_URL"`
	VsphereUsername   string `yaml:"vsphere_username" mapstructure:"VSPHERE_USERNAME"`
	VspherePassword   string `yaml:"vsphere_password" mapstructure:"VSPHERE_PASSWORD"`
	VaultAddress      string `yaml:"vault_address" mapstructure:"VAULT_ADDRESS"`
	VaultToken        string `yaml:"vault_token" mapstructure:"VAULT_TOKEN"`
	VaultAppRoleMount string `yaml:"vault_approle_mount" mapstructure:"VAULT_APPROLE_MOUNT"`
	TLSSkipVerify     bool   `yaml:"tls_skip_verify,optional" mapstructure:"TLS_SKIP_VERIFY"`
	WrapResponse      bool   `yaml:"wrap_response,optional" mapstructure:"WRAP_RESPONSE"`
}
