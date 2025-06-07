package main

import (
	"database"
	"database/sql"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"time"
	"strconv"
)

// | Date of change | By        | Comment |
// +----------------+-----------+---------+
// | 2025-06-05     | Polariusz | Created |
//
// # Structure:
// - {"BrokerId":<B>,"UserId":<U>}
//   - <B> : The ID of the Broker ROW matched from the BrokerId from PostCredentialsHandler()'s brokerId
//   - <U> : The ID of the User ROW matched from the BrokerId from PostCredentialsHandler()'s brokerId
//
// # Used in
// - type TopicsWrapper struct
//
// # Author
// - Polariusz
type BrokerUser struct {
	BrokerId int
	UserId int
}

// | Date of change | By        | Comment |
// +----------------+-----------+---------+
// | 2025-06-06     | Polariusz | Created |
//
// # Structure:
// - {"BrokerId":<B>,"UserId":<U>,"Message":<Message>}
//   - <B> : The ID of the Broker ROW matched from the BrokerId from PostCredentialsHandler()'s brokerId
//   - <U> : The ID of the User ROW matched from the BrokerId from PostCredentialsHandler()'s brokerId
//   - <M> : Published Message
//
// # Used in
// - messageBuilder()
//
// # Author
// - Polariusz
type JsonPublishMessage struct {
	ClientId string
	Message string
}

// # Author
// - Polariusz
const BADJSON = "I am nowt sowwy >:3. An expected! ewwow has happened. Youw weak json! iws of the wwongest fowmat thawt does nowt cowwespond tuwu the stwong awnd independent stwuct! >:P"

// | Date of change | By        | Comment              |
// +----------------+-----------+----------------------+
// | 2025-05-13     | Polariusz | Created              |
// | 2025-06-06     | Polariusz | Actually implemented |
//
// # Description
// - This method shall build a message containing the BrokerId, UserId and the Message that the messagePubHandler will be able to then differentiate and write into the database accordingly.
//
// # Author
// - Polariusz
func messageBuilder(clientId string, message string) []byte {
	fullMessage := JsonPublishMessage {
		clientId,
		message,
	}
	
	// ignoring error as the structure above cannot fail to marshal
	jsonMessage, _ := json.Marshal(fullMessage)

	return jsonMessage
}

// | Date of change | By        | Comment               |
// +----------------+-----------+-----------------------+
// |                | Polariusz | Created               |
// | 2025-05-13     | Polariusz | Documentation         |
// | 2025-05-16     | Polariusz | added AllKnownTopics  |
// | 2025-05-18     | Polariusz | added favouriteTopics |
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
	favouriteTopics []string // TODO: This should be in a database that we don't have yet.
	con *sql.DB
}

// | Date of change | By        | Comment                     |
// +----------------+-----------+-----------------------------+
// |                | Polariusz | Created                     |
// | 2025-05-13     | Polariusz | Documentation               |
// | 2025-06-04     | Polariusz | Added Username and Password |
//
// # Structure:
// - {"Ip":<I>,"Port":"<Po>","ClientId":"<C>","Username":"<U>","Password":"<Pa>"}
//   - <I> : The IP of the MQTT-Broker.
//   - <Po>: The Port that the MQTT-Broker opened for the protocol.
//   - <C> : Client ID that functions as an username. It makes the users distinct.
//   - <U> : Username for the broker, it's optional
//   - <Pa>: Password for the broker, it's optional
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
	Username string
	Password string
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
	con, err := database.OpenDatabase()
	if err != nil {
		fmt.Printf("WARN: Running without database\nErr:%s\n", err)
	}
	if err := database.SetupDatabase(con); err != nil {
		fmt.Printf("WARN: Issue with setting db up!\nErr:%s\n", err)
	}

	server := fiber.New()
	var serverState ServerState
	serverState.receivedMessages = make(map[string][]string)
	serverState.con = con

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
	server.Get("/topic/all-known", GetTopicAllKnownHandler(serverState))
	server.Post("/topic/favourites/mark", PostTopicFavouritesMark(serverState))
	server.Post("/topic/favourites/unmark", PostTopicFavouritesUnmark(serverState))		// TODO: DB
	server.Get("/topic/favourites", GetTopicFavourites(serverState))
}

