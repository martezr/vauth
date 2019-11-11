package main

import (
	nats "github.com/nats-io/nats.go"
	"log"

	"github.com/jasonlvhit/gocron"
)

func main() {
	s := gocron.NewScheduler()
	var syncInterval uint64
	syncInterval = 60
	s.Every(syncInterval).Minutes().Do(sync)
	<-s.Start()
}

func sync() {
	// Connect to a server
	log.Println("Connecting to nats")
	nc, err := nats.Connect("nats://nats:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	log.Println("Connected to nats")
	if err := nc.Publish("sync", []byte("Sync VMs")); err != nil {
		log.Fatal(err)
	}

	log.Println("Syncing VMs")
}
