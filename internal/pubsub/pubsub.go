package pubsub

import (
	"context"
	"fmt"
	"log"

	"github.com/vrnvu/go-aws-localstack/internal/pkg/cloud"
)

func PubSub(client cloud.PubSubClient) {
	ctx := context.Background()

	tARN := create(ctx, client)
	listTopics(ctx, client)
	sARN := subscribe(ctx, client, tARN)
	listTopicSubscriptions(ctx, client, tARN)
	publish(ctx, client, tARN)
	unsubscribe(ctx, client, sARN)
}

func create(ctx context.Context, client cloud.PubSubClient) string {
	arn, err := client.Create(ctx, "welcome-email")
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("create: topic ARN:", arn)

	return arn
}

func listTopics(ctx context.Context, client cloud.PubSubClient) {
	topics, err := client.ListTopics(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("list topics:")
	for _, topic := range topics {
		fmt.Printf("%+v\n", topic)
	}
}

func subscribe(ctx context.Context, client cloud.PubSubClient, topicARN string) string {
	arn, err := client.Subscribe(ctx, "email@example.com", "email", topicARN)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("subscribe: subscription ARN:", arn)

	return arn
}

func listTopicSubscriptions(ctx context.Context, client cloud.PubSubClient, topicARN string) {
	subs, err := client.ListTopicSubscriptions(ctx, topicARN)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("list topic subscriptions:")
	for _, sub := range subs {
		fmt.Printf("%+v\n", sub)
	}
}

func publish(ctx context.Context, client cloud.PubSubClient, topicARN string) {
	id, err := client.Publish(ctx, "hello!", topicARN)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("publish: message ID:", id)
}

func unsubscribe(ctx context.Context, client cloud.PubSubClient, subARN string) {
	if err := client.Unsubscribe(ctx, subARN); err != nil {
		log.Fatalln(err)
	}
	log.Println("unsubscribe: ok")
}
