package minio

import (
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/Anton-Kraev/gopay"
)

type Client struct {
	bucketName string
	client     *minio.Client
}

func NewClient(ctx context.Context, config Config) (c Client, err error) {
	const op = "minio.NewClient"

	c.bucketName = config.BucketName

	c.client, err = minio.New(config.URL, &minio.Options{
		Creds: credentials.NewStaticV4(config.User, config.Password, ""),
	})
	if err != nil {
		return Client{}, fmt.Errorf("%s: %w", op, err)
	}

	ok, err := c.client.BucketExists(ctx, c.bucketName)
	if err != nil {
		return Client{}, fmt.Errorf("%s: %w", op, err)
	}

	if ok {
		return c, nil
	}

	if err = c.client.MakeBucket(ctx, c.bucketName, minio.MakeBucketOptions{}); err != nil {
		return Client{}, fmt.Errorf("%s: %w", op, err)
	}

	return c, nil
}

func (c Client) GetData(ctx context.Context, id gopay.ID) ([]byte, error) {
	obj, err := c.client.GetObject(ctx, c.bucketName, string(id)+".pdf", minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("minio.Client.GetData: %w", err)
	}

	file, err := io.ReadAll(obj)
	if err != nil {
		return nil, fmt.Errorf("minio.Client.GetData: %w", err)
	}

	return file, nil
}
