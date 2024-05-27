package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const (
	partMiBs = 10
	mByte    = 1024 * 1024
)

// S3 provides a simple interface to store and retrieve data.
type S3 struct {
	bucket           string
	externalEndpoint string
	hostURL          string
	downloader       *manager.Downloader
	uploader         *manager.Uploader
	cli              *s3.Client
}

var _ Storage = (*S3)(nil)

// NewS3 creates a new S3 instance.
func NewS3(
	ctx context.Context,
	hostURL,
	externalEndpoint,
	credentialKey,
	credentialSecret,
	bucket string,
) (*S3, error) {
	resolver := aws.EndpointResolverWithOptionsFunc(
		func(service, _ string, _ ...any) (aws.Endpoint, error) {
			if service == s3.ServiceID && len(hostURL) > 0 {
				// カスタムエンドポイント使用
				return aws.Endpoint{
					URL:               hostURL,
					HostnameImmutable: true,
				}, nil
			}
			// 通常のエンドポイント使用
			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
		})

	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithEndpointResolverWithOptions(resolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(credentialKey, credentialSecret, "")),
	)
	if err != nil {
		return &S3{}, fmt.Errorf("failed to load configuration, %w", err)
	}

	svc := s3.NewFromConfig(cfg, func(options *s3.Options) {
		options.UsePathStyle = true
	})
	downloader := manager.NewDownloader(svc, func(d *manager.Downloader) {
		d.PartSize = partMiBs * mByte
	})
	uploader := manager.NewUploader(svc, func(u *manager.Uploader) {
		u.PartSize = partMiBs * mByte
	})
	return &S3{
		bucket:           bucket,
		hostURL:          hostURL,
		externalEndpoint: externalEndpoint,
		downloader:       downloader,
		uploader:         uploader,
		cli:              svc,
	}, nil
}

// PutObject stores data and returns the key.
func (s *S3) PutObject(
	ctx context.Context, reader io.Reader, key string,
) (string, error) {
	_, err := s.cli.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   reader,
	})
	if err != nil {
		return "", fmt.Errorf("failed to put object, %w", err)
	}
	return fmt.Sprintf("%s/%s/%s", s.externalEndpoint, s.bucket, key), nil
}

// UploadObject stores data and returns the key.
func (s *S3) UploadObject(
	ctx context.Context, reader io.Reader, key string,
) (string, error) {
	_, err := s.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		Body:   reader,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload object, %w", err)
	}
	return fmt.Sprintf("%s/%s/%s", s.externalEndpoint, s.bucket, key), nil
}

// GetObject retrieves data by key.
func (s *S3) GetObject(ctx context.Context, key string) (io.ReadCloser, error) {
	o, err := s.cli.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object, %w", err)
	}
	return o.Body, nil
}

// DownloadObject retrieves data by key.
func (s *S3) DownloadObject(ctx context.Context, writer io.WriterAt, key string) error {
	_, err := s.downloader.Download(ctx, writer, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("failed to download object, %w", err)
	}
	return nil
}

// DeleteObjects deletes data by key.
func (s *S3) DeleteObjects(ctx context.Context, keys []string) error {
	var objectIDs []types.ObjectIdentifier
	for _, key := range keys {
		objectIDs = append(objectIDs, types.ObjectIdentifier{Key: aws.String(key)})
	}
	o, err := s.cli.DeleteObjects(ctx, &s3.DeleteObjectsInput{
		Bucket: aws.String(s.bucket),
		Delete: &types.Delete{Objects: objectIDs},
	})
	if err != nil {
		return fmt.Errorf("failed to delete object, %w", err)
	}
	if len(o.Deleted) != len(keys) {
		return fmt.Errorf("failed to delete object, %w", err)
	}

	return nil
}

// ExistsObject checks if data exists by key.
func (s *S3) ExistsObject(ctx context.Context, key string) (bool, error) {
	_, err := s.cli.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		var e *types.NotFound
		if !errors.As(err, &e) {
			return false, fmt.Errorf("failed to head object, %w", err)
		}
		return false, nil
	}
	return true, nil
}

// GetURLFromKey returns the URL from the key.
func (s *S3) GetURLFromKey(_ context.Context, key string) (string, error) {
	return fmt.Sprintf("%s/%s/%s", s.externalEndpoint, s.bucket, key), nil
}

// GetKeyFromURL returns the key from the URL.
func (s *S3) GetKeyFromURL(_ context.Context, url string) (string, error) {
	if !strings.HasPrefix(url, s.externalEndpoint) {
		return "", fmt.Errorf("invalid URL")
	}
	return strings.TrimPrefix(url, fmt.Sprintf("%s/%s/", s.externalEndpoint, s.bucket)), nil
}
