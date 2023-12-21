# RSS FEED AGGREGATOR

This is a simple aggregator service that collects and organizes content from multiple rss feeds in one place, allowing users to subscribe to various sources such as blogs, news website or podcasts and combines themm into a simple feed.

This program is written in Go with postgres as the databsae, to run this program you need to have an .env file with the variables like so

```
PORT =
DB_URL = postgres://username:password@localhost:5432/dbname?sslmode=disable
```

to start the server run `make run` or `go run main.go` if you dont have make installed
