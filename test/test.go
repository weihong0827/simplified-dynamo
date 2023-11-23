package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	puturl := "http://127.0.0.1:8080/put?key=foo&value=bar"
	putresp, err := http.NewRequest("PUT", puturl, nil)
	if err != nil {
		fmt.Print(err.Error())
	}

	// log resp
	log.Println(putresp)
	

	geturl := "http://127.0.0.1:8080/get?key=foo"
    getresp, err := http.NewRequest("GET", geturl, nil)
    if err != nil {
        fmt.Print(err.Error())
    }

	// log resp
	log.Println(getresp)
}