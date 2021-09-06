package aws

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/vrnvu/go-aws-localstack/internal/pkg/cloud"
)

var _ cloud.PubSubClient = SNS{}

type SNS struct {
	timeout time.Duration
	client  *sns.SNS
}

func NewSNS(session *session.Session, timeout time.Duration) SNS {
	return SNS{
		timeout: timeout,
		client:  sns.New(session),
	}
}

// Creates a new topic and returns its ARN.

func (s SNS) Create(ctx context.Context, topic string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	res, err := s.client.CreateTopicWithContext(ctx, &sns.CreateTopicInput{
		Name: aws.String(topic),
	})

	if err != nil {
		return "", fmt.Errorf("create: %w", err)
	}

	return *res.TopicArn, nil
}

// Lists all topics.
func (s SNS) ListTopics(ctx context.Context) ([]*cloud.Topic, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	res, err := s.client.ListTopicsWithContext(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("list topics: %w", err)
	}

	topics := make([]*cloud.Topic, len(res.Topics))
	for i, topic := range res.Topics {
		topics[i] = &cloud.Topic{
			ARN: *topic.TopicArn,
		}
	}

	return topics, nil
}

// Subscribes a user (e.g. email, phone) to a topic and returns subscription ARN.
func (s SNS) Subscribe(ctx context.Context, endpoint string, protocol string, topicARN string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	res, err := s.client.SubscribeWithContext(ctx, &sns.SubscribeInput{
		Protocol:              aws.String(protocol),
		Endpoint:              aws.String(endpoint),
		TopicArn:              aws.String(topicARN),
		ReturnSubscriptionArn: aws.Bool(true),
	})
	if err != nil {
		return "", fmt.Errorf("subscribe: %w", err)
	}

	return *res.SubscriptionArn, nil
}

// Lists all subscriptions for a topic.
func (s SNS) ListTopicSubscriptions(ctx context.Context, topicARN string) ([]*cloud.Subscription, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	res, err := s.client.ListSubscriptionsByTopicWithContext(ctx, &sns.ListSubscriptionsByTopicInput{
		NextToken: nil,
		TopicArn:  aws.String(topicARN),
	})

	if err != nil {
		return nil, fmt.Errorf("list topic subscriptions: %w", err)
	}

	subs := make([]*cloud.Subscription, len(res.Subscriptions))
	for i, sub := range res.Subscriptions {
		subs[i] = &cloud.Subscription{
			ARN:      *sub.SubscriptionArn,
			TopicARN: *sub.TopicArn,
		}
	}
	return subs, nil
}

// Publishes a message to all subscribers of a topic and returns its message ID.
func (s SNS) Publish(ctx context.Context, message string, topicARN string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	res, err := s.client.PublishWithContext(ctx, &sns.PublishInput{
		Message:  &message,
		TopicArn: aws.String(topicARN),
	})
	if err != nil {
		return "", fmt.Errorf("publish: %w", err)
	}
	return *res.MessageId, nil
}

// Unsubscribes a topic subscription.
func (s SNS) Unsubscribe(ctx context.Context, subscriptionARN string) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if _, err := s.client.UnsubscribeWithContext(ctx, &sns.UnsubscribeInput{
		SubscriptionArn: aws.String(subscriptionARN),
	}); err != nil {
		return fmt.Errorf("unsubscribe: %w", err)
	}
	return nil
}
