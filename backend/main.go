package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	nats "github.com/nats-io/nats.go"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type VM struct {
	Name       string `json:"Name"`
	Datacenter string `json:"Datacenter"`
	Folder     string `json:"Folder"`
	Secretkey  string `json:"Secretkey"`
}

var vms []VM

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "vAuth 0.0.1")
}

func createVMRecord(w http.ResponseWriter, r *http.Request) {
	VMName := mux.Vars(r)["name"]
	fmt.Println(VMName)
	var vm VM
	_ = json.NewDecoder(r.Body).Decode(&vm)
	json.NewEncoder(w).Encode(vm)
	// Connect to the database
	log.Print("Connecting to the database")
	db, err := sql.Open("postgres", "postgresql://root@db:26257/vauth?sslmode=disable")
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	defer db.Close()
	log.Print("Successfully connected to the database")

	fmt.Fprint(w, "POST done")

	// Insert two rows into the "vms" table.
	if _, err := db.Exec(
		"UPSERT INTO vms (name, datacenter, secretkey, folder) VALUES ('" + vm.Name + "' , '" + vm.Datacenter + "' , '" + vm.Secretkey + "' , '" + vm.Folder + "')"); err != nil {
		log.Fatal(err)
	}

}

func getVMRecord(w http.ResponseWriter, r *http.Request) {
	VMName := mux.Vars(r)["name"]
	log.Printf("The record for %s was requested", VMName)
	for _, vm := range vms {
		if vm.Name == VMName {
			json.NewEncoder(w).Encode(vm)
		}
	}
}

func deleteVMRecord(w http.ResponseWriter, r *http.Request) {
	VMName := mux.Vars(r)["name"]

	for _, vm := range vms {
		if vm.Name == VMName {
			json.NewEncoder(w).Encode(vm)
		}
	}
}

func getAllVMs(w http.ResponseWriter, r *http.Request) {
	// Connect to the database
	log.Print("Connecting to the database")
	db, err := sql.Open("postgres", "postgresql://root@db:26257/vauth?sslmode=disable")
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	defer db.Close()
	log.Print("Successfully connected to the database")
	rows, err := db.Query("SELECT name, datacenter, secretkey, folder FROM vms")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	vms = nil
	for rows.Next() {
		var name, datacenter, secretkey, folder string
		if err := rows.Scan(&name, &datacenter, &secretkey, &folder); err != nil {
			log.Fatal(err)
		}
		vms = append(vms, VM{Name: name, Datacenter: datacenter, Secretkey: secretkey, Folder: folder})
	}
	json.NewEncoder(w).Encode(vms)
	log.Print("User requested all VMs")
}

func syncVMs(w http.ResponseWriter, r *http.Request) {
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
	log.Print("Syncing VMs")
}

func main() {
	// Connect to the database
	log.Print("Connecting to the database")
	db, err := sql.Open("postgres", "postgresql://root@db:26257/vauth?sslmode=disable")
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	defer db.Close()
	log.Print("Successfully connected to the database")

	// Create the "vauth" database if it doesn't exist
	if _, err := db.Exec(
		"CREATE DATABASE IF NOT EXISTS vauth;"); err != nil {
		log.Fatal(err)
	}

	// Create the "vms" table if it doesn't exist
	if _, err := db.Exec(
		"CREATE TABLE IF NOT EXISTS vms (name CHAR (50) PRIMARY KEY, datacenter CHAR (50), secretkey VARCHAR (64), folder CHAR (50))"); err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/vms", getAllVMs).Methods("GET")
	router.HandleFunc("/sync", syncVMs).Methods("GET")
	router.HandleFunc("/vm/{name}", getVMRecord).Methods("GET")
	router.HandleFunc("/vm/{name}", createVMRecord).Methods("POST")
	log.Fatal(http.ListenAndServe(":8090", router))
}
