## 第六章 并行

### 6.1 Golang scheduler 原理

![scheduler 原理](/ch6/images/scheduler.png)

如上图所示：

- 每个`M` 对应着一个计算机核心线程
- 每个`P` 绑定一个 `M` ，基于自己的队列完成基于 `G` 的调度 
- 每个`G` 是一个协程，类似于线程，但是更轻量级，这里这要理解，是`Golang` 应用程序调度的最小单位即可。

#### go routine 状态 & context switching

- Waiting : go routine处于此种状态下代表在等待某种资源，比如mutex等。
- Runnable: 在此状态下的go routine具备执行状态，在等待golang sheduler 调度到M上绑定计算资源以运行代码逻辑。
- Excuting: 此状态的goroutine处于正在执行状态，不需要等待任何资源，同时绑定了具体的计算资源，在运行代码逻辑。

只有处于`Runnable`的`go-routine`才可以被调度到计算资源，进而执行系统逻辑。

操作系统调度的最小单位是线程，golang context调度的最小粒度是go-routeine。golang 引用程序在执行过程中，会伴随着go-routine频繁抢占计算资源完成调度，调度过程中并不是无损的，需要设置各类寄存器，移动程序运行指针等。操作系统调度是基于Thread的，需要内核台和用户态的频繁切换，大概需要12K-18K的系统指令时间，Golang scheduler的调度是基于Golang的MGP模型的，只在用户态执行，因此效率会高很多，大概需要2.4K左右的系统指令时间。不管怎样，应该尽量避免频繁的`context swithing`，尤其在计算密集型应用；但是对于IO密集型应用这种`context switching`反而可以提高资源利用率。需要做好应用类型区分，合理利用golang scheduler的特性。


`Golang` scheduler 做了如下优化：

#### 网络IO多M复用

![Net puller goroutine调度](/ch6/images/net_puller_scheduler.png)

`G1` 调用网络IO 操作时会将其挂到`Net Puller`的`P`上，现有`P`可以调度其他`G`继续执行，避免CPU浪费，同时当`G1`网络IO结束后，自动会调度到活动的`P`上，完成指定工作；

#### 同步System call

![同步system call](/ch6/images/sync_system_call.png)

`G2` 进行非`异步system call`时，会新建出一个新的`M2`，对应的P会绑定到`M2`上，`P`基于`M2`继续执行代码，当`G2`完成`system call`后，将其放入某个`P`的`LRQ`上；

#### WorkStealing

![golang scheduler work stealing](/ch6/images/work_stealing.png)

当某个`P`的LRQ的所有goroutine都完成调度后，就开始从其他`P`中获取待完成的工作。


#### 用户态调度

![用户态调度](/ch6/images/user_scheduler.png)

`Golang`的Scheduler是用户态的调度，所以没有内核态和用户态之间的频繁切换；所以基于`go routine`的调度，context switching的消耗小很多。操作系统context switching 需要消耗1000~1500ns 大概消耗12K-18K的系统指令执行时间；golang的调度时间只有200ns，消耗只有2.4K左右的系统指令；




### 6.2 并发基础

可以使用 ```runtime.GOMAXPROCS(1)``` 将设置golang 应用只使用一个操作系统线程。使用```g := runtime.GOMAXPROCS(0)```可以设置golang应用使用尽量多的操作系统线程，这对于`container`环境下有很大意义，因为可以对比返回的线程数量和容器环境下分配给应用的cpu线程数量是否相等，如果不相等需要做一定调整，以达到系统最佳运行效果。


### 6.3 抢占scheduler

`Golang scheduler` 是基于抢占式的，所以在没有相关同步机制的情况下，不能对`go-routine`调度策略做任何的假设。比如：

```go

func main() {
    var wg sync.WaitGroup
    wg.Add(2)
    go func() {
        printHashes("A")
        wg.Done()
  }()

  go func() {
        printHashes("B")
        wg.Done()
  }()
  fmt.Println("Waiting To Finish")
  wg.Wait()
  fmt.Println("\nTerminating Program")
}

func printHashes(prefix string) {
    for i := 1; i <= 50000; i++ {
          num := strconv.Itoa(i)
          sum := sha1.Sum([]byte(num))
          fmt.Printf("%s: %05d: %x\n", prefix, i, sum)
  }
  fmt.Println("Completed", prefix)
}

```

统计context switching 次数：

```bash
$ ./example2 | cut -c1 | grep '[AB]' | uniq
B
A
B
A
B
A
B
A
B
A  9 Context Switches
$ ./example2 | cut -c1 | grep '[AB]' | uniq
B
A
B
A  3 Context Switches

```


可以发现每次运行代码，结果都不一样，因此不能对逻辑做太多假设。


### 6.4 数据竟态

