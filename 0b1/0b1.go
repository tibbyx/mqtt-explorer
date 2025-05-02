package main

import (
	"fmt"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

type MqttUserOpt struct {
	ClientId string
	Ip string
	Topic string
}

func printMqttUserCon(mqttUserOpt *MqttUserOpt) {
	fmt.Println("clientId : ", mqttUserOpt.ClientId)
	fmt.Println("ip       : ", mqttUserOpt.Ip)
	fmt.Println("topic    : ", mqttUserOpt.Topic)
}

func connectToMqtt(mqttUserOpt MqttUserOpt) int {
	printMqttUserCon(&mqttUserOpt)
	if mqttUserOpt.ClientId == "1" {
		return 69
	}
	return 0
}

func main() {
	app := app.New()
	window := app.NewWindow("MQTT-Explorer.lab")

	clientIdEntry := widget.NewEntry()
	ipEntry := widget.NewEntry()
	topicEntry := widget.NewEntry()
	mqttUserOptResponsiveLabel := widget.NewLabel("")

	mqttUserOptForm := &widget.Form {
		Items: []*widget.FormItem {
			{Text: "Your ID", Widget: clientIdEntry},
			{Text: "IP to connect to", Widget: ipEntry},
			{Text: "Topic to sub to", Widget: topicEntry},
		},
		OnSubmit: func() {
			res := connectToMqtt(MqttUserOpt{clientIdEntry.Text, ipEntry.Text, topicEntry.Text})
			if res == 69 {
				mqttUserOptResponsiveLabel.SetText("You fool!!")
			} else {
				mqttUserOptResponsiveLabel.SetText("")
			}
		},
	}

	grid := container.New(layout.NewGridLayout(2), mqttUserOptForm, mqttUserOptResponsiveLabel)

	window.SetContent(grid)

	window.ShowAndRun()
}
