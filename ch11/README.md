# 第11章 优化在线代码

本章将讲述如何优化在线代码。


## 11.1  样例代码

调优使用的样例代码如下所示：

![server.png](/ch11/images/server.svg)

如上图所示：

- 服务一个查询新闻的http服务程序
- 用户输入查询关键字后，先从缓存层中查询数据，如果缓存没有数据则从各个新闻源中获取数据；

由于网络和其他原因，对代码做了改进。不再从新闻源获取数据，而是随机生成一堆一些数据。

```go
const MessageLength = 100

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func mkItems() []Item {
	var results []Item
	for i :=0 ; i < 100; i++ {
		results = append(results,Item{
			Title:   randStringRunes(MessageLength),
			Link:    randStringRunes(MessageLength),
			Description : randStringRunes(MessageLength),
		})
	}
	return results
}
```


```go
//  查询新闻信息
func rssSearch(uid, term, engine, uri string) ([]Result, error) {
	var mu *sync.Mutex
	fetch.Lock()
	{
		var found bool
		mu, found = fetch.m[uri]
		if !found {
			mu = &sync.Mutex{}
			fetch.m[uri] = mu
		}
	}
	fetch.Unlock()

	var d Document
	mu.Lock()
	{
		// 从cache中获取数据
		v, found := cache.Get(uri)

		// 基于cache查询结果
		switch {
		case found: // 如果查询到直接使用
			d = v.(Document)

		default: //否则就重新生成一堆结果
			d.Channel.Items = mkItems()
			// Save this document into the cache.
			cache.Set(uri, d, expiration)
			log.Println("reloaded cache", uri)
		}
	}
	mu.Unlock()

	// Create an empty slice of results.
	results := []Result{}
	// 筛选结果
	for _, item := range d.Channel.Items {
	    // 忽略大小写
		if strings.Contains(strings.ToLower(item.Description), strings.ToLower(term)) {
			results = append(results, Result{
				Engine:  engine,
				Title:   item.Title,
				Link:    item.Link,
				Content: item.Description,
			})
		}
	}

	return results, nil
}
```

## 11.2 产生垃圾回收Trace

```bash
cd ch11
GODEBUG=gctrace=1  go run ./main.go
```


```
2022/04/24 11:07:04.186279 rss.go:95: reloaded cache http://feeds.bbci.co.uk/news/rss.xml
2022/04/24 11:07:04.188875 rss.go:95: reloaded cache http://rss.nytimes.com/services/xml/rss/nyt/HomePage.xml
2022/04/24 11:07:04.189291 rss.go:95: reloaded cache http://rss.cnn.com/rss/cnn_topstories.rss
2022/04/24 11:07:04.189740 rss.go:95: reloaded cache http://feeds.bbci.co.uk/news/world/rss.xml
2022/04/24 11:07:04.192021 rss.go:95: reloaded cache http://feeds.bbci.co.uk/news/politics/rss.xml
2022/04/24 11:07:04.195546 rss.go:95: reloaded cache http://rss.cnn.com/rss/cnn_world.rss
2022/04/24 11:07:04.197026 rss.go:95: reloaded cache http://rss.nytimes.com/services/xml/rss/nyt/US.xml
2022/04/24 11:07:04.199415 rss.go:95: reloaded cache http://rss.nytimes.com/services/xml/rss/nyt/Politics.xml
2022/04/24 11:07:04.203931 rss.go:95: reloaded cache http://rss.nytimes.com/services/xml/rss/nyt/Business.xml
2022/04/24 11:07:04.204262 rss.go:95: reloaded cache http://rss.cnn.com/rss/cnn_us.rss
2022/04/24 11:07:04.204952 rss.go:95: reloaded cache http://feeds.bbci.co.uk/news/world/us_and_canada/rss.xml
2022/04/24 11:07:04.206056 rss.go:95: reloaded cache http://rss.cnn.com/rss/cnn_allpolitics.rss
gc 1 @24.985s 0%: 0.053+5.9+0.006 ms clock, 0.21+0.32/1.2/0.26+0.026 ms cpu, 4->4->0 MB, 4 MB goal, 0 MB stacks, 0 MB globals, 4 P
gc 2 @61.101s 0%: 0.23+0.58+0.030 ms clock, 0.95+0.66/0.46/0.021+0.12 ms cpu, 4->4->0 MB, 4 MB goal, 0 MB stacks, 0 MB globals, 4 P
gc 3 @61.116s 0%: 0.13+0.78+0.082 ms clock, 0.53+0.30/0.51/0.29+0.32 ms cpu, 4->4->0 MB, 4 MB goal, 0 MB stacks, 0 MB globals, 4 P
gc 4 @61.131s 0%: 0.076+0.47+0.029 ms clock, 0.30+0.41/0.28/0.16+0.11 ms cpu, 4->4->0 MB, 4 MB goal, 0 MB stacks, 0 MB globals, 4 P
gc 5 @61.146s 0%: 0.081+0.55+0.020 ms clock, 0.32+0.30/0.26/0.10+0.083 ms cpu, 4->4->0 MB, 4 MB goal, 0 MB stacks, 0 MB globals, 4 P
gc 6 @61.160s 0%: 0.088+0.47+0.037 ms clock, 0.35+0.25/0.37/0+0.15 ms cpu, 4->4->0 MB, 4 MB goal, 0 MB stacks, 0 MB globals, 4 P
gc 7 @61.173s 0%: 0.072+0.34+0.016 ms clock, 0.29+0.23/0.26/0.14+0.066 ms cpu, 4->4->0 MB, 4 MB goal, 0 MB stacks, 0 MB globals, 4 P

```


