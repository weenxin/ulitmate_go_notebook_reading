package section1_basic

import "testing"

func TestBasic(t *testing.T) {
	t.Run("test integer slice", func(t *testing.T) {
		Print[int]([]int{1, 2, 3, 4, 5, 6, 7})
		t.Logf("int function address : %p", Print[int])
	})
	t.Run("test string slice", func(t *testing.T) {
		Print([]string{"one", "two", "three", "four", "five", "six", "seven"})
		t.Logf("string function address : %p", Print[string])
	})
	t.Run("test float64 slice", func(t *testing.T) {
		Print[float64]([]float64{1, 2, 3, 4, 5, 6, 7})
		t.Logf("float64 function address : %p", Print[float64])
	})

	t.Run("function address", func(t *testing.T) {
		f1 := Print[int]
		f2 := Print[float64]
		f3 := Print[string]
		//f4 := Print //will panic cannot use generic function Print without instantiation

		t.Logf("int function address : %p, float64 function address : %p , string function address : %p\n", f1, f2, f3)
	})
}