数据竟态主要发生在当多个cpu线程对同一片数据同时进行读写时，会导致数据竟态。对于数据竟态如果没有合理的同步机制，会导致代码逻辑异常。即使引入了同步机制后，也可能存在多个cpu同时访问同一片`cache line`，导致多个cpu对于`cache-line`的频繁同步，导致性能下降。


### 6.5 数据竟态实例

数据竟态问题是非常难以察觉和修复的。

```go
var counter int
func main() {
    const grs = 2
    var wg sync.WaitGroup
    wg.Add(grs)
    for g := 0; g < grs; g++ {
        go func() {
            for i := 0; i < 2; i++ {
                value := counter   //对全局变量进行了一次读，如果其他cpu对于counter已经有变更了，那么需要同步cache line，性能降低
                value++           //此时如果调度到其他的go-routine，会导致其他go-routine读到原来的counter，出现“幻读/幻写”
                counter = value  //进行了一次设置，此时其他的cpu的`cache-line` 被设置为dirty 
            }
            wg.Done()
        }()
    }
    wg.Wait()
    fmt.Println("Counter:", counter)
}
```

多次运行程序，结果都一样，couter结果都是4，但是并不代表代码逻辑是没有问题的。
```go

var counter int
func main() {
    const grs = 2
    var wg sync.WaitGroup
    wg.Add(grs)
    for g := 0; g < grs; g++ {
        go func() {
            for i := 0; i < 2; i++ {
                value := counter   //对全局变量进行了一次读，如果其他cpu对于counter已经有变更了，那么需要同步cache line，性能降低
                value++           //此时如果调度到其他的go-routine，会导致其他go-routine读到原来的counter，出现“幻读/幻写”
                log.Println("logging") // 此处增加了日志打印，有IO操作，此时会发生调度
                counter = value  //进行了一次设置，此时其他的cpu的`cache-line` 被设置为dirty 
            }
            wg.Done()
        }()
    }
    wg.Wait()
    fmt.Println("Counter:", counter)
}
```

此时在运行代码，发现结果是2，并不是4，触发了代码中的"幻读"。


### 6.6 数据竟态探测

如上所示的例子，数据竟态的问题探测是十分重要的，那么如何探测呢？

```bash
  go build -race
  ./example 
```

以上编译流程会导致程序的运行效率降低20%，但是却具有数据竟态检测能力。

``` bash
2021/02/01 17:30:52 logging
2021/02/01 17:30:52 logging
2021/02/01 17:30:52 logging
==================
WARNING: DATA RACE
Write at 0x000001278d88 by goroutine 8:
  main.main.func1()
      /data_race/example1/example1.go:41 +0xa6
Previous read at 0x000001278d88 by goroutine 7:
  main.main.func1()
      /data_race/example1/example1.go:38 +0x4a
Goroutine 8 (running) created at:
  main.main()
      /data_race/example1/example1.go:36 +0xaf
Goroutine 7 (finished) created at:
  main.main()
      /data_race/example1/example1.go:36 +0xaf
==================
2021/02/01 17:30:52 logging
Final Counter: 2
Found 1 data race(s)
```

可以发现已经发现了数据竟态问题。

![data_race_dect](/ch6/images/data_race_dect.png)

如上所示，38行和第41行有数据竟态问题。

### 6.7 atomics

如上所示的例子如何修复呢 ？可以使用atomic(原子操作)，来解决。`atomics`的实现方式有很多种，很多实现都是基于硬件，比如CAS或者基于cache-line实现等。`atomics` 的优势在于`go-routine`处于自旋状态，不会发生`context switching`，这对于异常短暂的操作，可以有效提高计算资源利用效率。

```go

var counter int32
func main() {
    const grs = 2
    var wg sync.WaitGroup
    wg.Add(grs)
    for g := 0; g < grs; g++ {
        go func() {
            for i := 0; i < 2; i++ {
                atomic.AddInt32(&counter, 1) // 增加atomic，防止数据竟态问题
            wg.Done()
        }()
  }
  wg.Wait()
      fmt.Println("Counter:", counter)
  }
}
```


### 6.8 Mutex

如上所示的例子，还可以通过`mutex`来解决。
```go
func main() {
  const grs = 2
  var wg sync.WaitGroup
  wg.Add(grs)
  var mu sync.Mutex
  for g := 0; g < grs; g++ {
      go func() {
          for i := 0; i < 2; i++ {
              mu.Lock() //加锁，保护数据
              {
                  value := counter
                  value++
                  counter = value
              }
              mu.Unlock()//释放锁
          wg.Done()
      }()
  }
  wg.Wait()
      fmt.Println("Counter:", counter)
  }
```

如上所示，可以说过`mutex`加锁保护数据，`mutex`与`automic`的区别在于`mutex`可以保护一段代码，但是`automic`能保证正常访问/读写数据。