// | Date of change | By        | Comment               |
// +----------------+-----------+-----------------------+
// |                | Polariusz | Created               |
// | 2025-05-13     | Polariusz | Documentation         |
// | 2025-06-04     | Polariusz | Integrated DB         |
// | 2025-06-05     | Polariusz | Updated documentation |
// | 2025-06-06     | Polariusz | Added auto subs       |
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
//   - {"goodJson":"Connecting to `Ip`:`Port` succeded", "brokerId":"<B>", "userId":"<U>", "subscribedTopics":[<ST>]}
//     - <B>  : This is the ID of the ROW from table Broker. The client needs to remember it and use it for the other functions.
//     - <U>  : This is the ID of the ROW from table User. The client needs to remember it and use it for the other functions.
//     - <ST> : It's the result from SelectSubscribedTopics() matched to data arguments <B> and <U>. Please take a look at `database.SelectTopic` struct.
// - 400 (Bad Request): JSON
//   - {"badJson":`const BADJSON`}
//   - {"badJson":`errorMessage`}
// - 404 (Not Found): JSON
//   - {"badMqtt":"Connecting to `Ip`:`Port` failed"}
// - 500 (Internal Server Error): JSON
//   - {"InternalServerError" : "Error while inserting in the <T> table", "Error" : "<E>"}
//     - <T> : It can be Broker or User
//     - <E> : SQL Error message
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

		// NOTE: I do this before to get the brokerId for the createMessageHandler.
		// Skipping err, as this should be validated in the validation function.
		port, _ := strconv.Atoi(userCreds.Port)
		brokerId, err := database.InsertNewBroker(serverState.con, database.InsertBroker{userCreds.Ip, port})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"InternalServerError" : "Error while inserting in the Broker table",
				"Error" : err.Error(),
			})
		}

		// test.mosquitto.org
		mqttOpts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s:%s", userCreds.Ip, userCreds.Port)).SetClientID(userCreds.ClientId)
		mqttOpts.SetKeepAlive(2 * time.Second)
		mqttOpts.SetPingTimeout(1 * time.Second)

		mqttOpts.SetDefaultPublishHandler(createMessageHandler(serverState, brokerId))

		mqttClient := mqtt.NewClient(mqttOpts)
		serverState.mqttClient = mqttClient

		if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"badJson": fmt.Sprintf("Connecting to %s:%s failed\n%s", userCreds.Ip, userCreds.Port, token.Error()),
			})
		}

		userId, err := database.InsertNewUser(serverState.con, database.InsertUser{brokerId, userCreds.ClientId, userCreds.Username, userCreds.Password, false})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"InternalServerError" : "Error while inserting in the User table",
				"Error" : err.Error(),
			})
		}

		topicList, err := database.SelectSubscribedTopics(serverState.con, brokerId, userId)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"InternalServerError" : "Error while selecting subscribed topics",
				"Error" : err.Error(),
			})
		}
		
		for _, topicToSub := range topicList {
			if token := serverState.mqttClient.Subscribe(topicToSub.Topic, 0, nil); token.Wait() && token.Error() != nil {
				fmt.Printf("ERROR: Subscribtion to topic %s failed!\n", topicToSub)
			}
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"goodJson" : fmt.Sprintf("Connecting to %s:%s succeded", userCreds.Ip, userCreds.Port),
			"brokerId" : brokerId,
			"userId" : userId,
			"subscribedTopics" : topicList,
		})
	}
}

