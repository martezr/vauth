package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	nats "github.com/nats-io/nats.go"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/vim25/mo"
	"log"
	"net/url"
	"strings"
)

func main() {
	time.Sleep(10 * time.Second)
	// Connect to a server
	nc, err := nats.Connect("nats://nats:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Use a WaitGroup to wait for a message to arrive
	wg := sync.WaitGroup{}
	wg.Add(1)

	// Subscribe to the "sync" topic
	if _, err := nc.Subscribe("sync", func(m *nats.Msg) {

		// Run the syncVM function to fetch a list of VMs to sync
		vms := syncVM()
		for _, name := range vms {
			if err := nc.Publish("updates", []byte(name)); err != nil {
				log.Fatal(err)
			}
		}
	}); err != nil {
		log.Fatal(err)
	}
	// Wait for a message to come in
	wg.Wait()
}

func syncVM() (vmnames []string) {
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
	if err != nil {
		fmt.Println(err)
	}

	finder.SetDatacenter(vsphereDatacenter)

	machines, err := finder.VirtualMachineList(ctx, "*")
	if err != nil {
		fmt.Println(err)
	}

	for _, vmdata := range machines {
		machine := machines[0]
		var props mo.VirtualMachine
		machine.Properties(ctx, vmdata.Reference(), nil, &props)
		var vmconfig []string
		// Evaluate if a virtual machine object is not a template and is powered on
		if props.Summary.Config.Template == false && props.Summary.Runtime.PowerState == "poweredOn" {
			for _, v := range props.Config.ExtraConfig {
				if strings.HasPrefix(v.GetOptionValue().Key, "guestinfo.vault.") {
					vmconfig = append(vmconfig, v.GetOptionValue().Key)
				}
			}
			if len(vmconfig) < 1 {
				log.Printf("%s is missing identity data and will be synced", vmdata.Name())
				vmnames = append(vmnames, vmdata.Name())
			}
		} else if props.Summary.Config.Template == true {
			log.Printf("%s has been marked as a template and will be skipped", vmdata.Name())
		} else if props.Summary.Runtime.PowerState == "poweredOff" {
			log.Printf("%s has a power state of %s and will be skipped", vmdata.Name(), props.Summary.Runtime.PowerState)
		}
	}
	return vmnames
}
