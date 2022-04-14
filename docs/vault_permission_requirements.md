# HashiCorp Vault Permission Requirements

The vAuth platform integrates with HashiCorp Vault to automate the process of injecting approle credentials to vSphere workloads. This means that the vAuth platform must authenticate to vault and have the appropriate permissions to perform the necessary operations. 

|Vault Path|Description|Operations|
|----------|-----------|----------|
|sys/auth|The vAuth platform reads the configured auth methods to evaluate if the configured approle backend has been configured|read|
|auth/approle/*||read, list|
|auth/approle/role/+/role-id||read|
|auth/approle/role/+/secret-id||update|


The following ex

```hcl
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