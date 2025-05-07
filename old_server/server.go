package main

import "github.com/gofiber/fiber/v2"
import "github.com/eclipse/paho.mqtt.golang"
import "fmt"
import "time"

type MqttUserConfig struct {
	Ip string
	ClientId string
	Topic string
}

func printMqttUserConfig(mqttUserConfig *MqttUserConfig) {
	fmt.Println("ip       : ", mqttUserConfig.Ip)
	fmt.Println("clientId : ", mqttUserConfig.ClientId)
	fmt.Println("topic    : ", mqttUserConfig.Topic)
}

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

func PostConnectHandler(mqttClient *mqtt.Client, mqttUserConfig *MqttUserConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var mqttOpts *mqtt.ClientOptions 

		if err := c.BodyParser(&mqttUserConfig); err != nil {
			return c.SendString("I am nowt sowwy >:3, but an expected! ewwow has happened. Youw weak json! iws of the wwongest fowmat thawt does nowt cowwespond tuwu the stwong awnd independent stwuct! >:P\n")
		}

		printMqttUserConfig(mqttUserConfig)

		// TODO: It would be nice to validate the IP before this

		//mqttOpts = mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s:1883", mqttUserConfig.Ip)).SetClientID(mqttUserConfig.ClientId)
		mqttOpts = mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://test.mosquitto.org:1883")).SetClientID(mqttUserConfig.ClientId)
		mqttOpts.SetKeepAlive(2 * time.Second)
		mqttOpts.SetPingTimeout(1 * time.Second)

		// TODO: Figure out how to use the handler.
		mqttOpts.SetDefaultPublishHandler(messagePubHandler)

		*mqttClient = mqtt.NewClient(mqttOpts)

		if token := (*mqttClient).Connect(); token.Wait() && token.Error() != nil {
			return c.SendString(fmt.Sprintf("Failure at connecting to IP: %s Topic: %s.\nError Message: %s\n",	mqttUserConfig.Ip, mqttUserConfig.Topic, token.Error()))
		}

		return c.SendString("Connected to the Json config\n")
	}
}

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

func main() {
	server := fiber.New()

	var mqttUserConfig MqttUserConfig
	var mqttClient mqtt.Client

	server.Post("/mqtt/connect", PostConnectHandler(&mqttClient, &mqttUserConfig))
	server.Post("/mqtt/subscribe", PostSubscribeHandler(&mqttClient, &mqttUserConfig))
	server.Post("/mqtt/message/:message", PostMessageHandler(&mqttClient, &mqttUserConfig))
	server.Post("/mqtt/disconnect", PostDisconnectHandler(&mqttClient, &mqttUserConfig))

	server.Listen(":3000")
}