// | Date of change | By        | Comment                  |
// +----------------+-----------+--------------------------+
// |                | Polariusz | Created                  |
// | 2025-05-13     | Polariusz | Documentation            |
// | 2025-06-05     | Polariusz | Improved Port validation |
//
// # Method-Type
// - Validator
//
// # Description
// - The method shall validate the argument `userCreds *MqttCredentials`.
//   - TODO: The validation need to be improved. Right now it only checks if the argument userCreds has empty strings.
//           The Port has now a better validation.
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
	port, err := strconv.Atoi(userCreds.Port);
	if err != nil {
		if errorMessage != nil {
			*errorMessage = "PORT cannot be converted to int"
		}
		return 2;
	}

	if port < 0 || port > 65535 {
		if errorMessage != nil {
			*errorMessage = "PORT is not between 0 and 65535"
		}
		return 2;
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

// | Date of change | By        | Comment          |
// +----------------+-----------+------------------+
// |                | Polariusz | Created          |
// | 2025-05-13     | Polariusz | Documentation    |
// | 2025-06-05     | Polariusz | Added BrokerUser |
//
// # Structure:
// - {"BrokerUserIds":{"BrokerId":<B>, "UserId":<U>},"Topics":[<T>]}
//   - <B> : The ID of the Broker ROW matched from the BrokerId from PostCredentialsHandler()'s brokerId
//   - <U> : The ID of the User ROW matched from the BrokerId from PostCredentialsHandler()'s brokerId
//   - <T> : String array of topics
//
// # Used in
// - PostTopicSubscribeHandler()
// - PostTopicUnsubscribeHandler()
// - PostTopicMarkFavourites()
// - PostTopicUnmarkFavourites()
//
// # Author
// - Polariusz
type TopicsWrapper struct {
	BrokerUserIDs BrokerUser
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
// - PostTopicMarkFavourites()
// - PostTopicUnmarkFavourites()
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
// | 2025-06-06     | Polariusz | Integrated Database    |
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
//   - {"result":{<TOPIC-N>:{"Status":"Fine","Message":"Subscribed to the topic"}}}
// - 207 (Multi Status): JSON
//   - {"result":{<TOPIC-N>:{"Status":<STATUS-N>,"Message":<MESSAGE-N>}}}
// - 400 (Bad Request): JSON
//   - {"badJson":`const BADJSON`}
// - 401 (Unauthorized): JSON
//   - {"Unauthorized":"The MQTT-Client is not connected to any brokers."}
// - 500 (Internal Server Error): JSON
//   - {"InternalServerError":"Error while selecting topics from the database","Error":"<E>"}
//     - <E> : SQL-Error message
//
// # Author
// - Polariusz
func PostTopicSubscribeHandler(serverState *ServerState) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if serverState.mqttClient == nil || !serverState.mqttClient.IsConnected() {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"Unauthorized": "The MQTT-Client is not connected to any brokers.",
			})
		}

		var subscribeTopics TopicsWrapper

		if err := c.BodyParser(&subscribeTopics); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"badJson": BADJSON,
			})
		}

		dbTopicList, err := database.SelectTopicsByBrokerId(serverState.con, subscribeTopics.BrokerUserIDs.BrokerId)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"InternalServerError": "Error while selecting topics from the database",
				"Error": err.Error(),
			})
		}

		dbSubscribedTopicList, err := database.SelectSubscribedTopics(serverState.con, subscribeTopics.BrokerUserIDs.BrokerId, subscribeTopics.BrokerUserIDs.UserId)

		topicResult := make(map[string]TopicResult)
		atLeastOneBadTopic := false

		for _, toSubTopic := range subscribeTopics.Topics {
			isSubscribed := false
			for _, SubTopic := range dbSubscribedTopicList {
				if toSubTopic == SubTopic.Topic {
					// TOPIC IS SUBSCRIBED
					isSubscribed = true
				}
			}
			
			isKnown := false
			knownTopicId := -1
			for _, dbTopic := range dbTopicList {
				if toSubTopic == dbTopic.Topic {
					// TOPIC IS KNOWN
					isKnown = true
					knownTopicId = dbTopic.Id
				}
			}

			if !isKnown {
				// SUBSCRIBE
				if token := serverState.mqttClient.Subscribe(toSubTopic, 0, nil); token.Wait() && token.Error() != nil {
					fmt.Printf("ERROR: Subscribtion to topic %s failed!\n", toSubTopic)
					continue
				}
				// INSERT TO TOPIC
				topicId, err := database.InsertNewTopic(serverState.con, database.InsertTopic{subscribeTopics.BrokerUserIDs.BrokerId, toSubTopic})
				if err != nil {
					fmt.Printf("Error in InsertNewTopic\n")
					atLeastOneBadTopic = true
					topicResult[toSubTopic] = TopicResult{"BigError", err.Error()}
					continue
				}
				// INSERT TO USERTOPICSUBSCRIBED
				if err := database.SubscribeTopic(serverState.con, subscribeTopics.BrokerUserIDs.BrokerId, subscribeTopics.BrokerUserIDs.UserId, topicId); err != nil {
					fmt.Printf("Error in SubscribeTopic\n")
					atLeastOneBadTopic = true
					topicResult[toSubTopic] = TopicResult{"BigError", err.Error()}
					continue
				}
				topicResult[toSubTopic] = TopicResult{"Fine", "Subscribed to the topic"}
			} else if !isSubscribed {
				// SUBSCRIBE
				if token := serverState.mqttClient.Subscribe(toSubTopic, 0, nil); token.Wait() && token.Error() != nil {
					fmt.Printf("ERROR: Subscribtion to topic %s failed!\n", toSubTopic)
					continue
				}
				// INSERT TO USERTOPICSUBSCRIBED
				err := database.SubscribeTopic(serverState.con, subscribeTopics.BrokerUserIDs.BrokerId, subscribeTopics.BrokerUserIDs.UserId, knownTopicId)
				if err != nil {
					fmt.Printf("Error in SubscribeTopic\n")
					atLeastOneBadTopic = true
					topicResult[toSubTopic] = TopicResult{"BigError", err.Error()}
					continue
				}
				topicResult[toSubTopic] = TopicResult{"Fine", "Subscribed to the topic"}
			} else {
				// WHAT
				atLeastOneBadTopic = true
				topicResult[toSubTopic] = TopicResult{"What", "The topic is already subscribed"}
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
// | 2025-06-05     | Polariusz | Integrated database    |
// | 2025-06-07     | Polariusz | UserTopicSubscribed    |
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
//   - {"result":{<TOPIC-N>:{"Status":"Fine","Message":"Unsubscribed successfully"}}}
// - 207 (Multi Status): JSON
//   - {"result":{<TOPIC-N>:{"Status":<STATUS-N>,"Message":<MESSAGE-N>}}}
// - 400 (Bad Request): JSON
//   - {"badJson":`const BADJSON`}
// - 401 (Unauthorized): JSON
//   - {"401":"You fool!"}
// - 500 (Internal Server Error): JSON
//   - {"InternalServerError":"Error while selecting topics from database","Error":"<SQL-ERROR-MESSAGE>"}
//
// # Author
// - Polariusz
func PostTopicUnsubscribeHandler(serverState *ServerState) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if serverState.mqttClient == nil || !serverState.mqttClient.IsConnected() {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"Unauthorized": "The MQTT-Client is not connected with the Broker.",
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

		dbSubscribedTopicList, err := database.SelectSubscribedTopics(serverState.con, unsubscribeTopics.BrokerUserIDs.BrokerId, unsubscribeTopics.BrokerUserIDs.UserId)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"InternalServerError" : "Error while selecting topics from database",
				"Error" : err.Error(),
			})
		}

		for _, toUnsubTopic := range unsubscribeTopics.Topics {
			isSubscribed := false
			for _, subTopic := range dbSubscribedTopicList {
				if toUnsubTopic == subTopic.Topic {
					isSubscribed = true
					if err := database.UnsubscribeTopic(serverState.con, unsubscribeTopics.BrokerUserIDs.BrokerId, unsubscribeTopics.BrokerUserIDs.UserId, subTopic.Id); err != nil {
						atLeastOneBadTopic = true
						topicResult[toUnsubTopic] = TopicResult{"BigError", err.Error()}
						continue
					}
					topicResult[toUnsubTopic] = TopicResult{"Fine", "Unsubscribed to the topic"}
				}
			}
			if !isSubscribed {
				atLeastOneBadTopic = true
				topicResult[toUnsubTopic] = TopicResult{"What", "The topic is not subscribed"}
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

// | Date of change | By        | Comment       |
// +----------------+-----------+---------------+
// |                | Polariusz | Created       |
// | 2025-05-13     | Polariusz | Documentation |
// | 2025-06-06     | Polariusz | Integrated DB |
//
// # Method-Type
// - Handler
//
// # Description
// - The method shall be a handler that allows to get a list of subscribed topics. These topics will be strings or simply string array.
// - The method shall return a 200 (Ok) with the subscribed topics.
// - The method shall accept a jsonified structure that follows the struct BrokerUser.
//
// # Usage
// - Call declared by the routing method addRoutes() URL with the GET-Method.
//
// # Returns
// - 200 (Ok): JSON
//   - {"topics":`[]database.SelectTopic`}
// - 400 (Bad Request): JSON
//   - {"badJson":`const BADJSON`}
// - 401 (Unauthorized): JSON
//   - The go server was never connected to the MQTT-Broker.
//   - {"Unauthorized":"The MQTT-Client is not connected to any brokers."}
// - 500 (Internal Server Error): JSON
//   - {"InternalServerError":"Error while selecting subscribed topics from database","Error":"<SQL-ERROR-MESSAGE>"}
//
// # Author
// - Polariusz
func GetTopicSubscribedHandler(serverState *ServerState) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if serverState.mqttClient == nil || !serverState.mqttClient.IsConnected() {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"Unauthorized": "The MQTT-Client is not connected to any brokers.",
			})
		}

		var brokerUser BrokerUser

		if err := c.BodyParser(&brokerUser); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"badJson": BADJSON,
			})
		}

		topicList, err := database.SelectSubscribedTopics(serverState.con, brokerUser.BrokerId, brokerUser.UserId)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"InternalServerError" : "Error while selecting subscribed topics from database",
				"Error" : err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"topics": topicList,
		})
	}
}

