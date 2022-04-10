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

### 基本使用

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

```
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

```
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

```
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

```
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

```
/* === INVALID === */
var _ = Describe("book", func() {
  var book *Book
  Expect(book.Title()).To(BeFalse()) // NO!  Place in a setup node instead.

  It("tests something", func() {...})
})
```

#### 为了避免测试用例污染：不要在容器节点中初始化变量

```
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

```
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

```
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

```
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

```
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

```
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

```
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


