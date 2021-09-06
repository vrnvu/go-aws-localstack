package cloud

import (
	"context"
	"io"
)

type MessageClient interface {
	// Creates a new long polling queue and returns its URL.
	CreateQueue(ctx context.Context, queueName string, isDLX bool) (string, error)
	// Get a queue ARN.
	QueueARN(ctx context.Context, queueURL string) (string, error)
	// Binds a DLX queue to a normal queue.
	BindDLX(ctx context.Context, queueURL, dlxARN string) error
	// Send a message to queue and returns its message ID.
	Send(ctx context.Context, req *SendRequest) (string, error)
	// Long polls given amount of messages from a queue.
	Receive(ctx context.Context, queueURL string) (*Message, error)
	// Deletes a message from a queue.
	Delete(ctx context.Context, queueURL, rcvHandle string) error
}

type BucketClient interface {
	// Creates a new bucket.
	Create(ctx context.Context, bucket string) error
	// Upload a new object to a bucket and returns its URL to view/download.
	UploadObject(ctx context.Context, bucket, key string, body io.Reader) (string, error)
	// Downloads an existing object from a bucket.
	DownloadObject(ctx context.Context, bucket, key string, body io.WriterAt) error
	// Deletes an existing object from a bucket.
	DeleteObject(ctx context.Context, bucket, key string) error
	// Lists all objects in a bucket.
	ListObjects(ctx context.Context, bucket string) ([]*Object, error)
}

type PubSubClient interface {
	// Creates a new topic and returns its ARN.
	Create(ctx context.Context, topic string) (string, error)
	// Lists all topics.
	ListTopics(ctx context.Context) ([]*Topic, error)
	// Subscribes a user (e.g. email, phone) to a topic and returns subscription ARN.
	Subscribe(ctx context.Context, endpoint, protocol, topicARN string) (string, error)
	// Lists all subscriptions for a topic.
	ListTopicSubscriptions(ctx context.Context, topicARN string) ([]*Subscription, error)
	// Publishes a message to all subscribers of a topic and returns its message ID.
	Publish(ctx context.Context, message, topicARN string) (string, error)
	// Unsubscribes a topic subscription.
	Unsubscribe(ctx context.Context, subscriptionARN string) error
}
