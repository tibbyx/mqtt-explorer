package main

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"slices"
	"strings"
	"github.com/valyala/fasthttp"
	"time"
	"bufio"
)

// # Author
// - Polariusz
const BADJSON = "I am nowt sowwy >:3. An expected! ewwow has happened. Youw weak json! iws of the wwongest fowmat thawt does nowt cowwespond tuwu the stwong awnd independent stwuct! >:P"

// | Date of change | By        | Comment |
// +----------------+-----------+---------+
// | 2025-05-13     | Polariusz | Created |
//
// # Description
// - This method shall build a message containing the Client-ID and the Message that the messagePubHandler will be able to then differentiate and write into the database accordingly.
//
// # Author
// - Polariusz
func messageBuilder(creds MqttCredentials, message string) string {
	// TODO: We need to think of a structure for categorising messages with the clientIds (Usernames).
	return message
}

// | Date of change | By        | Comment              |
// +----------------+-----------+----------------------+
// |                | Polariusz | Created              |
// | 2025-05-13     | Polariusz | Documentation        |
// | 2025-05-16     | Polariusz | added AllKnownTopics |
//
// # Description
//
// - The structure allows for all Handlers to have a common state.
//   - In this project this is fine as only one client shall have one server.
//     - The client functions as the frontend gui for the go-server.
//
// # Used in
// - All function handlers.
//
// # Author
// - Polariusz
type ServerState struct {
	userCreds MqttCredentials
	mqttClient mqtt.Client
	subscribedTopics []string // TODO: This probably has to be a struct array of a pair, a pair of topic and epoch time.
	allKnownTopics []string // TODO: This probably has to be a struct array of a pair, a pair of topic and epoch time.
	receivedMessages map[string][]string // TODO: This should be in a database that we don't have yet.
}

// | Date of change | By        | Comment       |
// +----------------+-----------+---------------+
// |                | Polariusz | Created       |
// | 2025-05-13     | Polariusz | Documentation |
//
// # Structure:
// - {"Ip":<I>,"Port":"<P>","ClientId":<C>}
//   - <I>: The IP of the MQTT-Broker.
//   - <P>: The Port that the MQTT-Broker opened for the protocol.
//   - <C>: Client ID that functions as an username. It makes the users distinct.
//
// # Used in
// - struct ServerState
//   - Therefore it is in scope in all function handlers.
// - PostCredentialsHandler()
// - validateCredentials()
//
// # Author
// - Polariusz
type MqttCredentials struct {
	Ip string
	Port string // TODO: It would be probably nice to store it as a numeric.
	ClientId string
}

// # Author
// - Polariusz
func (mc MqttCredentials) dump() {
	fmt.Printf("ip       : %s", mc.Ip)
	fmt.Printf("port     : %s", mc.Port)
	fmt.Printf("clientId : %s", mc.ClientId)
}

// # Author
// - Polariusz
func main() {
	server := fiber.New()
	var serverState ServerState
	serverState.receivedMessages = make(map[string][]string)

	addRoutes(server, &serverState)

	// need to build ui via 'npm run build' in client first
	// server.Static("/", "../client/dist")
	server.Listen(":3000")
}

// | Date of change | By        | Comment       |
// +----------------+-----------+---------------+
// |                | Polariusz | Created       |
// | 2025-05-13     | Polariusz | Documentation |
//
// # Method-Type
// - Routing
//
// # Description
// - The method shall assign URLs to function handlers.
//
// # Author
// - Polariusz
func addRoutes(server *fiber.App, serverState *ServerState) {
	server.Post("/credentials", PostCredentialsHandler(serverState))
	server.Post("/disconnect", PostDisconnectFromBrokerHandler(serverState))
	server.Post("/topic/subscribe", PostTopicSubscribeHandler(serverState))
	server.Post("/topic/unsubscribe", PostTopicUnsubscribeHandler(serverState))
	server.Get("/topic/subscribed", GetTopicSubscribedHandler(serverState))
	server.Post("/topic/send-message", PostTopicSendMessageHandler(serverState))
	server.Get("/topic/messages", GetTopicMessagesHandler(serverState))
	server.Get("/ping", GetPingHandler(serverState))
	server.Post("/write", writeStuff())
	server.Get("/sse", SseHandler())
	server.Get("/topic/all-known", GetTopicAllKnownHandler(serverState))
}

