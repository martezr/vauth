package main

import (
	//	"github.com/gorilla/mux"
	// "log"
	"net/http"
)

func main() {
	//	router := mux.NewRouter().StrictSlash(true)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)
	//	http.ListenAndServe(":3000", nil)
	http.ListenAndServeTLS(":443", "cert.pem", "key.pem", nil)
	//	router.Handle("/", fs)
	//	log.Fatal(http.ListenAndServeTLS(":443", "cert.pem",
	//		"key.pem", router))
}
