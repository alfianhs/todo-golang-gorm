package s3repo

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type s3Repo struct {
	bucketName string
	client     *s3.Client
}

func NewS3Repository(contextTimeout time.Duration) S3Repo {
	ctx, cancel := context.WithTimeout(context.Background(), contextTimeout)
	defer cancel()

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID:     viper.GetString("S3_ACCESS_KEY"),
				SecretAccessKey: viper.GetString("S3_SECRET_KEY"),
				SessionToken:    "",
			},
		}),
		config.WithRegion(viper.GetString("S3_REGION")),
		config.WithBaseEndpoint(viper.GetString("S3_ENDPOINT")),
	)
	if err != nil {
		logrus.Error("Load S3 Config Error : ", err)
		return nil
	}

	client := s3.NewFromConfig(cfg)
	return &s3Repo{
		client:     client,
		bucketName: viper.GetString("S3_BUCKET_NAME"),
	}
}

type S3Repo interface {
	UploadFile(ctx context.Context, objectName string, data []byte, mimeType string) (string, error)
}

func (r *s3Repo) UploadFile(ctx context.Context, objectName string, data []byte, mimeType string) (string, error) {
	// upload to bucket s3
	_, err := r.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(r.bucketName),
		Key:         aws.String(objectName),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(mimeType),
		ACL:         types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s/%s/%s", viper.GetString("S3_ENDPOINT"), r.bucketName, objectName)
	return url, nil
}