// | Date of change | By        | Comment       |
// +----------------+-----------+---------------+
// |                | Polariusz | Created       |
// | 2025-05-13     | Polariusz | Documentation |
//
// # Method-Type
// - Handler
//
// # Description
// - The method shall be a handler that allows to authenticate the user to the MQTT-Broker as the protocol must have a ClientId. 
// - The method shall accept a jsonified structure that follows the struct MqttCredentials.
// - The method shall return a 200 (Ok) if credentials are valid and the connection with the MQTT-Broker was estabilished.
// - The method shall return a 400 (Bad Request) if the data from the client does not match that one of the struct MqttCredentials.
// - The method shall return a 404 (Service Unavailable) if the connection to the MQTT-Broker failed.
//
// # Usage
// - Call declared by the routing method addRoutes() URL with the POST-Method.
// - The data must have a json structure that matches the struct MqttCredentials.
//
// # Returns
// - 200 (Ok): JSON
//   - {"goodJson":"Connecting to `Ip`:`Port` succeded"}
// - 400 (Bad Request): JSON
//   - {"badJson":`const BADJSON`}
//   - {"badJson":`errorMessage`}
// - 404 (Not Found): JSON
//   - {"badMqtt":"Connecting to `Ip`:`Port` failed"}
//
// # Author
// - Polariusz
func PostCredentialsHandler(serverState *ServerState) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var userCreds MqttCredentials

		if err := c.BodyParser(&userCreds); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"badJson": BADJSON,
			})
		}

		{
			errorMessage := ""
			if validateCredentials(&errorMessage, &userCreds) != 0 {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"badJson": errorMessage,
				})
			}
		}

		// test.mosquitto.org
		mqttOpts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s:%s", userCreds.Ip, userCreds.Port)).SetClientID(userCreds.ClientId)
		mqttOpts.SetKeepAlive(2 * time.Second)
		mqttOpts.SetPingTimeout(1 * time.Second)

		mqttOpts.SetDefaultPublishHandler(createMessageHandler(serverState))

		mqttClient := mqtt.NewClient(mqttOpts)

		if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"badJson": fmt.Sprintf("Connecting to %s:%s failed\n%s", userCreds.Ip, userCreds.Port, token.Error()),
			})
		}

		serverState.userCreds = userCreds
		serverState.mqttClient = mqttClient

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"goodJson": fmt.Sprintf("Connecting to %s:%s succeded", userCreds.Ip, userCreds.Port),
		})
	}
}

// | Date of change | By        | Comment       |
// +----------------+-----------+---------------+
// |                | Polariusz | Created       |
// | 2025-05-13     | Polariusz | Documentation |
//
// # Method-Type
// - Validator
//
// # Description
// - The method shall validate the argument `userCreds *MqttCredentials`.
//   - TODO: The validation need to be improved. Right now it only checks if the argument userCreds has empty strings.
//
// # Usage
// - Call the method with argument errorMessage if you want to know a more detailed error message and the userCreds that the user has inputted when logging in.
// - The method provides an int that defines what has happened.
// - The method's argument `errorMessage *string` provides a more detailed message on what went wrong with not valid variable of the argument `userCreds`.
// - If nessesary, The first argument can be simply ignored. 
//
// # Returns
// - 0: No error
// - 1: Ip was deemed incorrect
// - 2: Port was deemed incorrect
// - 3: ClientId was deemed incorrect
//
// # Author
// - Polariusz
func validateCredentials(errorMessage *string, userCreds *MqttCredentials) int {
	// VALIDATE IP
	if userCreds.Ip == "" {
		if errorMessage != nil {
			*errorMessage = "IP is incomprehensible"
		}
		return 1
	}

	// VALIDATE PORT 
	if userCreds.Port == "" {
		if errorMessage != nil {
			*errorMessage = "PORT is incomprehensible"
		}
		return 2
	}

	// VALIDATE CLIENTID
	if userCreds.ClientId == "" {
		if errorMessage != nil {
			*errorMessage = "CLIENT-ID is incomprehensible"
		}
		return 3
	}

	return 0
}

