package session

import (
	"errors"
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
	// 判断是否为指针类型
	destSlice := reflect.Indirect(reflect.ValueOf(values))
	// 获取指针指向的元素信息
	destType := destSlice.Type().Elem()
	// 将结构体和数据库表进行映射
	table := s.Model(reflect.New(destType).Type().Elem()).RefTable()

	// 组装SQL语句
	s.clause.Set(clause.SELECT, table.Name, table.FieldNames)
	sql, vars := s.clause.Build(clause.SELECT, clause.WHERE, clause.ORDERBY, clause.LIMIT)
	// 进行与数据库交互
	rows, err := s.Raw(sql, vars...).Query()
	if err != nil {
		return err
	}

	for rows.Next() {
		// 获取指针指向的元素信息
		dest := reflect.New(destType).Elem()
		// 结构体字段
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

func (s *Session) Update(kv ...interface{}) (int64, error) {
	m, ok := kv[0].(map[string]interface{})
	if !ok {
		m = make(map[string]interface{})
		for i := 0; i < len(kv); i += 2 {
			m[kv[i].(string)] = kv[i+1]
		}
	}
	s.clause.Set(clause.UPDATE, s.refTable.Name, m)
	sql, vars := s.clause.Build(clause.UPDATE, clause.WHERE)
	res, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (s *Session) Delete() (int64, error) {
	s.clause.Set(clause.DELETE, s.refTable.Name)
	sql, vars := s.clause.Build(clause.DELETE, clause.WHERE)
	res, err := s.Raw(sql, vars...).Exec()
	if err != nil {
		return 0, err
	}
	return res.RowsAffected()
}

func (s *Session) Count() (int64, error) {
	s.clause.Set(clause.COUNT, s.refTable.Name)
	sql, vars := s.clause.Build(clause.COUNT, clause.WHERE)
	row := s.Raw(sql, vars...).QueryRow()
	var tmp int64
	if err := row.Scan(&tmp); err != nil {
		return 0, err
	}
	return tmp, nil
}

func (s *Session) Limit(num int) *Session {
	s.clause.Set(clause.LIMIT, num)
	return s
}

func (s *Session) Where(desc string, args ...interface{}) *Session {
	var vars []interface{}
	s.clause.Set(clause.WHERE, append(append(vars, desc), args...))
	return s
}

func (s *Session) OrderBy(desc string) *Session {
	s.clause.Set(clause.ORDERBY, desc)
	return s
}

func (s *Session) First(value interface{}) error {
	dest := reflect.Indirect(reflect.ValueOf(value))
	destSlice := reflect.New(reflect.SliceOf(dest.Type())).Elem()
	if err := s.Limit(1).Find(destSlice.Addr().Interface()); err != nil {
		return err
	}
	if destSlice.Len() == 0 {
		return errors.New("NOT FIND")
	}
	dest.Set(destSlice.Index(0))
	return nil
}
