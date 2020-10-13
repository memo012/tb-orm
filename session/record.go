package session

import (
	"reflect"
	"tborm/clause"
)

func (s *Session) Insert(values ...interface{}) (int64, error) {
	recordValues := make([]interface{}, 0)
	for _, value := range values {
		table := s.Model(value).RefTable()
		s.clause.Set(clause.INSERT, table.Name, table.FieldNames)
		recordValues = append(recordValues, table.RecordValues(value))
	}

	s.clause.Set(clause.VALUES, recordValues...)
	sql, vars := s.clause.Build(clause.INSERT, clause.VALUES)
	result, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Find 的代码实现比较复杂，主要分为以下几步：
//1) destSlice.Type().Elem() 获取切片的单个元素的类型 destType
//	 使用 reflect.New() 方法创建一个 destType 的实例，作为 Model() 的入参 映射出表结构 RefTable()
//2）根据表结构，使用 clause 构造出 SELECT 语句 查询到所有符合条件的记录 rows
//3）遍历每一行记录，利用反射创建 destType 的实例 dest 将 dest 的所有字段平铺开 构造切片 values
//4）调用 rows.Scan() 将该行记录每一列的值依次赋值给 values 中的每一个字段
//5）将 dest 添加到切片 destSlice 中 循环直到所有的记录都添加到切片 destSlice 中
func (s *Session) Find(values interface{}) error {
	destSlice := reflect.Indirect(reflect.ValueOf(values))
	destType := destSlice.Type().Elem()
	table := s.Model(reflect.New(destType).Type().Elem()).RefTable()

	s.clause.Set(clause.SELECT, table.Name, table.FieldNames)
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	rows, err := s.Raw(sql, vars...).Query()
	if err != nil {
		return err
	}

	for rows.Next() {
		dest := reflect.New(destType).Elem()
		var values []interface{}
		for _, name := range table.FieldNames {
			values = append(values, dest.FieldByName(name).Addr().Interface())
		}
		if err := rows.Scan(values...); err != nil {
			return err
		}
		// 赋值
		destSlice.Set(reflect.Append(destSlice, dest))
	}
	return rows.Close()
}
