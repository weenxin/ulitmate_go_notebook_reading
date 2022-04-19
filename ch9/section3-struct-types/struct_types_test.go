package section3_struct_types

import (
	"testing"
)

type user struct {
	name string
}

func TestAdd(t *testing.T) {

	var users list[user]
	n1 := users.add(user{name: "weenxin"})
	n2 := users.add(user{name: "stone"})
	t.Log(n1.data.name, n2.data.name)

	var pUsers list[*user]
	n3 := pUsers.add(&user{name: "zhansan"})
	n4 := pUsers.add(&user{name: "lisi"})
	t.Log(n3.data.name, n4.data.name)

}
