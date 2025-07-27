//go:generate mockgen -source=file_storage.go -destination=../../adapter/repository/mocks/file_storage_mock.go -package=mocks -typed

package port

import (
	"bytes"
	"context"
)

type FileStorage interface {
	Save(ctx context.Context, key string, buf bytes.Buffer) (string, error)
	Remove(ctx context.Context, key string) error
	Get(ctx context.Context, key string) ([]byte, error)
}
