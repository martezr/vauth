package main

import (
	"context"
	"fmt"
	nats "github.com/nats-io/nats.go"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/event"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/vim25/types"
	"log"
	"net/url"
	"os"
	"reflect"
	"time"
)

func main() {
	time.Sleep(20 * time.Second)
	// Creating a connection context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	vsphereUsername := os.Getenv("VSPHERE_USERNAME")
	vspherePassword := os.Getenv("VSPHERE_PASSWORD")
	vsphereServer := os.Getenv("VSPHERE_SERVER")
	datacenter := os.Getenv("VSPHERE_DATACENTER")

	vcenterURL, err := url.Parse(fmt.Sprintf("https://%v/sdk", vsphereServer))
	if err != nil {
		log.Println(err)
	}
	credentials := url.UserPassword(vsphereUsername, vspherePassword)
	vcenterURL.User = credentials

	// Connecting to vCenter
	log.Print("Connecting to vCenter")
	client, err := govmomi.NewClient(ctx, vcenterURL, true)
	if err != nil {
		log.Println(err)
	}
	log.Printf("Connected to vCenter: %s", vsphereServer)
	finder := find.NewFinder(client.Client, true)

	dc, err := finder.DatacenterOrDefault(ctx, datacenter)

	finder.SetDatacenter(dc)
	refs := []types.ManagedObjectReference{dc.Reference()}

	// Setting up the event manager
	eventManager := event.NewManager(client.Client)
	err = eventManager.Events(ctx, refs, 10, true, false, handleEvent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
}

func handleEvent(ref types.ManagedObjectReference, events []types.BaseEvent) (err error) {
	// Connect to a server
	nc, err := nats.Connect("nats://nats:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	for _, event := range events {
		eventType := reflect.TypeOf(event).String()
		log.Printf("Event found of type %s\n", eventType)
		if eventType == "*types.VmPoweredOnEvent" {
			vmName := event.GetEvent().Vm.Name
			log.Printf("Detected power on event for %s", vmName)
			if err := nc.Publish("updates", []byte(vmName)); err != nil {
				log.Fatal(err)
			}
		}
	}
	return nil
}
