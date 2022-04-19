package section5_type_as_contraint

import "testing"

func TestAdd(t *testing.T) {
	t.Run("add interge", func(t *testing.T) {
		t.Log(Add(1.2, 2.1))
	})
}

func TestMatch(t * testing.T) {
	t.Run("test person", func(t *testing.T) {
		peoples := []Person{{name:"weenxin"},{name:"stone"}}
		index := match(peoples,Person{name: "stone"})
		t.Logf("find : %v ", index)
	})

	// Other 不在matcher列表中
	//t.Run("test others",func(t *testing.T) {
	//	others := []Other{{name:"weenxin"},{name:"stone"}}
	//	index := match(others,Other{name: "stone"})
	//	t.Logf("find : %v ", index)
	//})

	//Food 没有实现方法
	//t.Run("test foods",func(t *testing.T) {
	//	foods := []Food{{name:"weenxin"},{name:"stone"}}
	//	index := match(foods,Food{name: "stone"})
	//	t.Logf("find : %v ", index)
	//})
}
