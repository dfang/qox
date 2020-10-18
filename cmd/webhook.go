package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	faktory "github.com/contribsys/faktory/client"
	"github.com/rs/zerolog/log"
)

// StartWebhookd start webhookd
// just run `go StartWebhookd()` in main.go
func StartWebhookd() {
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
	log.Info().Msgf("Connected to %s %s", svr["description"], svr["faktory_version"])

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Info().Msg("webhook received a request")
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fail(err)
		}

		msg := string(body)
		log.Debug().Msg(msg)

		job := faktory.NewJob("upsert_order", msg)
		if err := cl.Push(job); err != nil {
			log.Info().Msgf("push job to faktory failed")
		}

		log.Info().Msgf("pushed a job with job id: %s, job type: %s, queue: %s", job.Jid, job.Type, job.Queue)
		fmt.Fprintln(w, "ok")
	})

	port := os.Getenv("WEBHOOKD_PORT")
	if port == "" {
		port = "9876"
	}

	log.Info().Msg(fmt.Sprintf("Webhook listens on %s.", port))
	log.Fatal().Msg(http.ListenAndServe(fmt.Sprintf(":%s", port), nil).Error())
}
