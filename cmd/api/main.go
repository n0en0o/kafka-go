package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"kafka-lab/internal/kafka"
	"kafka-lab/internal/models"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

var producer sarama.SyncProducer

func main() {

    config := kafka.NewConfig()
    var err error
    producer, err = sarama.NewSyncProducer([]string{kafka.Broker}, config)
    if err != nil {
        log.Fatal(err)
    }
    defer producer.Close()

    fmt.Println("Producer connected to Kafka")

    r := gin.Default()
    r.POST("/orders", createOrder)
    r.GET("/health", healthCheck)
    
    fmt.Println("API server starting on :8081")
    r.Run(":8081")
}

func createOrder(c *gin.Context) {

    order := generateTestOrder()
    
    sendToKafka(kafka.TopicOrders, order.OrderID, order)
    
    notification := models.NotificationEvent{
        Type:      "order_created",
        UserID:    order.UserID,
        UserEmail: order.UserEmail,
        Message:   "Your order has been created",
        OrderID:   order.OrderID,
        CreatedAt: time.Now(),
    }
    
    sendToKafka(kafka.TopicNotifications, order.OrderID, notification)
    
    c.JSON(http.StatusOK, gin.H{
        "order_id": order.OrderID,
        "status":   "created",
        "total":    order.TotalAmount,
    })
}

func sendToKafka(topic, key string, data interface{}) {
    jsonData, _ := json.Marshal(data)
    
    msg := &sarama.ProducerMessage{
        Topic: topic,
        Key:   sarama.StringEncoder(key),
        Value: sarama.ByteEncoder(jsonData),
    }
    
    partition, offset, err := producer.SendMessage(msg)
    if err != nil {
        log.Printf("Send error: %v", err)
    } else {
        log.Printf("Sent to %s: partition=%d offset=%d", topic, partition, offset)
    }
}


func generateTestOrder() models.OrderEvent {

    orderID := fmt.Sprintf("ORD-%d", time.Now().Unix())
    userID := fmt.Sprintf("USER-%d", rand.Intn(100))
    
    items := []models.Item{
        {ProductID: "P1", Name: "Laptop", Quantity: 1, Price: 1000},
        {ProductID: "P2", Name: "Mouse", Quantity: 1, Price: 50},
    }
    
    total := 0.0
    for _, item := range items {
        total += item.Price * float64(item.Quantity)
    }
    
    return models.OrderEvent{
        OrderID:     orderID,
        UserID:      userID,
        UserName:    "Test User",
        UserEmail:   "test@example.com",
        Items:       items,
        TotalAmount: total,
        Status:      "created",
        CreatedAt:   time.Now(),
    }
}

func healthCheck(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"status": "ok"})
}