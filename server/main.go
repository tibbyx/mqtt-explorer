package main

import (
	"database"

	"github.com/gofiber/fiber/v2"
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"slices"
	"time"
)

// # Author
// - Polariusz
const BADJSON = "I am nowt sowwy >:3. An expected! ewwow has happened. Youw weak json! iws of the wwongest fowmat thawt does nowt cowwespond tuwu the stwong awnd independent stwuct! >:P"

// | Date of change | By        | Comment              |
// +----------------+-----------+----------------------+
// |                | Polariusz | Created              |
// | 2025-05-13     | Polariusz | Documentation        |
// | 2025-05-16     | Polariusz | added AllKnownTopics |
// | 2025-05-20     | Nicolas   | added Database       |
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
	db *database.Database // Database connection
	currentConnectionID int64 // Tracks the current connection ID for disconnection logging
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

// | Date of change | By        | Comment       |
// +----------------+-----------+---------------+
// |                | Polariusz | Created       |
// | 2025-05-13     | Polariusz | Documentation |
//
// # Method-Type
// - Routing
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

// | Date of change | By        | Comment       |
// +----------------+-----------+---------------+
// |                | Polariusz | Created       |
// | 2025-05-13     | Polariusz | Documentation |
// | 2025-05-20     | Nicolas   | Added database logging |
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
		mqttOpts.SetKeepAlive(60 * time.Second)
		mqttOpts.SetPingTimeout(15 * time.Second)
		mqttOpts.SetCleanSession(true)
		mqttOpts.SetAutoReconnect(true)
		mqttOpts.SetConnectRetry(true)

		mqttOpts.SetDefaultPublishHandler(createMessageHandler(serverState))

		mqttClient := mqtt.NewClient(mqttOpts)

		if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"badJson": fmt.Sprintf("Connecting to %s:%s failed\n%s", userCreds.Ip, userCreds.Port, token.Error()),
			})
		}

		serverState.userCreds = userCreds
		serverState.mqttClient = mqttClient

		// Log connection to database if available
		if serverState.db != nil {
			connectionID, err := serverState.db.LogConnection(userCreds.ClientId, userCreds.Ip, userCreds.Port)
			if err != nil {
				fmt.Printf("Failed to log connection to database: %v\n", err)
			} else {
				serverState.currentConnectionID = connectionID
			}
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"goodJson": fmt.Sprintf("Connecting to %s:%s succeeded", userCreds.Ip, userCreds.Port),
		})
	}
}

// | Date of change | By        | Comment       |
// +----------------+-----------+---------------+
// |                | Polariusz | Created       |
// | 2025-05-13     | Polariusz | Documentation |
// | 2025-05-20     | Nicolas   | Added database topic tracking |
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
		if serverState.mqttClient == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				// TODO: Explain the message a bit more
				"401": "You fool!",
			})
		}

		var topics TopicsWrapper

		if err := c.BodyParser(&topics); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"badJson": BADJSON,
			})
		}

		var results map[string]TopicResult = make(map[string]TopicResult)
		allSuccessful := true

		for _, topic := range topics.Topics {
			if slices.Contains(serverState.subscribedTopics, topic) {
				results[topic] = TopicResult{
					Status:  "What",
					Message: "Already subscribed to this topic",
				}
				continue
			}

			token := serverState.mqttClient.Subscribe(topic, 0, nil)
			token.Wait() // blocking
			if token.Error() != nil {
				allSuccessful = false
				results[topic] = TopicResult{
					Status:  "Error",
					Message: fmt.Sprintf("Error subscribing to the topic: %v", token.Error()),
				}
			} else {
				serverState.subscribedTopics = append(serverState.subscribedTopics, topic)

				// Add to allKnownTopics if not already there
				if !slices.Contains(serverState.allKnownTopics, topic) {
					serverState.allKnownTopics = append(serverState.allKnownTopics, topic)
				}

				// Store topic in database if available
				if serverState.db != nil {
					if err := serverState.db.AddOrUpdateTopic(topic, true); err != nil {
						fmt.Printf("Failed to save topic to database: %v\n", err)
					}
				}

				results[topic] = TopicResult{
					Status:  "Fine",
					Message: "Subscribed to the topic",
				}
			}
		}

		if allSuccessful {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"goodJson": "Subscribed to requested topics",
			})
		} else {
			return c.Status(fiber.StatusMultiStatus).JSON(fiber.Map{
				"result": results,
			})
		}
	}
}

