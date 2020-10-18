package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	faktory "github.com/contribsys/faktory/client"
)

// just run `go startWebhookd()` in main.go
// startWebhookd start webhookd
func startWebhookd() {
	// - FAKTORY_URL=tcp://:admin@faktory:7419
	if os.Getenv("FAKTORY_URL") == "" {
		panic("Please set FAKTORY_URL")
	}

	cl, err := faktory.Open()
	if err != nil {
		fail(err)
	}
	defer cl.Close()

	data, err := cl.Info()
	if err != nil {
		fail(err)
	}
	svr := data["server"].(map[string]interface{})
	fmt.Printf("Connected to %s %s\n", svr["description"], svr["faktory_version"])

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fail(err)
		}
		log.Println(string(body))
		fmt.Fprint(w, string(body))

		job := faktory.NewJob("upsert_order", string(body))
		cl.Push(job)
	})

	port := os.Getenv("WEBHOOKD_PORT")
	if port == "" {
		port = "9876"
	}

	log.Printf("Handling HTTP requests on %s.", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), nil))
}

func fail(err error) {
	fmt.Println(err.Error())
	os.Exit(-1)
}
