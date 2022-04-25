package ch10

import "testing"

func TestStringEquals(t *testing.T) {
	data := []byte{'a','b','c','d','e'}

	arg1 := string(data)

	data[1] = 'a'

	arg2 := string(data)

	t.Logf("arg1 =%q,  arg2=%q", arg1,arg2)
}