比如：

```
gc 4 @61.131s 0%: 0.076+0.47+0.029 ms clock, 0.30+0.41/0.28/0.16+0.11 ms cpu, 4->4->0 MB, 4 MB goal, 0 MB stacks, 0 MB globals, 4 P
```


**gc 4 @61.131s 0%**

- 第4次 垃圾回收
- 程序已经运行了61.131s
- gc时间占比是0%

**0.076+0.47+0.029 ms clock**

- MarkUp setup(STW) :0.076ms
- Concurrent Marking :   0.47ms
- Mark Termination(STW): 0.029 ms

**0.30+0.41/0.28/0.16+0.11 ms cpu**

- 0.30 MarkUp setup(STW)
- 0.41 Concurrent Mark - 协助gc时间
- 0.28 Concurrent Mark - Background GC time
- 0.16 Concurrent Mark - Idle GC time
- 0.11 Mark Termination

**4->4->0 MB**

- 4M Heap memory in-use before the Marking started
- 4M  Heap memory in-use after the Marking finished
- 0M Heap memory marked as live after the Marking finished

**4 MB goal**

- Collection goal for heap memory in-use after Marking finished；内存达到这个值时进行再次垃圾回收

**4 P**

- 4个线程


## 11.3 增加负载

运行  `hey -m POST -c 1 -n 100 "http://localhost:5000/search?term=biden&cnn=on&bbc=on&nyt=on"` 命令对系统增加负载。

```
Summary:
  Total:        0.0963 secs
  Slowest:      0.0141 secs
  Fastest:      0.0007 secs
  Average:      0.0010 secs
  Requests/sec: 1038.4396


Response time histogram:
  0.001 [1]     |
  0.002 [97]    |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.003 [1]     |
  0.005 [0]     |
  0.006 [0]     |
  0.007 [0]     |
  0.009 [0]     |
  0.010 [0]     |
  0.011 [0]     |
  0.013 [0]     |
  0.014 [1]     |


Latency distribution:
  10% in 0.0007 secs
  25% in 0.0007 secs
  50% in 0.0008 secs
  75% in 0.0009 secs
  90% in 0.0010 secs
  95% in 0.0014 secs
  99% in 0.0141 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0001 secs, 0.0007 secs, 0.0141 secs
  DNS-lookup:   0.0001 secs, 0.0000 secs, 0.0095 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0003 secs
  resp wait:    0.0008 secs, 0.0006 secs, 0.0020 secs
  resp read:    0.0000 secs, 0.0000 secs, 0.0002 secs

Status code distribution:
  [200] 100 responses

```

## 11.4 增加调优入口

```go
package main
import (
    _ "net/http/pprof" // 引入pprof包，初始化
)
```

在 `src/net/http/pprof/pprof.go` 中可以看到引入pprof包后会增加以下http接口。

```go

func init() {
    http.HandleFunc("/debug/pprof/", Index)
    http.HandleFunc("/debug/pprof/cmdline", Cmdline)
    http.HandleFunc("/debug/pprof/profile", Profile)
    http.HandleFunc("/debug/pprof/symbol", Symbol)
    http.HandleFunc("/debug/pprof/trace", Trace)
}
```


访问如下url： `http://localhost:5000/debug/pprof`

```
heap profile: 8: 11536 [2923: 750016] @ heap/1048576
```

