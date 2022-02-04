Vault Agent vAuth Helper
=======

The Vault Agent vAuth helper application for automatically generating HashiCorp Vault AppRole credential files using credentials fetching from vAuth.

# Getting Started

1. Download the latest release of the vault-agent-vauth-helper from the Github releases.

2. Create a `config.hcl` configuration file in the directory that the `vault-agent-vauth-helper` 

```
sync_interval       = 10
role_id_file_path   = /tmp/role.txt
secret_id_file_path = /tmp/secret.txt
```

|Setting|Description|
|--|--|
|sync_interval||
|role_id_file_path||
|secret_id_file_path||



3. Create a `vauth-agent`

