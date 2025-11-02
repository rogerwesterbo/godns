package valkeyinterface

import "context"

type ValkeyInterface interface {
	GetData(ctx context.Context, key string) (string, error)
	SetData(ctx context.Context, key string, data string) error
	DeleteData(ctx context.Context, key string) error
	ListKeys(ctx context.Context) ([]string, error)
}