golang对于mutex的实现做了优化，mutex有两种模式：

- 正常模式，对于mutex加锁或进行自旋，当某个go-routine等待时间超过1M ns后，转换为饥饿模式
- 饥饿模式，对于mutex加锁直接加入到队列尾端，即使是当前锁是空闲的，也需要排队



### 6.9 读写锁

对于读多写少的情况，读写锁更适合。 读写锁原理为：

- 加读锁，对于没有写锁的情况下，直接对于已有读计数增加1，
  - 如果已经有写锁，读锁等待计数+1， 排队等待
  - 如果有读锁等待，则读锁等待计数+1
  - 如果没有写锁，直接正在读锁计数+1 
  - 释放读锁时，如果没有正在执行的读锁，并且有写锁等待，则唤醒写锁
- 加写锁
  - 如果已经有读锁，将读锁计数-maxlockcount，标示有写锁等待
  - 如果有写锁，则等待
  - 如果没有写锁则加锁成功
  - 释放锁时，如果有读锁等待，则唤醒读锁

### 6.10 channel 语义

[chan源码解读](/ch6/chan_source_code.md)


应该理解chan为一种信息的同步机制，而不是一个数据结构。这样在使用chan的时候考虑更多的是发送/接收信号，而不是读取/写入信息。如果在使用chan的时候没有信号语义的话，应该考虑是否应该使用chan来解决问题。

在使用chan时，一般需要考虑以下三个问题：
- 是否需要确保信号一定被收到
- 信号是否需要携带数据
- chan的状态

对于信号是否能被收到，是比较难以得到保证的。因为一个信号发送出去后，一方便需要一个接收者来接收，接收者何时出现是一个问题。另一方面，信号何时被接收的延时也是一个问题。

如果需要在信号传输过程中携带数据，则需要感知，信号的传输是信号和接收者1对1的，所以如果想多个goroutine接收到信号，需要在发一个信号。但是如果不希望接收数据，则可以同时通知多个goroutine，这对于取消工作，关闭工作流程场景十分重要的。

信号的状态有三种： 
- nil chan，对这类型的chan 发送和接收都会永久阻塞，关闭nil chan会panic
- open状态，一个open的chan，对于没有缓存的chan，发送和接收者需要同时存在，才能建立信号发送接收的流程。对于buffer的chan，可以先发送，或者先接收。当chan空时，接收者会阻塞，当chan缓存满时，发送者会阻塞。
- closed状态。使用`close(ch)`可以关闭chan，应该为了状态转换以关闭chan，而不是为了释放资源关闭chan（因为把chan当作信号机制用，而不是作为一个结构用）。给一个关闭的chan发送信号会导致`panic`，从一个关闭的chan接收信号会立刻收到一个`0状态`的信号。

### 6.11 channel 模式

chan的主要用例模式为以下几种：

#### 6.11.1 等待结果

```go
func waitForResult() {
    ch := make(chan string)
    go func() { //开启子协程工作
        time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
        ch <- "data"
        fmt.Println("child : sent signal")
    }()
    //等待结果完成
    d := <-ch
    fmt.Println("parent : recv'd signal :", d)
    time.Sleep(time.Second)
    fmt.Println("-------------------------------------------------")
}
```


#### 6.11.2 扇入扇出

```go
func fanOut() {
    children := 2000
    ch := make(chan string, children)

    for c := 0; c < children; c++ { //开启多个goroutine工作
        go func(child int) {
            time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
            ch <- "data"
            fmt.Println("child : sent signal :", child)
    }(c) }

    for children > 0 { //等待指定数量的工作完成
        d := <-ch
        children--
        fmt.Println(d)
        fmt.Println("parent : recv'd signal :", children)
    }

    time.Sleep(time.Second)
    fmt.Println("-------------------------------------------------")
}
```

#### 6.11.3 等待工作

```go
func waitForTask() {
    ch := make(chan string)
    go func() { // 等待工作
        d := <-ch
        fmt.Println("child : recv'd signal :", d)
    }()
    time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
    ch <- "data" //发送工作
    fmt.Println("parent : sent signal")
    time.Sleep(time.Second)
    fmt.Println("-------------------------------------------------")
}
```


#### 6.11.4 运行池

```go
func pooling() {
    ch := make(chan string)
    g := runtime.GOMAXPROCS(0) //获取使用最大的thread数量
    //开启一个协程池
    for c := 0; c < g; c++ {
        go func(child int) {
            //所有协程都消费工作内容，开始完成工作
            for d := range ch {
                fmt.Printf("child %d : recv'd signal : %s\n", child, d)
            }
            fmt.Printf("child %d : recv'd shutdown signal\n", child)
        }(c)
    }
    const work = 100
    //开始发布工作
    for w := 0; w < work; w++ {
        ch <- "data"
        fmt.Println("parent : sent signal :", w)
    }
    close(ch) //标示工作已经派发完成
    fmt.Println("parent : sent shutdown signal")
    time.Sleep(time.Second)
    fmt.Println("-------------------------------------------------")
}

```

