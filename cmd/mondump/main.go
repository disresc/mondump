package main

import (
	"fmt"
	drModels "github.com/disresc/lib/models"
	drReceiver "github.com/disresc/lib/receiver"
	"os"
)

func handle(event *drModels.Event) {
	fmt.Printf("New Event: %v", event)
}

func main() {
	name, found := os.LookupEnv("name")
	if !found {
		name = "transmitter"
	}
	/*topic, found := os.LookupEnv("topic")
	if !found {
		topic = "monitoring"
	}*/

	receiver := drReceiver.NewService(name)
	receiver.Start()
	for {
		event := <-receiver.EventChannel()
		handle(event)
	}
}
