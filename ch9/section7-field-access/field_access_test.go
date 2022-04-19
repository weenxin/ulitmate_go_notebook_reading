package section7_field_access

import "testing"

func TestInsertData(t *testing.T) {

	t.Run("insert user", func(t *testing.T) {
		u := &User{
			Id:    0,
			Name:  "weenxin",
			Email: "weenxin@123.com",
		}
		u , _ = InsertData(u)
		t.Logf("after insert data id : %d",u.Id)
	})

	t.Run("insert customer", func(t *testing.T) {
		c := &Customer{
			Id:    0,
			Name:  "weenxin",
			Email: "weenxin@123.com",
		}
		c , _ = InsertData(c)
		t.Logf("after insert data id : %d",c.Id)
	})

}