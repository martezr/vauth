# vauth

[![Build Status](https://img.shields.io/travis/martezr/packer-provisioner-puppet-bolt/master.svg)][travis]
[![GoReportCard][report-badge]][report]
[![GitHub release](https://img.shields.io/github/release/martezr/vauth.svg)](https://github.com/martezr/vauth/releases/)
[![license](https://img.shields.io/github/license/martezr/vauth.svg)](https://github.com/martezr/vauth/blob/master/LICENSE)

[travis]: https://travis-ci.org/martezr/vauth
[report-badge]: https://goreportcard.com/badge/github.com/martezr/vauth
[report]: https://goreportcard.com/report/github.com/martezr/vauth

VMware vSphere VM Identity Platform

The vAuth Identity platform works in conjunction with the [vSphere Vault Auth Plugin](https://github.com/martezr/vault-plugin-auth-vsphere). The vAuth platform provides identity information to virtual machines similiar to the metadata provided by public cloud providers. The platform is built to work with [HashiCorp Vault](https://www.vaultproject.io/) to enable VMware vSphere to be used as a trusted platform similar to public cloud providers such as AWS and Azure.

# vAuth Architecture

The vAuth platform is composed of multiple containers that leverage cloud-native practices with high availability and resillency in mind. The diagram below displays how the different services interact with one another.

![](./vauth-architecture.png)

|Component|Description|
|---------|-----------|
| Scheduler| Schedules the synchronization process to ensure that all virtual machines have identity data in the event a real-time event was missed |
| Syncer   | Synchronize virtual machines without identity data |
| Watcher | Watches for relevant VMware vSphere events such as power on operations and custom attribute changes |
| Worker | Generate the identity data (VM Name, Datacenter, Role and Secret Key) and add that to the virtual machines attributes|
| Backend | REST API for Vault authentication validation and database interaction |
| NATS | Event bus for system transactions |
| DB | Cockroachdb for persistent virtual machine identity data presented by the backend in response to Vault auth validation requests |

## Libraries

The following third party libraries have been used to build the vAuth platform.

|Name|Description|
|----|-----------|
|govmomi| VMware vSphere Golang SDK |
|nats-io| NATS pub/sub platform Golang SDK|
|jasonlvhit/gocron||
|mux| HTTTP router |
|pq| Golang Postgres driver for SQL database interaction|

## vSphere Account Permissions

The vAuth platform requires access to VMware vSphere to perform various operations such as watch events, update virtual machine guest information and more. The following table details the permissions that the account used by the vAuth platform would need in vSphere.

**Privileged interaction**
The following operations require a privilege to be assigned to the vSphere account that the vAuth platform uses.

|Permission|Description|
|----|-----------|
|Virtual Machine > Change Configuration > Advanced Configuration | The account needs to have permission to update the advanced configuration of virtual machines to provide the identity data to the guest operating system|

**Non-privileged interaction**

The watcher service connects to the vSphere event manager to "watch" for power on events in order to ensure that virtual machines receive identity data when they are powered on.

The watcher service also watches for changes to virtual machine custom attributes which are used for assigning the role to virtual machines.

## Setup


```
docker-compose up -d
```