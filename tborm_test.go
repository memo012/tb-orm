package tborm

import "testing"

func TestEngine_NewSession(t *testing.T) {
	t.Helper()
	engine,err := NewEngine("mysql", "root:root@tcp(127.0.0.1:3306)/vblog?charset=utf8mb4")
	if err != nil {
		t.Fatal("failed to connect", err)
	}
	engine.Close()
}
