package aws

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/vrnvu/go-aws-localstack/internal/pkg/cloud"
)

// TODO I do not understand this line
var _ cloud.BucketClient = S3{}

type S3 struct {
	timeout    time.Duration
	client     *s3.S3
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
}

func NewS3(session *session.Session, timeout time.Duration) S3 {
	s3manager.NewUploader(session)
	return S3{
		timeout:    timeout,
		client:     s3.New(session),
		uploader:   s3manager.NewUploader(session),
		downloader: s3manager.NewDownloader(session),
	}
}

// Creates a new bucket.
func (s S3) Create(ctx context.Context, bucket string) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if _, err := s.client.HeadBucketWithContext(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	}); err == nil {
		log.Printf("create: bucket %v already exists\n", bucket)
		return nil
	}

	_, err := s.client.CreateBucketWithContext(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}

	if err := s.client.WaitUntilBucketExistsWithContext(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	}); err != nil {
		return fmt.Errorf("create: %w", err)
	}

	return nil
}

// Upload a new object to a bucket and returns its URL to view/download.
func (s S3) UploadObject(ctx context.Context, bucket string, fileName string, body io.Reader) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	output, err := s.uploader.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
		Body:   body,
	})
	if err != nil {
		return "", fmt.Errorf("upload: %w", err)
	}

	if err := s.client.WaitUntilObjectExistsWithContext(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	}); err != nil {
		return "", fmt.Errorf("upload: %w", err)
	}

	return output.Location, nil
}

// Downloads an existing object from a bucket.
func (s S3) DownloadObject(ctx context.Context, bucket string, fileName string, body io.WriterAt) error {
	if _, err := s.downloader.DownloadWithContext(ctx, body, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	}); err != nil {
		return fmt.Errorf("download: %w", err)
	}

	return nil
}

// Deletes an existing object from a bucket.
func (s S3) DeleteObject(ctx context.Context, bucket string, fileName string) error {
	if _, err := s.client.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	}); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	if err := s.client.WaitUntilObjectNotExists(&s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileName),
	}); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

// Lists all objects in a bucket.
func (s S3) ListObjects(ctx context.Context, bucket string) ([]*cloud.Object, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	res, err := s.client.ListObjectsV2WithContext(ctx, &s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
	})

	if err != nil {
		return nil, fmt.Errorf("list: %w", err)
	}

	objects := make([]*cloud.Object, len(res.Contents))
	for i, object := range res.Contents {
		objects[i] = &cloud.Object{
			Key:        *object.Key,
			Size:       *object.Size,
			ModifiedAt: *object.LastModified,
		}
	}
	return objects, nil
}
