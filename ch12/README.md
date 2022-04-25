# 第12本章 Tracing

本章将学习如何对在线应用tracing。

## 12.1 样例代码

在docs中寻找包含topic的文章列表；

```go
func freq(topic string, docs []string) int {
	var found int

	for _, doc := range docs {
		file := fmt.Sprintf("%s.xml", doc[:8])
		f, err := os.OpenFile(file, os.O_RDONLY, 0)
		if err != nil {
			log.Printf("Opening Document [%s] : ERROR : %v", doc, err)
			return 0
		}

		data, err := io.ReadAll(f)
		f.Close()
		if err != nil {
			log.Printf("Reading Document [%s] : ERROR : %v", doc, err)
			return 0
		}

		var d document
		if err := xml.Unmarshal(data, &d); err != nil {
			log.Printf("Decoding Document [%s] : ERROR : %v", doc, err)
			return 0
		}

		for _, item := range d.Channel.Items {
			if strings.Contains(item.Title, topic) {
				found++
				continue
			}

			if strings.Contains(item.Description, topic) {
				found++
			}
		}
	}

	return found
}
```

文章内容如下所示：

```go
type (
	item struct {
		XMLName     xml.Name `xml:"item"`
		Title       string   `xml:"title"`
		Description string   `xml:"description"`
	}

	channel struct {
		XMLName xml.Name `xml:"channel"`
		Items   []item   `xml:"item"`
	}

	document struct {
		XMLName xml.Name `xml:"rss"`
		Channel channel  `xml:"channel"`
	}
)
```

主程序：

```go
func main() {

	docs := make([]string, 4000)
	for i := range docs {
		docs[i] = fmt.Sprintf("newsfeed-%.4d.xml", i)
	}

	topic := "president"
	n := freq(topic, docs)

	log.Printf("Searching %d files, found %s %d times.", len(docs), topic, n)
}
```

运行应用

```bash
go build -o trace ./trace.go
time ./trace
```

运行效率：

```
2022/04/24 20:46:36 Searching 4000 files, found president 28000 times.
./trace  3.35s user 0.25s system 90% cpu 3.975 total
```

## 12.2 产生tracing

```go
func main() {
	trace.Start(os.Stdout)
	defer trace.Stop()

	docs := make([]string, 4000)
	for i := range docs {
		docs[i] = fmt.Sprintf("newsfeed-%.4d.xml", i)
	}

	topic := "president"
	n := freq(topic, docs)

	log.Printf("Searching %d files, found %s %d times.", len(docs), topic, n)
}
```


## 12.3 查看Tracing

查看trace `go tool trace t.out` 可以打开网页端

![1-trace](/ch12/images/1-trace.png)

如上图所示：

- Goroutines： Goroutines数量
- Heap： 在使用堆内存大小
- Treads： 操作系统线程
- GC： 开始和结束GC详细信息
- Syscalls： 系统调用
- Procs： 逻辑核


可以在gc位置框选时间，查看gc发生次数，和gc占比时间。

![2-tracing-gc](/ch12/images/2-tracing-gc.png)


## 12.4 并行运行

更改算法，并行打开文件，查看tracing。

```go
func freqConcurrent(topic string, docs []string) int {
	var found int32

	g := len(docs)
	var wg sync.WaitGroup
	wg.Add(g)

	for _, doc := range docs {
		go func(doc string) {
			defer func() {
				wg.Done()
			}()

			file := fmt.Sprintf("%s.xml", doc[:8])
			f, err := os.OpenFile(file, os.O_RDONLY, 0)
			if err != nil {
				log.Printf("Opening Document [%s] : ERROR : %v", doc, err)
				return
			}

			data, err := io.ReadAll(f)
			f.Close()
			if err != nil {
				log.Printf("Reading Document [%s] : ERROR : %v", doc, err)
				return
			}

			var d document
			if err := xml.Unmarshal(data, &d); err != nil {
				log.Printf("Decoding Document [%s] : ERROR : %v", doc, err)
				return
			}

			for _, item := range d.Channel.Items {
				if strings.Contains(item.Title, topic) {
					found++ // 第一个版本
					continue
				}

				if strings.Contains(item.Description, topic) {
					found++ // 第一个版本
				}
			}
		}(doc)
	}

	wg.Wait()
	return int(found)
}
```

