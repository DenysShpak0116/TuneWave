package repository

import (
	"bytes"
	"context"
	"log/slog"

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
	logger  *slog.Logger
}

func NewFileStorage(cfg *config.Config, logger *slog.Logger) port.FileStorage {
	options := s3.Options{
		Region:      cfg.AWS.Region,
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(cfg.AWS.AccessKey, cfg.AWS.SecretKey, "")),
	}

	client := s3.New(options)

	return &FileStorage{
		region:  cfg.AWS.Region,
		bucket:  cfg.AWS.Bucket,
		storage: client,
		logger:  logger,
	}
}

func (fs *FileStorage) Save(ctx context.Context, key string, buf bytes.Buffer) (string, error) {
	const op = "adapter.repository.FileStorage.Save"
	logger := fs.logger.With(
		slog.String("op", op),
	)

	if _, err := fs.storage.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(fs.bucket),
		Key:    aws.String(key),
		Body:   &buf,
	}); err != nil {
		logger.Error("failed to save file", "key", key)
		return "", err
	}

	songURL := "https://" + fs.bucket + ".s3." + fs.region + ".amazonaws.com/" + key

	logger.Info("succesfully saved file", "url", songURL)
	return songURL, nil
}

func (fs *FileStorage) Remove(ctx context.Context, key string) error {
	const op = "adapter.repository.FileStorage.Remove"
	logger := fs.logger.With(
		slog.String("op", op),
	)

	if _, err := fs.storage.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(fs.bucket),
		Key:    aws.String(key),
	}); err != nil {
		logger.Error("failed to remove file", "key", key)
		return err
	}

	logger.Info("succesfully removed file", "key", key)
	return nil
}

func (fs *FileStorage) Get(ctx context.Context, key string) ([]byte, error) {
	const op = "adapter.repository.FileStorage.Get"
	logger := fs.logger.With(
		slog.String("op", op),
	)

	resp, err := fs.storage.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(fs.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		logger.Error("failed to get file", "key", key)
		return nil, err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(resp.Body); err != nil {
		return nil, err
	}

	logger.Info("succesfully retrieved file", "key", key)
	return buf.Bytes(), nil
}
