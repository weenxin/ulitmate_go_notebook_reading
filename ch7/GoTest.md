# go test 命令

# `go help test`

`go test` 会自动运行引入的包。并产生如下的summary信息：

```
 ok   archive/tar   0.011s
 FAIL archive/zip   0.022s
 ok   compress/gzip 0.033s
```

'Go test' 会重新编译所有有`*_test.go`文件的包。这些文件可以包含功能测试，benchmark测试，模糊测试。**所有以_和.开头的文件都被忽略，包括_test.go**。

在*_test*开头的包内声明的测试文件，会作为一个单独的包测试，并会以主测试二进制编译，链接与运行。可以在每个测试包中加入`testdata`文件夹保存用于测试的文件。

`go test`前会运行`go vet`测试文件的签名问题。如果在`go vert`过程中发现问题，将停止运行`go test`。

`go test`两种模式：
- `go test` 编译运行当前文件测试文件。侧重模式下不`caching`。
- `go test .` , `go test math` , `go test ./...`

go test 运行时，工作目录为被测试包所在的位置。

当使用第二种`go test`模式时，会缓存所有已经测试成功的结果。再次运行测试用例，并不会真正运行测试，而是直接打印测试结果并且会带有(cached)字样。

cache的策略时：
- 编译出的二进制文件时一致的；
- 不包含"不cache"友好的flag，比如"-benchtime, -cpu, -list, -parallel, -run, -short, -timeout, -failfast, -v" 这些都是cache友好的，如果运行`go test`还有其他类型的命令，则不会cache结果。
  - 所以如果想要取消cache最简单的做法就是增加"-count=1"flag；
  - 环境变量，文件等的变更并不会引起cahce失效；

`go test` 支持一下flag：

- args : 所有在`-args`后的参数，都作为参数传递给`go test`编译出的二进制包；
- c ： 只编译pkg.test，但是并不运行二进制；
- exec： 运行二进制文件
- json： 按照json 方式输出测试结果;
- o : 编译二进制到指定文件内；


# `go help testflag`

这部分flag同时也可以作用于`go test`命令。其中部分flag可以输出用于`go tool pprof`的文件。

-bench regexp：只运行benchmark测试，默认情况下不运行benchmark测试。使用 '-bench .' 或者 '-bench=.' 可以运行所有性能测试。正则表达式按照'/'分割，如果只运行某个子报的benchmark的话，需要依次匹配。所有parents测试包，都会按照`b.N=1`运行来标示运行子benchmark。比如`-bench=X/Y`，X只会运行`b.N=1`，但是所有Y匹配的包会全量运行。


-benchtime t： 指定测试时间 , 比如1s指定运行1秒钟，100x运行100次

-count n : 指定运行多少次用例。如果-cpu指定了cpu数量，运行次数将会是 [n*cpu]

-cover: 输出代码测试覆盖率

-covermode set,count,atomic: 覆盖度的计算公式

-coverpkg pattern1,pattern2,pattern3 ： 需要做代码覆盖分析的包

-cpu 1,2,4 ： 指定cpu数量，会基于一个cpu跑一次，基于2个cpu跑一次，基于4个跑一次。对于benchmarks和模糊测试有意义；

-failfast： 一个测试用例失败后就不去测试其他的测试用例啦；

-json： 以json方式输出测试报告；

-parallel n：当在测试中调用`t.Parallel()`时，所使用的并行数量，默认为GOMAXPROCS；

-run regexp： 指定运行测试的规则

```go

func TestChildTest(t *testing.T) {

	t.Logf("TestChildTest running")

	t.Run("hello", func(t *testing.T) {
		t.Logf("TestChildTest/hello running")
	})
	t.Run("hello1", func(t *testing.T) {
		t.Logf("TestChildTest/hello1 running")
	})
	t.Run("hello2", func(t *testing.T) {
		t.Logf("TestChildTest/hello2 running")
	})
}
```

运行 `go test -v  -run='TestChildTest/hello' ./...`
会输出：

```
=== RUN   TestChildTest
    child_test.go:16: TestChildTest running
=== RUN   TestChildTest/hello
    child_test.go:19: TestChildTest/hello running
=== RUN   TestChildTest/hello1
    child_test.go:22: TestChildTest/hello1 running
=== RUN   TestChildTest/hello2
    child_test.go:25: TestChildTest/hello2 running
```

运行： `go test -v  -run='TestChildTest/hello3' ./...`
只会输出：

```
=== RUN   TestChildTest
    child_test.go:16: TestChildTest running
--- PASS: TestChildTest (0.00s)

```

只会跑`TestChildTest`,其他的子测试不匹配，所以也不运行。


-short： 长时间运行的测试，降低运行时间；

-shuffle off,on,N ： 测试用例顺序随机化，可以关闭，可以开启，也可以指定随机的种子；

-timeout 　d : 测试二进制文件，超时时间，超时则panic。指定：0 则表示没有限制，默认是10分钟；

-v： 输出Log和Logf的信息

-vet list ：指定vet检查的列表