// | Date of change | By        | Comment                |
// +----------------+-----------+------------------------+
// |                | Polariusz | Created                |
// | 2025-05-13     | Polariusz | Documentation          |
// | 2025-05-16     | Polariusz | Changed one 400 to 207 |
// | 2025-05-20     | Nicolas   | Added database update  |
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
		if serverState.mqttClient == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				// TODO: Explain the message a bit more
				"401": "You fool!",
			})
		}
		var topics TopicsWrapper

		if err := c.BodyParser(&topics); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"badJson": BADJSON,
			})
		}

		var results map[string]TopicResult = make(map[string]TopicResult)
		allSuccessful := true

		for _, topic := range topics.Topics {
			if !slices.Contains(serverState.subscribedTopics, topic) {
				allSuccessful = false
				results[topic] = TopicResult{
					Status:  "What",
					Message: "Not subscribed to topic, so can't unsubscribe",
				}
				continue
			}

			token := serverState.mqttClient.Unsubscribe(topic)
			token.Wait() // blocking
			if token.Error() != nil {
				allSuccessful = false
				results[topic] = TopicResult{
					Status:  "Error",
					Message: fmt.Sprintf("Error unsubscribing from the topic: %v", token.Error()),
				}
			} else {
				// This is a wasteful delete, in larger applications a better method should be used.
				serverState.subscribedTopics = slices.DeleteFunc(serverState.subscribedTopics, func(s string) bool {
					return s == topic
				})

				// Update topic status in database
				if serverState.db != nil {
					if err := serverState.db.AddOrUpdateTopic(topic, false); err != nil {
						fmt.Printf("Failed to update topic in database: %v\n", err)
					}
				}

				results[topic] = TopicResult{
					Status:  "Fine",
					Message: "Unsubscribed from the topic",
				}
			}
		}

		if allSuccessful {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"goodJson": "Unsubscribed from requested",
			})
		} else {
			return c.Status(fiber.StatusMultiStatus).JSON(fiber.Map{
				"result": results,
			})
		}
	}
}

// | Date of change | By        | Comment       |
// +----------------+-----------+---------------+
// |                | Polariusz | Created       |
// | 2025-05-13     | Polariusz | Documentation |
// | 2025-05-20     | Nicolas   | Added database lookup |
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
		if serverState.mqttClient == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				// TODO: Explain the message a bit more
				"401": "You fool!",
			})
		}

		// If we have a database connection, retrieve subscribed topics from there
		if serverState.db != nil {
			topics, err := serverState.db.GetSubscribedTopics()
			if err == nil && len(topics) > 0 {
				// Only use database topics if we successfully retrieved them
				return c.Status(fiber.StatusOK).JSON(fiber.Map{
					"topics": topics,
				})
			}
		}

		// Fall back to in-memory topics if database lookup failed or returned no results
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"topics": serverState.subscribedTopics,
		})
	}
}

// | Date of change | By        | Comment       |
// +----------------+-----------+---------------+
// |                | Polariusz | Created       |
// | 2025-05-13     | Polariusz | Documentation |
// | 2025-05-20     | Nicolas   | Added database message storage |
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
		if serverState.mqttClient == nil {
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

		// Build message with client ID info
		finalMessage := messageBuilder(serverState.userCreds, messageWrapper.Message)

		// Store the sent message in the database if available
		if serverState.db != nil {
			// Store topic in database if it's not already known
			if !slices.Contains(serverState.allKnownTopics, messageWrapper.Topic) {
				if err := serverState.db.AddOrUpdateTopic(messageWrapper.Topic, false); err != nil {
					fmt.Printf("Failed to save topic to database: %v\n", err)
				}
				serverState.allKnownTopics = append(serverState.allKnownTopics, messageWrapper.Topic)
			}

			// Store message in database
			if err := serverState.db.SaveMessage(serverState.userCreds.ClientId, messageWrapper.Topic, finalMessage); err != nil {
				fmt.Printf("Failed to save message to database: %v\n", err)
			}
		}

		token := serverState.mqttClient.Publish(messageWrapper.Topic, 0, false, finalMessage)
		token.Wait()

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"goodJson": "Message posted",
		})
	}
}

