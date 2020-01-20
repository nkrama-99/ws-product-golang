package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

var count = 0

func main() {
	for true {
		time.Sleep(500 * time.Millisecond)
		MakeRequest()
	}
}

func MakeRequest() {

	fmt.Print("Request: ", count, "\n")

	resp, err := http.Get("http://127.0.0.1:8080/stats")
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(string(body))
	count++
}
