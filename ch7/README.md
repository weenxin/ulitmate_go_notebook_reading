## 第7章 测试

### 其他资料

[BDD & ginkgo](/ch7/BDD&Ginkgo.md)
[gorm & sqlmock](/ch7/SQLMock.md)
[go test 命令与参数](/ch7/GoTest.md)

### 7.1 单元测试

golang中测试做的比较好的点在于定义了什么是一个单元。golang的一个包就是一个单元，对应于代码树中的一个文件夹。编译器将每个包编译成为一个独立的二进制文件，已完成单元测试。

当我们谈及单元测试时，我们聚焦在一个独立的包中。这其中可能并不会限制我们不能使用外部系统（如数据库等）。这与使用多个包的集成测试是有区别的。

golang的工具栈已经对测试提供了很好的支持，只要按照如下方式就可以编写一个测试用例：
- 新建一个"*_test.go*开头的包；
- 函数名以*Test*开头，接收一个testing.T指针参数；


sample_test.go文件内容如下所示：

```go
package sample_test
import ("testing" )

func TestDownload(t *testing.T) {}
func TestUpload(t *testing.T) {}
```

如上所示的代码，使用`sample_test`作为包名，代表着只会测试包导出的函数，如果想要测试包内的所有函数，可以使用包的原名。

testing.T具有如下方法列表：

```go
type T
    func (c *T) Cleanup(f func())
    func (t *T) Deadline() (deadline time.Time, ok bool)
    func (c *T) Error(args ...interface{})
    func (c *T) Errorf(format string, args ...interface{})
    func (c *T) Fail()
    func (c *T) FailNow()
    func (c *T) Failed() bool
    func (c *T) Fatal(args ...interface{})
    func (c *T) Fatalf(format string, args ...interface{})
    func (c *T) Helper()
    func (c *T) Log(args ...interface{})
    func (c *T) Logf(format string, args ...interface{})
    func (c *T) Name() string
    func (t *T) Parallel()
    func (t *T) Run(name string, f func(t *T)) bool
    func (c *T) Skip(args ...interface{})
    func (c *T) SkipNow()
    func (c *T) Skipf(format string, args ...interface{})
    func (c *T) Skipped() bool
    func (c *T) TempDir() string
```

如上所示的代码，报告错误有两种方式：
- t.Log ： 输出日志，  使用`go test -v`输出相关日志；
- t.Fail ： 设置本测试用例为错误，继续执行测试用例；
- t.Error ： 设置本测试用例为错误，打印日志，并继续执行测试用例；
- t.FailNow ： 设置本测试用例为错误，并停止测试；
- t.Fatal ： 设置本测试用用例为错误，打印日志，并停止本测试用例执行；


```go
package sample_test
import (
    "testing"
    "http"
)
func TestDownload(t *testing.T) {
    url := "https://www.ardanlabs.com/blog/index.xml"
    statusCode := 200
    resp, err := http.Get(url) //请求资源
    if err != nil { //如果不是成功
        t.Fatalf("unable to issue GET on URL: %s: %s", url, err) //设置失败，打印日志，结束后续执行
    }
    defer resp.Body.Close() // 关闭链接
    if resp.StatusCode != statusCode { //判断状态码
       t.Log("exp:", statusCode)
       t.Log("got:", resp.StatusCode)
       t.Fatal("status codes don’t match")
} }
```

### 7.2 Table单元测试

很多单元测试，有很多分支，可以将多个分支组成一个表，每一行都是：输入，期望结果。这样可以复用一段代码，完成所有case的测试。

```go
package sample_test
import (
    "testing"
    "http"
)
func TestDownload(t *testing.T) {
    tt := []struct {
        url string //请求地址
        statusCode int //目标状态
    }{
        {"https://www.ardanlabs.com/blog/index.xml", http.StatusOK},
        {"http://rss.cnn.com/rss/cnn_topstorie.rss", http.StatusNotFound},
    }

    for _, test := range tt {
        resp, err := http.Get(test.url)
        if err != nil {
            t.Fatalf("unable to issue GET on URL: %s: %s", test.url, err)
        }
        defer resp.Body.Close()
        if resp.StatusCode != test.statusCode {
           t.Log("exp:", test.statusCode)
           t.Log("got:", resp.StatusCode)
           t.Fatal("status codes don’t match")
        }
    }
}
```

### 7.3 Mock Web调用

```go
package sample_test
import (
    "testing"
    "http"
    "httptest"
)

var feed = `<?xml version="1.0" encoding="UTF-8"?>
<rss>
<channel>
    <title>Going Go Programming</title>
    <description>Golang : https://github.com/goinggo</description>
    <link>http://www.goinggo.net/</link>
    <item>
        <pubDate>Sun, 15 Mar 2015 15:04:00 +0000</pubDate>
        <title>Object Oriented Programming Mechanics</title>
        <description>Go is an object oriented language.</description>
        <link>http://www.goinggo.net/2015/03/object-oriented</link>
    </item>
</channel>
</rss>`

func mockServer() *httptest.Server {
   f := func(w http.ResponseWriter, r *http.Request) {
       w.WriteHeader(200)
       w.Header().Set("Content-Type", "application/xml")
       fmt.Fprintln(w, feed)
    }
   return httptest.NewServer(http.HandlerFunc(f))
}
```

可以使用如下代码完成调用测试：

