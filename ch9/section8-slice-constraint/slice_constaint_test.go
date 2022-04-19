package section8_slice_constraint

import "testing"


func double(value int) int {
	return value * 2
}

type Numbers []int

func TestEvery(t *testing.T) {
	t.Run("test integer", func(t *testing.T) {
		items := Numbers{1,2,3,4,5,6}
		Every(items,double)
		t.Logf("values : %v, Type : %T", items,items)
	})
	t.Run("test slice function", func(t *testing.T) {
		items := Numbers{1,2,3,4,5,6}
		EverySlice(items,double)
		t.Logf("values : %v, Type: %T", items,items )
	})
}
