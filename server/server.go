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

func PostConnectHandler(mqttClient *mqtt.Client, mqttUserConfig *MqttUserConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var mqttOpts *mqtt.ClientOptions 

		if err := c.BodyParser(&mqttUserConfig); err != nil {
			return c.SendString("I am nowt sowwy >:3, but an expected! ewwow has happened. Youw weak json! iws of the wwongest fowmat thawt does nowt cowwespond tuwu the stwong awnd independent stwuct! >:P\n")
		}

		printMqttUserConfig(mqttUserConfig)

		// TODO: It would be nice to validate the IP before this

		mqttOpts = mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s:1883", mqttUserConfig.Ip)).SetClientID(mqttUserConfig.ClientId)
		mqttOpts.SetKeepAlive(2 * time.Second)
		mqttOpts.SetPingTimeout(1 * time.Second)

		// TODO: Figure out how to use the handler.
		// mqttState.mqttOpts.SetDefaultPublishHandler(<INSERT HANDLER HERE>)

		*mqttClient = mqtt.NewClient(mqttOpts)

		if token := (*mqttClient).Connect(); token.Wait() && token.Error() != nil {
			return c.SendString(fmt.Sprintf("Failure at connecting to IP: %s Topic: %s.\nError Message: %s\n",	mqttUserConfig.Ip, mqttUserConfig.Topic, token.Error()))
		}

		return c.SendString("Connected to the Json config\n")
	}
}

func PostDisconnectHandler(mqttClient *mqtt.Client, mqttUserConfig *MqttUserConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if token := (*mqttClient).Unsubscribe(mqttUserConfig.Topic); token.Wait() && token.Error() != nil {
			return c.SendString(fmt.Sprintf("Failure at unsubscribing from topic %s\nError Message: %s\n", mqttUserConfig.Topic, token.Error()))
		}

		(*mqttClient).Disconnect(250)
		return c.SendString("Disconnected")
	}
}

type Number struct {
	Id int
}

func main() {
	server := fiber.New()

	var mqttUserConfig MqttUserConfig
	var mqttClient mqtt.Client

	server.Post("/mqtt/connect", PostConnectHandler(&mqttClient, &mqttUserConfig))
	server.Post("/mqtt/disconnect", PostDisconnectHandler(&mqttClient, &mqttUserConfig))

	server.Listen(":3000")
}