// | Date of change | By        | Comment          |
// +----------------+-----------+------------------+
// |                | Polariusz | Created          |
// | 2025-05-13     | Polariusz | Documentation    |
// | 2025-06-06     | Polariusz | Added BrokerUser |
//
// # Structure:
// - {"BrokerUserIds":{"BrokerId":<B>,"UserId":<U>},"Topic":<T>,"Message":"<M>"}
//   - <T>: Topic
//   - <M>: Message
//
// # Used in
// - PostTopicSendMessageHandler()
//
// # Author
// - Polariusz
type MessageWrapper struct {
	BrokerUserIds BrokerUser
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
		if serverState.mqttClient == nil || !serverState.mqttClient.IsConnected() {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"Unauthorized": "The MQTT-Client is not connected to any brokers.",
			})
		}

		var messageWrapper MessageWrapper

		if err := c.BodyParser(&messageWrapper); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"badJson": BADJSON,
			})
		}

		user, err := database.SelectUserById(serverState.con, messageWrapper.BrokerUserIds.UserId)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"Internal Server Error": "Error while selecting User by ID",
				"Error" : err.Error(),
			})
		}

		// TODO: Validate topic and message!

		// TODO: This can be changed to check if the MQTT-Broker responds! Publish() method returns a token, and the token has method Wait() that waits for the respose and Error() that has either nil or an actual error.
		serverState.mqttClient.Publish(messageWrapper.Topic, 0, false, messageBuilder(user.ClientId, messageWrapper.Message))

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
//   - {"ServiceUnavailable":"The MQTT-Client is not connected to any broker"}
//
// # Author
// - Polariusz
func GetPingHandler(serverState *ServerState) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Is it connected?
		if serverState.mqttClient == nil || serverState.mqttClient.IsConnected() {
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
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"ServiceUnavailable": "The MQTT-Client is not connected to any broker",
			})
		}
	}
}