编译 `go build -race`

运行程序 ` ./trace > t.out`

```
==================
WARNING: DATA RACE
Read at 0x00c000158a5c by goroutine 16:
  main.freqConcurrent.func1()
      /Users/weenxin/go/src/github.com/weenxin/ulitmate_go_notebook_reading/ch12/trace/trace.go:145 +0x644
  main.freqConcurrent.func2()
      /Users/weenxin/go/src/github.com/weenxin/ulitmate_go_notebook_reading/ch12/trace/trace.go:149 +0x58

Previous write at 0x00c000158a5c by goroutine 9:
  main.freqConcurrent.func1()
      /Users/weenxin/go/src/github.com/weenxin/ulitmate_go_notebook_reading/ch12/trace/trace.go:145 +0x657
  main.freqConcurrent.func2()
      /Users/weenxin/go/src/github.com/weenxin/ulitmate_go_notebook_reading/ch12/trace/trace.go:149 +0x58

Goroutine 16 (running) created at:
  main.freqConcurrent()
      /Users/weenxin/go/src/github.com/weenxin/ulitmate_go_notebook_reading/ch12/trace/trace.go:110 +0x2ba
  main.main()
      /Users/weenxin/go/src/github.com/weenxin/ulitmate_go_notebook_reading/ch12/trace/trace.go:54 +0x17c

Goroutine 9 (finished) created at:
  main  main.main()
      /Users/weenxin/go/src/github.com/weenxin/ulitmate_go_notebook_reading/ch12/trace/trace.go:54 +0x17c
==================
2022/04/24 21:53:54 Searching 4000 files, found president 28000 times.
Found 1 data race(s)
```

 发现有数据并发访问的问题。更改使用atomic访问。

 ```go

 func freqConcurrent(topic string, docs []string) int {
 	var found int32

 	g := len(docs)
 	var wg sync.WaitGroup
 	wg.Add(g)

 	for _, doc := range docs {
 		go func(doc string) {
 			defer func() {
 				wg.Done()
 			}()

 			file := fmt.Sprintf("%s.xml", doc[:8])
 			f, err := os.OpenFile(file, os.O_RDONLY, 0)
 			if err != nil {
 				log.Printf("Opening Document [%s] : ERROR : %v", doc, err)
 				return
 			}

 			data, err := io.ReadAll(f)
 			f.Close()
 			if err != nil {
 				log.Printf("Reading Document [%s] : ERROR : %v", doc, err)
 				return
 			}

 			var d document
 			if err := xml.Unmarshal(data, &d); err != nil {
 				log.Printf("Decoding Document [%s] : ERROR : %v", doc, err)
 				return
 			}

 			for _, item := range d.Channel.Items {
 				if strings.Contains(item.Title, topic) {
 					atomic.AddInt32(&found,1) //第二个版本，解决并发问题
 					continue
 				}

 				if strings.Contains(item.Description, topic) {
 					atomic.AddInt32(&found,1) //第二个版本解决并发问题
 				}
 			}
 		}(doc)
 	}

 	wg.Wait()
 	return int(found)
 }

 ```

 编译：` go build -race ` 运行： `./trace > t.out`

```
2022/04/24 22:00:25 Searching 4000 files, found president 28000 times.
```

## 12.5 cache 友好

我们开启了多个goroutine同时运行代码，代码中由于多个线程核心同时，频繁更改和读取`found`变量导致L1Cache会持续处于Dirty中。

![cache-friendly](/ch12/images/3-cache-friendly.png)