// | Date of change | By        | Comment                |
// +----------------+-----------+------------------------+
// |                | Polariusz | Created                |
// | 2025-05-13     | Polariusz | Documentation          |
// | 2025-05-19     | Polariusz | Updated ping behaviour |
//
// # Method-Type
// - Handler
//
// # Description
// - The method shall be a handler that allows to check if the go-server is connected to the MQTT-Broker.
// - The method shall ignore any data being sent to the server, be it a json or any byte array.
// - The method shall return 200 (Ok) if the go-server is connected to the MQTT-Broker
// - The method shall return 200 (Ok) if the go-server is reconnecting to the MQTT-Broker
// - The method shall return 401 (Unauthorized) if the client never authenticated.
// - The method shall return 503 (Service Unavailable) if the MQTT-Broker does not respond back.
//
// # Usage
// - Call declared by the routing method addRoutes() URL with the GET-Method.
//
// # Returns
// - 200 (Ok): JSON
//   - MQTT-Broker responded, signifying that the connection to the Broker exists.
//     - {"OK":"Connection is active"}
//   - Paho is trying to reconnect to the MQTT-Broker, which should be fine
//     - {"Fine", "Reconnecting, but otherwise connected"}
// - 401 (Unauthorized): JSON
//   - The go server was never connected to the MQTT-Broker.
//   - {"Unauthorized":"Authenticate yourself first!"}
// - 503 (Service Unavailable): JSON
//   - The connection to the MQTT-Broker was severed.
//   - {"ServiceUnavailable":"The Credentials suggest that the server should be connected to a broker, but it isn't!","Ip":"<BROKER-IP>","Port":"<BROKER-PORT>","ClientId":"<CLIENT-ID>"}
//
// # Author
// - Polariusz
func GetPingHandler(serverState *ServerState) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if serverState.userCreds.Ip == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"Unauthorized": "Authenticate yourself first!",
			})
		}

		// Is it connected?
		if serverState.mqttClient.IsConnected() {
			// Is it really connected? (i.e not reconnecting)
			if serverState.mqttClient.IsConnectionOpen() {
				return c.Status(fiber.StatusOK).JSON(fiber.Map{
					"Ok": "Connection is active",
				})
			} else {
				// It is reconnecting
				return c.Status(fiber.StatusOK).JSON(fiber.Map{
					"Fine": "Reconnecting, but otherwise connected",
				})
			}
		} else {
			// It is not connected
			Ip := serverState.userCreds.Ip; serverState.userCreds.Ip = ""
			Port := serverState.userCreds.Port; serverState.userCreds.Port = ""
			ClientId := serverState.userCreds.ClientId; serverState.userCreds.ClientId = ""

			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"ServiceUnavailable": "The Credentials suggest that the server should be connected to a broker, but it isn't!",
				"Ip": Ip,
				"Port": Port,
				"ClientId": ClientId,
			})
		}
	}
}

// | Date of change | By     | Comment                 |
// +----------------+--------+-------------------------+
// | 2025-05-14     | Tibbyx | Created & Documentation |
// | 2025-05-20     | Nicolas | Added database message retrieval |
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

		// Store in-memory for backward compatibility
		if _, ok := serverState.receivedMessages[topic]; !ok {
			serverState.receivedMessages[topic] = []string{}
		}
		serverState.receivedMessages[topic] = append(serverState.receivedMessages[topic], payload)

		// Add to known topics if not already known
		if !slices.Contains(serverState.allKnownTopics, topic) {
			serverState.allKnownTopics = append(serverState.allKnownTopics, topic)

			// Store topic in database if available
			if serverState.db != nil {
				isSubscribed := slices.Contains(serverState.subscribedTopics, topic)
				if err := serverState.db.AddOrUpdateTopic(topic, isSubscribed); err != nil {
					fmt.Printf("Failed to save topic to database: %v\n", err)
				}
			}
		}

		// Store message in database if available
		if serverState.db != nil {
			if err := serverState.db.SaveMessage(serverState.userCreds.ClientId, topic, payload); err != nil {
				fmt.Printf("Failed to save message to database: %v\n", err)
			}
		}

		fmt.Printf("Received message on topic '%s': %s\n", topic, payload)
	}
}

