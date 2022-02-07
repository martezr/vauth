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
	"github.com/martezr/vauth/utils"
)

// Metadata defines the secret ID metadata payload for auditing
type Metadata struct {
	VirtualMachineName string `json:"name"`
}

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

// FetchAppRole fetches the approle credentials
func FetchAppRole(config utils.Config, vaultAddr string, token string, rolename string, vmname string) (status string, roleid string, secretid string, secretidttl string) {
	if config.VaultTLSSkipVerify {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}

		httpClient.Transport = tr
	}
	client, err := api.NewClient(&api.Config{Address: vaultAddr, HttpClient: httpClient})
	if err != nil {
		log.Println(err)
	}
	client.SetToken(token)

	// Verify the mount point exists
	authMethods, authmethoderr := client.Sys().ListAuth()
	if authmethoderr != nil {
		log.Println(authmethoderr)
	}

	pathExists := false
	for k := range authMethods {
		path := strings.Replace(k, "/", "", -1)
		if path == config.VaultAppRoleMount {
			pathExists = true
		}
	}

	if !pathExists {
		log.Printf("Unable to find %s mount point", config.VaultAppRoleMount)
		return "role not found", "", "", ""
	}

	// Evaluate if the role associated with the VM exists
	path := fmt.Sprintf("auth/%s/role", config.VaultAppRoleMount)
	roles, roleerr := client.Logical().List(path)
	if roleerr != nil {
		log.Println(roleerr)
	}

	if roles == nil {
		log.Printf("%s contains no roles", config.VaultAppRoleMount)
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
		rolepath := fmt.Sprintf("auth/%s/role/%s/role-id", config.VaultAppRoleMount, rolename)
		data, err := client.Logical().Read(rolepath)
		if err != nil {
			log.Println(err)
		}

		metadataPayload := map[string]string{"virtual_machine_name": vmname}
		metadataJSON, _ := json.Marshal(metadataPayload)

		options := map[string]interface{}{
			"metadata": string(metadataJSON),
		}
		secretpath := fmt.Sprintf("auth/%s/role/%s/secret-id", config.VaultAppRoleMount, rolename)

		if config.VaultWrapResponse {
			client.SetWrappingLookupFunc(func(operation, path string) string {
				return "5m"
			})
		}

		secretdata, newerr := client.Logical().Write(secretpath, options)
		if newerr != nil {
			log.Println(newerr)
		}

		tokenTTL := secretdata.Data["secret_id_ttl"]

		log.Printf("Token TTL: %v", tokenTTL)

		if config.VaultWrapResponse {
			if secretdata == nil {
				log.Println("nil secret")
			}
			if secretdata.WrapInfo == nil {
				log.Println("nil wrap info")
			}

			token := secretdata.WrapInfo.Token

			return "role found", data.Data["role_id"].(string), token, ""
		}
		return "role found", data.Data["role_id"].(string), secretdata.Data["secret_id"].(string), string(secretdata.Data["secret_id_ttl"].(json.Number))
	}
	return "role not found", "", "", ""
}