// | Date of change | By        | Comment                 |
// +----------------+-----------+-------------------------+
// | 2025-05-14     | Tibbyx    | Created & Documentation |
// | 2025-06-06     | Polariusz | Integrated with DB      |
//
// # Method-Type
// - MQTT Handler Factory
//
// # Description
// - The method shall create and return an MQTT message handler.
// - The handler processes incoming MQTT messages from subscribed topics.
// - The handler uses the JsonPublishString structure for messages
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
func createMessageHandler(serverState *ServerState, brokerId int) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		topic := msg.Topic()
		payload := msg.Payload()
		qos := msg.Qos()
		topicId := -1

		var jsonPublishMessage JsonPublishMessage
		if err := json.Unmarshal(payload, &jsonPublishMessage); err != nil {
			jsonPublishMessage.ClientId = "Unknown"
			jsonPublishMessage.Message = string(payload)
		}

		topicList, err := database.SelectTopicsByBrokerId(serverState.con, brokerId)
		if err != nil {
			fmt.Printf("Error while selecting topics by broker id and user id\nError: %s\n", err)
			return
		}
		for _, dbTopic := range topicList {
			if dbTopic.Topic == topic {
				topicId = dbTopic.Id
			}
		}
		// NOTE: The TopicId shouldn't be -1 still, because we only get messages with targetted topic when we are subscribed to these topics, so we know these topics!

		var userId int
		user, err := database.SelectUserByClientIdAndBrokerId(serverState.con, jsonPublishMessage.ClientId, brokerId)
		if err != nil {
			// No user found! Outsider!
			outsiderUserId, err := database.InsertNewUser(serverState.con, database.InsertUser{brokerId, jsonPublishMessage.ClientId, "", "", true})
			if err != nil {
				fmt.Printf("Error while inserting outsider.\nError: %s\n", err)
				return
			}
			userId = outsiderUserId
		} else {
			userId = user.Id
		}

		insertNewMessage := database.InsertMessage{userId, topicId, brokerId, qos, jsonPublishMessage.Message}

		fmt.Printf("Inserting into Message with arguments: %s", insertNewMessage)

		if err := database.InsertNewMessage(serverState.con, insertNewMessage); err != nil {
			// db error
			fmt.Printf("Error while inserting new message\nError: %s\n", err)
			return
		}
	}
}

