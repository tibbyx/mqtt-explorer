```
Currently there are four labs:

- old_server
    - gofiber as back-end and currently no frontend
- 0b1
    - fyne as both front-end and back-end.
    - I will probably abandon it though.
- 0b10
    - go-app as both front-end and back-end
    - it builds wasm and has the html + css flexibility
    - builds to web browsers
- 0b11
    - back-end          : gofiber
    - template engine   : django
    - front-end         : htmx
    - websocket         : yes
- 0b100
    - sse
    - I don't know how to make it work with htmx and I don't want to try React out.
- 0b101
    - Client shall ping the server for changes and if there is a change in state, it shall render it into the view.
    - This won't work with htmx. So React is the way.
    - For this to work smoothly, I will try to declare endpoints to the server.

Constants:
<IP>        : localhost
              127.0.0.1
                            
<PORT>      : 8080
<TABLE-CAP> : 500 // TODO: Talk about it with the comrades who pulled lists of anything from server.
                            
Terms:                      
<TOPIC>     : Definition:   <TOPIC> is an endpoint of the MQTT where clients can subscribe to and publish messages.
                            There can be any topic used. The working of the Topic is handled by the MQTT-Broker.
                            The <TOPIC> *shall* be mapped to all messages in order to sort them.
                            
<TABLE-ID>  : Definition:   The <TABLE-ID> *shall* be an ID, which is a number, matched to <MESSAGE>s.
                            If the <TABLE-ID> is 0, it shall be matched from the 0*<TABLE-CAP>+1 message up to the (0+1)*<TABLE-CAP> message.
                            If the <TABLE-ID> is 1, it shall be matched from the 1*<TABLE-CAP>+1 message up to the (1+1)*<TABLE-CAP> message.
                            If the <TABLE-ID> is n, it shall be matched from the n*<TABLE-CAP>+1 message up to the (n+1)*<TABLE-CAP> message.
                            
<MESSAGE>   : Definition:   The <MESSAGE> *shall* be a string that will be transmitted through the MQTT-Protocol to
                            MQTT-Broker.

JSON/RPC Endpoints:
    -   <IP>:<PORT>/get-history/<TOPIC>/<TABLE-ID>
        - The GET Method *shall* return matched to the <TOPIC> a list of <MESSAGE>s from other and self MQTT-Clients.
        - The GET Method *shall* return 404 if the <TOPIC> does not exist in the database.

    -   <IP>:<PORT>/get-new-messages/<TOPIC>
        - The GET Method *shall* return matched to the <TOPIC> a list of <MESSAGE>s from other MQTT-Clients that the <USER> has not yet seen.
        - The GET Method *shall* return 404 if the <TOPIC> does not exist in the database.
        - The functionality will be as such (the implementation may vary):
            - The User has subscribed to a <TOPIC>.
            - It subscribes to the <TOPIC> by the method "<IP>:<PORT>/get-history/<TOPIC>/<TABLE-ID>".
            - The server *shall* keep track of the last subscribed <TOPIC> matched to the <USER> and or <CLIENT-ID>
            - The server *shall* also keep track of the unix time of the subscribtion
            - When calling this method, the server *shall* query the message since the unix time of the subscribed <TOPIC>
                - This will always ensure that only the new messages will be queried
            - The server shall update the unix time of the subscription with the latest unix time from the query
            - The server returns the messages

    -   <IP>:<PORT>/post-credentials
        - The POST Method *shall* be used to post a <CLIENT-ID> for the MQTT-Client.
        - MQTT-Client *must* have a ClientID so that the MQTT-Broker to differenciate the MQTT-Clients (probably).
        - The <CLIENT-ID> *shall also* be used to associate the messages published through the MQTT protocol with the the <USER>.
        - The Method *shall* return an authentication cookie which will be the <CLIENT-ID>

    -   <IP>:<PORT>/post-message/<TOPIC>
        - The POST method *shall* write the <MESSAGE> to the database

The end functionality will resemble a chat server system where the information is not stored in a server, but rather in clients.
This has a side effect that if one client were to disconnect, it will not receive posted messages and will not be able to update them.
Either the client must be active at all times, or some sort of a database synchronisation mechanism must be implemented.

// TODO: Figure out if it is possible for ONE client to subscribe to MULTIPLE topics.

JSON Structures are to be defined // TODO: Finish this, it's almost 23:00 which is eepy time.
```
