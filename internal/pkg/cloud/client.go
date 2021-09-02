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
	UploadObject(ctx context.Context, bucket, fileName string, body io.Reader) (string, error)
	// Downloads an existing object from a bucket.
	DownloadObject(ctx context.Context, bucket, fileName string, body io.WriterAt) error
	// Deletes an existing object from a bucket.
	DeleteObject(ctx context.Context, bucket, fileName string) error
	// Lists all objects in a bucket.
	ListObjects(ctx context.Context, bucket string) ([]*Object, error)
}
