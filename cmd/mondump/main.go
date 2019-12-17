package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	drModels "github.com/disresc/lib/models"
	drReceiver "github.com/disresc/lib/receiver"
	"github.com/micro/go-micro/util/log"
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
	data, found := os.LookupEnv("data")
	if !found {
		data = "hosts.kvmtop-cpu.10;ves.kvmtop-cpu.10"
	}

	receiver := drReceiver.NewService(name)

	//receiver.RegisterData("hosts", "kvmtop-cpu", 10)
	//receiver.RegisterData("ves", "kvmtop-cpu", 10)
	dataLines := strings.Split(data, ";")
	for _, line := range dataLines {
		lineParts := strings.Split(line, ".")
		if len(lineParts) != 3 {
			log.Errorf("Invalid data line %s", line)
			return
		}
		source := lineParts[0]
		transmitter := lineParts[1]
		interval, err := strconv.Atoi(lineParts[2])
		if err != nil {
			log.Errorf("Cannot parse interval %s", lineParts[2])
			return
		}
		receiver.RegisterData(source, transmitter, interval)
	}

	receiver.Start()

	for {
		event := <-receiver.EventChannel()
		handle(event)
	}
}
