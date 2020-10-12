package session

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strings"
	"tborm/dialect"
	"tborm/log"
	"tborm/schema"
)

type Session struct {
	// 数据库引擎
	db *sql.DB
	dialect dialect.Dialect
	// SQL语句
	sql strings.Builder
	refTable *schema.Schema
	// SQL动态参数
	sqlValues []interface{}
}

func New(db *sql.DB, dialect dialect.Dialect) *Session {
	return &Session{db: db, dialect: dialect}
}

func (s *Session) Clear() {
	s.sql.Reset()
	s.sqlValues = nil
}

func (s *Session) DB() *sql.DB {
	return s.db
}

func (s *Session) Raw(sql string, values ...interface{}) *Session {
	s.sql.WriteString(sql)
	s.sql.WriteString(" ")
	s.sqlValues = append(s.sqlValues, values...)
	return s
}

func (s *Session) Exec() (result sql.Result, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlValues)
	if result, err = s.DB().Exec(s.sql.String(), s.sqlValues...); err != nil {
		log.Error(err)
	}
	return
}

func (s *Session) QueryRow() *sql.Row {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlValues)
	return s.DB().QueryRow(s.sql.String(), s.sqlValues...)
}

func (s *Session) Query() (rows *sql.Rows, err error) {
	defer s.Clear()
	log.Info(s.sql.String(), s.sqlValues)
	if rows, err = s.DB().Query(s.sql.String(), s.sqlValues...); err != nil {
		log.Error(err)
	}
	return
}
