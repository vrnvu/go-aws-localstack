package main

import (
	"fmt"
	"log"
	"time"

	"github.com/vrnvu/go-aws-localstack/internal/bucket"
	"github.com/vrnvu/go-aws-localstack/internal/message"
	"github.com/vrnvu/go-aws-localstack/internal/pkg/cloud/aws"
	"github.com/vrnvu/go-aws-localstack/internal/pubsub"
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

	fmt.Println("test: SQS")
	// Test message using SQS
	message.Message(aws.NewSQS(ses, time.Second*5))

	fmt.Println("test: s3")
	// Test bucket using S3
	bucket.Bucket(aws.NewS3(ses, time.Second*5))

	fmt.Println("test: pubsub")
	// Test pubsub using SNS
	pubsub.PubSub(aws.NewSNS(ses, time.Second*5))
}
