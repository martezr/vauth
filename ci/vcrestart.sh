export GOVC_URL='https://user:pass@localhost/sdk'
export GOVC_INSECURE=true

govc ls /DC0/vm | xargs govc vm.power -reset