# VMware vSphere Permission Requirements

The vAuth platform requires access to VMware vSphere to perform various operations such as watch events, update virtual machine guest information and more. The following table details the permissions that the account used by the vAuth platform would need in vSphere.

**Privileged interaction**
The following operations require a privilege to be assigned to the vSphere account that the vAuth platform uses.

|Permission|Description|
|----|-----------|
|Virtual Machine > Change Configuration > Advanced Configuration | The account needs to have permission to update the advanced configuration of virtual machines to provide the identity data to the guest operating system|