package main

import (
	"context"
	"fmt"
	"io/ioutil"
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

func main() {

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
	config.VsphereUsername = cfg.GetString("vsphere_username")
	config.VspherePassword = cfg.GetString("vsphere_password")
	config.VsphereURL = cfg.GetString("vsphere_server")
	config.VaultAddress = cfg.GetString("vault_address")
	config.VaultAppRoleMount = cfg.GetString("vault_app_role_mount")
	config.TLSSkipVerify = cfg.GetBool("tls_skip_verify")
	config.WrapResponse = cfg.GetBool("wrap_response")

	vcenterURL, err := url.Parse(fmt.Sprintf("https://%v/sdk", config.VsphereURL))
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
	log.Printf("Connected to vCenter: %s", config.VsphereURL)
	finder := find.NewFinder(vsphereClient.Client, true)

	dc, err := finder.DatacenterOrDefault(ctx, "")
	if err != nil {
		fmt.Println(err)
	}

	finder.SetDatacenter(dc)
	refs := []types.ManagedObjectReference{dc.Reference()}

	db = database.StartDB(config.DatDir)

	database.ListDBRecords(db)
	// Setting up the event manager
	eventManager := event.NewManager(vsphereClient.Client)
	go eventManager.Events(ctx, refs, 20, true, false, handleEvent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	//	http.Handle("/", clientHandler())

	// Start the server.
	uiport := cfg.GetString("ui_port")
	log.Println("UI listening on port", uiport)
	port := fmt.Sprintf(":%s", uiport)

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/vms", listVirtualMachines).Methods("POST")
	log.Fatal(http.ListenAndServe(port, router))
}

func listVirtualMachines(w http.ResponseWriter, r *http.Request) {
	namespace := mux.Vars(r)["namespace"]
	name := mux.Vars(r)["name"]
	provider := mux.Vars(r)["provider"]
	modulePayload := fmt.Sprintf("%s-%s-%s", namespace, name, provider)
	// Open our jsonFile
	jsonFile, err := os.Open(modulePayload + ".json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	w.Header().Set("Content-Type", "application/json")
	w.Write(byteValue)
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
