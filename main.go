package main

import (
	"fmt"
	"strings"

	"github.com/abiosoft/ishell"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func mqttClient(opts *MQTT.ClientOptions) MQTT.Client {
	var f MQTT.MessageHandler = func(client MQTT.Client, msg MQTT.Message) {
		fmt.Printf("MSG recieved : %s\n", msg.Payload())
	}

	topic := "nn/sensors"

	opts.OnConnect = func(c MQTT.Client) {
		if token := c.Subscribe(topic, 0, f); token.Wait() && token.Error() != nil {
			panic(token.Error())
		}
	}

	client := MQTT.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	} else {
		fmt.Printf("Connected to server\n")
	}
	return client
}

func main() {

	var mqttC MQTT.Client
	// create new shell.
	// by default, new shell includes 'exit', 'help' and 'clear' commands.
	shell := ishell.New()

	// display welcome info.
	shell.Println("MQTT Interactive Shell\n======================")

	shell.AddCmd(&ishell.Cmd{
		Name: "connect",
		Help: "Connect to a MQTT broker",
		Func: func(c *ishell.Context) {
			c.Println("Connecting: ")

			opts := MQTT.NewClientOptions().AddBroker("broker:port")
			opts.SetClientID("test")
			opts.SetPassword("passw0rd")
			opts.SetUsername("username")
			//opts.SetDefaultPublishHandler(f)

			mqttC = mqttClient(opts)
		},
	})

	shell.AddCmd(&ishell.Cmd{
		Name: "send",
		Help: "Send mqtt message (<topic> <message>",
		Func: func(c *ishell.Context) {
			c.Println("Sent: ", strings.Join(c.Args, " "))

			token := mqttC.Publish(c.Args[0], 0, false, strings.Join(c.Args[1:], " "))
			token.Wait()
		},
	})

	// run shell
	shell.Run()
}
