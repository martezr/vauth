package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"reflect"
	"strconv"
	"syscall"
	"time"

	"github.com/gorilla/mux"
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

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

var (
	db *bolt.DB
)

var config utils.Config
var vsphereClient *govmomi.Client

////go:embed frontend/dist/*
//var webUI embed.FS

//func clientHandler() http.Handler {
//	fsys := fs.FS(webUI)
//	contentStatic, _ := fs.Sub(fsys, "frontend/dist")
//	return http.FileServer(http.FS(contentStatic))
//}

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run the vAuth server",
	Long:  `Show this help output, or the help for a specified subcommand.`,
	Run: func(cmd *cobra.Command, args []string) {
		server()
	},
}

func server() {
	log.Println("Started vAuth 0.0.1")
	cfg := viper.New()
	cfg.AddConfigPath(".")
	cfg.SetConfigName("config")
	cfg.SetConfigType("yaml")

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
	config.VaultAddress = cfg.GetString("vault_address")
	config.VaultAppRoleMount = cfg.GetString("vault_approle_mount")
	config.VaultTLSSkipVerify = cfg.GetBool("vault_tls_skip_verify")
	config.VaultWrapResponse = cfg.GetBool("vault_wrap_response")

	vcenterURL, err := url.Parse(fmt.Sprintf("https://%v/sdk", config.VsphereServer))
	if err != nil {
		log.Println(err)
	}
	credentials := url.UserPassword(config.VsphereUsername, config.VspherePassword)
	vcenterURL.User = credentials

	// Connecting to vCenter
	log.Print("Connecting to vCenter")
	vsphereClient, err = govmomi.NewClient(ctx, vcenterURL, true)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Connected to vCenter: %s", config.VsphereServer)
	finder := find.NewFinder(vsphereClient.Client, true)

	dc, err := finder.DatacenterOrDefault(ctx, "")
	if err != nil {
		fmt.Println(err)
	}

	finder.SetDatacenter(dc)
	refs := []types.ManagedObjectReference{dc.Reference()}

	db = database.StartDB(config.DataDir)

	//database.ListDBRecords(db)
	// Setting up the event manager
	eventManager := event.NewManager(vsphereClient.Client)
	go eventManager.Events(ctx, refs, 50, true, false, handleEvent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	//	http.Handle("/", clientHandler())

	// Start the server.
	log.Println("UI listening on port", config.UIPort)
	port := fmt.Sprintf(":%s", config.UIPort)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/vms", listVirtualMachines).Methods("GET")
	log.Fatal(http.ListenAndServe(port, router))
}

func listVirtualMachines(w http.ResponseWriter, r *http.Request) {
	var output utils.VmRecords
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
			log.Printf("Detected power on event for %s", vmName)
			eventID := fmt.Sprintf("%d", event.GetEvent().ChainId)
			if isUnprocessedEvent(event) {
				database.AddDBRecord(db, vmName, eventID)
				updateVM(config.VaultAddress, config.VaultToken, vmName)
			}
		}
		// Detect VM custom attribute change
		if eventType == "*types.CustomFieldValueChangedEvent" {
			vmName := event.GetEvent().Vm.Name
			log.Printf("Detected custom attribute change event for %s", vmName)
			eventID := fmt.Sprintf("%d", event.GetEvent().ChainId)
			if isUnprocessedEvent(event) {
				database.AddDBRecord(db, vmName, eventID)
				updateVM(config.VaultAddress, config.VaultToken, vmName)
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
	evalID := database.ViewDBRecord(db, vmName)
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

func updateVM(vaultAddr string, token string, vmname string) {
	// Creating a connection context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	finder := find.NewFinder(vsphereClient.Client, true)

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
	attkey, _ := object.GetCustomFieldsManager(vsphereClient.Client)
	attID, _ := attkey.FindKey(ctx, "vauth-role")

	for _, vmdata := range machines {
		machine := machines[0]
		var props mo.VirtualMachine
		machine.Properties(ctx, vmdata.Reference(), nil, &props)
		if !props.Summary.Config.Template && props.Summary.Runtime.PowerState == "poweredOn" {
			customAttrs := make(map[string]interface{})
			role := ""
			if len(props.CustomValue) > 0 {
				for _, fv := range props.CustomValue {
					value := fv.(*types.CustomFieldStringValue).Value
					if value != "" {
						customAttrs[fmt.Sprint(fv.GetCustomFieldValue().Key)] = value
					}
					if fv.GetCustomFieldValue().Key == attID {
						log.Printf("Found the %s role associated with %s", value, vmname)
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
					authSpec := types.VirtualMachineConfigSpec{
						ExtraConfig: settings,
					}
					vmdata.Reconfigure(ctx, authSpec)
					log.Printf("Updated VM: %s", vmname)
				}
				if status == "role not found" {
					log.Printf("The %s role associated with %s does not exist in Vault", role, vmname)
				}
			} else {
				log.Printf("No role associated with %s", vmname)
			}
		}
	}
}

func SetupCloseHandler() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("\r- Ctrl+C pressed in Terminal")
		os.Exit(0)
	}()
}
