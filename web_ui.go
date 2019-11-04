package main

import (
	"time"

	"github.com/dfang/qor-demo/config/db"

	"github.com/gocraft/work/webui"
	"github.com/gomodule/redigo/redis"
)

var (
// redisHostPort  = flag.String("redis", "redis:6379", "redis hostport")
// redisDatabase  = flag.String("database", "0", "redis database")
// redisNamespace = flag.String("ns", "work", "redis namespace")
// webHostPort    = flag.String("listen", ":5040", "hostport to listen for HTTP JSON API")
)

// startWorkWebUI serves gocraft/work UI
// https://github.com/gocraft/work/blob/master/cmd/workwebui/main.go
func startWorkWebUI() {
	// flag.Parse()

	// redisDatabase = flag.String("database", "0", "redis database")

	// database, err := strconv.Atoi("0")
	// if err != nil {
	// 	fmt.Printf("Error: %v is not a valid database value", 0)
	// 	return
	// }

	// pool := newPool(*redisHostPort, database)
	// pool := newPool("redis:6379", database)

	// server := webui.NewServer("qor", pool, ":5040")
	server := webui.NewServer("qor", db.RedisPool, ":5040")
	server.Start()

	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt, os.Kill)

	// <-c

	// server.Stop()

	// fmt.Println("\nQuitting...")

}

func newPool(addr string, database int) *redis.Pool {
	return &redis.Pool{
		MaxActive:   3,
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.DialURL(addr, redis.DialDatabase(database))
		},
		Wait: true,
	}
}
