package s3 // Special thanks via https://docs.digitalocean.com/products/spaces/resources/s3-sdk-examples/

import (
	"bytes"
	"context"
	"errors"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"golang.org/x/exp/slog"

	c "github.com/LuchaComics/cps-backend/config"
	"github.com/LuchaComics/cps-backend/provider/uuid"
)

type S3Storager interface {
	UploadContent(ctx context.Context, path string, content []byte) (string, error)
	BucketExists(ctx context.Context, bucketName string) (bool, error)
}

type s3Storager struct {
	S3Client   *s3.Client
	UUID       uuid.Provider
	Logger     *slog.Logger
	BucketName string
}

// NewStorage connects to a specific S3 bucket instance and returns a connected
// instance structure.
func NewStorage(appConf *c.Conf, logger *slog.Logger, uuidp uuid.Provider) S3Storager {
	// DEVELOPERS NOTE:
	// How can I use the AWS SDK v2 for Go with DigitalOcean Spaces? via https://stackoverflow.com/a/74284205
	logger.Debug("s3 initializing...")

	// STEP 1: initialize the custom `endpoint` we will connect to.
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: appConf.AWS.Endpoint,
		}, nil
	})

	// STEP 2: Configure.
	sdkConfig, err := config.LoadDefaultConfig(
		context.TODO(), config.WithRegion(appConf.AWS.Region),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(appConf.AWS.AccessKey, appConf.AWS.SecretKey, "")),
	)
	if err != nil {
		log.Fatal(err) // We need to crash the program at start to satisfy google wire requirement of having no errors.
	}

	// STEP 3\: Load up s3 instance.
	s3Client := s3.NewFromConfig(sdkConfig)
	logger.Debug("s3 initialized")

	// Create our storage handler.
	s3Storage := &s3Storager{
		S3Client:   s3Client,
		Logger:     logger,
		UUID:       uuidp,
		BucketName: appConf.AWS.BucketName,
	}

	// STEP 4: Connect to the s3 bucket instance and confirm that bucket exists.
	doesExist, err := s3Storage.BucketExists(context.TODO(), appConf.AWS.BucketName)
	if err != nil {
		log.Fatal(err) // We need to crash the program at start to satisfy google wire requirement of having no errors.
	}
	if !doesExist {
		log.Fatal("bucket name does not exist") // We need to crash the program at start to satisfy google wire requirement of having no errors.
	}

	// Return our s3 storage handler.
	return s3Storage
}

func (s *s3Storager) UploadContent(ctx context.Context, path string, content []byte) (string, error) {
	objectKey := s.UUID.NewUUID()
	_, err := s.S3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.BucketName),
		Key:    aws.String(path),
		Body:   bytes.NewReader(content),
	})
	if err != nil {
		return "", err
	}
	return objectKey, nil
}

func (s *s3Storager) BucketExists(ctx context.Context, bucketName string) (bool, error) {
	// Note: https://docs.aws.amazon.com/code-library/latest/ug/go_2_s3_code_examples.html#actions

	_, err := s.S3Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	exists := true
	if err != nil {
		var apiError smithy.APIError
		if errors.As(err, &apiError) {
			switch apiError.(type) {
			case *types.NotFound:
				log.Printf("Bucket %v is available.\n", bucketName)
				exists = false
				err = nil
			default:
				log.Printf("Either you don't have access to bucket %v or another error occurred. "+
					"Here's what happened: %v\n", bucketName, err)
			}
		}
	}

	return exists, err
}
