package mongo

import (
	"context"

	"friberg/internal/consts"
	"friberg/utility/event.go"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ==================== 配置结构 ====================
type MongoConfig struct {
	Address  string // MongoDB URL
	Database string // 数据库名
}

// ==================== 客户端管理 ====================
var (
	mongoClient *mongo.Client // 内部使用
	dbRepo      DB            // 对外暴露的 Repository 接口
)

// ==================== 注册入口 ====================
func Register() {
	client, err := NewMongoClient()
	if err != nil {
		panic(err)
	}
	mongoClient = client

	// 初始化 Repository
	dbRepo = &dbImpl{
		collection: "subject", // 默认集合，可根据业务扩展更多集合
	}

	event.Event().Register(consts.EventServerClose, func(ctx context.Context, args ...interface{}) {
		g.Log().Debugf(ctx, "Mongo Event: %+v", args)
		if err := disconnect(ctx); err != nil {
			g.Log().Errorf(ctx, "Mongo Disconnect Error: %v", err)
		}
	})
}

// 外部获取 Repository 接口
func DBRepo() DB {
	if dbRepo == nil {
		panic("mongo not registered, call Register() first")
	}
	return dbRepo
}

// ==================== 配置加载 ====================
func MustLoadConfig(ctx context.Context) *MongoConfig {
	var cfg *MongoConfig
	if err := g.Cfg().MustGet(ctx, "mongo").Scan(&cfg); err != nil {
		panic(err)
	}
	if cfg == nil {
		panic(gerror.New("mongo config missing"))
	}
	return cfg
}

// ==================== Mongo 客户端初始化 ====================
func NewMongoClient() (*mongo.Client, error) {
	ctx := gctx.GetInitCtx()
	cfg := MustLoadConfig(ctx)

	g.Log().Debugf(ctx, "Mongo Config: %+v", cfg)

	clientOptions := options.Client().ApplyURI(cfg.Address)
	return mongo.Connect(ctx, clientOptions)
}

func disconnect(ctx context.Context) error {
	// 类型断言成 *dbImpl，或者在 DB 接口增加 FindByField
	if r, ok := dbRepo.(*dbImpl); ok {
		result, err := r.FindByField(ctx, "name", "Grand Theft Auto V")
		if err != nil {
			g.Log().Errorf(ctx, "FindByField Error: %v", err)
		} else {
			g.Dump(ctx, result)
		}
	}

	return mongoClient.Disconnect(ctx)
}

// ==================== Repository 接口 ====================
type DB interface {
	Create(ctx context.Context, doc interface{}) error
	FindByID(ctx context.Context, id string, result interface{}) error
	// 可以后续继续增加 MongoDB 接口方法，例如 Update, Delete, Find 等
}

// ==================== Repository 实现 ====================
type dbImpl struct {
	collection string
}

// 获取 Collection（内部使用）
func (r *dbImpl) coll(ctx context.Context) *mongo.Collection {
	cfg := MustLoadConfig(ctx)
	return mongoClient.Database(cfg.Database).Collection(r.collection)
}

// 创建文档
func (r *dbImpl) Create(ctx context.Context, doc interface{}) error {
	_, err := r.coll(ctx).InsertOne(ctx, doc)
	return err
}

// 根据 ID 查询文档
func (r *dbImpl) FindByID(ctx context.Context, id string, result interface{}) error {
	return r.coll(ctx).FindOne(ctx, bson.M{"_id": id}).Decode(result)
}

// 根据任意字段查询文档
func (r *dbImpl) FindByField(ctx context.Context, field string, value interface{}) (bson.M, error) {
	var result bson.M
	err := r.coll(ctx).FindOne(ctx, bson.M{field: value}).Decode(&result)
	return result, err
}
