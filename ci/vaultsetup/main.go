package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/hashicorp/vault/api"
)

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

const vaultToken = "vault"

func main() {

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	httpClient.Transport = tr

	client, err := api.NewClient(&api.Config{Address: "http://vault:8200", HttpClient: httpClient})
	if err != nil {
		panic(err)
	}
	client.SetToken(vaultToken)

	time.Sleep(10 * time.Second)

	// Enable approle authentication method
	client.Sys().EnableAuth("approle", "approle", "vauth approle backend")

	authMethods, authmethoderr := client.Sys().ListAuth()
	if authmethoderr != nil {
		panic(authmethoderr)
	}
	fmt.Println(authMethods)
	roles := [5]string{"app01", "app02", "app03", "app04", "app05"}
	for _, role := range roles {
		log.Println(role)
		createRole(client, role)
	}
}

func createRole(client *api.Client, rolename string) {
	options := map[string]interface{}{
		"secret_id_ttl":  "10m",
		"token_num_uses": 10,
		"token_ttl":      "60m",
	}

	_, newerr := client.Logical().Write("auth/approle/role/"+rolename, options)
	if newerr != nil {
		panic(newerr)
	}
}
