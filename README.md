# vauth
VMware vSphere VM Identity Platform

# vAuth Architecture

![](./vauth-architecture.png)

|Component|Description|
|---------|-----------|
| Scheduler| Schedule the synchronization process |
| Syncer   | Synchronize VMs without identity data|
| Watcher | Watch for relevant VMware vSphere events such as power on operations |
| Worker | Generate the identity data (VM Name, Datacenter, VM Folder and Secret Key) and add that to the VM attributes|
| Backend | REST API for Vault auth validation and Database interaction |
| NATS | Event bus for system transactions |
| DB | Cockroachdb for persistent VM identity data presented by the backend in response to Vault auth validation requests |

