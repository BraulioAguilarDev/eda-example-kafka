# Event-Driven Architecture with Go and Kafka

## Tools
- [Gin](https://gin-gonic.com)
- [ZooKeeper](https://zookeeper.apache.org/)
- [Kafka](https://kafka.apache.org/)

## Test code
```sh
$ git clone repo

# Start Kafka & zookeeper
$ docker-compose up

# Start producer REST
$ go run producer/producer.go

# Start consumer
$ go run worker/worker.go

# Send request
$ curl --location --request POST '0.0.0.0:8080/api/v1/comments' --header 'Content-Type: application/json' --data-raw '{"text":"my message"}'
```
