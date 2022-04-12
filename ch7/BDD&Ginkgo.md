## Ginkgo BDD测试

### 什么是BDD

BDD（Behavior driven development） 行为驱动测试，是对TDD（Test driven
development）的改进，BDD可以让开发聚焦在模型的行为上，有点像DDD聚焦在领域模型设计。软件开发的意义是生产有价值的产品，产品的特性是迭代的，用户需求也是变化的。所以如果程序模型的设计不能与产品模型，现实模型相对应，并跟随产品模型迭代，程序模型不断重构；非常容易形成所谓的"
屎山"。

### 如何开始BDD

- 制定模型行为规范
- 书写Given，When，Then 用例
- 按照规范开始写测试，写代码

### 规范

```yaml
given:
  - name: "a model"                                  # Given a model
    modules:
      - name: "get model catalog"                    # 对于book.Catalog函数测试
        behavior:
          - when: "pages count <= 300"              # 当pages count <= 300时
              - should: "be a short story"          # 应该是一个短故事

          - when: "pages count > 300"               # 当page count > 300时
              - should: "be a novel"                #应该是一个小说
```

制定如下所示的规范，可以完成如下代码书写

```go
package book

type Catalog int8

const (
	CategoryNovel      Catalog = iota
	CategoryShortStory Catalog = iota
)
const MaxShortStoryPages = 300

type Book struct {
	Title  string
	Author string
	Pages  int32
}

func (b Book) Catalog() Catalog {
	if b.Pages < MaxShortStoryPages {
		return CategoryShortStory
	} else {
		return CategoryNovel
	}
}
```

基于此可以完成如下测试用例的书写

```go
package book_test

var _ = ginkgo.Describe("Book", func() {

	//Get Catalog
	//基本用法
	ginkgo.Describe("get model catalog", func() {
		var foxInSocks, lesMis *book.Book
		ginkgo.BeforeEach(func() {
			lesMis = &book.Book{
				Title:  "Les Miserables",
				Author: "Victor Hugo",
				Pages:  2783,
			}
			foxInSocks = &book.Book{
				Title:  "Fox In Socks",
				Author: "Dr. Seuss",
				Pages:  24,
			}
		})

		ginkgo.Context("pages count <= 300", func() {
			ginkgo.It("be a short story", func() {
				gomega.Expect(foxInSocks.Catalog()).To(gomega.Equal(book.CategoryShortStory))
			})
		})
		ginkgo.Context("pages count > 300", func() {
			ginkgo.It("be a novel", func() {
				gomega.Expect(lesMis.Catalog()).To(gomega.Equal(book.CategoryNovel))
			})
		})
	})
}
```

### GinkGo基本使用