// | Date of change | By        | Comment |
// +----------------+-----------+---------+
// | 2025-06-06     | Polariusz | Created |
//
// # Structure:
// - {"BrokerUserIds":{"BrokerId":<B>, "UserId":<U>},"Topics":"<T>","Index":<I>}
//   - <B> : The ID of the Broker ROW matched from the BrokerId from PostCredentialsHandler()'s brokerId
//   - <U> : The ID of the User ROW matched from the BrokerId from PostCredentialsHandler()'s brokerId
//   - <T> : The Topic
//   - <I> : The query index, Select from `database.LIMIT_MESSAGES*index` to `database.LIMIT_MESSAGES*(1+index)` messages. If is less than 0, it will query all messages.
//
// # Used in
// - GetTopicMessagesHandler()
//
// # Author
// - Polariusz
type TopicWrapper struct {
	BrokerUserIDs BrokerUser
	Topic string
	Index int
}

// | Date of change | By        | Comment                 |
// +----------------+-----------+-------------------------+
// | 2025-05-14     | Tibbyx    | Created & Documentation |
// | 2025-06-06     | Polariusz | Integrated with DB      |
//
// # Method-Type
// - Handler
//
// # Description
// - The method shall be a handler that returns all stored messages for a specific topic.
// - The messages must have previously been received through an active MQTT subscription.
// - The topic is provided in the --data JSON
//
// # Usage
// - Call declared by the routing method addRoutes() URL with the GET-Method.
// - A data must be included that matches the structure of `TopicWrapper`.
//
// # Returns
// - 200 (Ok): JSON
//   - {"topic": "<topic-name>", "messages": [<`database.SelectMessage`>]}
// - 400 (Bad Request): JSON
//   - {"badJson":"`BADJSON`"}
//   - {"terribleJson":"The arguments in the json structure are missing"}
//   - {"terribleJson":"The argument `Topic` does not match the database."}
// - 401 (Unauthorized): JSON
//   - {"Unauthorized": "The MQTT-Client is not connected to any brokers."}
// - 500 (Internal Server Error): JSON
//   - {"InternalServerError":"Error while <Where> <Arguments>", "Error":"<SQL-Error>"}
//
// # Author
// - Tibbyx
func GetTopicMessagesHandler(serverState *ServerState) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if serverState.mqttClient == nil || !serverState.mqttClient.IsConnected() {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"Unauthorized": "The MQTT-Client is not connected to any brokers.",
			})
		}
		var topicWrapper TopicWrapper
		if err := c.BodyParser(&topicWrapper); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"badJson": BADJSON,
			})
		}
		if topicWrapper.BrokerUserIDs.BrokerId < 0 || topicWrapper.BrokerUserIDs.UserId < 0 || topicWrapper.Topic == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"terribleJson": "The arguments in the json structure are missing",
			})
		}

		topicId := -1
		topicList, err := database.SelectTopicsByBrokerId(serverState.con, topicWrapper.BrokerUserIDs.BrokerId)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"InternalServerError" : "Error while selecting topics matched with broker and user id",
				"Error" : err.Error(),
			})
		}
		for _, dbTopic := range topicList {
			if dbTopic.Topic == topicWrapper.Topic {
				topicId = dbTopic.Id
				break
			}
		}
		if topicId == -1 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"terribleJson": "The argument `Topic` does not match the database.",
			})
		}

		var messageList []database.SelectMessage
		if topicWrapper.Index < 0 {
			messageList, err = database.SelectMessagesByTopicIdAndBrokerId(serverState.con, topicId, topicWrapper.BrokerUserIDs.BrokerId)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"InternalServerError" : "Error while selecting messages matched with topic and broker",
					"Error" : err.Error(),
				})
			}
		} else {
			messageList, err = database.SelectMessagesByTopicIdBrokerIdAndIndex(serverState.con, topicId, topicWrapper.BrokerUserIDs.BrokerId, topicWrapper.Index)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"InternalServerError" : "Error while selecting messages matched with topic, broker and index",
					"Error" : err.Error(),
				})
			}
		}

		return c.JSON(fiber.Map{
			"topic": topicWrapper.Topic,
			"messages": messageList,
		})
	}
}

// | Date of change | By        | Comment       |
// +----------------+-----------+---------------|
// | 2025-05-16     | Polariusz | Created       |
// | 2025-06-07     | Polariusz | Integrated DB |
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
// - The client shall post a JSON that matches the structure of `BrokerUser`.
//
// # Returns
// - 200 (Ok): JSON
//   - {"Topics":[<TOPIC-N>]}
// - 400 (Bad Request): JSON
//   - {"badJson":`const BADJSON`}
// - 401 (Unauthorized): JSON
//   - {"Unauthorized":"The MQTT-Client is not connected to any brokers"}
// - 500 (Internal Server Error): JSON
//   - {"InternalServerError":"<SQL-ERROR>"}
//
// # Author
// - Polariusz
func GetTopicAllKnownHandler(serverState *ServerState) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if serverState.mqttClient == nil || !serverState.mqttClient.IsConnected() {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				// TODO: Explain the message a bit more
				"Unauthorized": "The MQTT-Client is not connected to any brokers",
			})
		}

		var brokerUser BrokerUser

		if err := c.BodyParser(&brokerUser); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"badJson": BADJSON,
			})
		}

		topicList, err := database.SelectTopicsByBrokerId(serverState.con, brokerUser.BrokerId)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"InternalServerError": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"Topics": topicList,
		})
	}
}

