package main

import (
	"context"
	"fmt"
	"log"

	"github.com/Popov-Dmitriy-Ivanovich/genmilk_backend/models"
	"github.com/segmentio/kafka-go"
)

func main() {
	test := models.Exterior{}
	fmt.Println("Starting server...", test)

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{"kafka:9092"},
		Topic:     "LactationLoadRequest",
		Partition: 0,
		MaxBytes:  10e6, // 10MB
	})
	r.SetOffset(0)

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
	}

	if err := r.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
}