#### 6.11.5 丢弃

```go

func drop() {
    const cap = 100
    ch := make(chan string, cap)
    go func() {
        // 接收工作，开始处理
        for p := range ch {
            fmt.Println("child : recv'd signal :", p)
        }
    }()
    const work = 2000
    for w := 0; w < work; w++ {
        select {
        case ch <- "data": //有能力处理工作
            fmt.Println("parent : sent signal :", w)
        default: //没有能力处理工作，丢弃一个work
            fmt.Println("parent : dropped data :", w)
        }
    }
    close(ch)
    fmt.Println("parent : sent shutdown signal")
    time.Sleep(time.Second)
    fmt.Println("-------------------------------------------------")
}

```

#### 6.11.6 取消

```go
func cancellation() {
    duration := 150 * time.Millisecond
    ctx, cancel := context.WithTimeout(context.Background(), duration)
    defer cancel()
    
    ch := make(chan string, 1)
    go func() {
        time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)
        ch <- "data"
    }()

    select {
    case d := <-ch: //在指定时间内完成了工作
        fmt.Println("work complete", d)
    case <-ctx.Done(): //在指定时间内，没有完成工作
        fmt.Println("work cancelled")
    }

    time.Sleep(time.Second)
    fmt.Println("-------------------------------------------------")
}
```



#### 6.11.7 信号限流

```go
func fanOutSem() {
    children := 2000
    ch := make(chan string, children)
    g := runtime.GOMAXPROCS(0) //获取可用线程数量
    sem := make(chan bool, g)

    for c := 0; c < children; c++ {
        go func(child int) {
          sem <- true //开始时抢占一个并行
               {
                t := time.Duration(rand.Intn(200)) * time.Millisecond
                time.Sleep(t)
                ch <- "data"
                fmt.Println("child : sent signal :", child)
           }
           <-sem //结束时，释放并行
       }(c)
    }

    for children > 0 {
        d := <-ch //接收结果
        children--
        fmt.Println(d)
        fmt.Println("parent : recv'd signal :", children)
    }

    time.Sleep(time.Second)
    fmt.Println("-------------------------------------------------")
}
```


#### 6.11.8 有限工作池

```go
func boundedWorkPooling() {
    work := []string{"paper", "paper", "paper", "paper", 2000: "paper"}

    g := runtime.GOMAXPROCS(0) //获取可用线程数量
    var wg sync.WaitGroup

    wg.Add(g)
    ch := make(chan string, g)
    for c := 0; c < g; c++ {
        go func(child int) { //根据线程数量建立指定大小的工作池
            defer wg.Done()
            for wrk := range ch {
                fmt.Printf("child %d : recv'd signal : %s\n", child, wrk)
            }
            fmt.Printf("child %d : recv'd shutdown signal\n", child)
        }(c)
}
    for _, wrk := range work {
        ch <- wrk //添加工作
    }
    close(ch)
    wg.Wait() //等待全部工作完成
    time.Sleep(time.Second)
    fmt.Println("-------------------------------------------------")
}
```


#### 6.11.9 超时重试

```go
func retryTimeout(ctx context.Context, retryInterval time.Duration, check func(ctx context.Context) error) {

    for {

        fmt.Println("perform user check call")
        if err := check(ctx); err == nil { //如果已经正常则直接返回
            fmt.Println("work finished successfully")
            return 
        }

        fmt.Println("check if timeout has expired")
        if ctx.Err() != nil { // 如果超过约定的timeout则返回错误
            fmt.Println("time expired 1 :", ctx.Err())
            return 
        }

        fmt.Printf("wait %s before trying again\n", retryInterval)
        t := time.NewTimer(retryInterval)

        select {
        case <-ctx.Done(): //已经超时了
            fmt.Println("timed expired 2 :", ctx.Err())
            t.Stop()
            return
        case <-t.C: //等待一段时间重试
            fmt.Println("retry again")
        } 
    }
}
```


#### 6.11.10 chan 结束context

```go
func channelCancellation(stop <-chan struct{}) {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    go func() {
        select {
        case <-stop: //如果收到结束信号则cancel context
            cancel()
        case <-ctx.Done():
        }
    }()

    func(ctx context.Context) error {
        //发起一个http请求，需要使用context来结束请求
        req, err := http.NewRequestWithContext(
            ctx,
            http.MethodGet,
            "https://www.ardanlabs.com/blog/index.xml",
            nil,
       )
       if err != nil {
          return err 
        }
        _ , err = http.DefaultClient.Do(req)
        if err != nil {
          return err 
        }

        return nil
    }(ctx) //基于context请求
}
```
