package session

import (
	"fmt"
	"reflect"
	"strings"
	"tborm/log"
	"tborm/schema"
)

//  将解析结果赋值给refTable
func (s *Session) Model(value interface{}) *Session {
	if s.refTable == nil ||
		reflect.TypeOf(value) != reflect.TypeOf(s.refTable.Model) {
		s.refTable = schema.Parse(value, s.dialect)
	}
	return s
}

// 将结构体和数据库字段解析的结果保存在成员变量 refTable
func (s *Session) RefTable() *schema.Schema {
	if s.refTable == nil {
		log.Error("Model is not set")
	}
	return s.refTable
}

func (s *Session) CreateTable() error {
	table := s.refTable
	var columns []string
	for _, field := range table.Fields {
		columns = append(columns, fmt.Sprintf("%s %s %s", field.Name, field.Type, field.Tag))
	}
	desc := strings.Join(columns, ",")
	_, err := s.Raw("CREATE TABLE %s (%s);", table.Name, desc).Exec()
	return err
}

func (s *Session) DropTable() error {
	_, err := s.Raw(fmt.Sprintf("DROP TABLE IF EXISTS %s", s.RefTable().Name)).Exec()
	return err
}

// todo 查看表是否存在
//func (s *Session) HasTable() bool {
//	sql, values := s.dialect.TableExitSQL(s.RefTable().Name)
//	row := s.Raw(sql, values...).QueryRow()
//	var tmp string
//	_ = row.Scan(&tmp)
//	return tmp == s.RefTable().Name
//}
