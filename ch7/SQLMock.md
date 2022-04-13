## mock SQL

很多时候我们系统依赖于各种MySQL等DB，此时如果想测试系统就需要
- 每个想要跑测试用例的同学都有一套可以完全运行应用的测试环境；
- DB连接信息等维护，放在仓库容易导致信息泄漏，人工维护又容易疏漏；

[go-sqlmock](https://github.com/DATA-DOG/go-sqlmock) 可以在不依赖DB环境的情况下测试相关代码。

主要的测试思想是：
- 执行指定语句后，你期望向后端发送一个什么样子的sql请求；
- 填写（mock）期望获得什么结果；
- 唯一索引等，通过mock结果实现；

当前大部分golang应用开发都不基于裸`database/sql`直接开发，大部分都会使用orm包，如下是一个使用go-orm与sqlmock配合的测试样例。

```go
    var manager *service.Manager
	var mock sqlmock.Sqlmock //全局变量方便访问
	var err error

	ginkgo.BeforeEach(func() { //每个测试用例运行之前都先运行这个代码
		var client *sql.DB
		client, mock, err = sqlmock.New() //每次都新建出一个mock
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		ginkgo.DeferCleanup(func() {
			defer func(db *sql.DB) {
				_ = db.Close() //结束的时候释放连接
			}(client)
		})

        //gorm与sql-mock联合使用
		db, err := gorm.Open(mysql.New(mysql.Config{SkipInitializeWithVersion: true, Conn: client}), &gorm.Config{})
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		//被测试对象使用的是sqlmock的对象
		manager = service.NewManager(db)
	})
```


### 增加对象

```go
ginkgo.Describe("save books to database", func() { //测试存储数据到db
		var b *model.Book
		ginkgo.Context("model with id = 0", func() { //当book的id为0
			ginkgo.BeforeEach(func() { //每个测试前都新建一个book
				b = &model.Book{
					Id:     0,
					Title:  "test save",
					Author: "test author",
					Pages:  100,
					Weight: 202,
				}
			})

			ginkgo.It("return no error & saved the model & model id should not be zero", func() {
			    //新建sql返回结果
				result := sqlmock.NewResult(1, 1)//返回影响1列，lastInsertedId为1
				//应该先开始一个事物
				mock.ExpectBegin()
				//接着执行insert语句
				mock.ExpectExec("^INSERT INTO `books`").
					WithArgs(b.Title, b.Author, b.Pages, b.Weight).
					WillReturnResult(result)
				//最后提交事务
				mock.ExpectCommit()

                //开始测试
				err = manager.AddBook(b)
				//期望结果应该是预期的
				gomega.Expect(err).To(gomega.BeNil())
				gomega.Expect(b.Id).To(gomega.Equal(int64(1)))
				//执行的sql应该与我们定义的匹配；
				gomega.Expect(mock.ExpectationsWereMet()).To(gomega.BeNil())
			})

		})
	})
```

### 查询数据

```go
ginkgo.Describe("get books from database", func() {
		var b *model.Book

		ginkgo.Context("model exists", func() { //待查询对象
			ginkgo.BeforeEach(func() {
				b = &model.Book{
					Id:     100,
					Title:  "test title",
					Author: "test author",
					Pages:  600,
					Weight: 400,
				}
			})
			ginkgo.It("return no error & return the model", func() {
			    //mock返回结果，返回一行，数据列和值
				result := sqlmock.NewRows([]string{"id", "title", "author", "pages", "weight"}).
					AddRow(b.Id, b.Title, b.Author, b.Pages, b.Weight)

				//应该下一个sql如下
				mock.ExpectQuery("SELECT (.+) FROM `books` WHERE `books`.`id` = ?").
					WithArgs(b.Id).
					WillReturnRows(result)

			    //开始测试对象
				returnedBook, err := manager.GetBook(b.Id)
				//结果应该与预期相符合
				gomega.Expect(err).To(gomega.BeNil())
				gomega.Expect(b.Id).To(gomega.Equal(returnedBook.Id))
				gomega.Expect(b.Title).To(gomega.Equal(returnedBook.Title))
				gomega.Expect(b.Author).To(gomega.Equal(returnedBook.Author))
				gomega.Expect(b.Pages).To(gomega.Equal(returnedBook.Pages))
				gomega.Expect(b.Weight).To(gomega.Equal(returnedBook.Weight))

                //下发的sql应该与预期符合
				gomega.Expect(mock.ExpectationsWereMet()).To(gomega.BeNil())

			})
		})
	})

```

### 删除数据

```go
ginkgo.Describe("delete books to database", func() {
		var b *model.Book
		ginkgo.Context("model exits ", func() {
			ginkgo.BeforeEach(func() {
				b = &model.Book{//待删除对象
					Id:     100,
				}
			})
			ginkgo.It("return no error & delete the model", func() {
			    //结果，影响函数1，lastInsertedId = 0
				result := sqlmock.NewResult(0, 1)
				//应该先开始一个事务
				mock.ExpectBegin()
				//删除数据,并且以id为参数
				mock.ExpectExec("^DELETE FROM `books` WHERE `books`.`id` = ?").
					WithArgs(b.Id).
					WillReturnResult(result) //返回目标结果
				//应该会提交事务
				mock.ExpectCommit()

				//开始测试
				err = manager.DeleteBook(b.Id)
				//期望没有错误发生
				gomega.Expect(err).To(gomega.BeNil())
				//期望sql与预期相符合
				gomega.Expect(mock.ExpectationsWereMet()).To(gomega.BeNil())

			})
		})
	})
```

```go
ginkgo.Describe("update books to database", func() {
		var b *model.Book
		ginkgo.Context("model exists ", func() {
			b = &model.Book{ //待更新对象
				Id:     100,
				Title:  "test title",
				Author: "test author",
				Pages:  100,
				Weight: 400,
			}
			ginkgo.It("return no error & updated the model", func() {
			    //mock结果，lastInsertedId 为0 ，影响1行
				result := sqlmock.NewResult(0, 1)
				//应该先开始一个事务
				mock.ExpectBegin()
				//然后下发sql，并且按照表达式下发参数
				mock.ExpectExec("^UPDATE `books` (.+)WHERE `id` = ?").
					WithArgs(b.Title, b.Author, b.Pages, b.Weight, b.Id).
					WillReturnResult(result) //返回目标结果
				//应该提交事务
				mock.ExpectCommit()

				//开始测试
				err = manager.UpdateBook(b)

				//不应该有错误发生
				gomega.Expect(err).To(gomega.BeNil())
				//sql与预期相符合
				gomega.Expect(mock.ExpectationsWereMet()).To(gomega.BeNil())

			})
		})
	})

```

