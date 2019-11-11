package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	nats "github.com/nats-io/nats.go"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"fmt"
	"net/url"

	"context"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

type vm struct {
	Name       string `json:"Name"`
	Datacenter string `json:"Datacenter"`
	Folder     string `json:"Folder"`
	Secretkey  string `json:"Secretkey"`
}

type extraConfig []types.BaseOptionValue

func main() {
	time.Sleep(10 * time.Second)
	// Connect to a server
	log.Println("Connecting to nats")
	nc, err := nats.Connect("nats://nats:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	log.Println("Connected to nats")

	// Use a WaitGroup to wait for a message to arrive
	wg := sync.WaitGroup{}
	wg.Add(1)

	// Subscribe
	if _, err := nc.Subscribe("updates", func(m *nats.Msg) {
		log.Printf("Tagging %s", string(m.Data))
		updateVM(string(m.Data))

	}); err != nil {
		log.Fatal(err)
	}

	// Wait for a message to come in
	wg.Wait()
}

func updateVM(name string) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	vsphereUsername := os.Getenv("VSPHERE_USERNAME")
	vspherePassword := os.Getenv("VSPHERE_PASSWORD")
	vsphereServer := os.Getenv("VSPHERE_SERVER")
	datacenter := os.Getenv("VSPHERE_DATACENTER")

	vcenterURL, err := url.Parse(fmt.Sprintf("https://%v/sdk", vsphereServer))
	if err != nil {
		fmt.Println(err)
	}
	credentials := url.UserPassword(vsphereUsername, vspherePassword)
	vcenterURL.User = credentials

	// Connecting to vCenter
	log.Print("Connecting to vCenter")
	client, err := govmomi.NewClient(ctx, vcenterURL, true)
	if err != nil {
		fmt.Println(err)
	}
	log.Printf("Connected to vCenter: %s", vsphereServer)
	finder := find.NewFinder(client.Client, true)

	vsphereDatacenter, err := finder.DatacenterOrDefault(ctx, datacenter)

	finder.SetDatacenter(vsphereDatacenter)

	machines, err := finder.VirtualMachineList(ctx, name)
	if err != nil {
		fmt.Println(err)
	}
	var vmsdata []vm

	for _, vmdata := range machines {
		machine := machines[0]
		var props mo.VirtualMachine
		machine.Properties(ctx, vmdata.Reference(), nil, &props)
		if props.Summary.Config.Template == false && props.Summary.Runtime.PowerState == "poweredOn" {
			inventoryPath := strings.Split(vmdata.InventoryPath, "/")
			folder := inventoryPath[len(inventoryPath)-2]
			datacenter := inventoryPath[1]
			token, _ := GenerateRandomString(64)

			name := vmdata.Name()
			vmsdata = append(vmsdata, vm{Name: name, Folder: folder, Datacenter: datacenter, Secretkey: token})

			var settings extraConfig
			settings = append(settings, &types.OptionValue{Key: "guestinfo.vault.vmname", Value: name})
			settings = append(settings, &types.OptionValue{Key: "guestinfo.vault.folder", Value: folder})
			settings = append(settings, &types.OptionValue{Key: "guestinfo.vault.datacenter", Value: datacenter})
			settings = append(settings, &types.OptionValue{Key: "guestinfo.vault.secretkey", Value: token})
			authSpec := types.VirtualMachineConfigSpec{
				ExtraConfig: settings,
			}
			vmdata.Reconfigure(ctx, authSpec)
			log.Printf("Updated VM: %s", name)

			// Post JSON payload
			payload, err := json.Marshal(vm{Name: name, Folder: folder, Datacenter: datacenter, Secretkey: token})
			if err != nil {
				log.Fatalln(err)
			}
			url := "http://backend:8090/vm/" + name
			resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
			if err != nil {
				log.Fatalln(err)
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatalln(err)
			}
			log.Println(string(body))
			log.Println("Updated Database VM Entry")

		} else {
			log.Printf("Skipped %s due to power state %s", vmdata.Name(), props.Summary.Runtime.PowerState)
		}
	}
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

func GenerateRandomString(n int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-/+="
	bytes, err := GenerateRandomBytes(n)
	if err != nil {
		return "", err
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes), nil
}