// | Date of change | By        | Comment                        |
// +----------------+-----------+--------------------------------+
// | 2025-05-16     | Polariusz | Created                        |
// | 2025-06-07     | Polariusz | Changed the connection checker |
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
		if serverState.mqttClient == nil || !serverState.mqttClient.IsConnected() {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"BadRequest": "The server isn't even connected to any MQTT-Brokers",
			})
		}

		serverState.mqttClient.Disconnect(250)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"Fine": "The MQTT-Client disconnected from the broker.",
		})
	}
}

// | Date of change | By        | Comment            |
// +----------------+-----------+--------------------+
// | 2025-05-18     | Polariusz | Created            |
// | 2025-06-07     | Polariusz | Integrated with DB |
//
// # Method-Type
// - Handler
//
// # Description
// - The fiber.Handler shall append a topic to the favourite list.
// - The function shall accept a json data which contains a list of Topics that the user wishes to mark as favourites.
//
// # Usage
// - Call declared by the routing method addRoutes() URL with the GET-Method.
// - The client shall post a JSON that matches the structure of `TopicsWrapper`.
//
// # Returns
// - 200 (OK): JSON
//   - {"result":{<TOPIC-N>:{"Status":"Fine","Message":"Marked topic as Favourite"}}}
// - 207 (Multi Status): JSON
//   - {"result":{<TOPIC-N>:{"Status":"<STATUS-N>","Message":"<MESSAGE-N>"}}}
// - 400 (Bad Request): JSON
//   - {"badJson": "<const BADJSON>"}
//   - {"terribleJson": "Arguments are not valid"}
// - 401 (Unauthorized): JSON
//   - {"Message": "Authenticate yourself first!"}
//
// # Author
// - Polariusz
func PostTopicFavouritesMark(serverState *ServerState) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if serverState.mqttClient == nil || !serverState.mqttClient.IsConnected() {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"Message": "Authenticate yourself first!",
			})
		}

		var markTopics TopicsWrapper

		if err := c.BodyParser(&markTopics); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"badJson": BADJSON,
			})
		}

		if markTopics.BrokerUserIDs.BrokerId < 0 || markTopics.BrokerUserIDs.UserId < 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"terribleJson": "Arguments are not valid",
			})
		}

		favTopicList, err := database.SelectFavouriteTopicsByBrokerIdAndUserId(serverState.con, markTopics.BrokerUserIDs.BrokerId, markTopics.BrokerUserIDs.UserId)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"InternalServerError": err.Error(),
			})
		}

		dbTopicList, err := database.SelectTopicsByBrokerId(serverState.con, markTopics.BrokerUserIDs.BrokerId)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"InternalServerError": err.Error(),
			})
		}

		topicResult := make(map[string]TopicResult)
		atLeastOneBadTopic := false

		for _, markTopic := range markTopics.Topics {
			isSubscribed := false
			for _, dbFavTopic := range favTopicList {
				if markTopic == dbFavTopic.Topic {
					atLeastOneBadTopic = false
					isSubscribed = true
					break;
				}
			}
			if isSubscribed {
				atLeastOneBadTopic = true
				topicResult[markTopic] = TopicResult{"What", "The topic is marked as favourite"}
				continue
			}

			isKnown := false
			for _, dbTopic := range dbTopicList {
				if markTopic == dbTopic.Topic {
					isKnown = true
					if err := database.InsertFavouriteTopic(serverState.con, markTopics.BrokerUserIDs.BrokerId, markTopics.BrokerUserIDs.UserId, dbTopic.Id); err != nil {
						atLeastOneBadTopic = true
						topicResult[markTopic] = TopicResult{"ServerError", err.Error()}
					} else {
						topicResult[markTopic] = TopicResult{"Fine", "Marked topic as favourite"}
					}
					break
				}
			}

			if !isKnown {
				atLeastOneBadTopic = true
				topicResult[markTopic] = TopicResult{"ClientError", "The topic is not known"}
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

// | Date of change | By        | Comment |
// +----------------+-----------+---------|
// | 2025-05-18     | Polariusz | Created |
//
// # Method-Type
// - Handler
//
// # Description
// - The fiber.Handler shall delete a topic from the favourite list.
// - The function shall accept a json data which contains a list of Topics that the user wishes to unmark from favourites.
//
// # Usage
// - Call declared by the routing method addRoutes() URL with the GET-Method.
// - The client shall post a JSON that matches the structure of `TopicsWrapper`.
//
// # Returns
// - 200 (OK): JSON
//   - {"result":{<TOPIC-N>:{"Status":"Fine","Message":"Unmarked topic from favourite list"}}}
// - 207 (Multi Status): JSON
//   - {"result":{<TOPIC-N>:{"Status":"<STATUS-N>","Message":"<MESSAGE-N>"}}}
// - 400 (Bad Request): JSON
//   - {"badJson": "<const BADJSON>"}
//   - {"terribleJSON":"Arguments are not valid"}
// - 401 (Unauthorized): JSON
//   - {"Message": "Authenticate yourself first!"}
// - 500 (Internal Server Error): JSON
//   - {"InternalServerError":"<SQL-ERROR>"}
//
// # Author
// - Polariusz
func PostTopicFavouritesUnmark(serverState *ServerState) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if serverState.mqttClient == nil || !serverState.mqttClient.IsConnected() {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"Message": "Authenticate yourself first!",
			})
		}

		var unmarkTopics TopicsWrapper

		if err := c.BodyParser(&unmarkTopics); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"badJson": BADJSON,
			})
		}

		if unmarkTopics.BrokerUserIDs.BrokerId < 0 || unmarkTopics.BrokerUserIDs.UserId < 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"terribleJson": "Arguments are not valid",
			})
		}

		favTopicList, err := database.SelectFavouriteTopicsByBrokerIdAndUserId(serverState.con, unmarkTopics.BrokerUserIDs.BrokerId, unmarkTopics.BrokerUserIDs.UserId)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"InternalServerError": err.Error(),
			})
		}

		topicResult := make(map[string]TopicResult)
		atLeastOneBadTopic := false

		for _, unmarkTopic := range unmarkTopics.Topics {
			topicFound := false
			for _, dbFavTopic := range favTopicList {
				if unmarkTopic == dbFavTopic.Topic {
					topicFound = true
					if err := database.DeleteFavouriteTopic(serverState.con, dbFavTopic.Id); err != nil {
						atLeastOneBadTopic = true
						topicResult[unmarkTopic] = TopicResult{"ServerError", err.Error()}
					} else {
						topicResult[unmarkTopic] = TopicResult{"Fine", "Unmarked topic from favourite list"}
					}
					break
				}
				if !topicFound {
					atLeastOneBadTopic = true
					topicResult[unmarkTopic] = TopicResult{"What", "The topics isn't on the favourite list"}
				}
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

// | Date of change | By        | Comment            |
// +----------------+-----------+--------------------+
// | 2025-05-18     | Polariusz | Created            |
// | 2025-06-07     | Polariusz | Integrated with DB |
//
// # Method-Type
// - Handler
//
// # Description
// - The fiber.Handler shall return favourite Topics.
//
// # Usage
// - Call declared by the routing method addRoutes() URL with the GET-Method.
// - The client shall post a JSON that matches the structure of `BrokerUser`.
//
// # Returns
// - 200 (OK): JSON
//   - {"Topics":[<SelectFavTopic-N>]}
// - 400 (Bad Request): JSON
//   - {"badJson":`BADJSON`}
//   - {"terribleJSON":"Arguments are not valid"}
// - 401 (Unauthorized): JSON
//   - {"Message": "Authenticate yourself first!"}
// - 500 (Internal Server Error): JSON
//   - {"InternalServerError":"<SQL-ERROR>"}
//
// # Author
// - Polariusz
func GetTopicFavourites(serverState *ServerState) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if serverState.mqttClient == nil || !serverState.mqttClient.IsConnected() {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"Message": "Authenticate yourself first!",
			})
		}

		var brokerUser BrokerUser

		if err := c.BodyParser(&brokerUser); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"badJson": BADJSON,
			})
		}

		if brokerUser.BrokerId < 0 || brokerUser.UserId < 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"terribleJson": "Arguments are not valid",
			})
		}

		favTopicList, err := database.SelectFavouriteTopicsByBrokerIdAndUserId(serverState.con, brokerUser.BrokerId, brokerUser.UserId)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"InternalServerError": err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"Topics": favTopicList,
		})
	}
}