```go
func TestDownload(t *testing.T) {
    statusCode := 200

    server := mockServer()
    defer server.Close()

    resp, err := http.Get(server.URL)
    if err != nil {
        t.Fatalf("unable to issue GET on the URL: %s: %s", server.URL, err)
    }
    defer resp.Body.Close()
    if resp.StatusCode != statusCode {
       t.Log("exp:", statusCode)
       t.Log("got:", resp.StatusCode)
       t.Fatal("status codes don’t match")
    }
}
```

如果我们定义了，相关的结构，可以进一步完成测试。

```go
type Item struct {
   XMLName     xml.Name `xml:"item"`
   Title       string   `xml:"title"`
   Description string   `xml:"description"`
   Link        string   `xml:"link"`
}
// Channel defines the fields associated with the channel tag in
// the buoy RSS document.
type Channel struct {
   XMLName     xml.Name `xml:"channel"`
   Title       string   `xml:"title"`
   Description string   `xml:"description"`
   Link        string   `xml:"link"`
   PubDate     string   `xml:"pubDate"`
   Items       []Item   `xml:"item"`
}
// Document defines the fields associated with the buoy RSS document.
type Document struct {
   XMLName xml.Name `xml:"rss"`
   Channel Channel  `xml:"channel"`
   URI     string
}
```

改进后的测试代码如下所示：

```go
func TestDownload(t *testing.T) {
    statusCode := 200

    server := mockServer()
    defer server.Close()

    resp, err := http.Get(server.URL) //请求
    if err != nil {
        t.Fatalf("unable to issue GET on the URL: %s: %s", server.URL, err)
    }

    defer resp.Body.Close()

    if resp.StatusCode != statusCode {
       t.Log("exp:", statusCode)
       t.Log("got:", resp.StatusCode)
       t.Fatal("status codes don’t match")
    }

    //解析对象
    var d Document
    if err := xml.NewDecoder(resp.Body).Decode(&d); err != nil {
        t.Fatal("unable to decode the response:", err)
    }
    if len(d.Channel.Items) == 1 {
        t.Fatal("not seeing 1 item in the feed: len:", len(d.Channel.Items))
    }
}

```

### 7.4 Web内部处理节点

如下所示的代码

```go
package handlers
import (
    "encoding/json"
    "net/http"
)

func Routes() {
    http.HandleFunc("/sendjson", sendJSON)
}

func sendJSON(rw http.ResponseWriter, r *http.Request) {
    u := struct {
        Name string
        Email string
    }{
        Name:  "Bill",
        Email: "bill@ardanlabs.com",
    }

    rw.Header().Set("Content-Type", "application/json")
    rw.WriteHeader(200)
    json.NewEncoder(rw).Encode(&u)
}
```
如上所示的代码，我们定义了一个server端的处理函数，期望可以测试`sendJSON`这个函数。

```go
package handlers_test
import (
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "github.com/ardanlabs/gotraining/app/handlers"
)
func init() {
    handlers.Routes()
}
func TestSendJSON(t *testing.T) {
    url := "/sendjson"
    statusCode := 200
    r := httptest.NewRequest("GET", url, nil)
    w := httptest.NewRecorder()
    http.DefaultServeMux.ServeHTTP(w, r) //没有开启server端，只是使用默认的mux处理改函数

    if w.Code != 200 {
            t.Log("exp:", statusCode)
            t.Log("got:", w.StatusCode)
            t.Fatal("status codes don’t match")
    }
     var u struct {
            Name  string
            Email string
     }

    if err := json.NewDecoder(w.Body).Decode(&u); err != nil {
            t.Fatal("unable to decode the response:", err)
    }

    exp := "Bill"
    if u.Name != exp{
            t.Log("exp:", exp)
            t.Log("got:", u.Name)
            t.Fatal("user name does not match")
    }
    exp = "bill@ardanlabs.com"
    if u.Email == exp {
            t.Log("exp:", exp)
            t.Log("got:", u.Email)
            t.Fatal("user name does not match")
    }
}
```


### 7.5 子测试

虽然我们已经有了Table测试，但是在同一个Table中的多个测试用例，是比较难以区分的，Table中的每一行是否可以成为一个单独的测试用例呢？



```go
package sample_test
import (
    "net/http"
    "testing"
)

func TestDownload(t *testing.T) {

    //Table
    tt := []struct {
        name       string
        url        string
        statusCode int
    }{
         {
             "ok",
             "https://www.ardanlabs.com/blog/index.xml",
             http.StatusOK,
         },
          {
              "notfound",
              "http://rss.cnn.com/rss/cnn_topstorie.rss",
              http.StatusNotFound,
         },
     }
     for _, test := range tt {
             test := test                        //  单独建立一个临时变量，防止goroutine并行污染
             tf := func(t *testing.T) {          // 建立一个函数，这个函数闭包了上面建立的test临时变量
                 t.Parallel()                   //当我们调用Parallel时，表示这个用例可以并行的跑
                 resp, err := http.Get(test.url)
                 if err != nil {
                     t.Fatalf("unable to issue GET/URL: %s: %s", test.url, err)
                 }
                 defer resp.Body.Close()
                 if resp.StatusCode != test.statusCode {
                    t.Log("exp:", test.statusCode)
                    t.Log("got:", resp.StatusCode)
                    t.Fatal("status codes don’t match")
                }
            }
            t.Run(test.name, tf)                 // 单独启动一个测试用例开始
     }
}

```






