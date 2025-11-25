package imongo

import "go.mongodb.org/mongo-driver/mongo"

// ----------------------------
// 表句柄类型（内部安全）
// ----------------------------
type TableHandle struct {
	coll *mongo.Collection // 私有字段，外部无法修改
}

// Coll 返回原生 *mongo.Collection（只读访问）
// 外部可以直接操作 Aggregate / Find / Insert 等方法
func (t *TableHandle) Coll() *mongo.Collection {
	return t.coll
}

// ----------------------------
// 对外唯一可访问的表
// ----------------------------
var tableSubject *TableHandle // 私有全局变量

// TableSubject 返回 subject 表的句柄（只读）
// 外部只能通过这个方法访问
func TableSubject() *TableHandle {
	return tableSubject
}

// ----------------------------
// 初始化表句柄（Register 调用）
// ----------------------------
func initTables() {
	db := getDatabase() // 内部函数，从注册模块获取数据库
	tableSubject = &TableHandle{
		coll: db.Collection("subject"),
	}
}
