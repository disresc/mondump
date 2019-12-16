package main

import (
	"fmt"
	"os"
	"time"

	drModels "github.com/disresc/lib/models"
	drReceiver "github.com/disresc/lib/receiver"
)

func handle(event *drModels.Event) {
	fmt.Printf("New Event: %v", event)
}

func main() {
	name, found := os.LookupEnv("name")
	if !found {
		name = "mondump"
	}

	receiver := drReceiver.NewService(name)
	receiver.Start()

	sendPeriodicRequest(receiver, "hosts", "kvmtop-cpu")
	sendPeriodicRequest(receiver, "ves", "kvmtop-cpu")

	for {
		event := <-receiver.EventChannel()
		handle(event)
	}
}

func sendPeriodicRequest(receiver *drReceiver.Service, source string, transmitter string) {
	ticker := time.NewTicker(5 * time.Second)
	//done := make(chan bool)
	sendRequest(receiver, source, transmitter)
	go func() {
		for {
			select {
			//case <-done:
			//	return
			case <-ticker.C:
				sendRequest(receiver, source, transmitter)
			}
		}
	}()
}

func sendRequest(receiver *drReceiver.Service, source string, transmitter string) {
	request := drModels.Request{
		Timeout:     15,
		Source:      source,
		Transmitter: transmitter,
		Interval:    10,
	}
	receiver.Request(&request)
}
