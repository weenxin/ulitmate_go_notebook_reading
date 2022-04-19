package section4_behavior_as_constraint

import "testing"

func TestStringfy(t *testing.T) {
	user := []User{{name: "weenxin"}, {name: "stone"}, {name: "zhangsan"}, {name: "lisi"}}
	values := stringfy(user)
	t.Log(values)
}