[官方文档](https://onsi.github.io/ginkgo/)

***安装***

```bash
$ go install -mod=mod github.com/onsi/ginkgo/v2/ginkgo
$ go get github.com/onsi/gomega/...
```

***产生test suite***

```bash
$ ginkgo bootstrap
Generating ginkgo test suite bootstrap for books in:
books_suite_test.go
```

***产生Test文件***

```bash
$ ginkgo generate model
Generating ginkgo test for Book in:
  book_test.go
```

***编写逻辑代码***

***写单元测试***

### Ginkgo 原理

- `Describe` ， `Context` , `When` 是语义等同的，称他们为构造语义，用来构造逻辑树，称这些节点为容器节点;
- `It` 是执行语义，称这些节点为运行节点；
- `BeforeEach` 会出现在构造语义中，完成运行代码的注册。其代码会在每个`It`执行前执行。

```go
var _ = Describe("Books", func() {
  var book *books.Book

  BeforeEach(func() {
    //Closure A
    book = &books.Book{
      Title: "Les Miserables",
      Author: "Victor Hugo",
      Pages: 2783,
    }
    Expect(book.IsValid()).To(BeTrue())
  })

  Describe("Extracting names", func() {
    When("author has both names", func() {
      It("extracts the last name", func() {
        //Closure B
        Expect(book.AuthorLastName()).To(Equal("Hugo"))
      })

      It("extracts the first name", func() {
        //Closure C
        Expect(book.AuthorFirstName()).To(Equal("Victor"))
      })
    })

    When("author has one name", func() {
      BeforeEach(func() {
        //Closure D
        book.Author = "Hugo"
      })

      It("extracts the last name", func() {
        //Closure E
        Expect(book.AuthorLastName()).To(Equal("Hugo"))
      })

      It("returns empty first name", func() {
        //Closure F
        Expect(book.AuthorFirstName()).To(BeZero())
      })
    })

  })
})
```

会分析成为如下所示的树行结构：

```go
Describe: "Books"
    |_BeforeEach: <Closure-A>
    |_Describe: "Extracting names"
        |_When: "author has both names"
            |_It: "extracts the last name", <Closure-B>
            |_It: "extracts the first name", <Closure-C>
    |_When: "author has one name"
        |_BeforeEach: <Closure-D>
            |_It: "extracts the last name", <Closure-E>
            |_It: "returns empty first name", <Closure-F>
```

根据如上所示的结构，会分析出如下几条测试case

```go
[
  {
    Texts: ["Books", "Extracting names", "author has both names", "extracts the last name"],
    Closures: <BeforeEach-Closure-A>, <It-Closure-B>
  },
  {
    Texts: ["Books", "Extracting names", "author has both names", "extracts the first name"],
    Closures: <BeforeEach-Closure-A>, <It-Closure-C>
  },
  {
    Texts: ["Books", "Extracting names", "author has one name", "extracts the last name"],
    Closures: <BeforeEach-Closure-A>, <BeforeEach-Closure-D>, <It-Closure-E>
  },
  {
    Texts: ["Books", "Extracting names", "author has one name", "returns empty first name"],
    Closures: <BeforeEach-Closure-A>, <BeforeEach-Closure-D>, <It-Closure-F>
  }
]
```

所以ginkgo会分为两个阶段

- 阶段1，梳理树形结构，此时不会执行任何测试代码，只会运行 `Describe` ， `Context` , `When`这些容器节点的代码
- 阶段2，运行测试用例，运行 `BeforeEach`, `JustBeforeEach`, `It`, `AfterEach`, `JustAfterEach`中的代码

理解了如上所说的流程，就可以很好的理解一些需要避免的点了。

#### 执行语义中不应该再出现容器节点

也就是说`BeforeEach`, `JustBeforeEach`, `It`, `AfterEach`, `JustAfterEach`这些节点代码中不应该出现`Describe` ， `Context` , `When`类的子节点。

例如如下所示的代码

```go
/* === INVALID === */
var _ = It("has a color", func() {
  Context("when blue", func() { // NO! Nodes can only be nested in containers
    It("is blue", func() { // NO! Nodes can only be nested in containers

    })
  })
})
```

#### 不要在容器节点中做任何Assert

因为容器节点内的代码只在测试树建立过程中运行，此时还没有运行真正的测试用例。

例如如下所示的代码 ：

```go
/* === INVALID === */
var _ = Describe("book", func() {
  var book *Book
  Expect(book.Title()).To(BeFalse()) // NO!  Place in a setup node instead.

  It("tests something", func() {...})
})
```

#### 为了避免测试用例污染：不要在容器节点中初始化变量

```go
/* === INVALID === */
var _ = Describe("book", func() {
  book := &books.Book{ // No!
    Title:  "Les Miserables",
    Author: "Victor Hugo",
    Pages:  2783,
  }

  It("is invalid with no author", func() {
    book.Author = "" // bam! we've changed the closure variable and it will never be reset.
    Expect(book.IsValid()).To(BeFalse())
  })

  It("is valid with an author", func() {
    Expect(book.IsValid()).To(BeTrue()) // this will fail if it runs after the previous test
  })
})
```

多个测试用例运行顺序是随机的，所以，如上所示的代码，第一个测试用例会影响到第二个测试用例的执行结果。应该将这部分代码放到`BeforeEach`中

```go
var _ = Describe("book", func() {
  var book *books.Book // declare in container nodes

  BeforeEach(func() {
    book = &books.Book {  //initialize in setup nodes
      Title:  "Les Miserables",
      Author: "Victor Hugo",
      Pages:  2783,
    }
  })

  It("is invalid with no author", func() {
    book.Author = ""
    Expect(book.IsValid()).To(BeFalse())
  })

  It("is valid with an author", func() {
    Expect(book.IsValid()).To(BeTrue())
  })
})
```

#### JustBeforeEach 分离配置与对象创建

`JustBeforeEach`的语句的运行顺序

- 会在每个所在容器节点的所有子测试用例都运行
- `BeforeEach`之后运行
- 并且在`It`之前运行

考虑到如下代码：

```go
Describe("some JSON decoding edge cases", func() {
  var book *books.Book
  var err error

  When("the JSON fails to parse", func() {
    BeforeEach(func() {
      book, err = NewBookFromJSON(`{  //创建和配置都在一起，多个测试用例都需要调用
        "title":"Les Miserables",
        "author":"Victor Hugo",
        "pages":2783oops
      }`)
    })

    It("returns a nil book", func() {
      Expect(book).To(BeNil())
    })

    It("errors", func() {
      Expect(err).To(MatchError(books.ErrInvalidJSON))
    })
  })

  When("the JSON is incomplete", func() {
    BeforeEach(func() {
      book, err = NewBookFromJSON(`{ //创建和配置都在一起，多个测试用例都需要调用
        "title":"Les Miserables",
        "author":"Victor Hugo",
      }`)
    })

    It("returns a nil book", func() {
      Expect(book).To(BeNil())
    })

    It("errors", func() {
      Expect(err).To(MatchError(books.ErrIncompleteJSON))
    })
  })
})
```

可以更改为如下代码：

```go
Describe("some JSON decoding edge cases", func() {
  var book *books.Book
  var err error
  var json string
  JustBeforeEach(func() {  //创建流程
    book, err = NewBookFromJSON(json)
    Expect(book).To(BeNil())
  })

  When("the JSON fails to parse", func() {
    BeforeEach(func() { //更改配置
      json = `{
        "title":"Les Miserables",
        "author":"Victor Hugo",
        "pages":2783oops
      }`
    })

    It("errors", func() {
      Expect(err).To(MatchError(books.ErrInvalidJSON))
    })
  })

  When("the JSON is incomplete", func() {
    BeforeEach(func() {   //更改配置
      json = `{
        "title":"Les Miserables",
        "author":"Victor Hugo",
      }`
    })

    It("errors", func() {
      Expect(err).To(MatchError(books.ErrIncompleteJSON))
    })
  })
})
```

#### AfterEach DeferCleanup 清理测试环境

如下所示的代码：

```go
Describe("Reporting book weight", func() {
  var book *books.Book

  BeforeEach(func() {
    book = &books.Book{
      Title: "Les Miserables",
      Author: "Victor Hugo",
      Pages: 2783,
      Weight: 500,
    }
  })

  Context("with no WEIGHT_UNITS environment set", func() {
    BeforeEach(func() {
      err := os.Clearenv("WEIGHT_UNITS")
      Expect(err).NotTo(HaveOccurred())
    })

    It("reports the weight in grams", func() {
      Expect(book.HumanReadableWeight()).To(Equal("500g"))
    })
  })

  Context("when WEIGHT_UNITS is set to oz", func() {
    BeforeEach(func() {
      err := os.Setenv("WEIGHT_UNITS", "oz")
      Expect(err).NotTo(HaveOccurred())
    })

    It("reports the weight in ounces", func() {
      Expect(book.HumanReadableWeight()).To(Equal("17.6oz"))
    })
  })

  Context("when WEIGHT_UNITS is invalid", func() {
    BeforeEach(func() {
      err := os.Setenv("WEIGHT_UNITS", "smoots")
      Expect(err).NotTo(HaveOccurred())
    })

    It("errors", func() {
      weight, err := book.HumanReadableWeight()
      Expect(weight).To(BeZero())
      Expect(err).To(HaveOccurred())
    })
  })
})

```

如上所示的测试用例可以正常运行，但是多次测试之间环境变量被更改，虽然测试用例正常，确实有一些环境变量的污染，同时在结束时，也没有恢复正常。可以通过如下方式更改：

```go
Describe("Reporting book weight", func() {
  var book *books.Book
  var originalWeightUnits string

  BeforeEach(func() {
    book = &books.Book{
      Title: "Les Miserables",
      Author: "Victor Hugo",
      Pages: 2783,
      Weight: 500,
    }
    originalWeightUnits = os.Getenv("WEIGHT_UNITS")
  })

  AfterEach(func() {
    err := os.Setenv("WEIGHT_UNITS", originalWeightUnits)
    Expect(err).NotTo(HaveOccurred())
  })
  ...
})
```

虽然这样可以解决问题，但是环境改变和reset在两个不同的函数中，结构上有些难以追踪。

```go
Describe("Reporting book weight", func() {
  var book *books.Book

  BeforeEach(func() {
    ...
    originalWeightUnits := os.Getenv("WEIGHT_UNITS")
    DeferCleanup(func() {
      err := os.Setenv("WEIGHT_UNITS", originalWeightUnits)
      Expect(err).NotTo(HaveOccurred())
    })
  })
  ...
})
```

可以将清理代码，使用`DeferCleanup`整理到`BeforeEach`中。`DeferCleanup`不被ginkgo认为是运行节点，所以不会之前的实践有冲突。

#### BeforeSuite , AfterSuite suite前的清理和创建

`BeforeSuite`和`AfterSuite`一般用在测试前的环境搭建合清理。比如project依赖与mysql，运行测试用例前应该先搭建mysql。

基于`BDD`的设计，我们将系统设计成为由多个模块拼接而成，每个模块都做好自己职责范围之内的事情。 `单元测试`毕竟能力是有限的，拼接过程也可能有代码异常。此时，可以引入`e2e`测试。在运行`e2e`
测试前需要将数据库表等相关依赖都创建出来。测试结束需要停止测试应用，并删除测试用的数据库表。

如下所示的代码：

```go
var _ = ginkgo.BeforeSuite(func() {

	ginkgo.By("initialzing tables") //创建表
	rootPassoword := os.Getenv("ROOT_DATABASE_PWD")
	gomega.Expect(rootPassoword).NotTo(gomega.BeEmpty())

	cmd := exec.Command("mysql", "-uroot", "-h127.0.0.1", fmt.Sprintf("-p%s", rootPassoword))
	cmd.Stdin = bytes.NewBuffer([]byte(`create database if not exists test; use test; CREATE TABLE IF NOT EXISTS books ( id INTEGER PRIMARY KEY AUTO_INCREMENT,     title varchar(255) NOT NULL,     author varchar(64) NOT NULL,     Pages int(10) not null,     weight int(10) not null );`))
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	ginkgo.By("initializing database")
	dsn := os.Getenv("DATABASE_DSN")
	gomega.Expect(dsn).NotTo(gomega.BeEmpty())

	ginkgo.By("start server") //在启动测试前，我们编译好最新版本的应用放在指定位置，此时我们运行应用
	go func() {
		server = exec.Command("./ch7-e2e-test", fmt.Sprintf("--dsn=%s", dsn), fmt.Sprintf("--address=%s", address))
		server.Stderr = os.Stderr
		server.Stdout = os.Stdout
		defer ginkgo.GinkgoRecover() //程序异常需要捕获
		err := server.Start()
		if err != nil {
			ginkgo.Fail(fmt.Sprintf("start server failed :%s", err))
		}

	}()
	//wait for server start
	time.Sleep(1 * time.Second)

})
```

```go
var _ = ginkgo.AfterSuite(func() {
	ginkgo.By("stop server") //停止运行最新版应用
	if err := server.Process.Kill(); err != nil {
		fmt.Println("Server Shutdown:", err)
	}

	ginkgo.By("clear database") //清理数据库
	rootPassoword := os.Getenv("ROOT_DATABASE_PWD")
	cmd := exec.Command("mysql", "-uroot", "-h127.0.0.1", fmt.Sprintf("-p%s", rootPassoword))
	cmd.Stdin = bytes.NewBuffer([]byte("drop database if exists test;"))
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
})
```

如上在全部测试用例开始前，

- 先初始化了应用所需要的环境，并且启动应用；
- 启动测试用例；
- 关闭应用程序，清理数据库环境；

测试用例如下所示：

```go
var _ = ginkgo.Describe("Api", func() { //测试对外API
	ginkgo.Describe("Normal use", func() { //当正常的使用用例
		var b *model.Book
		ginkgo.BeforeEach(func() {
			b = &model.Book{
				Id:     0,
				Title:  "test title",
				Author: "test author",
				Pages:  100,
				Weight: 100,
			}
		})
		ginkgo.It("smoking test", func() { //创建并获取
			ginkgo.By("creating book")
			urlCrate := fmt.Sprintf("http://%s/books/", address)
			content, err := json.Marshal(b)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			response, err := http.Post(urlCrate, "application/json", bytes.NewBuffer(content))
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer func(Body io.ReadCloser) {
				_ = Body.Close()
			}(response.Body)
			gomega.Expect(response.StatusCode).To(gomega.Equal(http.StatusOK))
			content, err = ioutil.ReadAll(response.Body)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			var bookInserted = struct {
				Status  string
				Message string
				Data    model.Book
			}{}
			err = json.Unmarshal(content, &bookInserted)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(bookInserted.Data.Id).NotTo(gomega.Equal(0))
			gomega.Expect(b.Title).To(gomega.Equal(bookInserted.Data.Title))
			gomega.Expect(b.Author).To(gomega.Equal(bookInserted.Data.Author))
			gomega.Expect(b.Pages).To(gomega.Equal(bookInserted.Data.Pages))
			gomega.Expect(b.Weight).To(gomega.Equal(bookInserted.Data.Weight))

			ginkgo.By("get book")
			url := fmt.Sprintf("http://%s/books/%d", address, bookInserted.Data.Id)
			resp, err := http.Get(url)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			defer func(Body io.ReadCloser) {
				err := Body.Close()
				if err != nil {

				}
			}(resp.Body)
			content, err = ioutil.ReadAll(resp.Body)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			var data = struct {
				Status  string
				Message string
				Data    model.Book
			}{}
			err = json.Unmarshal(content, &data)
			gomega.Expect(err).NotTo(gomega.HaveOccurred())
			gomega.Expect(bookInserted.Data.Id).To(gomega.Equal(data.Data.Id))
			gomega.Expect(b.Title).To(gomega.Equal(data.Data.Title))
			gomega.Expect(b.Author).To(gomega.Equal(data.Data.Author))
			gomega.Expect(b.Pages).To(gomega.Equal(data.Data.Pages))
			gomega.Expect(b.Weight).To(gomega.Equal(data.Data.Weight))
		})
	})
})
```

#### Ginkgo如何处理错误

- 用户可以使用`gomega`的`expect`来指定返回结果应该与期望一直，当然用户也可以直接调用`ginkgo.Fail`来直接结束用例；
- 在`BeforeEach`, `JustBeforeEach`, `It` 这些运行节点中调用失败函数时，会终止本用例的执行，如果期望只要有一个用例失败就停止所有，可以启用`ginkgo --fail-fast`这个flag
- 在其他函数（`BeforeSuite`,`AfterSuite`）或是单独启用一个goRoutine中调用错误函数时，会导致ginkgo panic，应该按照如下方式捕获异常；

```go
It("panics in a goroutine", func() {
  var c chan interface{}
  go func() {
    defer GinkgoRecover()
    Fail("boom")
    close(c)
  }()
  <-c
})
```

#### 日志输出

当有部分重要信息debug信息需要输出时，可以使用ginkgo的输出函数，默认情况下只有当测试用例失败时这些信息才会输出到stdout中，可以使用`ginkgo -v`将所有错误信息都输出。

- ginkgo.GinkgoWriter.Print(a ...interface{})
- ginkgo.GinkgoWriter.Println(a ...interface{})
- ginkgo.GinkgoWriter.Printf(format string, a ...interface{})

#### 通过By来书写复杂用例文档

```go
var _ = Describe("Browsing the library", func() {
  BeforeEach(func() {
    By("Fetching a token and logging in")

    authToken, err := authClient.GetToken("gopher", "literati")
    Expect(err).NotTo(HaveOccurred())

    Expect(libraryClient.Login(authToken)).To(Succeed())
  })

  It("should be a pleasant experience", func() {
    By("Entering an aisle")
    aisle, err := libraryClient.EnterAisle()
    Expect(err).NotTo(HaveOccurred())

    By("Browsing for books")
    books, err := aisle.GetBooks()
    Expect(err).NotTo(HaveOccurred())
    Expect(books).To(HaveLen(7))

    By("Finding a particular book")
    book, err := books.FindByTitle("Les Miserables")
    Expect(err).NotTo(HaveOccurred())
    Expect(book.Title).To(Equal("Les Miserables"))

    By("Checking a book out")
    Expect(libraryClient.CheckOut(book)).To(Succeed())
    books, err = aisle.GetBooks()
    Expect(err).NotTo(HaveOccurred())
    Expect(books).To(HaveLen(6))
    Expect(books).NotTo(ContainElement(book))
  })
})
```

如上所示的用例，如果有错误出现，会按照by出现的顺序依次输出，更容易定位问题出现在哪个环节。


#### TableTesting

TableTesting已经被很多项目所采用，ginkgo也支持Table Specs的。

```go
DescribeTable("Extracting the author's first and last name",
  func(author string, isValid bool, firstName string, lastName string) {
    book := &books.Book{
      Title: "My Book"
      Author: author,
      Pages: 10,
    }
    Expect(book.IsValid()).To(Equal(isValid))
    Expect(book.AuthorFirstName()).To(Equal(firstName))
    Expect(book.AuthorLastName()).To(Equal(lastName))
  },
  Entry("When author has both names", "Victor Hugo", true, "Victor", "Hugo"),
  Entry("When author has one name", "Hugo", true, "", "Hugo"),
  Entry("When author has a middle name", "Victor Marie Hugo", true, "Victor", "Hugo"),
  Entry("When author has no name", "", false, "", ""),
)
```
需要注意的是，TableTesting的Entry函数会在测试树构建过程调用，因此

```go
/* === INVALID === */
Describe("book", func() {
  var shelf map[string]*books.Book //Shelf is declared here

  BeforeEach(func() {
    shelf = map[string]*books.Book{ //...and initialized here
      "Les Miserables": &books.Book{Title: "Les Miserables", Author: "Victor Hugo", Pages: 2783},
      "Fox In Socks": &books.Book{Title: "Fox In Socks", Author: "Dr. Seuss", Pages: 24},
    }
  })

  DescribeTable("Categorizing books",
    func(book *books.Book, category books.Category) {
      Expect(book.Category()).To(Equal(category))
    },
    Entry("Novels", shelf["Les Miserables"], books.CategoryNovel),
    Entry("Short story", shelf["Fox in Socks"], books.CategoryShortStory),
  )
})
```
会发生panic，因为`shelf`会在`BeforeEach`也就是测试用例运行过程调用，所以在`Entry`访问就会panic。可以按照如下实现：

```go
Describe("book", func() {
  var shelf map[string]*books.Book //Shelf is declared here

  BeforeEach(func() {
    shelf = map[string]*books.Book{ //...and initialized here
      "Les Miserables": &books.Book{Title: "Les Miserables", Author: "Victor Hugo", Pages: 2783},
      "Fox In Socks": &books.Book{Title: "Fox In Socks", Author: "Dr. Seuss", Pages: 24},
    }
  })

  DescribeTable("Categorizing books",
    func(key string, category books.Category) {
      Expect(shelf[key]).To(Equal(category))
    },
    Entry("Novels", "Les Miserables", books.CategoryNovel),
    Entry("Novels", "Fox in Socks", books.CategoryShortStory),
  )
})
```

可以指定entry的测试用例描述：

```go
var _ = Describe("Math", func() {
  DescribeTable("addition",
    func(a, b, c int) {
      Expect(a+b).To(Equal(c))
    },
    Entry("1+2=3", 1, 2, 3),
    Entry("-1+2=1", -1, 2, 1),
    Entry("0+0=0", 0, 0, 0),
    Entry("10+100=101", 10, 100, 110), //OOPS TYPO
  )
})
```
方式比较机械，有typo的可能。

```go
var _ = Describe("Math", func() {
  DescribeTable("addition",
    func(a, b, c int) {
      Expect(a+b).To(Equal(c))
    },
    Entry(nil, 1, 2, 3),
    Entry(nil, -1, 2, 1),
    Entry(nil, 0, 0, 0),
    Entry(nil, 10, 100, 110),
  )
})
```
此时会输出 `Entry: 1, 2, 3`, `Entry: -1, 2, 1`, `Entry: 0, 0, 0`, `Entry: 10, 100, 110`。这样似乎也不是很好。

```go
var _ = Describe("Math", func() {
  DescribeTable("addition",
    func(a, b, c int) {
      Expect(a+b).To(Equal(c))
    },
    func(a, b, c int) string {
      return fmt.Sprintf("%d + %d = %d", a, b, c)
    }
    Entry(nil, 1, 2, 3),
    Entry(nil, -1, 2, 1),
    Entry(nil, 0, 0, 0),
    Entry(nil, 10, 100, 110),
  )
})
```

也有比较好的使用方式：

```go
ar _ = Describe("Math", func() {
  DescribeTable("addition",
    func(a, b, c int) {
      Expect(a+b).To(Equal(c))
    },
    EntryDescription("%d + %d = %d")
    Entry(nil, 1, 2, 3),
    Entry(nil, -1, 2, 1),
    Entry("zeros", 0, 0, 0),
    Entry(EntryDescription("%[3]d = %[1]d + %[2]d"), 10, 100, 110)
    Entry(func(a, b, c int) string {fmt.Sprintf("%d = %d", a + b, c)}, 4, 3, 7)
  )
})
```


### Ginkgo 运行测试用例

- Ginkgo假设所有测试用例之间都是独立的
- Ginkgo在运行测试前会对测试用例运行顺序随机处理；

#### 并行运行测试用例

可以通过`-p`参数指定测试用例并行数量；如`ginkgo -procs=N`使用方式。如下所示的代码：

```go
Describe("Bookmark", func() {
  var book *books.Book

  BeforeEach(func() {
    book = &books.Book{
      Title:  "Les Miserables",
      Author: "Victor Hugo",
      Pages:  2783,
    }
  })

  It("has no bookmarks by default", func() {
    Expect(book.Bookmarks()).To(BeEmpty())
  })

  It("can add bookmarks", func() {
    book.AddBookmark(173)
    Expect(book.Bookmarks()).To(ContainElement(173))
  })
})
```
并行运行两个测试用例，"似乎"会有并发问题，因为两个测试用例都在访问book字段，但是在`ginkgo`是没有这个问题的。
ginkgo 怎么运行测试用例的呢 ？
- 通过 ginkgo 先会调用`go test -c`编译出一个用于测试的可执行文件
- ginkgo调度多个process运行测试用例

可以通过`GinkgoParallelProcess()`获得当前用例所在的进程id。

```go
Describe("Storing books in an external database", func() {
  BeforeEach(func() {
    namespace := fmt.Sprintf("namespace-%d", GinkgoParallelProcess())
    Expect(dbClient.SetNamespace(namespace)).To(Succeed())
    DeferCleanup(dbClient.ClearNamespace, namespace)
  })

  It("returns empty when there are no books", func() {
    Expect(dbClient.Books()).To(BeEmpty())
  })

  Context("when a book is in the database", func() {
    var book *books.Book
    BeforeEach(func() {
      book = &books.Book{
        Title:  "Les Miserables",
        Author: "Victor Hugo",
        Pages:  2783,
      }
      Expect(dbClient.Store(book)).To(Succeed())
    })

    It("can fetch the book", func() {
      Expect(dbClient.Books()).To(ConsistOf(book))
    })

    It("can update the book", func() {
      book.Author = "Victor Marie Hugo"
      Expect(dbClient.Store(book)).To(Succeed())
      Expect(dbClient.Books()).To(ConsistOf(book))
    })

    It("can delete the book", func() {
      Expect(dbClient.Delete(book)).To(Succeed())
      Expect(dbClient.Books()).To(BeEmpty())
    })
  })
})
```

你也可以通过如下代码获得并行数。

```go
suiteConfig, _ := GinkgoConfiguration()
totalProcesses := suiteConfig.ParallelTotal
```


#### SynchronizedBeforeSuite 和 SynchronizedAfterSuite
为了支持并行：
- `BeforeSuite`前会运行`SynchronizedBeforeSuite`
- `AfterSuite`之后运行`SynchronizedAfterSuite`

并且`SynchronizedBeforeSuite`和`SynchronizedAfterSuite`只会被执行一次


#### 用例装饰

- `Serial` 保证能够按照代码先后顺序执行，但是并不一定在同一个Process执行；
- `Ordered` 保证顺序的同时，保证在统一额procsss执行；

```go
Describe("Something expensive", Serial, func() {
  It("is a resource hog that can't run in parallel", func() {
    ...
  })

  It("is another resource hog that can't run in parallel", func() {
    ...
  })
})
```

```go
Describe("checking out a book", Ordered, func() {
  var book *books.Book
  var err error

  It("can fetch a book from a library", func() {
    book, err = libraryClient.FetchByTitle("Les Miserables")
    Expect(err).NotTo(HaveOccurred())
    Expect(book.Title).To(Equal("Les Miserables"))
  })

  It("can check out the book", func() {
    Expect(library.CheckOut(book)).To(Succeed())
  })

  It("no longer has the book in stock", func() {
    book, err = libraryClient.FetchByTitle("Les Miserables")
    Expect(err).To(MatchError(books.NOT_IN_STOCK))
    Expect(book).To(BeNil())
  })
})
```


#### `BeforeAll`和`AfterAll`

`BeforeAll`和`AfterAll`可以在`Ordered`中使用。
- `BeforeAll`保证在所有用例前指定一次；
- `AfterAll`保证在所有用例后执行一次；


#### 筛选用例

- Pending 不会运行的用例
- Skip 会跳过用例
  - It 中跳过本次用例
  - BeforeEach 跳过作用的所有用例
  - BeforeAll 跳过`Ordered`跳过作用的所有用例
  - BeforeSuite 跳过测试`suite`
- Focus 用例中有Focus后，只会运行Focus的用例，其他用例都不跑

子用例没有Focus的，会执行所有子用例
```go
FDescribe("some specs you're debugging", func() {
  It("might be failing", func() { ... })//运行这个用例
  It("might also be failing", func() { ... })//运行这个用例
})
```

子用例有focused，就只跑Foucused的用例
```go
FDescribe("some specs you're debugging", func() {
  It("might be failing", func() { ... })
  FIt("might also be failing", func() { ... }) //只运行这个用例
})
```

#### 给测试用例增加Label

```go
Describe("Storing books", Label("integration", "storage"), func() {
  It("can save entire shelves of books to the central library", Label("network", "slow", "library storage"), func() {
    // has labels [integration, storage, network, slow, library storage]
  })

  It("cannot delete books from the central library", Label("network", "library storage"), func() {
    // has labels [integration, storage, network, library storage]
  })

  It("can check if a book is stored in the central library", Label("network", "slow", "library query"), func() {
    // has labels [integration, storage, network, slow, library query]
  })

  It("can save books locally", Label("local"), func() {
    // has labels [integration, storage, local]
  })

  It("can delete books locally", Label("local"), func() {
    // has labels [integration, storage, local]
  })
})
```

每个用例的Label是所有子用例的集合。

- `ginkgo --label-filter="integration"` 运行所有有`integration`的用例
- `ginkgo --label-filter="!slow"` 运行所有没有`slow`的用例
- `ginkgo --label-filter="network && !slow"` 运行所有会有network并且没有`slow`的用例
- `ginkgo --label-filter=/library/` 运行所有匹配`library`模式的用例


#### 按照文件filter

- `ginkgo --focus-file=foo` 匹配foo_test.go 和 /foo/bar_test.go
- `ginkgo --skip-file=foo` 匹配foo_test.go 和 /foo/bar_test.go
- `FILE_REGEX:LINE1-LINE2` 可以指定文件位置

