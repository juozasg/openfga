package tests

import (
	"context"
	"testing"

	openfgav1 "github.com/openfga/api/proto/openfga/v1"
	"github.com/openfga/openfga/cmd/run"
	"github.com/openfga/openfga/pkg/logger"
	"github.com/openfga/openfga/pkg/testfixtures/storage"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

// TestClientBootstrapper defines a client interface definition that can be used by tests
// to bootstrap OpenFGA resources (stores, models, relationship tuples, etc..) needed to
// execute tests.
type TestClientBootstrapper interface {
	CreateStore(ctx context.Context, in *openfgav1.CreateStoreRequest, opts ...grpc.CallOption) (*openfgav1.CreateStoreResponse, error)
	WriteAuthorizationModel(ctx context.Context, in *openfgav1.WriteAuthorizationModelRequest, opts ...grpc.CallOption) (*openfgav1.WriteAuthorizationModelResponse, error)
	Write(ctx context.Context, in *openfgav1.WriteRequest, opts ...grpc.CallOption) (*openfgav1.WriteResponse, error)
}

func StartServer(t testing.TB, cfg *run.Config) context.CancelFunc {
	logger := logger.MustNewLogger(cfg.Log.Format, cfg.Log.Level)
	serverCtx := &run.ServerContext{Logger: logger}
	return StartServerWithContext(t, cfg, serverCtx)
}

func StartServerWithContext(t testing.TB, cfg *run.Config, serverCtx *run.ServerContext) context.CancelFunc {
	container := storage.RunDatastoreTestContainer(t, cfg.Datastore.Engine)
	cfg.Datastore.URI = container.GetConnectionURI(true)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		err := serverCtx.Run(ctx, cfg)
		require.NoError(t, err)
	}()

	return cancel
}
