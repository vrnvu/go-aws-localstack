package bucket

import (
	"context"
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/vrnvu/go-aws-localstack/internal/pkg/cloud"
)

var bucket = "s3-bucket-test"

func Bucket(client cloud.BucketClient) {
	ctx := context.Background()

	name := "./assets/id.txt"

	// Creates a new bucket.
	create(ctx, client)
	// Upload a new object to a bucket and returns its URL to view/download.
	uploadObject(ctx, client, name, "id.txt")
	// Lists all objects in a bucket.
	listObjects(ctx, client)
	// Downloads an existing object from a bucket.
	downloadObject(ctx, client, "id.txt", "/tmp", "id.txt")
	// Deletes an existing object from a bucket.
	deleteObject(ctx, client, "id.txt")
	// Lists all objects in a bucket.
	listObjects(ctx, client)
}

func create(ctx context.Context, client cloud.BucketClient) {
	if err := client.Create(ctx, bucket); err != nil {
		log.Fatalln(err)
	}
	log.Println("create: ok")
}

func uploadObject(ctx context.Context, client cloud.BucketClient, name, key string) {
	file, err := os.Open(name)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	url, err := client.UploadObject(ctx, bucket, key, file)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("upload object:", url)

}

func createPath(path string) {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
}

func downloadObject(ctx context.Context, client cloud.BucketClient, key, path, toFileName string) {
	createPath(path)

	file, err := os.Create(filepath.Join(path, toFileName))
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	if err := client.DownloadObject(ctx, bucket, key, file); err != nil {
		log.Fatalln(err)
	}
	log.Println("download object: ok")
}

func deleteObject(ctx context.Context, client cloud.BucketClient, key string) {
	if err := client.DeleteObject(ctx, bucket, key); err != nil {
		log.Fatalln(err)
	}
	log.Println("delete object: ok")

}

func listObjects(ctx context.Context, client cloud.BucketClient) {
	objects, err := client.ListObjects(ctx, bucket)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("list objects:")
	for _, object := range objects {
		log.Printf("%+v\n", object)
	}
}
