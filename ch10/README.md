# 第十章 Profiling(诊断/分析)

在本章我们将学习如何对代码诊断和分析。相对来说比较枯燥，枯燥就会容易让人想起放弃，所以还是希望可以不那么枯燥的把这个聊一聊。以下是我的一些心里路程。

在尝试学习本章的时候，我曾想过放弃当前在做的事情。比如去做一些看起来更加`fancy（花哨）`的事情，比如
- 看看`golang http包`源码；
- 整理下`kubernetes`资料；
- 学习下`kubernetes client-go`；
- 学习`kubernetes kubelet`源码；
- 学习下Prometheus源码；
- 研究研究微服务；
- 写一个cmdb等等。

放弃的念头是十分容易出来的，但是真正做这个决定是比较难得。但是如果放弃去做其他的事情，那么：
- 我会怀疑，我是否是一个靠谱programer或者合格的"同学"吗？我是一个合格的engineer吗？为什么我没有坚持下去？
- 这件事没有坚持下来，我能坚持下来，做好另外一件事吗？
- 做这个事情，也许也有一些收益，对我未来的工作和学习有参考和意义。

我还是决定继续把手头上的事情做好，虽然很枯燥，能学到多少技术上的东西也不清楚，或者后面的工作能不能用，也不确定，但是还是把基础打好。这让我想起我的一个好朋友。他在做一件对于人生来说很重要的事情的时候，付出了100%努力，虽然结果并没有达到他预期的那么好。当我们回忆那段经历的时候，而他当时的想法并不是期望"逆袭"人生，而只是向自己证明自己是可以的，仅此而已，而他确实做到了。


摘自文章的一句话。

`Those who can make you believe absurdities can make you commit atrocities" - Voltaire`

`动辄让人忽悠信邪之辈，也难免受人教唆/鼓动为非作歹。—— 伏尔泰`



## 10.1 介绍

我们可以使用golang的工具集对程序诊断分析，诊断分析既是一次侦探，又是一次旅行。需要我们对运行的程序有一定的了解。诊断数据只是一堆干燥的数据，需要参与诊断的人赋予给这些数据语义。


### 10.1.1 诊断基础

**分析器是怎么工作的？**

一个分析器会启动运行我们的程序，并且会定时中断程序的运行。这是通过向调试程序发送`SIGPROF`信号来实现的，接收到这个信号会让待调试的程序暂停运行并且运行分析器。这个分析器会抓取每个线程的PC（program counter），然后再继续只从programer。

**分析器做什么，不做什么？**

在我开始做性能分析前，我需要一个稳定的环境，以得到一个可重复的结果。

- 运行Profiling的机器，必须是空闲的；不要在一个共享的机器做profiling；
- 小心省电模式和cpu温度？（power saving and thermal scaling）；
- 不要使用虚拟机或者共享云主机；

如果你能负担的起这个成本就这么搞测试环境，说关闭省电模式，不要更新软件。

### 10.1.2 性能分析类型

**CPU性能分析**

CPU性能调优是最常用的性能调优分析。当开启CPU性能调优时，golang的`runtime`会（没10ms）打断goroutine并基于goroutine的调用堆栈。但我们把心更难分析的数据保存起来，我们就可以对代码的执行`hot path`进行分析了。函数出现的次数越多，那么这个函数吃掉cpu的时间越多。

**内存分析**

开启内存分析，当程序在堆中申请内存时，会记录内存堆栈。内存分析和cpu调优一样都是基于采样的。默认情况下，每512Kb采样一次，采样率可以调整。对于栈内存分配可以认为是无损的，不会被追溯。由于内存分析是基于采样的，而且他追踪未被使用的内存分配，因此基于此评估程序总内存是不靠谱的。

**阻塞分析**

与CPU分析不同的是阻塞分析主要记录，某个goroutine浪费在等待状态的时间。这对于分析并发性能是十分重要的，阻塞分析对于大量goroutine可以并行执行时，但是却被阻塞了，这个场景的分析是十分重要的。

阻塞主要发生在以下环节：
- 从非缓存chan，发送和接收数据；
- 向一个满状态的chan发送数据，从空状态的chan接收数据；
- mutex锁；

阻塞分析应该是最后的工具，在此之前我们应该已经确定CPU和Mem都不是性能瓶颈。

