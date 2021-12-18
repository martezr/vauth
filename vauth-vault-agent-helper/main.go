package main

import (
	"log"
	"os"

	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/vmware/vmw-guestinfo/rpcvmx"
	"github.com/vmware/vmw-guestinfo/vmcheck"
)

type Config struct {
	SyncInterval     int    `hcl:"sync_interval"`
	RoleIDFilePath   string `hcl:"role_id_file_path"`
	SecretIDFilePath string `hcl:"secret_id_file_path"`
}

// WriteVaultFile writes the vault auth data to a file
func writeVaultFile(path string, data string) {
	f, err := os.Create(path)

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.WriteString(data)

	if err2 != nil {
		log.Fatal(err2)
	}
}

func main() {
	// Evaluate whether the workload is running on vSphere
	isVM, err := vmcheck.IsVirtualWorld()
	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	if !isVM {
		log.Fatalf("ERROR: not running on vSphere or VMware tools are not installed.")
	}

	var config Config
	hclErr := hclsimple.DecodeFile("config.hcl", nil, &config)
	if hclErr != nil {
		log.Fatalf("Failed to load configuration: %s", hclErr)
	}
	log.Printf("Configuration is %#v", config)

	vsphereConfig := rpcvmx.NewConfig()

	roledata := ""
	if out, err := vsphereConfig.String("guestinfo.vault.roleid", ""); err != nil {
		log.Fatalf("ERROR: String failed with %s", err)
	} else {
		roledata = out
	}
	writeVaultFile(config.RoleIDFilePath, roledata)

	secretdata := ""
	if out, err := vsphereConfig.String("guestinfo.vault.secretid", ""); err != nil {
		log.Fatalf("ERROR: String failed with %s", err)
	} else {
		secretdata = out
	}
	writeVaultFile(config.SecretIDFilePath, secretdata)
}
