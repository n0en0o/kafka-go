package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"kafka-lab/internal/kafka"
	"kafka-lab/internal/models"

	"github.com/IBM/sarama"
)

var totalOrders = 0
var totalRevenue = 0.0

func main() {
    config := kafka.NewConfig()
    consumer, err := sarama.NewConsumer([]string{kafka.Broker}, config)
    if err != nil {
        log.Fatal(err)
    }
    defer consumer.Close()

    partitionConsumer, err := consumer.ConsumePartition(kafka.TopicOrders, 0, sarama.OffsetNewest)
    if err != nil {
        log.Fatal(err)
    }
    defer partitionConsumer.Close()

    fmt.Println("Analytics Service started")

    sigchan := make(chan os.Signal, 1)
    signal.Notify(sigchan, os.Interrupt)

    ticker := time.NewTicker(60 * time.Second)
    defer ticker.Stop()

    go func() {
        for range ticker.C {
            fmt.Printf("\n--- Analytics ---\n")
            fmt.Printf("Total Orders: %d\n", totalOrders)
            fmt.Printf("Total Revenue: $%.2f\n", totalRevenue)
            if totalOrders > 0 {
                fmt.Printf("Average Order: $%.2f\n", totalRevenue/float64(totalOrders))
            }
        }
    }()

    for {
        select {
        case msg := <-partitionConsumer.Messages():
            var order models.OrderEvent
            if err := json.Unmarshal(msg.Value, &order); err != nil {
                log.Printf("Parse error: %v", err)
                continue
            }
            
            totalOrders++
            totalRevenue += order.TotalAmount
            
            fmt.Printf("Analytics: Order %s added $%.2f\n", order.OrderID, order.TotalAmount)

        case err := <-partitionConsumer.Errors():
            log.Printf("Error: %v", err)

        case <-sigchan:
            fmt.Println("Shutting down")
            return
        }
    }
}