package basic

import (
	"math"
	"runtime"
	"sync"
)

func merge(l, r []int) []int {
	var result []int
	for {
		switch {
		case len(l) == 0:
			result = append(result, r...)
			return result
		case len(r) == 0:
			result = append(result, l...)
			return result
		case l[0] < r[0]:
			result = append(result, l[0])
			l = l[1:]
		default:
			result = append(result, r[0])
			r = r[1:]
		}
	}
}

func Single(n []int) []int {
	if len(n) <= 1 {
		return n
	}
	mid := len(n) / 2
	return merge(Single(n[:mid]), Single(n[mid:]))
}

func Unlimited(n []int) []int {
	if len(n) <= 1 {
		return n
	}
	mid := len(n) / 2
	var l []int
	var r []int

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		l = Unlimited(n[:mid])
		wg.Done()
	}()

	go func() {
		r = Unlimited(n[mid:])
		wg.Done()
	}()
	wg.Wait()

	return merge(l, r)
}

func NumCpu(n []int) []int {
	if len(n) <= 1 {
		return n
	}

	maxLevel := int(math.Log2(float64(runtime.GOMAXPROCS(0))))

	return numCpu(n, maxLevel)

}

func numCpu(n []int, maxLevel int) []int {
	if len(n) <= 1 {
		return n
	}
	mid := len(n) / 2

	if maxLevel > 0 {
		var l, r []int
		var wg sync.WaitGroup
		{
		}
		wg.Add(2)
		go func() {
			l = numCpu(n[:mid], maxLevel-1)
			wg.Done()
		}()
		go func() {
			r = numCpu(n[mid:], maxLevel-1)
			wg.Done()
		}()
		wg.Wait()
		return merge(l, r)
	}

	return append(merge(numCpu(n[:mid], maxLevel-1), numCpu(n[mid:], maxLevel-1)))
}
