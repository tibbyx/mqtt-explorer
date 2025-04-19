# mqtt-explorer

This is an mqtt client. It connects to the mqtt broker.

Compile and run:

```bash
go build server.go
./server
```

Or... like this if you want.

```bash
go run server.go
```

Examples:

```bash
curl --header "Content-Type: application/json"  --request POST  --data '{"Ip":"localhost","ClientId":"Jotaro","Topic":"mlem/cat"}' localhost:3000/mqtt/connect
> Connected to the Json config

curl -X POST localhost:3000/mqtt/subscribe
> Subscribed

curl -X POST localhost:3000/mqtt/message/<MESSAGE>
>

curl -X POST localhost:3000/mqtt/disconnect
> Disconnected
```

The /mqtt/message/<MESSAGE> post the message to the broker and the clients that are listening to it will receive the message.
