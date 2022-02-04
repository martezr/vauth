export GOVC_URL='https://user:pass@localhost/sdk'
export GOVC_INSECURE=true

govc fields.add -type VirtualMachine 'vauth-role'

govc ls /DC0/vm | xargs govc fields.set 'vauth-role' 'demo01'

sleep 5

govc ls /DC0/vm | xargs govc vm.power -reset