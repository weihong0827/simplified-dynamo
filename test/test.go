package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	puturl := "http://127.0.0.1:8080/put?key=foo&value=bar"
	putresp, err := http.NewRequest(http.MethodPut, puturl, nil)
	if err != nil {
		fmt.Print(err.Error())
	}

	client := &http.Client{}
	resp, err := client.Do(putresp)
	// log resp
	output, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}
	log.Println(output)
	log.Println(resp.Status)

	geturl := "http://127.0.0.1:8080/get?key=foo"
	getresp, err := http.Get(geturl)
	if err != nil {
		fmt.Print(err.Error())
	}

	getOutput, err := io.ReadAll(getresp.Body)
	if err != nil {
		fmt.Print(err.Error())
	}

	// log resp
	log.Println(getOutput)
	log.Println(getresp.Status)
}
