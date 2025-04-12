# mqtt-explorer

Compile and run:

```bash
go build main.go
./main
```

Or... like this if you want.

```bash
go run main.go
```

Examples:

```bash
curl -X POST localhost:3000/message/<MESSAGE>
> You fool, you utter buffoon.

curl -X POST localhost:3000/ip/<IP>
curl -X POST localhost:3000/message/<MESSAGE>
> You fool, you utter buffoon.

curl -X POST localhost:3000/reset
curl -X POST localhost:3000/topic/<TOPIC>
curl -X POST localhost:3000/message/<MESSAGE>
> You fool, you utter buffoon.

curl -X POST localhost:3000/ip/<IP>
curl -X POST localhost:3000/topic/<TOPIC>
curl -X POST localhost:3000/message/<MESSAGE>
> IP: <IP> Topic: <TOPIC> MSG: <MESSAGE>
```
