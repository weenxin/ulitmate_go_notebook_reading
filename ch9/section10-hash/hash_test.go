package section10_hash

import (
	"hash/fnv"
	"testing"
)

func TestHash(t *testing.T) {

	function := func(data string, buckets int) int {
		h := fnv.New32()
		h.Write([]byte(data))
		return int(h.Sum32())%buckets
	}
	function2 := func(data int, buckets int) int {
		return data%buckets
	}
	t1 := New[string,int](100,function)
	t2 := New[int,string](100,function2)

	values := map[string]int{"one":1,"two":2,"three":3,"four":4}
	for key,value := range values {
		t1.Insert(key,value)
		t2.Insert(value,key)
	}
	for key,value := range values {
		intValue , exists := t1.Get(key)
		if !exists {
			t.Log("not exists")
		}else{
			t.Logf("getting data: %s : %d",key, intValue)
		}
		stringValue , exists := t2.Get(value)
		if !exists {
			t.Log("not exists")
		}else{
			t.Logf("getting data: %d : %s",value, stringValue)
		}
	}
}
