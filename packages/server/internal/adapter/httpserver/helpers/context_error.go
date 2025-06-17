package helpers

import (
	"context"
)

type ctxKey string

const errorKey ctxKey = "api_error"

func SetAPIError(ctx context.Context, err error) context.Context {
	return context.WithValue(ctx, errorKey, err)
}

func GetAPIError(ctx context.Context) error {
	if err, ok := ctx.Value(errorKey).(error); ok {
		return err
	}
	return nil
}
