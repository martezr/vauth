package approle

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
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
	logError(err)

	client.SetToken(token)

	// Verify the mount point exists
	authMethods, authmethoderr := client.Sys().ListAuth()
	logError(authmethoderr)

	pathExists := false
	for k := range authMethods {
		path := strings.Replace(k, "/", "", -1)
		if path == config.VaultAppRoleMount {
			pathExists = true
		}
	}

	if !pathExists {
		hclog.Default().Info(fmt.Sprintf("unable to find %s mount point", config.VaultAppRoleMount))
		return "role not found", "", "", ""
	}

	// Evaluate if the role associated with the VM exists
	path := fmt.Sprintf("auth/%s/role", config.VaultAppRoleMount)
	roles, roleerr := client.Logical().List(path)
	logError(roleerr)

	if roles == nil {
		hclog.Default().Info(fmt.Sprintf("%s contains no roles", config.VaultAppRoleMount))
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
		logError(err)

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
		logError(newerr)

		//tokenData := secretdata

		if config.VaultWrapResponse {
			if secretdata == nil {
				hclog.Default().Named("vault").Warn("nil secret")
			}
			if secretdata.WrapInfo == nil {
				hclog.Default().Named("vault").Warn("nil wrap info")
			}

			token := secretdata.WrapInfo.Token

			return "role found", data.Data["role_id"].(string), token, ""
		}
		current_time := time.Now().UTC()
		val, _ := secretdata.Data["secret_id_ttl"].(json.Number).Int64()
		leaseTime := time.Second * time.Duration(val)

		tokenExpirationTime := current_time.Add(leaseTime)
		expirationDate := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02dZ",
			tokenExpirationTime.Year(), tokenExpirationTime.Month(), tokenExpirationTime.Day(),
			tokenExpirationTime.Hour(), tokenExpirationTime.Minute(), tokenExpirationTime.Second())
		return "role found", data.Data["role_id"].(string), secretdata.Data["secret_id"].(string), string(expirationDate)
	}
	return "role not found", "", "", ""
}

// logError logs error messages
func logError(errormessage error) {
	if errormessage != nil {
		hclog.Default().Named("vault").Error(errormessage.Error())
	}
}
