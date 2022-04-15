package basic

import (
	"reflect"
	"testing"
)

func TestSingle(t *testing.T) {
	tt := []struct {
		input  []int
		expect []int
	}{
		{
			[]int{1, 4, 2, 6, 9, 8},
			[]int{1, 2, 4, 6, 8, 9},
		},
		{
			[]int{4, 1, 2, 6, 9, 8},
			[]int{1, 2, 4, 6, 8, 9},
		},
		{
			[]int{4},
			[]int{4},
		},
		{
			[]int{},
			[]int{},
		},
	}

	for _, test := range tt {
		result := Single(test.input)
		if !reflect.DeepEqual(result, test.expect) {
			t.Fatalf("input : %v expect : %v , got: %v ", test.input, test.expect, result)
		}
	}
}

func TestUnlimited(t *testing.T) {
	tt := []struct {
		input  []int
		expect []int
	}{
		{
			[]int{1, 4, 2, 6, 9, 8},
			[]int{1, 2, 4, 6, 8, 9},
		},
		{
			[]int{4, 1, 2, 6, 9, 8},
			[]int{1, 2, 4, 6, 8, 9},
		},
		{
			[]int{4},
			[]int{4},
		},
		{
			[]int{},
			[]int{},
		},
	}

	for _, test := range tt {
		result := Unlimited(test.input)
		if !reflect.DeepEqual(result, test.expect) {
			t.Fatalf("input : %v expect : %v , got: %v ", test.input, test.expect, result)
		}
	}
}

func TestNumCpu(t *testing.T) {
	tt := []struct {
		input  []int
		expect []int
	}{
		{
			[]int{1, 4, 2, 6, 9, 8},
			[]int{1, 2, 4, 6, 8, 9},
		},
		{
			[]int{4, 1, 2, 6, 9, 8},
			[]int{1, 2, 4, 6, 8, 9},
		},
		{
			[]int{4, 1, 2, 6, 9, 8, 10, 22, 33, 4, 5, 6},
			[]int{1, 2, 4, 4, 5, 6, 6, 8, 9, 10, 22, 33},
		},
		{
			[]int{4},
			[]int{4},
		},
		{
			[]int{},
			[]int{},
		},
	}

	for _, test := range tt {
		result := NumCpu(test.input)
		if !reflect.DeepEqual(result, test.expect) {
			t.Fatalf("input : %v expect : %v , got: %v ", test.input, test.expect, result)
		}
	}
}

var n []int

func init() {
	for i := 0; i < 1_000; i++ {
		n = append(n, 1000-i)
	}
}
func BenchmarkMergeSingle(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Single(n)
	}
}
func BenchmarkMergeUnlimited(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Unlimited(n)
	}
}
func BenchmarkMergeNumCPU(b *testing.B) {
	for i := 0; i < b.N; i++ {
		numCpu(n, 0)
	}
}
