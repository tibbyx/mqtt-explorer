### Compile:
```bash
go build main.go
```

### Run:
```bash
go run main.go
```

### Or if you just compiled the code:
```bash
./main
```

### If you compiled on windows:
```powershell
.\main.exe
```

### To log in to the MQTT-Broker:

```bash
curl --request POST --header "Content-Type: application/json" --data '{"Ip" : "<BROKER-IP-HERE>", "Port" : "<BROKER-PORT-HERE>", "ClientId" : "<CLIENT-NAME-HERE>"}' localhost:3000/credentials
```

#### Or in other words, you need to POST into localhost:3000/credentials a JSON with this format:
```javascript
{
  "Ip" : "<BROKER-IP-HERE>",
  "Port" : "<BROKER-PORT-HERE>",
  "ClientId" : "<CLIENT-NAME-HERE>"
}
```

#### If the server cannot process the json, it will return a 400 (Bad request) with a JSON:
```javascript
{
  "badJson" : "I am nowt sowwy >:3. An expected! ewwow has happened. Youw weak json! iws of the wwongest fowmat thawt does nowt cowwespond tuwu the stwong awnd independent stwuct! >:P"
}
```

#### If the IP, Port or Client-Id are bad, the server will return a 400 (Bad request) with a JSON:
```javascript
{
  "badJson" : "IP is incomprehensible"
}
```

#### Or:
```javascript
{
  "badJson" : "PORT is incomprehensible"
}
```

#### Or:
```javascript
{
  "badJson" : "CLIENT-ID is incomprehensible"
}
```

#### If for some reason the server will not be able to connect to the MQTT-Broker (from several reasons, like the Broker IP or Port is invalid, or maybe the client-id is already in use), it will return a 400 (Bad request) with a JSON:
```javascript
{
  "badJson": "Connecting to <IP>:<PORT> failed\n<MQTT-ERROR>"
}
```

#### If everything will go well, the server will return a 200 (OK) with a JSON:

```javascript
{
  "goodJson" : "Connecting to <IP>:<PORT> succeded"
}
```

### To subscribe to a topic or multiple at once:

```bash
curl --request POST --header "Content-Type: application/json" --data '{"Topics":["<TOPIC-1>", "<TOPIC-2>", "<TOPIC-N>"]}' localhost:3000/topic/subscribe
```

#### Or in other words, you need to POST into localhost:3000/topic/subscribe a JSON with this format:

```javascript
{
  "Topics" :
  [
    "<TOPIC-1>",
    "<TOPIC-2>",
    "<TOPIC-N>"
  ]
}
```

#### If the client has not log in with the credentials, the server will return a 401 (Unauthorized) with a JSON:
```javascript
{
  "401" : "You fool!"
}
```

#### If the server cannot process the json, it will return a 400 (Bad request) with a JSON:
```javascript
{
  "badJson" : "I am nowt sowwy >:3. An expected! ewwow has happened. Youw weak json! iws of the wwongest fowmat thawt does nowt cowwespond tuwu the stwong awnd independent stwuct! >:P"
}
```

#### The server might not be able to subscribe to all topics, in this case it will return a 400 (Bad Request) (TODO: This status code is probably bad and needs to be replaced) with a JSON:
```javascript
{
  "badJson" : "Could not subscribe to these topics",
  "topics" : [
    "badTopic-1",
    "badTopic-2",
    "badTopic-N"
  ]
}
```

#### If everything will go well, the server will return a 200 (OK) with a JSON:
```javascript
{
  "goodJson" : "Subscribed to the topics"
}
```

### To unsubscribe to a topic or multiple at once:

```
curl --request POST --header "Content-Type: application/json" --data '{"Topics":["<TOPIC-1>", "<TOPIC-2>", "<TOPIC-N>"]}' localhost:3000/topic/unsubscribe
```

#### Or in other words, you need to POST into localhost:3000/topic/unsubscribe a JSON with this format:
```javascript
{
  "Topics" :
  [
    "<TOPIC-1>",
    "<TOPIC-2>",
    "<TOPIC-N>"
  ]
}
```

#### If the client has not log in with the credentials, the server will return a 401 (Unauthorized) with a JSON:
```javascript
{
  "401" : "You fool!"
}
```

#### If the client wants to unsubscribe from topics that weren't even subscribed, the server will return a 400 (Bad Request) (TODO: This status code is probably bad and needs to be replaced) with a JSON:
```javascript
{
  "badJson" : "Some topics were not even subscribed",
  "badTopics" :
  [
    "<TOPIC-1>",
    "<TOPIC-2>",
    "<TOPIC-N>"
  ],
  "unsubscribedTopics" : 
  [
    "<TOPIC-3>",
    "<TOPIC-4>",
    "<TOPIC-M>"
  ]
}
```

#### If eveything went well, the server will return a 200 (OK) with a JSON:

```javascript
{
  "goodJson" : "Unsubscribed from all"
}
```

### To get the subcribed list:

```bash
curl localhost:3000/topic/subscribed
```

#### Or in other words, you need to GET into localhost:3000/topic/unsubscribe
_There is no need for anything to send._

#### If the client has not log in with the credentials, the server will return a 401 (Unauthorized) with a JSON:
```javascript
{
  "401" : "You fool!"
}
```

#### If eveything went well, the server will return a 200 (OK) with a JSON:

```javascript
{
  "topics" :
  [
    "<TOPIC-1>",
    "<TOPIC-2>",
    "<TOPIC-n>",
  ]
}
```

### To send a message:

curl --request POST --header "Content-Type: application/json" --data '{"Topics" : ["<TOPIC-1>", "<TOPIC-2>", "<TOPIC-N>"]}' localhost:3000/topic/unsubscribe
```bash
curl --request POST --header "Content-Type: application/json" --data '{"Topic" : "<TOPIC>", "Message" : "<MESSAGE>"}' localhost:3000/topic/send-message
```

#### If the client has not log in with the credentials, the server will return a 401 (Unauthorized) with a JSON:
```javascript
{
  "401" : "You fool!"
}
```

#### If the server cannot process the json, it will return a 400 (Bad request) with a JSON:
```javascript
{
  "badJson" : "I am nowt sowwy >:3. An expected! ewwow has happened. Youw weak json! iws of the wwongest fowmat thawt does nowt cowwespond tuwu the stwong awnd independent stwuct! >:P"
}
```

#### If eveything went well, the server will return a 200 (OK) with a JSON:
```javascript
{
  "goodJson" : "Message posted"
}
```
