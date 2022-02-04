# vAuth

[![GoReportCard][report-badge]][report]
[![GitHub release](https://img.shields.io/github/release/martezr/vauth.svg)](https://github.com/martezr/vauth/releases/)
[![license](https://img.shields.io/github/license/martezr/vauth.svg)](https://github.com/martezr/vauth/blob/master/LICENSE)

[report-badge]: https://goreportcard.com/badge/github.com/martezr/vauth
[report]: https://goreportcard.com/report/github.com/martezr/vauth

VMware vSphere VM Identity Platform

The vAuth platform provides identity information to virtual machines similiar to the metadata provided by public cloud providers. The platform is built to work with [HashiCorp Vault](https://www.vaultproject.io/) to enable VMware vSphere to be used as a trusted platform similar to public cloud providers such as AWS and Azure.


## How vAuth Works

The vAuth platform queries the virtual machine

vAuth generates a new secret ID for the Vault approle role

## HashiCorp Vault Minimum Permissions

The vAuth platform requires the following minimum permissions to integrate with HashiCorp Vault.

* List all authentication methods
* Read and list all roles in the approle backend

The following is an example least privilege Vault policy. The policy assumes that `approle` is the name of the approle authentication method backend/path. 

```
# List auth methods
path "sys/auth" {
  capabilities = ["read"]
}

# List roles
path "auth/approle/*" {
  capabilities = [ "read", "list" ]
}

# Read the role IDs for all roles in the approle auth backend
path "auth/approle/role/+/role-id" {
   capabilities = [ "read" ]
}

# Generate secret IDs for all roles in the approle auth backend
path "auth/approle/role/+/secret-id" {
  capabilities = [ "update" ]  
}
```

## vSphere Account Permissions

The vAuth platform requires access to VMware vSphere to perform various operations such as watch events, update virtual machine guest information and more. The following table details the permissions that the account used by the vAuth platform would need in vSphere.

**Privileged interaction**
The following operations require a privilege to be assigned to the vSphere account that the vAuth platform uses.

|Permission|Description|
|----|-----------|
|Virtual Machine > Change Configuration > Advanced Configuration | The account needs to have permission to update the advanced configuration of virtual machines to provide the identity data to the guest operating system|

## Setup


```
---
ui_port: 8000
data_dir: .
vsphere_server: "localhost"
vsphere_username: "user"
vsphere_password: "pass"
vault_address: "http://localhost:8200"
vault_token: "vault"
vault_approle_mount: "approle"
wrap_response: true
tls_skip_verify: true
```

|Setting Name|Description|Type|Example|
|------------|-----------|---|----|
|ui_port     | The port on which the vAuth web UI will listen|string|8000|
|data_dir    | The path on the filesystem that will be used to store vAuth | string | /vauthdata|
| vsphere_server | The FQDN or IP address of the vCenter server that vAuth will connect to | string | vcenter.domain.local |
|vsphere_tls_skip_verify | Whether to skip the verification of the vCenter SSL certificate or not | boolean | false |
| vsphere_username | The username of the user account that vAuth will use to connect to vCenter | string | vauth@vsphere.local |
| vsphere_password | The password of the user account that vAuth will use to connect to vCenter | string | securepassword |
| vault_address | The URL of the HashiCorp Vault instance that vAuth will connect to | string | https://demo.domain.local:8200 |
| vault_token | The vault token that used by vAuth to authenticate to HashiCorp Vault | string | vaultpassword|
| vault_approle_mount | The name of the approle authentication backend used by vAuth to generate new approle role credentials | string | approle |
| vault_wrap_reponse | | boolean | true |
| vault_tls_skip_verify | Whether to | boolean | false |


### Binary Installation

The following

1. Download the vAuth binary from Github

2. Make the vAuth binary executable

```
chmod +x vauth
```

3. 

```
vauth server
```
### Docker Installation


```
docker run --name vauth -e UI_PORT=9000 -e DATA_DIR=/app -e VSPHERE_SERVER=grtvcenter01.grt.local -e VSPHERE_USERNAME=vauth@vsphere.local -e VSPHERE_PASSWORD="Password123#" -e VAULT_ADDRESS="https://grtmanage01.grt.local:8200" -e VAULT_TOKEN="s.ewdkchV1oqIwTxxI8G3INWVG" -e VAULT_APPROLE_MOUNT=approle public.ecr.aws/i4r5n0t9/vauth:1.0
```

### Kubernetes Installation

The vAuth platform can be deployed on a Kubernetes cluster with the following steps