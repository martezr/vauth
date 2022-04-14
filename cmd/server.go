package cmd

import (
	"context"
	"crypto/tls"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/martezr/vauth/approle"
	"github.com/martezr/vauth/database"
	"github.com/martezr/vauth/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/event"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	bolt "go.etcd.io/bbolt"
)

type extraConfig []types.BaseOptionValue

var (
	db *bolt.DB
)

var config utils.Config
var vsphereClient *govmomi.Client

//go:embed webui/*
var webUI embed.FS

func clientHandler() http.Handler {
	fsys := fs.FS(webUI)
	contentStatic, _ := fs.Sub(fsys, "webui")
	return http.FileServer(http.FS(contentStatic))
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the vAuth server",
	Long:  `Show this help output, or the help for a specified subcommand.`,
	Run: func(cmd *cobra.Command, args []string) {
		server()
	},
}

func server() {
	log.Println("starting vAuth 0.0.1")
	cfg := viper.New()
	if cfgFile != "" {
		cfg.SetConfigFile(cfgFile)
	} else {
		cfg.AddConfigPath(".")
		cfg.SetConfigName("config")
		cfg.SetConfigType("yaml")
	}

	cfg.AutomaticEnv()

	if err := cfg.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			log.Println("No config file found")
		} else {
			// Config file was found but another error was produced
		}
	}

	err := cfg.Unmarshal(&config)
	if err != nil {
		log.Fatalf("unable to decode into struct, %v", err)
	}

	// Creating a connection context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config.UIPort = cfg.GetString("ui_port")
	config.DataDir = cfg.GetString("data_dir")
	config.VsphereServer = cfg.GetString("vsphere_server")
	config.VsphereTLSSkipVerify = cfg.GetBool("vsphere_tls_skip_verify")
	config.VsphereUsername = cfg.GetString("vsphere_username")
	config.VspherePassword = cfg.GetString("vsphere_password")
	config.VsphereDatacenters = cfg.GetStringSlice("vsphere_datacenters")
	config.VaultAddress = cfg.GetString("vault_address")
	config.VaultAppRoleMount = cfg.GetString("vault_approle_mount")
	config.VaultTLSSkipVerify = cfg.GetBool("vault_tls_skip_verify")
	config.VaultToken = cfg.GetString("vault_token")
	config.VaultWrapResponse = cfg.GetBool("vault_wrap_response")

	vcenterURL, err := url.Parse(fmt.Sprintf("https://%v/sdk", config.VsphereServer))
	if err != nil {
		log.Println(err)
	}
	credentials := url.UserPassword(config.VsphereUsername, config.VspherePassword)
	vcenterURL.User = credentials

	// Connecting to vCenter
	log.Print("connecting to vCenter server")

	vsphereClient, err = govmomi.NewClient(ctx, vcenterURL, true)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("connected to vCenter: %s", config.VsphereServer)
	finder := find.NewFinder(vsphereClient.Client, true)
	db = database.StartDB(config.DataDir)

	for _, datacenter := range config.VsphereDatacenters {
		// Iterate through dcs
		dc, err := finder.DatacenterOrDefault(ctx, datacenter)
		if err != nil {
			log.Println(err)
		}
		log.Println(dc)

		finder.SetDatacenter(dc)
		refs := []types.ManagedObjectReference{dc.Reference()}
		//database.ListDBRecords(db)
		// Setting up the event manager
		eventManager := event.NewManager(vsphereClient.Client)
		go eventManager.Events(ctx, refs, 100, true, false, handleEvent)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %s\n", err)
			os.Exit(1)
		}
	}

	log.Println("ui listening on port", config.UIPort)
	port := fmt.Sprintf(":%s", config.UIPort)

	http.Handle("/", clientHandler())
	http.HandleFunc("/api/v1/vms", listVirtualMachines)
	http.HandleFunc("/api/v1/snapshot", BackupHandleFunc)
	http.HandleFunc("/api/v1/health", GetHealthStatus)

	log.Panic(
		http.ListenAndServe(port, nil),
	)
}

func GetVsphereHealthStatus() string {
	var httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	httpClient.Transport = tr
	vSphereURL := fmt.Sprintf("https://%s", config.VsphereServer)
	resp, err := httpClient.Get(vSphereURL)
	if err != nil {
		log.Println(err)
		return "unhealthy"
	}
	if resp.StatusCode == 200 {
		return "healthy"
	} else {
		return "unhealthy"
	}
}

func GetVaultHealthStatus() string {
	var httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	httpClient.Transport = tr
	VaultURL := fmt.Sprintf("%s/v1/sys/health", config.VaultAddress)
	resp, err := httpClient.Get(VaultURL)
	if err != nil {
		log.Println(err)
		return "unhealthy"
	}
	if resp.StatusCode == 200 {
		return "healthy"
	} else {
		return "unhealthy"
	}
}

// GetHealthStatus returns the health of the vAuth platform
func GetHealthStatus(w http.ResponseWriter, r *http.Request) {
	var output utils.HealthStatus
	output.Version = "0.0.1"
	output.VaultStatus = GetVaultHealthStatus()
	output.VsphereStatus = GetVsphereHealthStatus()
	jsonOutput, _ := json.MarshalIndent(output, "", " ")
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonOutput)
}

func listVirtualMachines(w http.ResponseWriter, r *http.Request) {
	var output utils.VMRecords
	data := database.ListDBRecords(db)
	output.Records = data
	output.Total = len(data)
	jsonOutput, _ := json.MarshalIndent(output, "", " ")
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonOutput)
}

