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

### To disconnect from the MQTT-Broker:
```bash
curl -X POST localhost:3000/disconnect
```

#### If the go-server wasn't connected to any MQTT-Broker, the server will return a 400 (Bad Request) with a JSON:
```javascript
{
  "BadRequest": "The server isn't even connected to any MQTT-Brokers"
}
```

#### If everything went well, the server will return a 200 (OK) with a JSON:
```javascript
{
  "fine" : "The MQTT-Client disconented from <IP>:<PORT> Broker"
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

#### The server might not be able to subscribe to all topics, in this case it will return a 207 (Multi Status) with a JSON:
```javascript
{
  "result" :
  {
    <TOPIC-1> :
    {
      "Status" : <STATUS-1>,
      "Message" : <MESSAGE-1>
    },
    <TOPIC-2> :
    {
      "Status" : <STATUS-2>,
      "Message" : <MESSAGE-2>
    },
    <TOPIC-N> :
    {
      "Status" : <STATUS-N>,
      "Message" : <MESSAGE-N>
    }
  }
}
```

#### If everything will go well, the server will return a 200 (OK) with a JSON:
```javascript
{
  "result" :
  {
    <TOPIC-1> :
    {
      "Status" : "Fine",
      "Message" : "Subscribed to the topic"
    },
    <TOPIC-2> :
    {
      "Status" : "Fine",
      "Message" : "Subscribed to the topic"
    },
    <TOPIC-N> :
    {
      "Status" : "Fine",
      "Message" : "Subscribed to the topic"
    }
  }
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

#### If the server cannot process the json, it will return a 400 (Bad request) with a JSON:
```javascript
{
  "badJson" : "I am nowt sowwy >:3. An expected! ewwow has happened. Youw weak json! iws of the wwongest fowmat thawt does nowt cowwespond tuwu the stwong awnd independent stwuct! >:P"
}
```

#### If the client wants to unsubscribe from topics that weren't even subscribed, the server will return a 207 (Multi Status) with a JSON:
```javascript
{
  "result" :
  {
    <TOPIC-1> :
    {
      "Status" : <STATUS-1>,
      "Message" : <MESSAGE-1>
    },
    <TOPIC-2> :
    {
      "Status" : <STATUS-2>,
      "Message" : <MESSAGE-2>
    },
    <TOPIC-N> :
    {
      "Status" : <STATUS-N>,
      "Message" : <MESSAGE-N>
    }
  }
}
```

#### If eveything went well, the server will return a 200 (OK) with a JSON:
```javascript
{
  "result" :
  {
    <TOPIC-1> :
    {
      "Status" : "Fine",
      "Message" : "Unsubscribed successfully"
    },
    <TOPIC-2> :
    {
      "Status" : "Fine",
      "Message" : "Unsubscribed successfully"
    },
    <TOPIC-N> :
    {
      "Status" : "Fine",
      "Message" : "Unsubscribed successfully"
    }
  }
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
    "<TOPIC-n>"
  ]
}
```

### To get all known topics:
```bash
curl localhost:3000/topic/all-known
```

#### If the client is not authenticated yet, the server will return a 401 (Unauthorized) with a JSON:
```javascript
{
  "Message": "Authenticate yourself first!"
}
```

#### If everything went well, the server will return a 200 (OK) with a JSON:
```javascript
{
  "Topics" : ["<TOPIC-1>", "<TOPIC-2>", "<TOPIC-N>"]
}
```

### To send a message:
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

### To get the messages matched to a topic:
```bash
curl localhost:3000/topic/messages?topic=<TOPIC>
```

#### If the <TOPIC> is empty, the server will return a 400 (Bad Request) with a JSON:
```javascript
{
  "error": "Missing topic query parameter"
}
```

#### If everything went well, the server will return a 200 (OK) with a JSON:
```javascript
{
  "topic": <TOPIC>,
  "messages": [<MESSAGE-1>, <MESSAGE-2>, <MESSAGE-N>]
}
```

### To check if the go server is still connected to the MQTT-Broker:
```bash
curl localhost:3000/ping
```

#### Or in other words, you need to GET into localhost:3000/ping
_There is no need for anything to send._

#### If the client has not log in with the credentials, the server will return a 401 (Unauthorized) with a JSON:
```javascript
{
  "401" : "You fool!"
}
```

#### If the go server does not retrieve a response from the MQTT-Broker, it will return a 503 (Service Unavailable) with a JSON:
```javascript
{
  "badMqtt" : "Big f!"
}
```

#### If everything went well, the server will return a 200 (OK) with a JSON:
```javascript
{
  "goodMqtt" : "pong"
}
```
