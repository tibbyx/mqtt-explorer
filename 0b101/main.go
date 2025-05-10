package main

import (
	"github.com/gofiber/fiber/v2"
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"time"
	"slices"
	"strings"
)

const BADJSON = "I am nowt sowwy >:3. An expected! ewwow has happened. Youw weak json! iws of the wwongest fowmat thawt does nowt cowwespond tuwu the stwong awnd independent stwuct! >:P"

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

type ServerState struct {
	userCreds MqttCredentials
	mqttClient mqtt.Client
	subscribedTopics []string 
}

type MqttCredentials struct {
	Ip string
	Port string
	ClientId string
}

func (mc MqttCredentials) dump() {
	fmt.Println("ip       : ", mc.Ip)
	fmt.Println("port     : ", mc.Port)
	fmt.Println("clientId : ", mc.ClientId)
}

func main() {
	server := fiber.New()
	var serverState ServerState

	addRoutes(server, &serverState)

	server.Listen(":3000")
}

func addRoutes(server *fiber.App, serverState *ServerState) {
	server.Post("/credentials", PostCredentialsHandler(serverState))
	server.Post("/topic/subscribe", PostTopicSubscribeHandler(serverState))
	server.Post("/topic/unsubscribe", PostTopicUnsubscribeHandler(serverState))
	server.Get("/topic/subscribed", GetTopicSubscribedHandler(serverState))
	server.Post("/topic/send-message", PostMessageHandler(serverState))
}

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

		mqttOpts.SetDefaultPublishHandler(messagePubHandler)

		mqttClient := mqtt.NewClient(mqttOpts)

		if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
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

/*
 * TODO: The validation need to be improved. Right now it only checks if the argument userCreds has empty strings.
 * Returns:
 * - 0: No error
 * - 1: Ip was deemed incorrect
 * - 2: Port was deemed incorrect
 * - 3: ClientId was deemed incorrect
*/
func validateCredentials(errorMessage *string, userCreds *MqttCredentials) int {
	// VALIDATE IP
	if userCreds.Ip == "" {
		*errorMessage = "IP is incomprehensible"
		return 1
	}

	// VALIDATE PORT 
	if userCreds.Port == "" {
		*errorMessage = "PORT is incomprehensible"
		return 2
	}

	// VALIDATE CLIENTID
	if userCreds.ClientId == "" {
		*errorMessage = "CLIENT-ID is incomprehensible"
		return 3
	}

	return 0
}

type TopicsWrapper struct {
	Topics []string
}

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

		var badTopics []string

		for _, topic := range subscribeTopics.Topics {
			if !slices.Contains(serverState.subscribedTopics, topic) {
				if token := serverState.mqttClient.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
					badTopics = append(badTopics, topic)
				} else {
					serverState.subscribedTopics = append(serverState.subscribedTopics, topic)
				}
			}
		}

		if len(badTopics) != 0 {
			// TODO: The status code is meh, as the function at this point would
			// subscribe to at least some of the requested topics.
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"badJson": "Could not subscribe to these topics",
				"topics": badTopics,
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"goodJson":"Subscribed to the topics",
		})
	}
}

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

		var badTopics []string
		var goodTopics []string

		for _, topic := range unsubscribeTopics.Topics {
			if !slices.Contains(serverState.subscribedTopics, topic) {
				badTopics = append(serverState.subscribedTopics, topic)
			} else {
				serverState.mqttClient.Unsubscribe(topic)
				goodTopics = append(goodTopics, topic)
				for index, stateTopic := range serverState.subscribedTopics {
					if strings.Compare(stateTopic, topic) == 0 {
						serverState.subscribedTopics = append(serverState.subscribedTopics[:index], serverState.subscribedTopics[index+1:]...)
						break
					}
				}
			}
		}

		if len(badTopics) != 0 {
			// TODO: The status code is meh, as the function at this point would
			// subscribe to at least some of the requested topics.
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"badJson":"Some topics were not even subscribed",
				"badTopics":badTopics,
				"unsubscribedTopics":goodTopics,
			})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"goodJson":"Unsubscribed from all",
		})
	}
}

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

type MessageWrapper struct {
	Topic string
	Message string
}

func PostMessageHandler(serverState *ServerState) fiber.Handler {
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

		serverState.mqttClient.Publish(messageWrapper.Topic, 0, false, messageWrapper.Message)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"goodJson":"Message posted",
		})
	}
}
