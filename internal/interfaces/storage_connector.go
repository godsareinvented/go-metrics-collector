package interfaces

import "context"

type StorageConnector interface {
	Ping(ctx context.Context) (bool, error)
	Close() error
}