// | Date of change | By     | Comment                 |
// +----------------+--------+-------------------------+
// | 2025-05-14     | Tibbyx | Created & Documentation |
// | 2025-05-20     | Nicolas | Added database message retrieval |
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

		// If database is available, try to get messages from there first
		if serverState.db != nil {
			dbMessages, err := serverState.db.GetMessagesByTopic(topic)
			if err == nil && len(dbMessages) > 0 {
				// Extract just the payload strings for the response
				var messagePayloads []string
				for _, msg := range dbMessages {
					messagePayloads = append(messagePayloads, msg.Payload)
				}

				return c.Status(fiber.StatusOK).JSON(fiber.Map{
					"topic":    topic,
					"messages": messagePayloads,
				})
			}
			// If database query failed or returned no results, fall back to in-memory messages
		}

		// Fall back to in-memory messages
		messages, exists := serverState.receivedMessages[topic]
		if !exists {
			messages = []string{} // Return empty array if no messages
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"topic":    topic,
			"messages": messages,
		})
	}
}

// | Date of change | By        | Comment |
// +----------------+-----------+---------|
// | 2025-05-16     | Polariusz | Created |
// | 2025-05-20     | Nicolas   | Added database lookup |
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
		if serverState.mqttClient == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				// TODO: Explain the message a bit more
				"Message": "Authenticate yourself first!",
			})
		}

		// If database is available, try to get topics from there first
		if serverState.db != nil {
			dbTopics, err := serverState.db.GetAllTopics()
			if err == nil && len(dbTopics) > 0 {
				return c.Status(fiber.StatusOK).JSON(fiber.Map{
					"Topics": dbTopics,
				})
			}
			// If database query failed or returned no results, fall back to in-memory topics
		}

		// Fall back to in-memory topics
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"Topics": serverState.allKnownTopics,
		})
	}
}

// | Date of change | By        | Comment |
// +----------------+-----------+---------|
// | 2025-05-16     | Polariusz | Created |
// | 2025-05-20     | Nicolas   | Added database disconnect logging |
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
// - 400 (BadRequest): JSON
//   - {"BadRequest":"The server isn't even connected to any MQTT-Brokers"}
//
// # Author
// - Polariusz
func PostDisconnectFromBrokerHandler(serverState *ServerState) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if serverState.mqttClient == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"BadRequest": "The server isn't even connected to any MQTT-Brokers",
			})
		}

		// Log disconnection in database if available
		if serverState.db != nil && serverState.currentConnectionID > 0 {
			if err := serverState.db.LogDisconnection(serverState.currentConnectionID); err != nil {
				fmt.Printf("Failed to log disconnection to database: %v\n", err)
			}
			// Reset connection ID
			serverState.currentConnectionID = 0
		}

		// Save data from previous connection for response
		ip := serverState.userCreds.Ip
		port := serverState.userCreds.Port

		// Unsubscribe from all topics
		for _, topic := range serverState.subscribedTopics {
			token := serverState.mqttClient.Unsubscribe(topic)
			token.Wait()
		}

		// Disconnect from client
		serverState.mqttClient.Disconnect(250)
		serverState.mqttClient = nil
		serverState.subscribedTopics = []string{}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"Fine": fmt.Sprintf("The MQTT-Client disconnected from %s:%s Broker", ip, port),
		})
	}
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
	server.Get("/topic/all-known", GetTopicAllKnownHandler(serverState))
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
	Status  string
	Message string
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
	Topic   string // TODO: This could be converted to a string array if you wish for the publich messages method to send the same message to multiple topics.
	Message string
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

func main() {
	db, err := database.SetupDatabase()
	if(err != nil) {
		fmt.Println("WARN: Running without database:", err)
	}

	server := fiber.New()
	var serverState ServerState
	serverState.receivedMessages = make(map[string][]string)
	serverState.db = db

	addRoutes(server, &serverState)

	// need to build ui via 'npm run build' in client first
	// server.Static("/", "../client/dist")
	server.Listen(":3000")
}