// | Date of change | By        | Comment       |
// +----------------+-----------+---------------+
// |                | Polariusz | Created       |
// | 2025-05-13     | Polariusz | Documentation |
//
// # Structure:
// - {"Topics":<T>}
//   - <T>: String array of topics
//
// # Used in
// - PostTopicSubscribeHandler()
// - PostTopicUnsubscribeHandler()
//
// # Author
// - Polariusz
type TopicsWrapper struct {
	Topics []string
}

// | Date of change | By        | Comment |
// +----------------+-----------+---------+
// | 2025-05-16     | Polariusz | Created |
//
// # JSON-Structure:
// - {"Status":<S>,"Message":<M>}
//   - <S>: Status from the methods. Currently it can be 'Fine', 'What' and 'Error'
//   - <M>: Defined string literals in the methods. These explain what has happened.
//
// # Used in
// - PostTopicSubscribeHandler()
// - PostTopicUnsubscribeHandler()
//
// # Author
// - Polariusz
type TopicResult struct {
	Status string
	Message string
}

// | Date of change | By        | Comment                |
// +----------------+-----------+------------------------+
// |                | Polariusz | Created                |
// | 2025-05-13     | Polariusz | Documentation          |
// | 2025-05-16     | Polariusz | Changed one 400 to 207 |
//
// # Method-Type
// - Handler
//
// # Description
// - The method shall be a handler that allows to subscribe the MQTT-Broker's topics.
// - The method shall accept a jsonified structure that follows the struct TopicsWrapper.
// - The method shall return a 200 (Ok) if all requested topics were subscribed.
// - The method shall return a 207 (Multi Status) if at least one topic was not subscribed.
// - The method shall return a 400 (Bad Request) if the data from the client does not match that one fo the struct TopicsWrapper.
//
// # Usage
//
// - Call declared by the routing method addRoutes() URL with the POST-Method.
// - The data must have a json structure that matches the struct TopicsWrapper.
//
// # Returns
// - 200 (Ok): JSON
//   - {"goodJson":"Subscribed to requested topics"}
// - 207 (Multi Status): JSON
//   - {"result":{<TOPIC-N>:{"Status":<STATUS-N>,"Message":<MESSAGE-N>}}}
// - 400 (Bad Request): JSON
//   - {"badJson":`const BADJSON`}
// - 401 (Unauthorized): JSON
//   - {"401":"You fool!"}
//
// # Author
// - Polariusz
func PostTopicSubscribeHandler(serverState *ServerState) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if serverState.userCreds.Ip == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				// TODO: Explain the message a bit more
				"401": "You fool!",
			})
		}

		var subscribeTopics TopicsWrapper

		if err := c.BodyParser(&subscribeTopics); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"badJson": BADJSON,
			})
		}

		// TODO: Validate topics (if they are empty)

		topicResult := make(map[string]TopicResult)
		atLeastOneBadTopic := false

		for _, topic := range subscribeTopics.Topics {
			if !slices.Contains(serverState.subscribedTopics, topic) {
				if token := serverState.mqttClient.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
					atLeastOneBadTopic = true
					topicResult[topic] = TopicResult{"Error", "Make sure that the topic conforms the MQTT-Broker configuration."}
				} else {
					topicResult[topic] = TopicResult{"Fine", "Subscribed to the topic"}
					serverState.subscribedTopics = append(serverState.subscribedTopics, topic)
					serverState.allKnownTopics = append(serverState.allKnownTopics, topic)
				}
			} else {
				atLeastOneBadTopic = true
				topicResult[topic] = TopicResult{"What", "The topic is already subscribed"}
			}
		}

		if atLeastOneBadTopic {
			return c.Status(fiber.StatusMultiStatus).JSON(fiber.Map{
				"result": topicResult,
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"result": topicResult,
		})
	}
}

