package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const (
	streamName   = "Messages"
	subjects     = "app.>"
	consumerName = "SimpleConsumer"
)

func main() {

	nc, err := nats.Connect(os.Getenv("NATS_URL"))
	if err != nil {
		log.Fatal(err)
	}

	defer nc.Close()

	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatal("create js client: ", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cfgStream := jetstream.StreamConfig{
		Replicas:    3,
		Name:        streamName,
		Subjects:    []string{subjects},
		Storage:     jetstream.FileStorage,
		Retention:   jetstream.WorkQueuePolicy,
		AllowDirect: true,
	}

	_, err = js.CreateOrUpdateStream(ctx, cfgStream)
	if err != nil {
		log.Fatal("create stream : ", err)
	}

	cfgConsu := jetstream.ConsumerConfig{
		Name:          consumerName,
		FilterSubject: subjects,
		Durable:       consumerName,
	}
	log.Printf("Registering consumer as %s\n", cfgConsu.Name)

	cons, err := js.CreateConsumer(ctx, cfgStream.Name, cfgConsu)
	if err != nil {
		log.Fatal(err)
	}

	processedMessages := 1

	cc, err := cons.Consume(func(msg jetstream.Msg) {

		err := msg.InProgress()
		if err != nil {
			log.Printf("Error setting message to in progress: %v", err)
		}

		log.Printf("✨ Received a message on topic %s : %s\n", msg.Subject(), string(msg.Data()))

		log.Println("Processing message...")
		time.Sleep(1 * time.Second)
		err = msg.Ack()

		log.Println("Message processed ! 🎉 - " + fmt.Sprint(processedMessages))
		if err != nil {
			log.Printf("Error acknowledging message: %v", err)
		}

		processedMessages++

	})
	if err != nil {
		log.Fatal(err)
	}

	defer cc.Drain()

	log.Println("🚀 Consumer is running")

	select {}

}
