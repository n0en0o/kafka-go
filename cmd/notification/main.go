package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"

	"kafka-lab/internal/kafka"
	"kafka-lab/internal/models"

	"github.com/IBM/sarama"
)

func main() {

    config := kafka.NewConfig()
    consumer, err := sarama.NewConsumer([]string{kafka.Broker}, config)
    if err != nil {
        log.Fatal(err)
    }
    defer consumer.Close()

    partitionConsumer, err := consumer.ConsumePartition(kafka.TopicNotifications, 0, sarama.OffsetNewest)
    if err != nil {
        log.Fatal(err)
    }
    defer partitionConsumer.Close()

    fmt.Println("Notification Service started")

    sigchan := make(chan os.Signal, 1)
    signal.Notify(sigchan, os.Interrupt)

    for {
        select {
        case msg := <-partitionConsumer.Messages():
            var notification models.NotificationEvent
            if err := json.Unmarshal(msg.Value, &notification); err != nil {
                log.Printf("Parse error: %v", err)
                continue
            }
            
            fmt.Printf("Notification: %s to %s - %s\n",
                notification.Type, notification.UserEmail, notification.Message)

        case err := <-partitionConsumer.Errors():
            log.Printf("Error: %v", err)

        case <-sigchan:
            fmt.Println("Shutting down")
            return
        }
    }
}