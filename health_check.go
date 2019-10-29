package main

import (
	"net/http"
	"time"

	"github.com/dfang/qor-demo/config/db"
	"github.com/heptiolabs/healthcheck"
)

// just run go startHealthCheck() in main.go
func startHealthCheck() {
	health := healthcheck.NewHandler()
	// Our app is not happy if we've got more than 100 goroutines running.
	health.AddLivenessCheck("goroutine-threshold", healthcheck.GoroutineCountCheck(100))
	// Our app is not ready if we can't resolve our upstream dependency in DNS.
	// health.AddReadinessCheck(
	// 	"upstream-dep-dns",
	// 	healthcheck.DNSResolveCheck("upstream.example.com", 50*time.Millisecond))

	// Our app is not ready if we can't connect to our database (`var db *sql.DB`) in <1s.
	health.AddReadinessCheck("database", healthcheck.DatabasePingCheck(db.DB.DB(), 5*time.Second))

	// go http.ListenAndServe("0.0.0.0:8086", health)
	http.ListenAndServe("0.0.0.0:8086", health)
}
