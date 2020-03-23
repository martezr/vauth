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
	Role       string `json:"Role"`
	Secretkey  string `json:"Secretkey"`
}

var vms []VM

func homeLink(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "vAuth 0.0.1")
}

func createVMRecord(w http.ResponseWriter, r *http.Request) {
	VMName := mux.Vars(r)["name"]
	log.Printf("Updating database record for %s", VMName)
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

	// Insert or update the VM info into the "vms" table.
	if _, err := db.Exec(
		"UPSERT INTO vms (name, datacenter, secretkey, role) VALUES ('" + vm.Name + "' , '" + vm.Datacenter + "' , '" + vm.Secretkey + "' , '" + vm.Role + "')"); err != nil {
		log.Fatal(err)
	}

}

func getVMRecord(w http.ResponseWriter, r *http.Request) {
	VMName := mux.Vars(r)["name"]
	log.Printf("The record for %s was requested", VMName)

	// Connect to the database
	log.Print("Connecting to the database")
	db, err := sql.Open("postgres", "postgresql://root@db:26257/vauth?sslmode=disable")
	if err != nil {
		log.Fatal("error connecting to the database: ", err)
	}
	defer db.Close()
	log.Print("Successfully connected to the database")
	rows, err := db.Query("SELECT name, datacenter, secretkey, role FROM vms WHERE name='" + VMName + "'")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var data VM
	for rows.Next() {
		var name, datacenter, secretkey, role string
		if err := rows.Scan(&name, &datacenter, &secretkey, &role); err != nil {
			log.Fatal(err)
		}
		data = VM{Name: name, Datacenter: datacenter, Secretkey: secretkey, Role: role}
	}
	json.NewEncoder(w).Encode(data)
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
	rows, err := db.Query("SELECT name, datacenter, secretkey, role FROM vms")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	vms = nil
	for rows.Next() {
		var name, datacenter, secretkey, role string
		if err := rows.Scan(&name, &datacenter, &secretkey, &role); err != nil {
			log.Fatal(err)
		}
		vms = append(vms, VM{Name: name, Datacenter: datacenter, Secretkey: secretkey, Role: role})
	}
	json.NewEncoder(w).Encode(vms)
	log.Print("User requested all VMs")
}

func syncVMs(w http.ResponseWriter, r *http.Request) {
	// Connect to NATS
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
	log.Print("Triggered VM sync")
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
		"CREATE TABLE IF NOT EXISTS vms (name CHAR (50) PRIMARY KEY, datacenter CHAR (50), secretkey VARCHAR (64), role CHAR (50))"); err != nil {
		log.Fatal(err)
	}

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", homeLink)
	router.HandleFunc("/vms", getAllVMs).Methods("GET")
	router.HandleFunc("/sync", syncVMs).Methods("GET")
	router.HandleFunc("/vm/{name}", getVMRecord).Methods("GET")
	router.HandleFunc("/vm/{name}", createVMRecord).Methods("POST")
	log.Fatal(http.ListenAndServeTLS(":443", "cert.pem",
		"key.pem", router))
}