```go
func freqConcurrent(topic string, docs []string) int {
	var found int32

	g := len(docs)
	var wg sync.WaitGroup
	wg.Add(g)

	for _, doc := range docs {
		go func(doc string) {
			var lFound int32  //版本3
			defer func() {
				atomic.AddInt32(&found, lFound) //版本3
				wg.Done()
			}()

			file := fmt.Sprintf("%s.xml", doc[:8])
			f, err := os.OpenFile(file, os.O_RDONLY, 0)
			if err != nil {
				log.Printf("Opening Document [%s] : ERROR : %v", doc, err)
				return
			}

			data, err := io.ReadAll(f)
			f.Close()
			if err != nil {
				log.Printf("Reading Document [%s] : ERROR : %v", doc, err)
				return
			}

			var d document
			if err := xml.Unmarshal(data, &d); err != nil {
				log.Printf("Decoding Document [%s] : ERROR : %v", doc, err)
				return
			}

			for _, item := range d.Channel.Items {
				if strings.Contains(item.Title, topic) {
					lFound ++  //版本3
					continue
				}

				if strings.Contains(item.Description, topic) {
					lFound ++  //版本3
				}
			}
		}(doc)
	}

	wg.Wait()
	return int(found)
}

```


## 12.6 并行运行结果

编译： `go build` 运行： ` time ./trace > t.out`

```
2022/04/24 22:49:13 Searching 4000 files, found president 28000 times.
./trace > t.out  6.31s user 0.32s system 257% cpu 2.577 total
```

查看tracing, `go tool trace t.out`

![4-tracing-concurrency](/ch12/images/4-tracing-concurrency.png)

GC时间如下所示：

![5-tracing-concurrency-gc](/ch12/images/5-tracing-concurrency-gc.png)

- 大约占比 （795,319,812 ns/2,425,732,028 ns）33%左右；
- 大约节省了（3.75-2.57时间）


![6-tracing-concurrency-gc-graph](/ch12/images/6-tracing-concurrency-gc-graph.png)

如上图所示可以看到GC发生的时间点，GoRoutine从1000多到2个。内存也逐渐释放。1000+的goroutine会频繁发生context-switching。

![7-context-switching](/ch12/images/7-context-switching.png)

如上图可看在Goroutine很多的时候，groutine在频繁的switching。

## 12.7 Goroutine 池

由于开启几千个goroutine会导致context-switching浪费时间，为了即获得足够多的并行能力，又降低context-switching时间，可以使用goroutine池，降低goroutine频繁新建和释放的gc压力。

```go
func freqProcessors(topic string, docs []string) int {
	var found int32

	g := runtime.GOMAXPROCS(0) //获取系统核心数量
	var wg sync.WaitGroup
	wg.Add(g)

	ch := make(chan string, g)

	for i := 0; i < g; i++ {
		go func() {
			var lFound int32
			defer func() {
				atomic.AddInt32(&found, lFound)
				wg.Done()
			}()

			for doc := range ch {
				file := fmt.Sprintf("%s.xml", doc[:8])
				f, err := os.OpenFile(file, os.O_RDONLY, 0)
				if err != nil {
					log.Printf("Opening Document [%s] : ERROR : %v", doc, err)
					return
				}

				data, err := io.ReadAll(f)
				f.Close()
				if err != nil {
					log.Printf("Reading Document [%s] : ERROR : %v", doc, err)
					return
				}

				var d document
				if err := xml.Unmarshal(data, &d); err != nil {
					log.Printf("Decoding Document [%s] : ERROR : %v", doc, err)
					return
				}

				for _, item := range d.Channel.Items {
					if strings.Contains(item.Title, topic) {
						lFound++
						continue
					}

					if strings.Contains(item.Description, topic) {
						lFound++
					}
				}
			}
		}()
	}

	for _, doc := range docs {
		ch <- doc
	}
	close(ch)

	wg.Wait()
	return int(found)
}
```
## 12.8 线程池结果

编译： `go build` 运行 :`time ./trace  > t.out`

```
2022/04/25 09:24:04 Searching 4000 files, found president 28000 times.
./trace > t.out  5.45s user 0.29s system 237% cpu 2.417 total
```

查看profiling文件： `go tool trace t.out`

![polling-tracing](/ch12/images/8-tracing-polling.png)