- 8 8个活跃对象
- 11536 活跃对象占用11536字节
- 2923 一共分配了2923个对象
- 750016 总共分配了750016字节



```
1: 112 [879: 98448] @ 0x10e6d2d 0x10e6d15 0x10e6d7f 0x13712b1 0x1370c96 0x1068f81
#	0x10e6d2c	strings.(*Builder).grow+0xec								/usr/local/go/src/strings/builder.go:67
#	0x10e6d14	strings.(*Builder).Grow+0xd4								/usr/local/go/src/strings/builder.go:81
#	0x10e6d7e	strings.ToLower+0x13e									/usr/local/go/src/strings/strings.go:600
#	0x13712b0	github.com/weenxin/ulitmate_go_notebook_reading/ch11/profiling/search.rssSearch+0x470	/Users/weenxin/go/src/github.com/weenxin/ulitmate_go_notebook_reading/ch11/profiling/search/rss.go:105
#	0x1370c95	github.com/weenxin/ulitmate_go_notebook_reading/ch11/profiling/search.NYT.Search+0x155	/Users/weenxin/go/src/github.com/weenxin/ulitmate_go_notebook_reading/ch11/profiling/search/nyt.go:25
```

如上所示是有内存分配的对象：

- 1 ： 总共 1 个活跃对象
- 112： 活跃对象占用112字节
- 879： 总共分配879个对象
- 98448： 总共占用98448字节

## 11.4 查看内存profile

运行 `go tool pprof -noinlines http://localhost:5000/debug/pprof/allocs`

输入： `top 15 -cum`


查看15个内存分配最多的调用（以cum视角）

```
(pprof) top 15 -cum
Showing nodes accounting for 1347.70MB, 91.94% of 1465.86MB total
Dropped 62 nodes (cum <= 7.33MB)
Showing top 15 nodes out of 30
      flat  flat%   sum%        cum   cum%
         0     0%     0%  1312.64MB 89.55%  github.com/weenxin/ulitmate_go_notebook_reading/ch11/profiling/search.rssSearch
 1307.14MB 89.17% 89.17%  1307.14MB 89.17%  strings.ToLower
         0     0% 89.17%   470.55MB 32.10%  github.com/weenxin/ulitmate_go_notebook_reading/ch11/profiling/search.BBC.Search
         0     0% 89.17%   441.55MB 30.12%  github.com/weenxin/ulitmate_go_notebook_reading/ch11/profiling/search.NYT.Search
         0     0% 89.17%   400.54MB 27.32%  github.com/weenxin/ulitmate_go_notebook_reading/ch11/profiling/search.CNN.Search
         0     0% 89.17%   149.21MB 10.18%  net/http.(*conn).serve
         0     0% 89.17%   133.20MB  9.09%  github.com/braintree/manners.(*gracefulHandler).ServeHTTP
   20.04MB  1.37% 90.54%   133.20MB  9.09%  github.com/weenxin/ulitmate_go_notebook_reading/ch11/profiling/service.handler
         0     0% 90.54%   133.20MB  9.09%  net/http.(*ServeMux).ServeHTTP
         0     0% 90.54%   133.20MB  9.09%  net/http.HandlerFunc.ServeHTTP
         0     0% 90.54%   133.20MB  9.09%  net/http.serverHandler.ServeHTTP
   19.52MB  1.33% 91.87%   101.15MB  6.90%  github.com/weenxin/ulitmate_go_notebook_reading/ch11/profiling/service.render
    0.50MB 0.034% 91.91%    81.63MB  5.57%  github.com/weenxin/ulitmate_go_notebook_reading/ch11/profiling/service.executeTemplate
         0     0% 91.91%    81.13MB  5.53%  html/template.(*Template).Execute
    0.50MB 0.034% 91.94%    81.13MB  5.53%  text/template.(*Template).execute

```

当然也可以使用 `web rssSearch`   通过Web查看内存分配情况。

从上面可以看出，`rssSearch`分配了非常多的内存。

```go
// Capture the data we need for our results if we find the search term.
	for _, item := range d.Channel.Items {
		if strings.Contains(strings.ToLower(item.Description), strings.ToLower(term)) {
			results = append(results, Result{
				Engine:  engine,
				Title:   item.Title,
				Link:    item.Link,
				Content: item.Description,
			})
		}
	}
```

原因是每次匹配的时候，都使用了`ToLower`忽略大小写。`ToLower`会基于入参返回一个新的字符串。这里会有很大的内存消耗。

