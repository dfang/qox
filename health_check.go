package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dfang/qor-demo/config"
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

	// check postgres
	// Our app is not ready if we can't connect to our database (`var db *sql.DB`) in <5s.
	health.AddReadinessCheck("postgres", healthcheck.DatabasePingCheck(db.DB.DB(), 5*time.Second))

	// check redis
	health.AddReadinessCheck("redis", healthcheck.TCPDialCheck(fmt.Sprintf("%s:%s", config.Config.Redis.Host, config.Config.Redis.Port), 5*time.Second))

	// check worker
	health.AddReadinessCheck("worker", healthcheck.HTTPGetCheck("http://localhost:5040/worker_pools", 5*time.Second))
	health.AddLivenessCheck("worker", healthcheck.HTTPGetCheck("http://localhost:5040/worker_pools", 5*time.Second))

	//check web is live
	health.AddLivenessCheck("web", healthcheck.HTTPGetCheck("http://localhost:7000/health", 5*time.Second))

	// go http.ListenAndServe("0.0.0.0:8086", health)
	http.ListenAndServe("0.0.0.0:8086", health)
}
