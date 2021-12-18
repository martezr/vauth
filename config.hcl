ui_port = "8000"
vsphere "grtlab" {
    vsphere_url      = "grtvcenter01.grt.local"
    vsphere_username = "vauth@vsphere.local"
    vsphere_password = "Password123#"
}

vault "grtmanage" {
    vault_address       = "https://grtmanage01.grt.local:8200"
    tls_skip_verify     = true
    vault_token         = "s.ewdkchV1oqIwTxxI8G3INWVG"
    vault_approle_mount = "approle"
    wrap_response       = true
}
