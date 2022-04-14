#!/bin/bash

mkdir -p /opt/vauth

# Create vAuth system user and group
id -u vauth &>/dev/null || useradd --system --no-create-home --user-group vauth

# Create vAuth service
cat << 'EOF' > /etc/systemd/system/vauth.service
[Unit]
Description="vAuth - A vSphere identity platform"
Documentation=https://martezr.github.io/vauth
Requires=network-online.target
After=network-online.target

[Service]
User=vauth
Group=vauth
ExecStart=/usr/bin/vauth server
ExecReload=/bin/kill --signal HUP $MAINPID
KillMode=process
KillSignal=SIGTERM
Restart=on-failure
LimitNOFILE=65536

[Install]
WantedBy=multi-user.target
EOF

cat << 'EOF' > /opt/vauth/config.yaml
---
ui_port: 8000
data_dir: /opt/vauth
vsphere_server: "localhost"
vsphere_tls_skip_verify: true
vsphere_username: "user"
vsphere_password: "pass"
vsphere_datacenters: ["DC0","DC1","DC2"]
vault_address: "http://localhost:8200"
vault_token: "vault"
vault_approle_mount: "approle"
vault_wrap_response: true
vault_tls_skip_verify: true
EOF

systemctl enable vauth
systemctl start vauth