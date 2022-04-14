# Configuration Settings

|Setting Name|Description|Type|Example|
|------------|-----------|---|----|
|ui_port     | The port on which the vAuth web UI will listen|string|8000|
|data_dir    | The path on the filesystem that will be used to store vAuth data | string | /vauthdata|
| vsphere_server | The FQDN or IP address of the vCenter server that vAuth will connect to | string | vcenter.domain.local |
|vsphere_tls_skip_verify | Whether to skip the verification of the vCenter SSL certificate or not | boolean | false |
| vsphere_username | The username of the user account that vAuth will use to connect to vCenter | string | vauth@vsphere.local |
| vsphere_password | The password of the user account that vAuth will use to connect to vCenter | string | securepassword |
| vsphere_datacenters | The vSphere datacenters to enable authentication on | []string | ["DC1","DC2] |
| vault_address | The URL of the HashiCorp Vault instance that vAuth will connect to | string | https://demo.domain.local:8200 |
| vault_token | The vault token that used by vAuth to authenticate to HashiCorp Vault | string | vaultpassword|
| vault_approle_mount | The name of the approle authentication backend used by vAuth to generate new approle role credentials | string | approle |
| vault_wrap_reponse | Whether to wrap the secret id returned from vault | boolean | true |
| vault_tls_skip_verify | Whether to skip the verification of the vault SSL certificate or not	 | boolean | false |