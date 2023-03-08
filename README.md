# vAuth

[![GoReportCard][report-badge]][report]
[![GitHub release](https://img.shields.io/github/release/martezr/vauth.svg)](https://github.com/martezr/vauth/releases/)
[![license](https://img.shields.io/github/license/martezr/vauth.svg)](https://github.com/martezr/vauth/blob/master/LICENSE)

[report-badge]: https://goreportcard.com/badge/github.com/martezr/vauth
[report]: https://goreportcard.com/report/github.com/martezr/vauth

VMware vSphere VM Identity Platform

The vAuth platform provides identity information to virtual machines similiar to the metadata provided by public cloud providers. The platform is built to work with [HashiCorp Vault](https://www.vaultproject.io/) to enable VMware vSphere to be used as a trusted platform similar to public cloud providers such as AWS and Azure.

## How vAuth Works

The following steps provide a high level overview of how the vAuth platform works and interacts with vSphere and HashiCorp Vault:

1. The vAuth platform listens for virtual machine power on and virtual machine custom attribute change events. 
2. When one of these events are detected the platform looks up the role associated with the virtual machine. The role is defined via the `vauth-role` custom attribute.
3. The vAuth platform then queries the configured HashiCorp Vault instance from the approle backend configured.
4. If a matching role is found then a role ID and secret ID are generated for that virtual machine and set in the virtual machine's VMware guest information.
5. Once the credentials have been set, the virtual machine guest operating system is able to query the credentials.

## HashiCorp Vault Minimum Permissions

The vAuth platform requires the following minimum permissions to integrate with HashiCorp Vault.

* List all authentication methods
* Read and list all roles in the configured approle backend

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

|Setting Name|Description|Type|Example|
|------------|-----------|---|----|
|ui_port     | The port on which the vAuth web UI will listen|string|8000|
|data_dir    | The path on the filesystem that will be used to store vAuth | string | /vauthdata|
| vsphere_server | The FQDN or IP address of the vCenter server that vAuth will connect to | string | vcenter.domain.local |
|vsphere_tls_skip_verify | Whether to skip the verification of the vCenter SSL certificate or not | boolean | false |
| vsphere_username | The username of the user account that vAuth will use to connect to vCenter | string | vauth@vsphere.local |
| vsphere_password | The password of the user account that vAuth will use to connect to vCenter | string | securepassword |
| vsphere_datacenters | The vSphere datacenters to enable authentication on | []string | ["DC1","DC2] |
| vault_address | The URL of the HashiCorp Vault instance that vAuth will connect to | string | https://demo.domain.local:8200 |
| vault_token | The vault token that used by vAuth to authenticate to HashiCorp Vault | string | vaultpassword|
| vault_approle_mount | The name of the approle authentication backend used by vAuth to generate new approle role credentials | string | approle |
| vault_wrap_reponse | Whether to wrap the response for the secret ID | boolean | true |
| vault_tls_skip_verify | Whether to skip the verification of the Vault SSL certificate or not | boolean | false |


### Binary Installation

The vAuth platform can be deployed using the vAuth binary on linux systems.

1. Download the vAuth binary from the latest Github release

```bash
export VAUTH_VERSION="0.0.2"
```

```bash
curl --silent --remote-name \
  https://github.com/martezr/vauth/releases/download/v${VAUTH_VERSION}/vauth_${VAUTH_VERSION}_linux_amd64.zip
```

2. Make the vAuth binary executable

```
chmod +x vauth
```

3. Start the vAuth service

```
vauth server
```

### Docker Installation

The vAuth platform can be deployed with Docker using the following command:

```
docker run --name vauth -e UI_PORT=9000 -e DATA_DIR=/app -e VSPHERE_SERVER=grtvcenter01.grt.local -e VSPHERE_USERNAME=vauth@vsphere.local -e VSPHERE_PASSWORD="Password123#" -e VAULT_ADDRESS="https://grtmanage01.grt.local:8200" -e VAULT_TOKEN="s.ewdkchV1oqIwTxxI8G3INWVG" -e VAULT_APPROLE_MOUNT=approle public.ecr.aws/i4r5n0t9/vauth:1.0
```

### Kubernetes Installation

The vAuth platform can be deployed to a Kubernetes cluster using the following manifest:

```yaml
apiVersion: v1
kind: Namespace
metadata:
  name: vauth
---
apiVersion: v1
kind: ConfigMap
metadata:
  name: vauth-config
  namespace: vauth
data:
  VSPHERE_SERVER: "grtvcenter01.grt.local"
  DATA_DIR: "/vauthdata"
  VSPHERE_TLS_SKIP_VERIFY: "true"
  VSPHERE_USERNAME: "vauth@vsphere.local"
  VSPHERE_PASSWORD: "Password123#"
  VSPHERE_DATACENTERS: ["DC1","DC2"]
  VAULT_ADDRESS: "https://10.0.0.202:8200"
  VAULT_TOKEN: "s.r5A9FBMiQyRzXcEh7Ab7ZE4K"
  VAULT_APPROLE_MOUNT: "approle"
  VAULT_WRAP_RESPONSE: "true"
  VAULT_TLS_SKIP_VERIFY: "true"
---
apiVersion: v1
kind: Pod
metadata:
  name: vauth
  namespace: vauth
  labels:
    name: vauth
spec:
  containers:
  - name: vauth
    image: public.ecr.aws/i4r5n0t9/vauth:1.0
    imagePullPolicy: Always
    envFrom:
      - configMapRef:
          name: vauth-config
    volumeMounts:
    - mountPath: /vauthdata
      name: cache-volume
  volumes:
  - name: cache-volume
    emptyDir: {}
```