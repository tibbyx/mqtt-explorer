package main


import (
	"github.com/gofiber/fiber/v2"
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"time"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

/*
func PostDisconnectHandler(mqttClient *mqtt.Client, mqttUserConfig *MqttUserConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if mqttUserConfig.Topic == "" {
			return c.SendString("It wasn't even connected\n")
		}
		if token := (*mqttClient).Unsubscribe(mqttUserConfig.Topic); token.Wait() && token.Error() != nil {
			return c.SendString(fmt.Sprintf("Failure at unsubscribing from topic %s\nError Message: %s\n", mqttUserConfig.Topic, token.Error()))
		}

		(*mqttClient).Disconnect(250)
		mqttUserConfig.Ip = ""
		mqttUserConfig.ClientId = ""
		mqttUserConfig.Topic = ""
		return c.SendString("Disconnected\n")
	}
}

func PostSubscribeHandler(mqttClient *mqtt.Client, mqttUserConfig *MqttUserConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if token := (*mqttClient).Subscribe(mqttUserConfig.Topic, 0, nil); token.Wait() && token.Error() != nil {
			return c.SendString(fmt.Sprintf("You fool! The topic you want to subscribe to is invalid! You utter buffoon!"))
		}
		return c.SendString("Subscribed\n")
	}
}

func PostMessageHandler(mqttClient *mqtt.Client, mqttUserConfig *MqttUserConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		(*mqttClient).Publish(mqttUserConfig.Topic, 0, false, c.Params("message"))
		return nil
	}
}
*/

type ServerState struct {
	userCreds MqttCredentials
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

	/*
	var userCreds MqttCredentials
	var mqttClient mqtt.Client
	*/

	addRoutes(server)

	server.Listen(":3000")
}

func addRoutes(server *fiber.App) {
	server.Post("/credentials", PostCredentialsHandler())
	//server.Post("/mqtt/subscribe", PostSubscribeHandler(&mqttClient, &mqttUserConfig))
	//server.Post("/mqtt/message/:message", PostMessageHandler(&mqttClient, &mqttUserConfig))
	//server.Post("/mqtt/disconnect", PostDisconnectHandler(&mqttClient, &mqttUserConfig))
}

func PostCredentialsHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {

		var userCreds MqttCredentials

		if err := c.BodyParser(&userCreds); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"badJson": "I am nowt sowwy >:3. An expected! ewwow has happened. Youw weak json! iws of the wwongest fowmat thawt does nowt cowwespond tuwu the stwong awnd independent stwuct! >:P",
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

		//return c.SendString(fmt.Sprintf("IP:%s;Port:%s;ClientID:%s", userCreds.Ip, userCreds.Port, userCreds.ClientId))

		mqttOpts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s:%s", userCreds.Ip, userCreds.Port)).SetClientID(userCreds.ClientId)
		//mqttOpts = mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://test.mosquitto.org:1883")).SetClientID(mqttUserConfig.ClientId)
		mqttOpts.SetKeepAlive(2 * time.Second)
		mqttOpts.SetPingTimeout(1 * time.Second)

		mqttOpts.SetDefaultPublishHandler(messagePubHandler)

		mqttClient := mqtt.NewClient(mqttOpts)

		if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"badJson": fmt.Sprintf("Connecting to %s:%s failed\n%s", userCreds.Ip, userCreds.Port, token.Error()),
			})
		}

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
