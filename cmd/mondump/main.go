package main

import (
	"fmt"
	"os"
	"time"

	drModels "github.com/disresc/lib/models"
	drReceiver "github.com/disresc/lib/receiver"
)

func handle(event *drModels.Event) {
	time := time.Unix(event.GetTimestamp(), 0)
	fmt.Printf("Event from %s at %s\n", event.GetSource(), time)
	for _, item := range event.GetItems() {
		fmt.Printf("\t%s\t%s\t%s\n", item.GetTransmitter(), item.GetMetric(), item.GetValue())
	}
	fmt.Printf("\n")
}

func main() {
	name, found := os.LookupEnv("name")
	if !found {
		name = "mondump"
	}

	receiver := drReceiver.NewService(name)
	receiver.RegisterData("hosts", "kvmtop-net", 10)
	//receiver.RegisterData("ves", "kvmtop-cpu", 10)
	receiver.Start()

	for {
		event := <-receiver.EventChannel()
		handle(event)
	}
}
