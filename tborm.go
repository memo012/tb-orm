package tborm

import (
	"database/sql"
	"tborm/dialect"
	"tborm/log"
	"tborm/session"
)

type Engine struct {
	db      *sql.DB
	dialect dialect.Dialect
}

func NewEngine(driverName, source string) (e *Engine, err error) {
	// step1: 初始化一个sql.DB参数
	// 不会立即建立一个数据库的网络连接 也不会对数据库链接参数进行校验
	db, err := sql.Open(driverName, source)
	if err != nil {
		log.Error(err)
		return
	}

	// step2: 测试数据库是否连接成功
	if err = db.Ping(); err != nil {
		log.Error(err)
		return
	}

	// step3: 确保数据库存在
	dial, ok := dialect.GetDialect(driverName)
	if !ok {
		log.Errorf("dialect %s Not Found", driverName)
		return
	}

	// step4: 连接成功
	e = &Engine{db: db, dialect: dial}
	log.Info("Connect database success")
	return
}

func (e *Engine) Close() {
	if err := e.db.Close(); err != nil {
		log.Error("Failed to close database")
	}
	log.Info("Close database success")
}

func (e *Engine) NewSession() *session.Session {
	return session.New(e.db, e.dialect)
}

type TxFunc func(s *session.Session) (interface{}, error)

func (e *Engine) Transaction(f TxFunc) (result interface{}, err error) {
	s := e.NewSession()
	if err := s.Begin(); err != nil {
		return nil, err
	}
	defer func() {
		if p := recover(); p != nil {
			_ = s.Rollback()
			panic(p)
		} else if err != nil {
			_ = s.Rollback()
		} else {
			err = s.Commit()
		}
	}()
	return f(s)
}
