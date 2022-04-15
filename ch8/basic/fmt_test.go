package basic

import (
	"fmt"
	"testing"
)

var gs string

var a []int

func BenchmarkSprint(b *testing.B) {
	var s string
	a = append(a, b.N)
	for i := 0; i < b.N; i++ {
		s = fmt.Sprint("hello")
	}
	if len(a) > 4 {
		fmt.Println(a)
	}
	gs = s
}
func BenchmarkSprintf(b *testing.B) {
	var s string
	for i := 0; i < b.N; i++ {
		s = fmt.Sprintf("hello")
	}
	gs = s
}

var gp People
var gpHeap *People

type People struct {
	id [20480]byte
}

var peoples []People

func BenchmarkArrayValueLoopNoReceiveValue(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, item := range peoples { //get data , it will do 2 times copy
			_ = item
		}
		//but compiler will do some optimization, maybe it will not loop the peoples
	}
}
func BenchmarkArrayPointerLoopNoReceiveData(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for index := range peoples { //pointer loop , will only copy data once , gp = peoples[index]
			_ = peoples[index]
		}
		//but compiler will do some optimization, maybe it will only loop the peoples
	}
}

func BenchmarkArrayValueLoopReceiveValue(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, item := range peoples { //get data , it will do 2 times copy
			gp = item
		}
	}
}

func BenchmarkArrayPointerLoopReceiveData(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for index := range peoples { //pointer loop , will only copy data once , gp = peoples[index]
			gp = peoples[index]
		}
	}
}

func BenchmarkArrayPointerLoopOnHeapReceiveValue(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for index := range peoples { //pointer loop , will only copy data once ,but will malloc data in heap
			data := new(People)
			*data = peoples[index]
			gpHeap = data
		}
	}
}

func init() {
	peoples = make([]People, 100)
	for index := range peoples {
		for j := range peoples[index].id {
			peoples[index].id[j] = byte(j)
		}
	}
}
