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
