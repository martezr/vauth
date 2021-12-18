package approle

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/vault/api"
	"github.com/martezr/vsphere-vauth/utils"
)

type Metadata struct {
	VirtualMachineName string `json:"name"`
}

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

func FetchAppRole(config utils.Config, vaultAddr string, token string, rolename string, vmname string) (status string, roleid string, secretid string, secretidttl string) {
	if config.Vault[0].TLSSkipVerify {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		httpClient.Transport = tr
	}
	client, err := api.NewClient(&api.Config{Address: vaultAddr, HttpClient: httpClient})
	if err != nil {
		panic(err)
	}
	client.SetToken(token)

	// Verify mount point exists
	authMethods, authmethoderr := client.Sys().ListAuth()
	if authmethoderr != nil {
		panic(authmethoderr)
	}

	pathExists := false
	for k := range authMethods {
		path := strings.Replace(k, "/", "", -1)
		if path == config.Vault[0].VaultAppRoleMount {
			log.Printf("Found %s mount point", config.Vault[0].VaultAppRoleMount)
			pathExists = true
		}
	}

	if !pathExists {
		log.Printf("Unable to find %s mount point", config.Vault[0].VaultAppRoleMount)
		return "role not found", "", "", ""
	}

	// Evaluate if the role associated with the VM exists
	path := fmt.Sprintf("auth/%s/role", config.Vault[0].VaultAppRoleMount)
	roles, roleerr := client.Logical().List(path)
	if roleerr != nil {
		panic(roleerr)
	}

	if roles == nil {
		log.Printf("%s contains no roles", config.Vault[0].VaultAppRoleMount)
		return "role not found", "", "", ""
	}

	output := roles.Data["keys"]

	list := output.([]interface{})
	var outroles []string
	for _, role := range list {
		outroles = append(outroles, role.(string))
	}

	// Fetch role ID from Vault
	if utils.Contains(outroles, rolename) {
		rolepath := fmt.Sprintf("auth/%s/role/%s/role-id", config.Vault[0].VaultAppRoleMount, rolename)
		data, err := client.Logical().Read(rolepath)
		if err != nil {
			panic(err)
		}

		metadataPayload := map[string]string{"virtual_machine_name": vmname}
		metadataJSON, _ := json.Marshal(metadataPayload)

		options := map[string]interface{}{
			"metadata": string(metadataJSON),
		}
		secretpath := fmt.Sprintf("auth/%s/role/%s/secret-id", config.Vault[0].VaultAppRoleMount, rolename)

		if config.Vault[0].WrapResponse {
			client.SetWrappingLookupFunc(func(operation, path string) string {
				return "5m"
			})
		}

		secretdata, newerr := client.Logical().Write(secretpath, options)
		if newerr != nil {
			panic(newerr)
		}

		if config.Vault[0].WrapResponse {
			if secretdata == nil {
				log.Fatal("nil secret")
			}
			if secretdata.WrapInfo == nil {
				log.Fatal("nil wrap info")
			}

			token := secretdata.WrapInfo.Token

			return "role found", data.Data["role_id"].(string), token, ""
		}
		return "role found", data.Data["role_id"].(string), secretdata.Data["secret_id"].(string), string(secretdata.Data["secret_id_ttl"].(json.Number))
	}
	return "role not found", "", "", ""
}