// | Date of change | By        | Comment                |
// +----------------+-----------+------------------------+
// |                | Polariusz | Created                |
// | 2025-05-13     | Polariusz | Documentation          |
// | 2025-05-16     | Polariusz | Changed one 400 to 207 |
//
// # Method-Type
// - Handler
//
// # Description
// - The method shall be a handler that allows to unsubscribe from subscribed topics.
// - The method shall accept a jsonified structure that follows the struct TopicsWrapper.
// - The method shall return a 200 (Ok) if all the topics from the converted to type TopicsWrapper have been successfully unsubscribed.
// - The method shall return a 207 (Multi Status) if at least one topic could not be unsubscribed from.
// - The method shall return a 400 (Bad Request) if the data from the client does not match that one fo the struct TopicsWrapper.
//
// # Usage
//
// - Call declared by the routing method addRoutes() URL with the POST-Method.
// - The data must have a json structure that matches the struct TopicsWrapper.
//
// # Returns
// - 200 (Ok): JSON
//   - {"goodJson":"Unsubscribed from requested"}
// - 207 (Multi Status): JSON
//   - {"result":{<TOPIC-N>:{"Status":<STATUS-N>,"Message":<MESSAGE-N>}}}
// - 400 (Bad Request): JSON
//   - {"badJson":`const BADJSON`}
// - 401 (Unauthorized): JSON
//   - {"401":"You fool!"}
//
// # Author
// - Polariusz
func PostTopicUnsubscribeHandler(serverState *ServerState) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if serverState.userCreds.Ip == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				// TODO: Explain the message a bit more
				"401": "You fool!",
			})
		}

		var unsubscribeTopics TopicsWrapper

		if err := c.BodyParser(&unsubscribeTopics); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"badJson": BADJSON,
			})
		}

		// TODO: Validate topics (if they are empty)

		topicResult := make(map[string]TopicResult)
		atLeastOneBadTopic := false

		for _, topic := range unsubscribeTopics.Topics {
			if !slices.Contains(serverState.subscribedTopics, topic) {
				topicResult[topic] = TopicResult{"What", "The topic wasn't even subscribed"}
				atLeastOneBadTopic = true
			} else {
				// unsubscribe from mqtt-broker
				serverState.mqttClient.Unsubscribe(topic)
				// add fine to the map
				topicResult[topic] = TopicResult{"Fine", "Unsubscribed successfully"}
				// Find the topic in the array and remove it.
				for index, stateTopic := range serverState.subscribedTopics {
					if strings.Compare(stateTopic, topic) == 0 {
						serverState.subscribedTopics = append(serverState.subscribedTopics[:index], serverState.subscribedTopics[index+1:]...)
						break
					}
				}
			}
		}

		if atLeastOneBadTopic {
			// TODO: The status code is meh, as the function at this point would
			// subscribe to at least some of the requested topics.
			return c.Status(fiber.StatusMultiStatus).JSON(fiber.Map{
				"result": topicResult,
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"result": topicResult,
		})
	}
}

