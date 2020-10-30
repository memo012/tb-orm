package dialect

import "reflect"

// 类说明  该类旨于tb-orm框架兼容多个数据库
// 不同的数据库支持的数据类型有所差异
// 将差异的部分提取出来 每一种数据库分别实现 实现最大程度的复用和耦合

var dialectMap = map[string]Dialect{}

type Dialect interface {
	// 将Go语言的类型转换为该数据库的数据类型
	DataTypeOf(typ reflect.Value) string
	// 返回某个表是否存在的SQL语句 参数是表名
	TableExitSQL(tableName string) (string, []interface{})
}

// 注册dialect实例
func RegisterDialect(name string, dialect Dialect)  {
	dialectMap[name] = dialect
}

// 获取dialect实例
func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectMap[name]
	return
}
