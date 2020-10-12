package session

import "testing"

type User struct {
	Name string `tborm:"PRIMARY KEY"`
	Age  int
}

func TestSession_CreateTable(t *testing.T) {
	session := NewSession().Model(&User{})
	_ = session.CreateTable()
}
