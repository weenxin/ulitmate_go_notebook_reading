# 第八章 benchmark测试

在本章，我们将一起学习在golang中如何写基准测试，同时说明当我们在整个测试结果时，遇到的"基准测试谎言"问题。

## 8.1 基准测试基础

基准测试不靠谱！于此同时，在我们没有做基准测试前，这些都是我们的猜测。golang的一个优点就是你不需要对任何事情做猜测。

golang的标准库中，已对基准测试做了很好的支持。要想开始基准测试，你需要：
- 新建一个`*_tets.go`文件
- 新建benchmark函数，函数以`Benchmark`开头，紧跟着测试名称，**测试名称的第一个字母应该是大写的**，函数接收一个`*testing.B`对象。

如下所示的代码，开启了两个基准测试函数，分别用来测试上传和下载的效率。

```go
sample_test.go
package sample
import (
    "testing"
)
func BenchmarkDownload(b *testing.B) {}
func BenchmarkUpload(b *testing.B) {}
```

如下是一个有意思的基准测试例子；

```go

import (
	"fmt"
	"testing"
)

var gs string

func BenchmarkSprint(b *testing.B) {
	var s string
	for i := 0; i < b.N; i++ {
		s = fmt.Sprint("hello")
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
```

如上两个测试用例，想要测试``Sprint`和`Sprintf`哪个函数的运行效率更高。大部分同学可能认为应该是`Sprint`的运行效率，因为没有`format`的流程。但是实际上并不是这样。比如在我的mac上运行结果如下所示：

运行 ` go test -bench . -benchtime 3s ./ch8/basic/...` 命令，结果如下所示：

```
goos: darwin
goarch: amd64
pkg: github.com/weenxin/ulitmate_go_notebook_reading/ch8/basic
cpu: Intel(R) Core(TM) i5-7360U CPU @ 2.30GHz
BenchmarkSprint-4       54031044                65.10 ns/op
BenchmarkSprintf-4      76011352                49.57 ns/op
PASS
ok      github.com/weenxin/ulitmate_go_notebook_reading/ch8/basic       9.962s

```
可以看到`Sprintf`效率反而更高。这与我们的预期是不一样的。



如上所示的代码，我们指定了只能跑3秒钟，但是实际上我们的代码是按照次数计量的，那么他们之间是如何换算的呢？

```go
for i := 0; i < b.N; i++ {
		s = fmt.Sprint("hello")
	}
```

运行的逻辑是：
- 先跑一次用例
- 基于次计算是否可以继续跑100次
- 在计算是否可以继续乘100，以此类推

使用如下所示的代码可以验证结论：

```go
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
```

会输出：

```
goos: darwin
goarch: amd64
pkg: github.com/weenxin/ulitmate_go_notebook_reading/ch8/basic
cpu: Intel(R) Core(TM) i5-7360U CPU @ 2.30GHz
BenchmarkSprint
[1 100 10000 1000000 16437782]
[1 100 10000 1000000 16437782 20209514]
BenchmarkSprint-4       20209514                59.80 ns/op
BenchmarkSprintf
BenchmarkSprintf-4      22030516                48.94 ns/op
PASS
ok      github.com/weenxin/ulitmate_go_notebook_reading/ch8/basic       4.036s
```
可以看到benchmark测试过程中，持续在迭代计算次数。

另外一个有关基准测试的点是有关编译器的。编译器会将测试代码编译出一个二进制应用程序。需要注意在for循环中的代码应该与在生产环境中使用的方式是一致的。一个微小的改动都可能影响到编译器的行为，进而基准测试与生产环境不一致的问题。

```go
func BenchmarkSprintf(b *testing.B) {
	var s string
	for i := 0; i < b.N; i++ {
		s = fmt.Sprintf("hello") // 应该捕获返回值，如果不捕获可能与正常使用有差异，导致编译器将这里优化
	}
	gs = s
}
```
如上所示的代码，捕获`fmt.Sprintf`的返回值是十分必要的,如果因为涉及到数据的拷贝。换另外一个用例：

运行 `go test -v -bench BenchmarkArray  ./ch8/basic/...` 指令。

```go

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


