package imongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

// TableHandle 表对象（只包含表名）
type TableHandle struct {
	name string
}

// Table 入口方法（外部调用）
func Table(name string) *TableHandle {
	return &TableHandle{name: name}
}

// Collection 内部安全包装，外部获取 *mongo.Collection 的唯一方式
func (t *TableHandle) Collection(ctx context.Context) *mongo.Collection {
	cli := mustClient() // 内部可访问 Client
	cfg := mustLoadConfig(ctx)
	return cli.Database(cfg.Database).Collection(t.name)
}

var (
	TableSubject = Table("subject")
)