**一次只进行一种类型性能调试**

不要一次进行多种类型性能分析的，因为每种类型的性能分析都是有损的。尤其当我们加大内存分配的采样频率时。

### 10.1.3 合适停止性能分析

当我发现有很长的时间浪费在`runtime.mallocgc`这个函数时，我们的程序可能在频繁的做一些小对象的申请和释放。我们的profile会告诉我们内存分配发生在哪里。

如果大部分时间浪费在了channel的接收发送以及mutex等的`sync`操作上，某些资源可能是瓶颈，需要重构代码，以减少共享资源的访问，一般可以通过sharding，buffer，batch，copy-on-write技术来解决。

如果你的程序浪费了很多实践中在`syscall.Read/Write`，应用程序可能在频繁的调用大量的读写操作。`bufio`包裹`os.File`或者`net.Conn`可以解决问题。

如果大量时间话费在`GC`上，那么要么是你的应用程序在分配大量小对象，要么你的堆内存太小。

- 大内存分配影响到GC的pacing时间，但是大量小对象会影响到标记时间；
- 将大量小对象合并为大对象，可以提高内存分配效率，同时降低gc时间；
- 没有包含指针的对象是不会被垃圾回收程序扫描的，如果减少指针feild会大大降低垃圾回收扫描时间。

### 10.1.4 性能规则

1. 永远不要猜测性能；
2. 测量必须是相关的；
3. 在决定某个环节是性能之前，先做Profiling；
4. 测试之后发现我是正确的。


### 10.1.5 Go 和 OS 工具

`time`工具可以给你有关程序运行速度的基本概念。

`perf` 如果的应用程序运行在linux，可以使用perf工具。

## 10.2 样例代码

```go
var data = []struct {
	input  []byte
	output []byte
}{
	{[]byte("abc"), []byte("abc")},
	{[]byte("elvis"), []byte("Elvis")},
	{[]byte("aElvis"), []byte("aElvis")},
	{[]byte("abcelvis"), []byte("abcElvis")},
	{[]byte("eelvis"), []byte("eElvis")},
	{[]byte("aelvis"), []byte("aElvis")},
	{[]byte("aabeeeelvis"), []byte("aabeeeElvis")},
	{[]byte("e l v i s"), []byte("e l v i s")},
	{[]byte("aa bb e l v i saa"), []byte("aa bb e l v i saa")},
	{[]byte(" elvi s"), []byte(" elvi s")},
	{[]byte("elvielvis"), []byte("elviElvis")},
	{[]byte("elvielvielviselvi1"), []byte("elvielviElviselvi1")},
	{[]byte("elvielviselvis"), []byte("elviElvisElvis")},
}
```

组合一个输入流

```go
func assembleInputStream() []byte {
	var in []byte
	for _, d := range data {
		in = append(in, d.input...)
	}
	return in
}
```

算法1

```go
func algOne(data []byte, find []byte, repl []byte, output *bytes.Buffer) {
	input := bytes.NewBuffer(data)
	size := len(find)
	buf := make([]byte, len(find)
	end := size - 1
	if n, err := io.ReadFull(input, buf[:end]); err != nil { //把input当作一个io.Reader接口来用，此时会导致input逃逸
		output.Write(buf[:n])
		return
	}
	for {
		if _, err := io.ReadFull(input, buf[end:]); err != nil {//把input当作一个io.Reader接口来用，此时会导致input逃逸
			output.Write(buf[:end])
			return
		}

		if bytes.Equal(buf, find) { //3 转换成为一个string，会有一次cop，性能降低
			output.Write(repl)
			if _, err := io.ReadFull(input, buf[end:]); err != nil { //当作io.Reader来用，input会逃逸
				output.Write(buf[:end])
				return
			}

			continue
		}
		output.WriteByte(buf[0])
		copy(buf, buf[1:])
	}
}
```

算法2

```go
func algTwo(data []byte, find []byte, repl []byte, output *bytes.Buffer) {
	input := bytes.NewReader(data)
	size := len(find)
	idx := 0
	for {
		b, err := input.ReadByte()
		if err != nil {
			break
		}
		if b == find[idx] {
			idx++
			if idx == size {
				output.Write(repl)
				idx = 0
			}
			continue
		}
		if idx != 0 {
			output.Write(find[:idx])
			input.UnreadByte()
			idx = 0
			continue
		}
		output.WriteByte(b)
		idx = 0
	}
}
```


