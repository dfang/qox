# QOR example application

This is an example application to show and explain features of [QOR](http://getqor.com).

Chat Room: [![Join the chat at https://gitter.im/qor/qor](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/qor/qor?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

## Quick Started

### Go version: 1.11+

````shell
# 1. Run docker-compose to setup database
$ docker-compose up -d

# 2. Config database.yml
> db:
    adapter: postgres
    name: qor_example
    user: adminer
    password: adminer
    host: 192.168.66.182


# 3. Config env for qiniu or s3 (config/config.go)
export QOR_QINIU_ACCESS_ID=<QOR_QINIU_ACCESS_ID>
export QOR_QINIU_ACCESS_KEY=<QOR_QINIU_ACCESS_KEY>
export QOR_QINIU_BUCKET=<QOR_QINIU_BUCKET>
export QOR_QINIU_REGION=<QOR_QINIU_REGION>
export QOR_QINIU_ENDPOINT=QOR_QINIU_ENDPOINT>

# 3. Seed

```go
$ go run config/db/seeds/main.go config/db/seeds/seeds.go
```

# 4. Run

```go
$ cd $GOPATH/src/github.com/dfang/qor-demo
$ go run main.go
```

## Admin Management Interface

[Qor Example admin configuration](https://github.com/dfang/qor-demo/blob/master/config/admin/admin.go)

Online Demo Website: [demo.getqor.com/admin](http://demo.getqor.com/admin)

## RESTful API

[Qor Example API configuration](https://github.com/dfang/qor-demo/blob/master/config/api/api.go)

Online Example APIs:

- Users: [http://demo.getqor.com/api/users.json](http://demo.getqor.com/api/users.json)
- User 1: [http://demo.getqor.com/api/users/1.json](http://demo.getqor.com/api/users/1.json)
- User 1's Orders [http://demo.getqor.com/api/users/1/orders.json](http://demo.getqor.com/api/users/1/orders.json)
- User 1's Order 1 [http://demo.getqor.com/api/users/1/orders/1.json](http://demo.getqor.com/api/users/1/orders/1.json)
- User 1's Orders 1's Items [http://demo.getqor.com/api/users/1/orders.json](http://demo.getqor.com/api/users/1/orders/1/items.json)
- Orders: [http://demo.getqor.com/api/orders.json](http://demo.getqor.com/api/orders.json)
- Products: [http://demo.getqor.com/api/products.json](http://demo.getqor.com/api/products.json)

## License

Released under the MIT License.

[@QORSDK](https://twitter.com/qorsdk)
```
````

## How to run this locally

```
go run main.go background_job.go health_check.go web_ui.go

or

gin

or

fresh
```

- [gin](https://github.com/codegangsta/gin)
- [fresh](https://github.com/gravityblast/fresh)
- [godev](https://github.com/zephinzer/godev)

## About health check

[liveness and readiness probe](https://github.com/heptiolabs/healthcheck#http-endpoints)

- /live: liveness endpoint (HTTP 200 if healthy, HTTP 503 if unhealthy)
- /ready: readiness endpoint (HTTP 200 if healthy, HTTP 503 if unhealthy)
Pass the ?full=1 query parameter to see the full check results as JSON. These are omitted by default for performance.




## Gocraft/work Web UI

```
workwebui -redis="redis:6379" -ns="qor" -listen=":5040"
```
