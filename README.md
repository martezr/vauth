# vAuth

[![GoReportCard][report-badge]][report]
[![GitHub release](https://img.shields.io/github/release/martezr/vauth.svg)](https://github.com/martezr/vauth/releases/)
[![license](https://img.shields.io/github/license/martezr/vauth.svg)](https://github.com/martezr/vauth/blob/master/LICENSE)

[report-badge]: https://goreportcard.com/badge/github.com/martezr/vauth
[report]: https://goreportcard.com/report/github.com/martezr/vauth

VMware vSphere VM Identity Platform

The vAuth platform provides identity information to virtual machines similiar to the metadata provided by public cloud providers. The platform is built to work with [HashiCorp Vault](https://www.vaultproject.io/) to enable VMware vSphere to be used as a trusted platform similar to public cloud providers such as AWS and Azure.

## HashiCorp Vault Minimum Permissions

The vAuth platform requires the following minimum permissions.

* List all authentication methods
* Read and list all roles in the approle backend

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
docker run --name vauth -e VSPHERE_SERVER=grtvcenter01.grt.local -e VSPHERE_USERNAME=vauth@vsphere.local -e VSPHERE_PASSWORD="Password123#" -e VAULT_ADDRESS="https://grtmanage01.grt.local:8200" -e VAULT_TOKEN="s.ewdkchV1oqIwTxxI8G3INWVG" -e VAULT_APPROLE_MOUNT=approle public.ecr.aws/i4r5n0t9/vauth:1.0
```
