package section2_underlying_types

import "testing"

func TestVectors(t *testing.T) {

	t.Run("integer", func(t *testing.T) {
		vGenericInter := vector[int]{1, 2, 3, 4, 5, 6}
		t.Log(vGenericInter.last())
		//vGenericInter = vector{1, 2, 3, 4, 5, 6} //cannot use generic type vector[T any] without instantiation
		//t.Log(vGenericInter.last())
	})

	t.Run("string", func(t *testing.T) {
		//vGenericInter := vector{"one", "two", "three", "four", "five", "six"}
		//编译报错： cannot use generic type vector[T any] without instantiation

		//vGenericInter := (vector[string])([]string{"one", "two", "three", "four", "five", "six"}) //这样是可以的
		vGenericInter := vector[string]{"one", "two", "three", "four", "five", "six"}
		t.Log(vGenericInter.last())
		//vGenericInter = vector{"one", "two", "three", "four", "five", "six"} //cannot use generic type vector[T any] without instantiation
		//t.Log(vGenericInter.last())
	})
}
