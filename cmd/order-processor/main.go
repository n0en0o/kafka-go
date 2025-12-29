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
    consumer, err := sarama.NewConsumer([]string{"localhost:9092"}, config)
    if err != nil {
        log.Fatal(err)
    }
    defer consumer.Close()

    partitionConsumer, err := consumer.ConsumePartition(kafka.TopicOrders, 0, sarama.OffsetNewest)
    if err != nil {
        log.Fatal(err)
    }
    defer partitionConsumer.Close()

    fmt.Println("Order Processor started")

    sigchan := make(chan os.Signal, 1)
    signal.Notify(sigchan, os.Interrupt)

    for {
        select {
        case msg := <-partitionConsumer.Messages():
            var order models.OrderEvent
            if err := json.Unmarshal(msg.Value, &order); err != nil {
                log.Printf("Parse error: %v", err)
                continue
            }
            
            fmt.Printf("Processing order: %s, User: %s, Total: $%.2f\n",
                order.OrderID, order.UserID, order.TotalAmount)
            
            order.Status = "processed"
            fmt.Printf("Order %s status: %s\n", order.OrderID, order.Status)

        case err := <-partitionConsumer.Errors():
            log.Printf("Error: %v", err)

        case <-sigchan:
            fmt.Println("Shutting down")
            return
        }
    }
}