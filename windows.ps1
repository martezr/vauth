$VMwareToolsExe = "C:\Program Files\VMware\VMware Tools\vmtoolsd.exe"
$VMName = & $VMwareToolsExe --cmd "info-get guestinfo.vault.vmname" | Out-String  
$Role = & $VMwareToolsExe --cmd "info-get guestinfo.vault.role" | Out-String
$Datacenter = & $VMwareToolsExe --cmd "info-get guestinfo.vault.datacenter" | Out-String
$Secretkey = & $VMwareToolsExe --cmd "info-get guestinfo.vault.secretkey" | Out-String
$VMName = $VMname.replace("`n"," ")
$Role = $Role.replace("`n"," ")
$Datacenter = $Datacenter.replace("`n"," ")
$Secretkey = $Secretkey.replace("`n"," ")

$test = @{}
$test['vmname'] = $VMName.Replace("`r","").Trim()
$test['datacenter'] = $Datacenter.Replace("`r","").Trim()
$test['secretkey'] = $Secretkey.Replace("`r","").Trim()
$test['role'] = $Role.Replace("`r","").Trim()
$test = $test | ConvertTo-Json

$vault_output = Invoke-WebRequest -Uri http://10.0.0.70:8200/v1/auth/vsphere/login -ContentType "application/json" -Method POST -Body $test
$vault_output.Content |  ConvertFrom-Json | ConvertTo-Json

 