// | Date of change | By        | Comment       |
// +----------------+-----------+---------------+
// |                | Polariusz | Created       |
// | 2025-05-13     | Polariusz | Documentation |
//
// # Method-Type
// - Handler
//
// # Description
// - The method shall be a handler that allows to get a list of subscribed topics. These topics will be strings or simply string array.
// - The method shall return a 200 (Ok) with the subscribed topics.
//
// # Usage
// - Call declared by the routing method addRoutes() URL with the GET-Method.
//
// # Returns
// - 200 (Ok): JSON
//   - {"topics":`serverState.subscribedTopics`}
// - 401 (Unauthorized): JSON
//   - The go server was never connected to the MQTT-Broker.
//   - {"401":"You fool!"}
//
// # Author
// - Polariusz
func GetTopicSubscribedHandler(serverState *ServerState) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if serverState.userCreds.Ip == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				// TODO: Explain the message a bit more
				"401": "You fool!",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"topics": serverState.subscribedTopics,
		})
	}
}

// | Date of change | By        | Comment       |
// +----------------+-----------+---------------+
// |                | Polariusz | Created       |
// | 2025-05-13     | Polariusz | Documentation |
//
// # Structure:
// - {"Topic":<T>,"Message":"<M>"}
//   - <T>: Topic
//   - <M>: Message
//
// # Used in
// - PostTopicSendMessageHandler()
//
// # Author
// - Polariusz
type MessageWrapper struct {
	Topic string // TODO: This could be converted to a string array if you wish for the publich messages method to send the same message to multiple topics.
	Message string
}

// | Date of change | By        | Comment       |
// +----------------+-----------+---------------+
// |                | Polariusz | Created       |
// | 2025-05-13     | Polariusz | Documentation |
//
// # Method-Type
// - Handler
//
// # Description
// - The method shall be a handler that allows to send messages to the MQTT-Broker.
// - The method shall accept a jsonified structure that follows the struct MessageWrapper.
// - The method shall return a 200 (Ok) if the go-server publishes a message.
// - The method shall return a 400 (Bad Request) if the data from the client does not match that one fo the struct MessageWrapper.
//
// # Usage
// - Call declared by the routing method addRoutes() URL with the POST-Method.
// - The data must have a json structure that matches the struct MessageWrapper.
//
// # Returns
// - 200 (Ok): JSON
//   - {"goodJson":"Message posted"}
// - 400 (Bad Request): JSON
//   - {"badJson":`const BADJSON`}
// - 401 (Unauthorised): JSON
//   - {"401":"You fool!"}
//
// # Author
// - Polariusz
func PostTopicSendMessageHandler(serverState *ServerState) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if serverState.userCreds.Ip == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				// TODO: Explain the message a bit more
				"401": "You fool!",
			})
		}

		var messageWrapper MessageWrapper

		if err := c.BodyParser(&messageWrapper); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"badJson": BADJSON,
			})
		}

		// TODO: Validate topic and message!

		// TODO: This can be changed to check if the MQTT-Broker responds! Publish() method returns a token, and the token has method Wait() that waits for the respose and Error() that has either nil or an actual error.
		serverState.mqttClient.Publish(messageWrapper.Topic, 0, false, messageBuilder(serverState.userCreds, messageWrapper.Message))

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"goodJson": "Message posted",
		})
	}
}

// | Date of change | By        | Comment       |
// +----------------+-----------+---------------+
// |                | Polariusz | Created       |
// | 2025-05-13     | Polariusz | Documentation |
//
// # Method-Type
// - Handler
//
// # Description
// - The method shall be a handler that allows to check if the go-server is connected to the MQTT-Broker.
// - The method shall ignore any data being sent to the server, be it a json or any byte array.
// - The method shall return 200 (Ok) if the go-server is connected to the MQTT-Broker
// - The method shall return 501 (Service Unavailable) if the MQTT-Broker does not respond back.
//
// # Usage
// - Call declared by the routing method addRoutes() URL with the GET-Method.
//
// # Returns
// - 200 (Ok): JSON
//   - MQTT-Broker responded, signifying that the connection to the Broker exists.
//   - {"goodMqtt":"pong"}
// - 401 (Unauthorized): JSON
//   - The go server was never connected to the MQTT-Broker.
//   - {"401":"You fool!"}
// - 503 (Service Unavailable): JSON
//   - The connection to the MQTT-Broker was severed.
//   - {"badMqtt":"Big f!"}
//
// # Author
// - Polariusz
func GetPingHandler(serverState *ServerState) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if serverState.userCreds.Ip == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				// TODO: Explain the message a bit more
				"401": "You fool!",
			})
		}

		if token := serverState.mqttClient.Publish("ping", 0, false, '0'); token.Wait() && token.Error() != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"badMqtt": "Big f!",
			})
		}
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"goodMqtt": "pong",
			})
	}
}

