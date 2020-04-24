package utils

import (
	"context"
	"time"

	"cloud.google.com/go/storage"
	"github.com/samedguener/ImageService/errors"
	"github.com/sirupsen/logrus"
)

// InitGCPCloudStorageBucket ...
func InitGCPCloudStorageBucket() {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		logrus.Fatalf("failed to create client: %v", err)
	}

	bucket := client.Bucket(BucketName.Value)
	if _, err = bucket.Attrs(ctx); err == nil {
		logrus.Infof("bucket with name '%s' already exists", BucketName.Value)
		return
	}

	bucketAttrs := &storage.BucketAttrs{
		StorageClass: "STANDARD",
		Location:     "eu",
	}

	ctx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	if err := bucket.Create(ctx, GCPProjectID.Value, bucketAttrs); err != nil {
		logrus.Fatalf("failed to create bucket with name '%s': '%v'", BucketName.Value, err)
	}

	logrus.Infof("bucket %s, created at %s, is located in %s with storage class %s",
		BucketName.Value, bucketAttrs.Created, bucketAttrs.Location, bucketAttrs.StorageClass)

}

// GetBucket ...
func GetBucket() (*storage.BucketHandle, error) {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, errors.Internal.Wrapf(err, "failed to create client")
	}

	bucket := client.Bucket(BucketName.Value)

	return bucket, err
}