如上图所示
- gc时间占比： 20%
- goroutine最多4个
- 堆内存分配也比较稳定

![polling-tracing-gc](/ch12/images/9-tracing-polling-gc.png)

## 12.9 垃圾回收大小

运行： `time GOGC=1000 ./trace > t.out`

```
2022/04/25 09:46:12 Searching 4000 files, found president 28000 times.
GOGC=1000 ./trace > t.out  5.21s user 0.21s system 296% cpu 1.829 total
```

默认情况下GOGC是100，这代表着如果当前堆内存使用大小为4MB，当堆内存达到8MB时，就会发生垃圾回收。如果`GOGC`为1000则当内存达到40MB时，才会发生垃圾回收。

![10-gogc](/ch12/images/10-gogc.png)

如上图所示，可以发现

- 堆内存大大增加
- gc发生次数大大降低


![11-gogc-time-percent](/ch12/images/11-gogc-time-percent.png)

可以发现gc时间占比大大降低。


## 12.10 Task and Region

可以针对某个文件或者某段代码tracing查看各个阶段程序运行状态。

```go
func freqProcessors(topic string, docs []string) int {
	var found int32

	g := runtime.GOMAXPROCS(0)
	var wg sync.WaitGroup
	wg.Add(g)

	ch := make(chan string, g)

	for i := 0; i < g; i++ {
		go func() {
			var lFound int32
			defer func() {
				atomic.AddInt32(&found, lFound)
				wg.Done()
			}()

			for doc := range ch {
				ctx, task := trace.NewTask(context.Background(),doc) // 增加Task

				reg := trace.StartRegion(ctx, "OpenFile") // 增加Region
				file := fmt.Sprintf("%s.xml", doc[:8])
				f, err := os.OpenFile(file, os.O_RDONLY, 0)
				if err != nil {
					log.Printf("Opening Document [%s] : ERROR : %v", doc, err)
					return
				}
				reg.End() // 结束Region

				reg = trace.StartRegion(ctx, "ReadAll") // 增加Region
				data, err := io.ReadAll(f)
				f.Close()
				if err != nil {
					log.Printf("Reading Document [%s] : ERROR : %v", doc, err)
					return
				}
				reg.End() // 结束Region

				reg = trace.StartRegion(ctx, "Unmarshal") // 增加Region
				var d document
				if err := xml.Unmarshal(data, &d); err != nil {
					log.Printf("Decoding Document [%s] : ERROR : %v", doc, err)
					return
				}
				reg.End() // 结束Region

				reg = trace.StartRegion(ctx, "Contains") // 增加Region
				for _, item := range d.Channel.Items {
					if strings.Contains(item.Title, topic) {
						lFound++
						continue
					}

					if strings.Contains(item.Description, topic) {
						lFound++
					}
				}
				reg.End()// 结束Region
				task.End() // 结束Task
			}
		}()
	}

	for _, doc := range docs {
		ch <- doc
	}
	close(ch)

	wg.Wait()
	return int(found)
}
```

编译： `go build` 运行： `time ./trace > t.out`

查看profiling ： `go tool trace t.out`

点击 `User-defined tasks` 可以查看task的运行情况：

![task-and-region](/ch12/images/12-task-and-region.png)


点击具体的tas可以查看各个region运行情况：

```

When	Elapsed	Goroutine ID	Events
0.000807833s	5.279889ms		Task 2 (goroutine view) (complete)
0.000807833	 .         	22	task newsfeed-0000.xml (id 2, parent 0) created
0.000810166	 .     2333	22	region OpenFile started (duration: 37.584µs)
0.000864055	 .    53889	22	region ReadAll started (duration: 741.75µs)
0.001607638	 .   743583	22	region Unmarshal started (duration: 4.469612ms)
0.006079833	 .  4472195	22	region Contains started (duration: 6µs)
0.006087722	 .     7889	22	task end
GC:0s
```
点击链接可以查看`goroutine`视角

![goroutine-tracing](/ch12/images/13-go-routine-tracing.png)




