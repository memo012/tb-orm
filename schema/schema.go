package schema

import (
	"go/ast"
	"reflect"
	"tborm/dialect"
)

// 该类 旨于实现 对象和表的转换

// 表名 -- 结构名
// 字段名和字段类型 -- 成员变量和类型
// 额外的约束条件(例外非空，主键等) -- 成员变量的Tag

type Filed struct {
	// 字段名
	Name string
	// 字段类型
	Type string
	// 约束条件
	Tag string
}

// Schema 主要包含被映射的对象(Model) 表名(Name) 字段(Fields)
// FieldNames 包含所有的字段名(列名)
// fieldMap 记录字段名和 Field 的映射关系 方便之后直接使用 无需遍历 Fields
type Schema struct {
	Model      interface{}
	Name       string
	Fields     []*Filed
	FieldNames []string
	fieldMap   map[string]*Filed
}

// 通过字段名 获取 field对象
func (s *Schema) GetField(name string) *Filed {
	return s.fieldMap[name]
}

func Parse(dest interface{}, d dialect.Dialect) *Schema {
	// 获取指针指向的实例类型对象
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()
	schema := &Schema{
		Model:    dest,
		Name:     modelType.Name(),
		fieldMap: make(map[string]*Filed),
	}

	// 遍历 dest 中的字段 属性
	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		// 忽略匿名字段和私有字段
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Filed{
				Name: p.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}
			// 获取字段上标签
			if v, ok := p.Tag.Lookup("tborm"); ok {
				field.Tag = v
			}
			//
			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, p.Name)
			schema.fieldMap[p.Name] = field
		}
	}
	return schema
}

func (s *Schema) RecordValues(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var fieldValues []interface{}
	for _, field := range s.Fields {
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}
	return fieldValues
}
