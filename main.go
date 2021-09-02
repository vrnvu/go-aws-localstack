package main

import (
	"log"
	"time"

	"github.com/vrnvu/go-aws-localstack/internal/message"
	"github.com/vrnvu/go-aws-localstack/internal/pkg/cloud/aws"
)

func main() {
	// Create a session instance.
	ses, err := aws.New(aws.Config{
		Address: "http://localhost:4566",
		Region:  "eu-west-1",
		Profile: "localstack",
		ID:      "test",
		Secret:  "test",
	})
	if err != nil {
		log.Fatalln(err)
	}

	// Test message
	message.Message(aws.NewSQS(ses, time.Second*5))
}
