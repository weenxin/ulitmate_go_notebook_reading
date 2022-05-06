package service_test

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/weenxin/ulitmate_go_notebook_reading/ch7/model"
	"github.com/weenxin/ulitmate_go_notebook_reading/ch7/service"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var _ = ginkgo.Describe("manager to get book", func() {

	var manager *service.Manager
	var mock sqlmock.Sqlmock
	var err error

	ginkgo.BeforeEach(func() {
		var client *sql.DB
		client, mock, err = sqlmock.New()
		gomega.Expect(err).NotTo(gomega.HaveOccurred())

		ginkgo.DeferCleanup(func() {
			defer func(db *sql.DB) {
				_ = db.Close()
			}(client)
		})

		db, err := gorm.Open(mysql.New(mysql.Config{SkipInitializeWithVersion: true, Conn: client}), &gorm.Config{})
		gomega.Expect(err).NotTo(gomega.HaveOccurred())
		manager = service.NewManager(db)
	})

	ginkgo.Describe("save books to database", func() {
		var b *model.Book
		ginkgo.Context("model with id = 0", func() {
			ginkgo.BeforeEach(func() {
				b = &model.Book{
					Title:  "test save",
					Author: "test author",
					Pages:  100,
					Weight: 202,
				}
			})

			ginkgo.It("return no error & saved the model & model id should not be zero", func() {

				result := sqlmock.NewResult(1, 1)
				mock.ExpectBegin()
				mock.ExpectExec("^INSERT INTO `books`").
					WithArgs(sqlmock.AnyArg(),sqlmock.AnyArg(),sqlmock.AnyArg(),b.Title, b.Author, b.Pages, b.Weight).
					WillReturnResult(result)
				mock.ExpectCommit()

				err = manager.AddBook(b)
				gomega.Expect(err).To(gomega.BeNil())
				gomega.Expect(b.ID).To(gomega.Equal(uint(1)))
				gomega.Expect(mock.ExpectationsWereMet()).To(gomega.BeNil())

			})

		})
	})

	ginkgo.Describe("get books from database", func() {
		var b *model.Book

		ginkgo.Context("model exists", func() {
			ginkgo.BeforeEach(func() {
				b = &model.Book{
					Model: gorm.Model{
						ID:     100,
					},
					Title:  "test title",
					Author: "test author",
					Pages:  600,
					Weight: 400,
				}
			})
			ginkgo.It("return no error & return the model", func() {
				result := sqlmock.NewRows([]string{"id", "title", "author", "pages", "weight"}).
					AddRow(b.ID, b.Title, b.Author, b.Pages, b.Weight)
				mock.ExpectQuery("SELECT \\* FROM `books` WHERE `books`\\.`id` = \\? AND `books`\\.`deleted_at` IS NULL ORDER BY `books`\\.`id` LIMIT 1").
					WithArgs(b.ID).
					WillReturnRows(result)
				returnedBook, err := manager.GetBook(b.ID)
				gomega.Expect(err).To(gomega.BeNil())
				gomega.Expect(b.ID).To(gomega.Equal(returnedBook.ID))
				gomega.Expect(b.Title).To(gomega.Equal(returnedBook.Title))
				gomega.Expect(b.Author).To(gomega.Equal(returnedBook.Author))
				gomega.Expect(b.Pages).To(gomega.Equal(returnedBook.Pages))
				gomega.Expect(b.Weight).To(gomega.Equal(returnedBook.Weight))

				gomega.Expect(mock.ExpectationsWereMet()).To(gomega.BeNil())

			})
		})
	})

	ginkgo.Describe("delete books to database", func() {
		var b *model.Book
		ginkgo.Context("model exits ", func() {
			ginkgo.BeforeEach(func() {
				b = &model.Book{ //待删除的对象
					Model: gorm.Model{
						ID: 100,
					},
				}
			})
			ginkgo.It("return no error & delete the model", func() {
				//删除结果
				result := sqlmock.NewResult(0, 1)
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `books` SET `deleted_at`=\\? WHERE `books`.`id` = \\? AND `books`.`deleted_at` IS NULL").
					WithArgs(sqlmock.AnyArg(),b.ID).
					WillReturnResult(result)
				mock.ExpectCommit()
				err = manager.DeleteBook(b.ID)
				gomega.Expect(err).To(gomega.BeNil())

				gomega.Expect(mock.ExpectationsWereMet()).To(gomega.BeNil())

			})
		})
	})

	ginkgo.Describe("update books to database", func() {
		var b *model.Book
		ginkgo.Context("model exists ", func() {
			b = &model.Book{
				Model: gorm.Model{
					ID: 100,
				},
				Title:  "test title",
				Author: "test author",
				Pages:  100,
				Weight: 400,
			}
			ginkgo.It("return no error & updated the model", func() {
				result := sqlmock.NewResult(0, 1)
				mock.ExpectBegin()
				mock.ExpectExec("^UPDATE `books` (.+)WHERE `books`.`deleted_at` IS NULL AND `id` = \\?").
					WithArgs(sqlmock.AnyArg(),b.Title, b.Author, b.Pages, b.Weight, b.ID).
					WillReturnResult(result)
				mock.ExpectCommit()
				err = manager.UpdateBook(b)
				gomega.Expect(err).To(gomega.BeNil())

				gomega.Expect(mock.ExpectationsWereMet()).To(gomega.BeNil())
			})
		})
	})
})
