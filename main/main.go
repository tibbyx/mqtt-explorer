package main

import "github.com/gofiber/fiber/v2"

var ip = ""
var topic = ""

var getIndex = func(c *fiber.Ctx) error {
	return c.SendString("Just Monika.\n")
}

var postIp = func(c *fiber.Ctx) error {
	// TODO: Validate the IP
	ip = c.Params("ip")
	return nil
}

var postTopic = func(c *fiber.Ctx) error {
	// TODO: Validate the topic(??????)
	topic = c.Params("topic")
	return nil
}

var postMessage = func(c *fiber.Ctx) error {
	// TODO: Probably send something else.
	if(ip == "" || topic == "") {
		return c.SendString("You fool, you utter buffoon.\n")
	}
	return c.SendString("IP: " + ip + " Topic: " + topic + " MSG: " + c.Params("message") + "\n")
}

var postReset = func(c *fiber.Ctx) error {
	ip = ""
	topic = ""
	return nil
}

func main() {
	server := fiber.New()

	server.Get("/", getIndex)
	server.Post("/ip/:ip", postIp)
	server.Post("/topic/:topic", postTopic)
	server.Post("/message/:message", postMessage)
	server.Post("/reset", postReset)

	server.Listen(":3000")
}