## 10.3 性能测试

```go
var output bytes.Buffer
var in = assembleInputStream()
var find = []byte("elvis")
var repl = []byte("Elvis")
func BenchmarkAlgorithmOne(b *testing.B) {
	for i := 0; i < b.N; i++ {
		output.Reset()
		algOne(in, find, repl, &output)
	}
}
func BenchmarkAlgorithmTwo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		output.Reset()
		algTwo(in, find, repl, &output)
	}
}
```

运行 `go test -bench BenchmarkAlgorithm  -benchtime 3s -benchmem ./ch10/...`

结果如下所示：

`
goos: darwin
goarch: amd64
pkg: github.com/weenxin/ulitmate_go_notebook_reading/ch10
cpu: Intel(R) Core(TM) i5-7360U CPU @ 2.30GHz
BenchmarkAlgorithmOne-4          1784953              2292 ns/op              53 B/op          2 allocs/op
BenchmarkAlgorithmTwo-4          6776229               474.9 ns/op             0 B/op          0 allocs/op
PASS
ok      github.com/weenxin/ulitmate_go_notebook_reading/ch10    13.598s
`

## 10.4 内存分析

执行 ：  `go test -bench . -benchtime 3s -benchmem -memprofile p.out`

会产生两个文件：`p.out` 和 `ch10.test`

`
-rw-r--r--   1 weenxin  staff   188B Apr 19 17:55 README.md
drwxr-xr-x   5 weenxin  staff   160B Apr 21 16:44 ch10/
-rwxr-xr-x   1 weenxin  staff   2.9M Apr 21 16:44 ch10.test*
drwxr-xr-x   4 weenxin  staff   128B Apr 20 15:12 ch11/
drwxr-xr-x   5 weenxin  staff   160B Apr 11 16:00 ch6/
drwxr-xr-x  13 weenxin  staff   416B Apr 13 21:06 ch7/
drwxr-xr-x   4 weenxin  staff   128B Apr 16 17:25 ch8/
drwxr-xr-x  13 weenxin  staff   416B Apr 19 18:07 ch9/
-rw-r--r--   1 weenxin  staff   1.6K Apr 20 15:24 go.mod
-rw-r--r--   1 weenxin  staff    15K Apr 20 15:24 go.sum
-rw-r--r--   1 weenxin  staff   1.1K Apr 21 16:45 p.out //产
`

- `p.out` : 是产生的profile文件；
- `ch10.test` : 运行benchmark的二进制文件。

这些是用来做profiling时需要的。

运行： `go tool pprof ch10.test p.out`

进入命令行，可以使用list命令查看具体代码的内存分配情况。

![1-listOne](/ch10/images/1-memory.png)

如上图所示：

- flat： 本行分配的内存
- cum： 函数调用分配的内存

一共有两次堆内存分配
- `buf := make([]byte, size)` 由于size是一个不定长度，所以在堆内存分配
- `input := bytes.NewBuffer(data)` 返回bytes.Buffer指针对象
   - `if n, err := io.ReadFull(input, buf[:end]); err != nil` 此时将input转换为一个io.Reader interface，input逃逸到堆


改进代码：

```go
func algOne(data []byte, find []byte, repl []byte, output *bytes.Buffer) {
	input := bytes.NewBuffer(data)
	size := len(find)
	buf := make([]byte, 5)
	end := size - 1
	if n, err := input.Read(buf[:end]); err != nil { //直接调用对象的方法，没有隐式转换，所以不会逃逸
	//if n, err := io.ReadFull(input, buf[:end]); err != nil { //把input当作一个io.Reader接口来用，此时会导致input逃逸
		output.Write(buf[:n])
		return
	}
	for {
		if _, err := input.Read(buf[end:]); err != nil { //直接调用对象的方法，没有隐式转换，所以不会逃逸
		//if _, err := io.ReadFull(input, buf[end:]); err != nil {//把input当作一个io.Reader接口来用，此时会导致input逃逸
			output.Write(buf[:end])
			return
		}

		//if equal(buf, find) { //3 占用时间比较多，因此可以优化下
		if bytes.Equal(buf, find) { //3 转换成为一个string，会有一次cop，性能降低
		//	output.Write(repl)
		//	if _, err := io.ReadFull(input, buf[end:]); err != nil { //当作io.Reader来用，input会逃逸
		//		output.Write(buf[:end])
		//		return
		//	}

			data, err := input.ReadByte() //　调用方法，没有隐式转换
			buf[end:][0] = data
			if  err != nil {
				output.Write(buf[:end])
				return
			}

			continue
		}
		output.WriteByte(buf[0])
		copy(buf, buf[1:])
	}
}
```

