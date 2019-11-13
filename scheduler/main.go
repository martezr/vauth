package main

import (
	"github.com/jasonlvhit/gocron"
	nats "github.com/nats-io/nats.go"
	"log"
)

func main() {
	s := gocron.NewScheduler()
	s.Every(60).Minutes().Do(sync)
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
	log.Println("Scheduling VM synchronization process")
}
