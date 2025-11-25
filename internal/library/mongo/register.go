package imongo

import (
	"context"
	"fmt"
	"sync"

	"friberg/internal/consts"
	"friberg/utility/event"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoConfig struct {
	Address  string
	Database string
}

var (
	mongoClient *mongo.Client
	mongoCfg    *mongoConfig
	clientLock  sync.RWMutex
)

// Register 初始化 MongoClient 并加载配置
func Register() {
	ctx := gctx.GetInitCtx()
	// 1. 加载配置
	fmt.Println("mongo register")
	cfg := &mongoConfig{}
	if err := g.Cfg().MustGet(ctx, "mongo").Scan(cfg); err != nil {
		panic(err)
	}
	if cfg == nil {
		panic(gerror.New("mongo config missing"))
	}
	mongoCfg = cfg

	// 2. 创建客户端
	raw, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.Address))
	if err != nil {
		panic(err)
	}

	clientLock.Lock()
	mongoClient = raw
	clientLock.Unlock()
	initTables() // 初始化表句柄
	// 3. 注册服务器关闭事件
	event.Event().Register(consts.EventServerClose, func(ctx context.Context, args ...interface{}) {
		_ = disconnect(ctx)
	})
}

// mustClient 返回单例 Client
func mustClient() *mongo.Client {
	clientLock.RLock()
	defer clientLock.RUnlock()
	if mongoClient == nil {
		panic("mongo client not initialized")
	}
	return mongoClient
}

// disconnect 优雅关闭 MongoClient
func disconnect(ctx context.Context) error {
	clientLock.RLock()
	defer clientLock.RUnlock()
	if mongoClient == nil {
		return nil
	}
	return mongoClient.Disconnect(ctx)
}

// 获取全局数据库（仅包内使用）
func getDatabase() *mongo.Database {
	return mustClient().Database(mongoCfg.Database)
}
