package utils

type Config struct {
	UIPort  string          `hcl:"ui_port"`
	Vsphere []VsphereConfig `hcl:"vsphere,block"`
	Vault   []VaultConfig   `hcl:"vault,block"`
}

type VsphereConfig struct {
	Name            string `hcl:"name,label"`
	VsphereURL      string `hcl:"vsphere_url"`
	VsphereUsername string `hcl:"vsphere_username"`
	VspherePassword string `hcl:"vsphere_password"`
}

type VaultConfig struct {
	Name              string `hcl:"name,label"`
	VaultAddress      string `hcl:"vault_address"`
	VaultToken        string `hcl:"vault_token"`
	VaultAppRoleMount string `hcl:"vault_approle_mount"`
	TLSSkipVerify     bool   `hcl:"tls_skip_verify,optional"`
	WrapResponse      bool   `hcl:"wrap_response,optional"`
}