// TODO: I've written this to play around with SSE a bit more. This doesn't do anything with the mqtt part.
//       It simply needs to be recycled for a way of the client of getting new messages.
func writeStuff() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var messageWrapper MessageWrapper

		if err := c.BodyParser(&messageWrapper); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"badJson": BADJSON,
			})
		}

		active = append(active, ActiveMessageWrapper{
			ClientId: "kurwaNieWiem",
			Topic: messageWrapper.Topic,
			Message: messageWrapper.Message,
			epoch: time.Now().Unix(),
		})
		return nil
	}
}

// TODO: I've written this to play around with SSE a bit more. This doesn't do anything with the mqtt part.
//       It simply needs to be recycled for a way of the client of getting new messages.
var active []ActiveMessageWrapper

// TODO: I've written this to play around with SSE a bit more. This doesn't do anything with the mqtt part.
//       It simply needs to be recycled for a way of the client of getting new messages.
type ActiveMessageWrapper struct {
	ClientId string
	Topic string
	Message string
	epoch int64
}

// TODO: I've written this to play around with SSE a bit more. This doesn't do anything with the mqtt part.
//       It simply needs to be recycled for a way of the client of getting new messages.
func SseHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// setting headers to make the handler be the SSE.
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("Transfer-Encoding", "chunked")

		// Here are some things to unpack..:
		// Status() is simple to understand, we simply tell that whatever we return is 200.
		// With content, we define a stream writer. That stream writer will be a fasthttp.StreamWriter.
		// fasthttp.StreamWriter wants a function that we write below.
		// The argument w *bufio.Writer abstracts the socket connection for us and we can use it to write any things
		// to the client part of the project.
		c.Status(fiber.StatusOK).Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
			for {
				//var msg string

				if len(active) > 0 {
					aaaaa := struct{Messages []ActiveMessageWrapper}{active,}
					msg, _ := json.Marshal(aaaaa)
					//msg = fmt.Sprintf("%d - message received: %s", i, active[0])
					// To conform to the sse protocol, it must begin with `data:` and end with `\n\n`.
					fmt.Fprintf(w, "data: Message: %s\n\n", msg)
					// remove the message from the buffer
					active = active[:0]
					// w.Flush() writes anything in the println immediately out. With the err := and the check below, we check if we got any error.
					err := w.Flush()
					if err != nil {
						// Refreshing page in web browser will establish a new
						// SSE connection, but only (the last) one is alive, so
						// dead connections must be closed here.
						fmt.Printf("Error while flushing: %v. Closing http connection.\n", err)

						break
					}
				}
				time.Sleep(10 * time.Second)
			}
		}))
		return nil
	}
}

// | Date of change | By     | Comment                 |
// +----------------+--------+-------------------------+
// | 2025-05-14     | Tibbyx | Created & Documentation |
//
// # Method-Type
// - MQTT Handler Factory
//
// # Description
// - The method shall create and return an MQTT message handler.
// - The handler processes incoming MQTT messages from subscribed topics.
// - The payloads are stored in the in-memory map receivedMessages within the ServerState.
// - The topic string is used as key, and messages are appended as values (string slices).
//
// # Usage
// - Used in PostCredentialsHandler() to assign the MQTT client’s default message handler.
// - Requires a reference to ServerState to access the receivedMessages map.
//
// # Returns
// - mqtt.MessageHandler: A function that handles and stores MQTT messages.
//
// # Author
// - Tibbyx
func createMessageHandler(serverState *ServerState) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		topic := msg.Topic()
		payload := string(msg.Payload())

		fmt.Printf("Received message: %s from topic: %s\n", payload, topic)

		serverState.receivedMessages[topic] = append(serverState.receivedMessages[topic], payload)
	}
}