由于有缓存，所以在从新闻源获取之后，对数据处理，降低字符串处理频率。

```go
        //优化后的代码
		results = append(results,Item{
			Title:   strings.ToLower(randStringRunes(MessageLength)) ,
			Link:    strings.ToLower(randStringRunes(MessageLength)),
			Description : strings.ToLower(randStringRunes(MessageLength)),
		})
```

将字符串处理流程前置，降低内存分配。

```go
    itemLower := strings.ToLower(term)

	// Capture the data we need for our results if we find the search term.
	for _, item := range d.Channel.Items {
		if strings.Contains(item.Description, itemLower) {
			results = append(results, Result{
				Engine:  engine,
				Title:   item.Title,
				Link:    item.Link,
				Content: item.Description,
			})
		}
	}
```

重新运行应用

- `GODEBUG=gctrace=1  go run ./main.go` :执行程序
- `hey -m POST -c 100 -n 10000 "http://localhost:5000/search?term=biden&cnn=on&bbc=on&nyt=on"`  增加压力


运行效率有提高：

```

Summary:
  Total:        1.1577 secs
  Slowest:      0.0887 secs
  Fastest:      0.0002 secs
  Average:      0.0109 secs
  Requests/sec: 8637.8386


Response time histogram:
  0.000 [1]     |
  0.009 [5316]  |■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■■
  0.018 [2800]  |■■■■■■■■■■■■■■■■■■■■■
  0.027 [1273]  |■■■■■■■■■■
  0.036 [407]   |■■■
  0.044 [106]   |■
  0.053 [49]    |
  0.062 [19]    |
  0.071 [20]    |
  0.080 [7]     |
  0.089 [2]     |


Latency distribution:
  10% in 0.0017 secs
  25% in 0.0040 secs
  50% in 0.0084 secs
  75% in 0.0152 secs
  90% in 0.0235 secs
  95% in 0.0284 secs
  99% in 0.0442 secs

Details (average, fastest, slowest):
  DNS+dialup:   0.0001 secs, 0.0002 secs, 0.0887 secs
  DNS-lookup:   0.0000 secs, 0.0000 secs, 0.0088 secs
  req write:    0.0000 secs, 0.0000 secs, 0.0053 secs
  resp wait:    0.0107 secs, 0.0002 secs, 0.0886 secs
  resp read:    0.0001 secs, 0.0000 secs, 0.0090 secs

Status code distribution:
  [200] 10000 responses

```

查看内存分配，运行`go tool pprof -noinlines http://localhost:5000/debug/pprof/allocs` ；


```
(pprof) top 15 -cum
Showing nodes accounting for 295.02MB, 70.74% of 417.05MB total
Dropped 34 nodes (cum <= 2.09MB)
Showing top 15 nodes out of 44
      flat  flat%   sum%        cum   cum%
         0     0%     0%   408.55MB 97.96%  net/http.(*conn).serve
         0     0%     0%   359.53MB 86.21%  github.com/braintree/manners.(*gracefulHandler).ServeHTTP
   67.15MB 16.10% 16.10%   359.53MB 86.21%  github.com/weenxin/ulitmate_go_notebook_reading/ch11/profiling/service.handler
         0     0% 16.10%   359.53MB 86.21%  net/http.(*ServeMux).ServeHTTP
         0     0% 16.10%   359.53MB 86.21%  net/http.HandlerFunc.ServeHTTP
         0     0% 16.10%   359.53MB 86.21%  net/http.serverHandler.ServeHTTP
   43.54MB 10.44% 26.54%   258.38MB 61.95%  github.com/weenxin/ulitmate_go_notebook_reading/ch11/profiling/service.render
    3.50MB  0.84% 27.38%   214.84MB 51.51%  github.com/weenxin/ulitmate_go_notebook_reading/ch11/profiling/service.executeTemplate
         0     0% 27.38%   211.34MB 50.67%  html/template.(*Template).Execute
       3MB  0.72% 28.10%   211.34MB 50.67%  text/template.(*Template).execute
         0     0% 28.10%   208.34MB 49.96%  text/template.(*state).walk
         0     0% 28.10%   177.84MB 42.64%  bytes.(*Buffer).Write
         0     0% 28.10%   177.84MB 42.64%  bytes.(*Buffer).grow
  177.84MB 42.64% 70.74%   177.84MB 42.64%  bytes.makeSlice
         0     0% 70.74%   103.25MB 24.76%  fmt.Fprint

```

可以看到 `rssSearch` 已经没有分配了。
