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
					Id:     0,
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
					WithArgs(b.Title, b.Author, b.Pages, b.Weight).
					WillReturnResult(result)
				mock.ExpectCommit()

				err = manager.AddBook(b)
				gomega.Expect(err).To(gomega.BeNil())
				gomega.Expect(b.Id).To(gomega.Equal(int64(1)))
				gomega.Expect(mock.ExpectationsWereMet()).To(gomega.BeNil())

			})

		})
	})

	ginkgo.Describe("get books from database", func() {
		var b *model.Book

		ginkgo.Context("model exists", func() {
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
				result := sqlmock.NewRows([]string{"id", "title", "author", "pages", "weight"}).
					AddRow(b.Id, b.Title, b.Author, b.Pages, b.Weight)
				mock.ExpectQuery("SELECT (.+) FROM `books` WHERE `books`.`id` = ?").
					WithArgs(b.Id).
					WillReturnRows(result)
				returnedBook, err := manager.GetBook(b.Id)
				gomega.Expect(err).To(gomega.BeNil())
				gomega.Expect(b.Id).To(gomega.Equal(returnedBook.Id))
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
					Id: 100,
				}
			})
			ginkgo.It("return no error & delete the model", func() {
				//删除结果
				result := sqlmock.NewResult(0, 1)
				mock.ExpectBegin()
				mock.ExpectExec("^DELETE FROM `books` WHERE `books`.`id` = ?").
					WithArgs(b.Id).
					WillReturnResult(result)
				mock.ExpectCommit()
				err = manager.DeleteBook(b.Id)
				gomega.Expect(err).To(gomega.BeNil())

				gomega.Expect(mock.ExpectationsWereMet()).To(gomega.BeNil())

			})
		})
	})

	ginkgo.Describe("update books to database", func() {
		var b *model.Book
		ginkgo.Context("model exists ", func() {
			b = &model.Book{
				Id:     100,
				Title:  "test title",
				Author: "test author",
				Pages:  100,
				Weight: 400,
			}
			ginkgo.It("return no error & updated the model", func() {
				result := sqlmock.NewResult(0, 1)
				mock.ExpectBegin()
				mock.ExpectExec("^UPDATE `books` (.+)WHERE `id` = ?").
					WithArgs(b.Title, b.Author, b.Pages, b.Weight, b.Id).
					WillReturnResult(result)
				mock.ExpectCommit()
				err = manager.UpdateBook(b)
				gomega.Expect(err).To(gomega.BeNil())

				gomega.Expect(mock.ExpectationsWereMet()).To(gomega.BeNil())

			})
		})
	})
})
