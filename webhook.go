package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	faktory "github.com/contribsys/faktory/client"
	"github.com/rs/zerolog/log"
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
		log.Info().Msg("webhook received data")
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fail(err)
		}

		msg := string(body)
		log.Debug().Msg(msg)
		// fmt.Fprint(w, string(body))

		job := faktory.NewJob("upsert_order", msg)
		cl.Push(job)
	})

	port := os.Getenv("WEBHOOKD_PORT")
	if port == "" {
		port = "9876"
	}

	log.Info().Msg(fmt.Sprintf("Handling HTTP requests on %s.", port))
	log.Fatal().Msg(http.ListenAndServe(fmt.Sprintf(":%s", port), nil).Error())
}
