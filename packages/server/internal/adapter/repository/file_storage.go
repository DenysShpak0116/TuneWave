package repository

import (
	"bytes"
	"context"

	"github.com/DenysShpak0116/TuneWave/packages/server/internal/adapter/config"
	"github.com/DenysShpak0116/TuneWave/packages/server/internal/core/port"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type FileStorage struct {
	region  string
	bucket  string
	storage *s3.Client
}

func NewFileStorage(cfg *config.Config) port.FileStorage {
	options := s3.Options{
		Region:      cfg.AWS.Region,
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(cfg.AWS.AccessKey, cfg.AWS.SecretKey, "")),
	}

	client := s3.New(options)

	return &FileStorage{
		region:  cfg.AWS.Region,
		bucket:  cfg.AWS.Bucket,
		storage: client,
	}
}

func (fs *FileStorage) Save(ctx context.Context, key string, buf bytes.Buffer) (string, error) {
	_, err := fs.storage.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(fs.bucket),
		Key:    aws.String(key),
		Body:   &buf,
	})
	if err != nil {
		return "", err
	}

	songURL := "https://" + fs.bucket + ".s3." + fs.region + ".amazonaws.com/" + key
	return songURL, nil
}

func (fs *FileStorage) Remove(ctx context.Context, key string) error {
	_, err := fs.storage.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(fs.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}
	return nil
}

func (fs *FileStorage) Get(ctx context.Context, key string) ([]byte, error) {
	resp, err := fs.storage.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(fs.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
