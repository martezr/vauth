package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"time"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

func main() {
	// Creating a connection context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	log.Println("sleeping for 20 seconds...")
	time.Sleep(20 * time.Second)
	vcenterURL, err := url.Parse(fmt.Sprintf("https://%v/sdk", "vcsim"))
	if err != nil {
		log.Println(err)
	}
	credentials := url.UserPassword("user", "pass")
	vcenterURL.User = credentials

	// Connecting to vCenter
	log.Print("connecting to vCenter server")

	vsphereClient, err := govmomi.NewClient(ctx, vcenterURL, true)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("connected to vCenter: %s", "vcsim")
	finder := find.NewFinder(vsphereClient.Client, true)

	log.Println("adding vauth custom attribute")
	m, err := object.GetCustomFieldsManager(vsphereClient.Client)
	if err != nil {
		log.Panicln(err)
	}

	def, err := m.Add(ctx, "vauth-role", "VirtualMachine", nil, nil)
	if err != nil {
		log.Panicln(err)
	}
	attID, _ := m.FindKey(ctx, "vauth-role")

	log.Printf("%s - %d - %s\n", def.Name, def.Key, def.Type)

	datacenters := []string{"DC0", "DC1", "DC2"}
	for {
		for _, datacenter := range datacenters {
			// Iterate through dcs
			dc, err := finder.DatacenterOrDefault(ctx, datacenter)
			if err != nil {
				log.Println(err)
			}
			log.Println(dc)

			finder.SetDatacenter(dc)
			machines, err := finder.VirtualMachineList(ctx, "*")
			if err != nil {
				log.Println(err)
			}
			for _, vm := range machines {
				//			vm.PowerOff(ctx)
				vm.Reset(ctx)
				approles := []string{"app01", "app02", "app03", "app04", "app05"}
				randomIndex := rand.Intn(len(approles))
				pick := approles[randomIndex]

				m.Set(ctx, vm.Reference(), attID, pick)
				var props mo.VirtualMachine
				vm.Properties(ctx, vm.Reference(), nil, &props)
				for _, fv := range props.CustomValue {
					value := fv.(*types.CustomFieldStringValue).Value
					if value != "" {
						customAttrs := make(map[string]interface{})
						customAttrs[fmt.Sprint(fv.GetCustomFieldValue().Key)] = value
					}
					vm.Name()
					log.Printf("name: %s role: %s", vm.Name(), value)
				}
				time.Sleep(2 * time.Second)
			}
		}
	}
}