// | Date of change | By     | Comment                 |
// +----------------+--------+-------------------------+
// | 2025-05-14     | Tibbyx | Created & Documentation |
//
// # Method-Type
// - Handler
//
// # Description
// - The method shall be a handler that returns all stored messages for a specific topic.
// - The messages must have previously been received through an active MQTT subscription.
// - The topic is provided as a query parameter.
//
// # Usage
// - Call declared by the routing method addRoutes() URL with the GET-Method.
// - URL must include ?topic=<topic-name>
//
// # Returns
// - 200 (Ok): JSON
//   - {"topic": "<topic-name>", "messages": ["msg1", "msg2", ...]}
// - 400 (Bad Request): JSON
//   - {"error": "Missing topic query parameter"}
//
// # Author
// - Tibbyx
func GetTopicMessagesHandler(serverState *ServerState) fiber.Handler {
	return func(c *fiber.Ctx) error {
		topic := c.Query("topic")
		if topic == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Missing topic query parameter",
			})
		}

		messages, ok := serverState.receivedMessages[topic]
		if !ok {
			messages = []string{}
		}

		return c.JSON(fiber.Map{
			"topic": topic,
			"messages": messages,
		})
	}
}

// | Date of change | By        | Comment |
// +----------------+-----------+---------|
// | 2025-05-16     | Polariusz | Created |
//
// # Method-Type
// - Handler
//
// # Description
// - The method shall return a list of all previously subscribed Topics in a JSON format
// - The method shall return a 200 (Ok) with the list if user is authenticated
// - The method shall return a 401 (Unauthorized) if the user is not authenticated
//
// # Usage
// - Call declared by the routing method addRoutes() URL with the GET-Method.
//
// # Returns
// - 200 (Ok): JSON
//   - {"Topics":[<TOPIC-N>]}
// - 401 (Unauthorized): JSON
//   - {"Message":"Authenticate yourself first!"}
//
// # Author
// - Polariusz
func GetTopicAllKnownHandler(serverState *ServerState) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if serverState.userCreds.Ip == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				// TODO: Explain the message a bit more
				"Message": "Authenticate yourself first!",
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"Topics": serverState.allKnownTopics,
		})
	}
}

// | Date of change | By        | Comment |
// +----------------+-----------+---------|
// | 2025-05-16     | Polariusz | Created |
//
// # Method-Type
// - Handler
//
// # Description
// - The method shall disconnect the from the argument serverstate mqttClient from the MQTT-Broker
// - The method shall return a 200 (Ok) if the user is authenticated
// - The method shall return a 401 (Unauthorized) if the user is not authenticated
//
// # Usage
// - Call declared by the routing method addRoutes() URL with the GET-Method.
//
// # Returns
// - 200 (Ok): JSON
//   - {"Fine":"The MQTT-Client disconnected from <IP>:<PORT> Broker"}
// - 401 (Unauthorized): JSON
//   - {"BadRequest":"The server isn't even connected to any MQTT-Brokers"}
//
// # Author
// - Polariusz
func PostDisconnectFromBrokerHandler(serverState *ServerState) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if(serverState.userCreds.Ip == "") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"BadRequest": "The server isn't even connected to any MQTT-Brokers",
			})
		}

		serverState.mqttClient.Disconnect(250)
		message := fmt.Sprintf("The MQTT-Client disconented from %s:%s Broker", serverState.userCreds.Ip, serverState.userCreds.Port)

		serverState.userCreds.Ip = ""
		serverState.userCreds.Port = ""
		serverState.userCreds.ClientId = ""
		serverState.subscribedTopics = nil

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"Fine": message,
		})
	}
}
