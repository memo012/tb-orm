package session

import "testing"

var (
	u1    = &User{Name: "Tom", Age: 18}
	u2    = &User{Name: "Sam", Age: 25}
	user3 = &User{"Jack", 25}
)

func testRecordInit(t *testing.T) *Session {
	t.Helper()
	s := NewSession()
	//_, err := s.Insert(u1, u2)
	//if err != nil {
	//	t.Fatal(err)
	//}
	return s
}

func TestSession_Insert(t *testing.T) {
	s := testRecordInit(t)
	affected, err := s.Insert(user3)
	if err != nil || affected != 1 {
		t.Fatal("failed to create record")
	}
}

func TestSession_Find(t *testing.T) {
	s := testRecordInit(t)
	var users []User
	if err := s.Find(&users); err != nil || len(users) != 2 {
		t.Fatal("failed to query all")
	}
}