运行 `go test -bench . -benchtime 3s -benchmem -memprofile p.out ./ch10/...`


```
goos: darwin
goarch: amd64
pkg: github.com/weenxin/ulitmate_go_notebook_reading/ch10
cpu: Intel(R) Core(TM) i5-7360U CPU @ 2.30GHz
BenchmarkAlgorithmOne-4          2723008              1328 ns/op               0 B/op          0 allocs/op
BenchmarkAlgorithmTwo-4          8369998               430.0 ns/op             0 B/op          0 allocs/op
PASS
ok      github.com/weenxin/ulitmate_go_notebook_reading/ch10    9.922s

```

已经没有堆内存分配了

## 10.5 CPU分析

运行： `go test -bench . -benchtime 3s -cpuprofile p.out ./ch10/...`

```
goos: darwin
goarch: amd64
pkg: github.com/weenxin/ulitmate_go_notebook_reading/ch10
cpu: Intel(R) Core(TM) i5-7360U CPU @ 2.30GHz
BenchmarkAlgorithmOne-4          2689377              1305 ns/op
BenchmarkAlgorithmTwo-4          8294487               483.1 ns/op
PASS
ok      github.com/weenxin/ulitmate_go_notebook_reading/ch10    9.559s
```

运行 ` go tool pprof p.out`

返回如下结果

![2-cpu](/ch10/images/2-cpu.png)

可以看到` if bytes.Equal(buf, find)`，这一行调用`Equal`函数用了非常多的时间。

```go
// Equal reports whether a and b
// are the same length and contain the same bytes.
// A nil argument is equivalent to an empty slice.
func Equal(a, b []byte) bool {
	// Neither cmd/compile nor gccgo allocates for these string conversions.
	return string(a) == string(b)
}
```

可以看到`Equal`做了一次转换。每一次转换都做了一次内存拷贝。具体见： https://github.com/golang/go/issues/25484

```go
func equal(a , b []byte) bool {
	return  (*(*string)(unsafe.Pointer(&a))) ==  (*(*string)(unsafe.Pointer(&b)))
}
```

改下代码：
```
if equal(buf, find) { //3 占用时间比较多，因此可以优化下
//if bytes.Equal(buf, find) { //3 转换成为一个string，会有一次cop，性能降低
```

运行：  `go test -bench . -benchtime 3s -cpuprofile p.out ./ch10/...`

结果如下所示：

```
goos: darwin
goarch: amd64
pkg: github.com/weenxin/ulitmate_go_notebook_reading/ch10
cpu: Intel(R) Core(TM) i5-7360U CPU @ 2.30GHz
BenchmarkAlgorithmOne-4          2454951              1600 ns/op
BenchmarkAlgorithmTwo-4          6060313               513.8 ns/op
PASS
ok      github.com/weenxin/ulitmate_go_notebook_reading/ch10    13.544s
```
优化下代码：

```go
func equal(a , b []byte) bool {

	if len(a) != len(b) {
		return false
	}
	for index := range a {
		if a[index] != b[index] {
			return  false
		}
	}
	return true
}
```

得到如下结果

```
goos: darwin
goarch: amd64
pkg: github.com/weenxin/ulitmate_go_notebook_reading/ch10
cpu: Intel(R) Core(TM) i5-7360U CPU @ 2.30GHz
BenchmarkAlgorithmOne-4          3829674               942.7 ns/op
BenchmarkAlgorithmTwo-4          8493554               429.6 ns/op
PASS
ok      github.com/weenxin/ulitmate_go_notebook_reading/ch10    9.152s

```

确实和想象的不一样。






