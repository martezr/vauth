package function

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/hashicorp/vault/api"
	handler "github.com/openfaas/templates-sdk/go-http"
	toml "github.com/pelletier/go-toml"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

type extraConfig []types.BaseOptionValue

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

// cloudEvent is a subsection of a Cloud Event.
type cloudEvent struct {
	Data    types.Event `json:"data"`
	Source  string      `json:"source"`
	Subject string      `json:"subject"`
}

const secretPath = "/var/openfaas/secrets/vebaconfig"

// vcConfig represents the toml vcconfig file
type vcConfig struct {
	VCenter struct {
		Server   string
		User     string
		Password string
		Insecure bool
	}
	Vault struct {
		Server string
		Token  string
	}
}

// Handle parses the event, interacts with Vault and populates the virtual machine's guest info
func Handle(req handler.Request) (handler.Response, error) {
	// Parse the vSphere event
	cloudEvt, err := parseCloudEvent(req.Body)
	if err != nil {
		return errRespondAndLog(fmt.Errorf("parsing cloud event data: %w", err))
	}

	// Load config every time, to ensure the most updated version is used.
	cfg, err := loadTomlCfg(secretPath)
	if err != nil {
		return errRespondAndLog(fmt.Errorf("loading of vcconfig: %w", err))
	}

	vsphereUsername := cfg.VCenter.User
	vspherePassword := cfg.VCenter.Password
	vsphereServer := cfg.VCenter.Server

	vcenterURL, err := url.Parse(fmt.Sprintf("https://%v/sdk", vsphereServer))
	if err != nil {
		log.Println(err)
	}
	credentials := url.UserPassword(vsphereUsername, vspherePassword)
	vcenterURL.User = credentials

	vaultAddr := cfg.Vault.Server
	staticToken := cfg.Vault.Token

	updateVM(vaultAddr, staticToken, cloudEvt.Data.Vm.Name, vcenterURL)

	return handler.Response{
		Body:       req.Body,
		StatusCode: http.StatusOK,
	}, nil
}

func errRespondAndLog(err error) (handler.Response, error) {
	log.Println(err.Error())

	return handler.Response{
		Body:       []byte(err.Error()),
		StatusCode: http.StatusInternalServerError,
	}, err
}

func parseCloudEvent(req []byte) (cloudEvent, error) {
	var event cloudEvent

	err := json.Unmarshal(req, &event)
	if err != nil {
		return cloudEvent{}, fmt.Errorf("unmarshalling json: %w", err)
	}

	return event, nil
}

func loadTomlCfg(path string) (*vcConfig, error) {
	var cfg vcConfig

	secret, err := toml.LoadFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to load vcconfig.toml: %w", err)
	}

	err = secret.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal vcconfig.toml: %w", err)
	}

	err = validateConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("insufficient information in vcconfig.toml: %w", err)
	}

	return &cfg, nil
}

// ValidateConfig ensures the bare minimum of information is in the config file.
func validateConfig(cfg vcConfig) error {
	reqFields := map[string]string{
		"vcenter server":   cfg.VCenter.Server,
		"vcenter user":     cfg.VCenter.User,
		"vcenter password": cfg.VCenter.Password,
		"vault server":     cfg.Vault.Server,
		"vault token":      cfg.Vault.Token,
	}

	// Multiple fields may be missing, but err on the first encountered.
	for k, v := range reqFields {
		if v == "" {
			return errors.New("required field(s) missing, including " + k)
		}
	}

	return nil
}

func contains(slice []string, inputValue string) bool {
	for _, sliceValue := range slice {
		if sliceValue == inputValue {
			return true
		}
	}
	return false
}

func fetchAppRole(vaultAddr string, token string, backend string, rolename string) (status string, roleid string, secretid string) {

	client, err := api.NewClient(&api.Config{Address: vaultAddr, HttpClient: httpClient})
	if err != nil {
		panic(err)
	}

	client.SetToken(token)

	roles, roleerr := client.Logical().List("auth/approle/role")
	if roleerr != nil {
		panic(roleerr)
	}

	output := roles.Data["keys"]

	list := output.([]interface{})
	var outroles []string
	for _, role := range list {
		outroles = append(outroles, role.(string))
	}

	if contains(outroles, rolename) {
		rolepath := fmt.Sprintf("auth/approle/role/%s/role-id", rolename)
		data, err := client.Logical().Read(rolepath)
		if err != nil {
			panic(err)
		}

		options := map[string]interface{}{}
		secretpath := fmt.Sprintf("auth/approle/role/%s/secret-id", rolename)
		secretdata, newerr := client.Logical().Write(secretpath, options)
		if newerr != nil {
			panic(newerr)
		}

		return "role found", data.Data["role_id"].(string), secretdata.Data["secret_id"].(string)
	}

	return "role not found", "", ""
}

func updateVM(vaultAddr string, token string, vmname string, vcurl *url.URL) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Connecting to vCenter
	log.Print("Connecting to vCenter")
	client, err := govmomi.NewClient(ctx, vcurl, true)
	if err != nil {
		fmt.Println(err)
	}

	finder := find.NewFinder(client.Client, true)

	vsphereDatacenter, err := finder.DatacenterOrDefault(ctx, "")
	if err != nil {
		fmt.Println(err)
	}

	finder.SetDatacenter(vsphereDatacenter)

	machines, err := finder.VirtualMachineList(ctx, vmname)
	if err != nil {
		fmt.Println(err)
	}
	// Fetch IAM Role information
	attkey, _ := object.GetCustomFieldsManager(client.Client)
	attID, _ := attkey.FindKey(ctx, "vauth-role")

	for _, vmdata := range machines {
		machine := machines[0]
		var props mo.VirtualMachine
		machine.Properties(ctx, vmdata.Reference(), nil, &props)
		if props.Summary.Config.Template == false && props.Summary.Runtime.PowerState == "poweredOn" {
			customAttrs := make(map[string]interface{})
			role := ""
			if len(props.CustomValue) > 0 {
				for _, fv := range props.CustomValue {
					value := fv.(*types.CustomFieldStringValue).Value
					if value != "" {
						customAttrs[fmt.Sprint(fv.GetCustomFieldValue().Key)] = value
					}
					if fv.GetCustomFieldValue().Key == attID {
						log.Printf("vauth-role value: %s", value)
						role = value
					}
				}
			}

			status, roleid, secretid := fetchAppRole(vaultAddr, token, "approle", role)
			if status == "role found" {
				fmt.Print(status, roleid, secretid)

				var settings extraConfig
				settings = append(settings, &types.OptionValue{Key: "guestinfo.vault.role", Value: role})
				settings = append(settings, &types.OptionValue{Key: "guestinfo.vault.roleid", Value: roleid})
				settings = append(settings, &types.OptionValue{Key: "guestinfo.vault.secretid", Value: secretid})
				authSpec := types.VirtualMachineConfigSpec{
					ExtraConfig: settings,
				}
				vmdata.Reconfigure(ctx, authSpec)
				log.Printf("Updated VM: %s", vmname)
			}
			if status == "role found" {
				log.Printf("The %s role associated with %s does not exist in Vault", role, vmname)
			}
		}
	}
}
