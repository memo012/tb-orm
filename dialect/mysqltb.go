package dialect

import (
	"fmt"
	"reflect"
	"time"
)

type mysqltb struct{}

var _ Dialect = (*mysqltb)(nil)

func init() {
	RegisterDialect("mysql", &mysqltb{})
}

// go数据类型与MySQL数据库类型转换
func (m *mysqltb) DataTypeOf(typ reflect.Value) string {
	switch typ.Kind() {
	case reflect.Bool:
		return "bool"
	case reflect.Int, reflect.Int32, reflect.Uint, reflect.Uint32:
		return "int(11)"
	case reflect.Int8, reflect.Uint8:
		return "TINYINT"
	case reflect.Int16, reflect.Uint16:
		return "SMALLINT"
	case reflect.Int64, reflect.Uint64:
		return "BIGINT"
	case reflect.Float32, reflect.Float64:
		return "DOUBLE"
	case reflect.String:
		return "varchar(255)"
	case reflect.Slice, reflect.Array:
		return "blob"
	case reflect.Struct:
		if _, ok := typ.Interface().(time.Time); ok {
			return "datetime"
		}
	}
	panic(fmt.Sprintf("invalid sql type %s (%s)", typ.Type().Name(), typ.Kind()))
}

func (m *mysqltb) TableExitSQL(tableName string) (string, []interface{}) {
	args := []interface{}{tableName}
	return "SELECT name FROM sqlite_master WHERE type='table' and name = ?", args
}
