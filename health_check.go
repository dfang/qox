package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/dfang/qor-demo/config"
	"github.com/dfang/qor-demo/config/db"
	"github.com/heptiolabs/healthcheck"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// just run go startHealthCheck() in main.go
func startHealthCheck() {
	registry := prometheus.NewRegistry()
	health := healthcheck.NewMetricsHandler(registry, "qor")

	// Our app is not happy if we've got more than 100 goroutines running.
	health.AddLivenessCheck("goroutine-threshold", healthcheck.GoroutineCountCheck(100))
	// Our app is not ready if we can't resolve our upstream dependency in DNS.
	// health.AddReadinessCheck(
	// 	"upstream-dep-dns",
	// 	healthcheck.DNSResolveCheck("upstream.example.com", 50*time.Millisecond))

	// check postgres
	// Our app is not ready if we can't connect to our database (`var db *sql.DB`) in <5s.
	health.AddLivenessCheck("postgres", healthcheck.DatabasePingCheck(db.DB.DB(), 5*time.Second))

	// check redis
	health.AddLivenessCheck("redis", healthcheck.TCPDialCheck(fmt.Sprintf("%s:%s", config.Config.Redis.Host, config.Config.Redis.Port), 5*time.Second))

	// check worker
	health.AddLivenessCheck("worker-ui", healthcheck.HTTPGetCheck("http://localhost:5040/worker_pools", 5*time.Second))

	health.AddLivenessCheck("faktory", healthcheck.HTTPGetCheck("http://localhost:7420", 5*time.Second))

	//check web is live
	health.AddLivenessCheck("web", healthcheck.HTTPGetCheck("http://localhost:7000/health", 5*time.Second))

	adminMux := http.NewServeMux()
	adminMux.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	// Expose a liveness check on /live
	adminMux.HandleFunc("/live", health.LiveEndpoint)

	// Expose a readiness check on /ready
	adminMux.HandleFunc("/ready", health.ReadyEndpoint)
	go http.ListenAndServe("0.0.0.0:9402", adminMux)

	// Expose the /live and /ready endpoints over HTTP (on port 8086):
	// go http.ListenAndServe("0.0.0.0:8086", health)
	http.ListenAndServe("0.0.0.0:8086", health)
}