func init() {
	peoples = make([]People, 100)
	for index := range peoples {
		for j := range peoples[index].id {
			peoples[index].id[j] = byte(j)
		}
	}
}
```

产生结果如下所示：

```
goos: darwin
goarch: amd64
pkg: github.com/weenxin/ulitmate_go_notebook_reading/ch8/basic
cpu: Intel(R) Core(TM) i5-7360U CPU @ 2.30GHz
BenchmarkArrayValueLoopNoReceiveValue
BenchmarkArrayValueLoopNoReceiveValue   32290458                35.70 ns/op
BenchmarkArrayPointerLoopNoReceiveData
BenchmarkArrayPointerLoopNoReceiveData  18846804                62.29 ns/op
```

如上所示的代码，按照我们的理解，指针遍历要比值遍历更快（因为不需要copy People.id数组），但是实际上的测试却不是这样。这是因为编译器做了很多优化，换成如下用例：

```go
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
```
结果如下所示：
```
goos: darwin
goarch: amd64
pkg: github.com/weenxin/ulitmate_go_notebook_reading/ch8/basic
cpu: Intel(R) Core(TM) i5-7360U CPU @ 2.30GHz
BenchmarkArrayValueLoopReceiveValue
BenchmarkArrayValueLoopReceiveValue        10000            101690 ns/op
BenchmarkArrayPointerLoopReceiveData
BenchmarkArrayPointerLoopReceiveData       23486             52781 ns/op
```

如果我们开启 `-benchmem` 这个flag，可以看到如下结果如下所示：

运行 `go test -v -benchtime 1s  -bench  BenchmarkArray -cpu 1 -benchmem  ./ch8/basic/...` 命令，会得到：

```go
goos: darwin
goarch: amd64
pkg: github.com/weenxin/ulitmate_go_notebook_reading/ch8/basic
cpu: Intel(R) Core(TM) i5-7360U CPU @ 2.30GHz
BenchmarkArrayValueLoopNoReceiveValue
BenchmarkArrayValueLoopNoReceiveValue           32451716                35.72 ns/op            0 B/op          0 allocs/op
BenchmarkArrayPointerLoopNoReceiveData
BenchmarkArrayPointerLoopNoReceiveData          18987433                62.22 ns/op            0 B/op          0 allocs/op
BenchmarkArrayValueLoopReceiveValue
BenchmarkArrayValueLoopReceiveValue                10000            104908 ns/op               0 B/op          0 allocs/op
BenchmarkArrayPointerLoopReceiveData
BenchmarkArrayPointerLoopReceiveData               22872             52563 ns/op               0 B/op          0 allocs/op
BenchmarkArrayPointerLoopOnHeapReceiveValue
BenchmarkArrayPointerLoopOnHeapReceiveValue         4010            288683 ns/op         2048000 B/op        100 allocs/op
PASS
ok      github.com/weenxin/ulitmate_go_notebook_reading/ch8/basic       7.763s
```

如上所示可以看到，对于内存分配统计应该只有堆的统计，栈上内存分配是不统计的，如下代码可以让内存分配在堆中。同时我们也注意到，同样是两次分配内存和赋值，但是在栈上的效率要高的多。

```go
func BenchmarkArrayPointerLoopOnHeapReceiveValue(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for index := range peoples { //pointer loop , will only copy data once ,but will malloc data in heap
			data := new(People)
			*data = peoples[index]
			gpHeap = data
		}
	}
}
```

## 8.2 子基准测试
有时候我们希望将相同类型的基准测试统一管理起来，也许你希望使用TableTest的方式运行你的测试用例，子基准测试可以满足你的需求。

```go
package basic
import (
    "fmt"
    "testing"
)
var gs string
func BenchmarkSprint(b *testing.B) {
    b.Run("none", benchSprint)
    b.Run("format", benchSprintf)
}
```

如上所示的代码，可以通过如下方式运行和晒讯测试用例。
```bash
$ go test -bench .
$ go test -bench BenchmarkSprint/none
$ go test -bench BenchmarkSprint/format
```

## 8.3 验证基准测试

本章的开始，已经阐述基准测试"不靠谱"的问题。考虑合并排序。

```go

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
```

单线程跑：

```go
func Single(n []int) []int {
	if len(n) <= 1 {
		return n
	}
	mid := len(n) / 2
	return merge(Single(n[:mid]), Single(n[mid:]))
}
```

无限并行

```go
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

	return merge(Unlimited(l), Unlimited(r))
}
```

指定cpu数量

```go

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
```
运行功能测试  ` go test -v  ./ch8/basic/...`

```
=== RUN   TestSingle
--- PASS: TestSingle (0.00s)
=== RUN   TestUnlimited
--- PASS: TestUnlimited (0.00s)
=== RUN   TestNumCpu
--- PASS: TestNumCpu (0.00s)
PASS
ok      github.com/weenxin/ulitmate_go_notebook_reading/ch8/basic       0.127s
```

指定基准测试

```go

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

```


运行 `go test -benchtime 5s  -bench BenchmarkMerge  ./ch8/basic/...`

```
goos: darwin
goarch: amd64
pkg: github.com/weenxin/ulitmate_go_notebook_reading/ch8/basic
cpu: Intel(R) Core(TM) i5-7360U CPU @ 2.30GHz
BenchmarkMergeSingle-4             45496            129516 ns/op
BenchmarkMergeUnlimited-4             36         170992436 ns/op
BenchmarkMergeNumCPU-4             44782            130402 ns/op
PASS
ok      github.com/weenxin/ulitmate_go_notebook_reading/ch8/basic       21.329s

```

单独运行： `go test -benchtime 5s  -bench BenchmarkMergeN  ./ch8/basic/...` 比之前快了一点。可能数组大小有关系，或者电脑性能不行。

```go
goos: darwin
goarch: amd64
pkg: github.com/weenxin/ulitmate_go_notebook_reading/ch8/basic
cpu: Intel(R) Core(TM) i5-7360U CPU @ 2.30GHz
BenchmarkMergeNumCPU-4             40107            160717 ns/op
PASS
ok      github.com/weenxin/ulitmate_go_notebook_reading/ch8/basic       8.664s
```

理论上由于`BenchmarkMergeUnlimited`产生了大量的`go-routeine`然后会做很多的垃圾回收工作，进而拖慢后续的测试用例效果，也会导致性能下降。








