package imongo

import (
	"context"
	"sync"

	"friberg/internal/consts"
	"friberg/utility/event.go"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoConfig struct {
	Address  string
	Database string
}

var (
	mongoClient *mongo.Client
	clientLock  sync.RWMutex
)

func Register() {
	raw, err := newClient()
	if err != nil {
		panic(err)
	}

	clientLock.Lock()
	mongoClient = raw
	clientLock.Unlock()

	event.Event().Register(consts.EventServerClose, func(ctx context.Context, args ...interface{}) {
		_ = disconnect(ctx)
	})
}

func mustClient() *mongo.Client {
	clientLock.RLock()
	defer clientLock.RUnlock()

	if mongoClient == nil {
		panic("mongo client not initialized")
	}
	return mongoClient
}

func mustLoadConfig(ctx context.Context) *MongoConfig {
	var cfg *MongoConfig
	if err := g.Cfg().MustGet(ctx, "mongo").Scan(&cfg); err != nil {
		panic(err)
	}
	if cfg == nil {
		panic(gerror.New("mongo config missing"))
	}
	return cfg
}

func newClient() (*mongo.Client, error) {
	ctx := gctx.GetInitCtx()
	cfg := mustLoadConfig(ctx)
	return mongo.Connect(ctx, options.Client().ApplyURI(cfg.Address))
}

func disconnect(ctx context.Context) error {
	clientLock.RLock()
	defer clientLock.RUnlock()

	if mongoClient == nil {
		return nil
	}
	return mongoClient.Disconnect(ctx)
}