func handleEvent(ref types.ManagedObjectReference, events []types.BaseEvent) (err error) {
	for _, event := range events {
		eventType := reflect.TypeOf(event).String()
		// Detect VM power on events
		if eventType == "*types.VmPoweredOnEvent" {
			vmName := event.GetEvent().Vm.Name
			log.Printf("detected power on event for %s", vmName)
			eventID := fmt.Sprintf("%d", event.GetEvent().ChainId)
			if isUnprocessedEvent(event) {
				role := updateVM(config.VaultAddress, config.VaultToken, vmName, event.GetEvent().Datacenter.Name)
				var vmData utils.VMRecord
				vmData.LatestEventId = eventID
				vmData.Name = vmName
				vmData.Role = role
				out, _ := json.Marshal(vmData)
				database.AddDBRecord(db, vmName, string(out))
			}
		}
		// Detect VM custom attribute change
		if eventType == "*types.CustomFieldValueChangedEvent" {
			vmName := event.GetEvent().Vm.Name
			log.Printf("detected custom attribute change event for %s", vmName)
			eventID := fmt.Sprintf("%d", event.GetEvent().ChainId)
			if isUnprocessedEvent(event) {
				role := updateVM(config.VaultAddress, config.VaultToken, vmName, event.GetEvent().Datacenter.Name)
				var vmData utils.VMRecord
				vmData.LatestEventId = eventID
				vmData.Name = vmName
				vmData.Role = role
				out, _ := json.Marshal(vmData)
				database.AddDBRecord(db, vmName, string(out))
			}
		}
		// Detect VM removal events
		if eventType == "*types.VmRemovedEvent" {
			vmName := event.GetEvent().Vm.Name
			log.Printf("detected remove event for %s", vmName)
			if isUnprocessedEvent(event) {
				log.Printf("Delete VM: %s", vmName)
				database.DeleteDBRecord(db, vmName)
			}
		}
	}
	return nil
}

func isUnprocessedEvent(event types.BaseEvent) (response bool) {
	vmName := event.GetEvent().Vm.Name
	eventID := fmt.Sprintf("%d", event.GetEvent().ChainId)
	eventIDInt, err := strconv.Atoi(eventID)
	if err != nil {
		log.Println(err)
	}
	evalData := database.ViewDBRecord(db, vmName)
	var evalParse utils.VMRecord
	var evalID string
	if evalData != "" {
		err = json.Unmarshal([]byte(evalData), &evalParse)
		if err != nil {
			log.Println(err)
		}
		evalID = evalParse.LatestEventId
	} else {
		evalID = ""
	}
	var evalIDInt int
	if evalID == "" {
		evalIDInt = 0
	} else {
		evalIDInt, err = strconv.Atoi(evalID)
		if err != nil {
			log.Println(err)
		}
	}

	if eventIDInt > evalIDInt {
		return true
	}
	return false
}

func updateVM(vaultAddr string, token string, vmname string, datacenter string) (role string) {
	// Creating a connection context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	finder := find.NewFinder(vsphereClient.Client, true)

	vsphereDatacenter, err := finder.DatacenterOrDefault(ctx, datacenter)
	if err != nil {
		log.Println(err)
	}

	finder.SetDatacenter(vsphereDatacenter)

	machines, err := finder.VirtualMachineList(ctx, vmname)
	if err != nil {
		log.Println(err)
	}
	// Fetch IAM Role information
	attkey, _ := object.GetCustomFieldsManager(vsphereClient.Client)
	attID, _ := attkey.FindKey(ctx, "vauth-role")

	for _, vmdata := range machines {
		machine := machines[0]
		var props mo.VirtualMachine
		machine.Properties(ctx, vmdata.Reference(), nil, &props)
		if !props.Summary.Config.Template && props.Summary.Runtime.PowerState == "poweredOn" {
			customAttrs := make(map[string]interface{})
			role = ""
			if len(props.CustomValue) > 0 {
				for _, fv := range props.CustomValue {
					value := fv.(*types.CustomFieldStringValue).Value
					if value != "" {
						customAttrs[fmt.Sprint(fv.GetCustomFieldValue().Key)] = value
					}
					if fv.GetCustomFieldValue().Key == attID {
						log.Printf("found the %s role associated with %s", value, vmname)
						role = value
					}
				}
			}

			if role != "" {
				status, roleid, secretid, secretidttl := approle.FetchAppRole(config, vaultAddr, token, role, vmname)
				if status == "role found" {
					var settings extraConfig
					settings = append(settings, &types.OptionValue{Key: "guestinfo.vault.role", Value: role})
					settings = append(settings, &types.OptionValue{Key: "guestinfo.vault.roleid", Value: roleid})
					settings = append(settings, &types.OptionValue{Key: "guestinfo.vault.secretid", Value: secretid})
					settings = append(settings, &types.OptionValue{Key: "guestinfo.vault.secretidttl", Value: secretidttl})
					settings = append(settings, &types.OptionValue{Key: "guestinfo.vault.vaultaddr", Value: config.VaultAddress})
					authSpec := types.VirtualMachineConfigSpec{
						ExtraConfig: settings,
					}
					vmdata.Reconfigure(ctx, authSpec)
					log.Printf("updated VM: %s", vmname)
				}
				if status == "role not found" {
					log.Printf("the %s role associated with %s does not exist in Vault", role, vmname)
				}
			} else {
				log.Printf("no role associated with %s", vmname)
			}
		}
	}
	return role
}
