package main

import (
	"fmt"
	"log"

	"github.com/IBM/sarama"
)

func main() {
    
    config := sarama.NewConfig()
    config.Version = sarama.V2_8_1_0

    client, err := sarama.NewClient([]string{"localhost:9092"}, config)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    admin, err := sarama.NewClusterAdminFromClient(client)
    if err != nil {
        log.Fatal(err)
    }
    defer admin.Close()

    topics := []string{"orders", "notifications"}
    
    for _, topic := range topics {
        err := admin.CreateTopic(topic, &sarama.TopicDetail{
            NumPartitions:     3,
            ReplicationFactor: 1,
        }, false)
        
        if err != nil {
            fmt.Printf("Topic %s: %v\n", topic, err)
        } else {
            fmt.Printf("Topic %s created\n", topic)
        }
    }
}
