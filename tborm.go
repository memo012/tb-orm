package tborm

import (
	"TbORM/log"
	"TbORM/session"
	"database/sql"
)

type Engine struct {
	db *sql.DB
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

	// step3: 连接成功
	e = &Engine{db: db}
	log.Info("Connect database success")
	return
}

func (e *Engine) Close()  {
	if err := e.db.Close(); err != nil {
		log.Error("Failed to close database")
	}
	log.Info("Close database success")
}

func (e *Engine) NewSession() *session.Session {
	return session.New(e.db)
}