package model_test

import (
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"os"

	"github.com/weenxin/ulitmate_go_notebook_reading/ch7/model"
)

var _ = ginkgo.Describe("Book", func() {

	//Get Catalog
	//基本用法
	ginkgo.Describe("get model catalog", func() {
		var foxInSocks, lesMis *model.Book
		ginkgo.BeforeEach(func() {
			lesMis = &model.Book{
				Title:  "Les Miserables",
				Author: "Victor Hugo",
				Pages:  2783,
			}
			foxInSocks = &model.Book{
				Title:  "Fox In Socks",
				Author: "Dr. Seuss",
				Pages:  24,
			}
		})

		ginkgo.Context("pages count <= 300", func() {
			ginkgo.It("be a short story", func() {
				gomega.Expect(foxInSocks.Catalog()).To(gomega.Equal(model.CategoryShortStory))
			})
		})
		ginkgo.Context("pages count > 300", func() {
			ginkgo.It("be a novel", func() {
				gomega.Expect(lesMis.Catalog()).To(gomega.Equal(model.CategoryNovel))
			})
		})
	})

	//Extract author name
	//这里主要演示BeforeEach是在每运行每个用例前之前执行的
	ginkgo.Describe("extract author name", func() {
		var theBook *model.Book
		ginkgo.BeforeEach(func() {
			theBook = &model.Book{
				Title:  "Les Miserables",
				Author: "",
				Pages:  2783,
			}
		})
		ginkgo.Context("author name has first name and last name , example: 'Victor Hugo'", func() {
			ginkgo.BeforeEach(func() {
				theBook.Author = "Victor Hugo"
			})
			ginkgo.It("author last name is 'Hugo'", func() {
				gomega.Expect(theBook.LastName()).To(gomega.Equal("Hugo"))
			})
			ginkgo.It("author first name is 'Victor'", func() {
				gomega.Expect(theBook.FirstName()).To(gomega.Equal("Victor"))
			})
		})
		ginkgo.Context("author name has only last name, example: 'Hugo'", func() {
			ginkgo.BeforeEach(func() {
				theBook.Author = "Hugo"
			})
			ginkgo.It("author only has last name , example: 'Hugo'", func() {
				gomega.Expect(theBook.LastName()).To(gomega.Equal("Hugo"))
			})
			ginkgo.It("author first name is ''", func() {
				gomega.Expect(theBook.FirstName()).To(gomega.Equal(""))
			})
		})
		ginkgo.Context("author has middle name , example: 'Victor Marie Hugo'", func() {
			ginkgo.BeforeEach(func() {
				theBook.Author = "Victor Marie Hugo"
			})
			ginkgo.It("author last name is 'Hugo'", func() {
				gomega.Expect(theBook.LastName()).To(gomega.Equal("Hugo"))
			})
			ginkgo.It("author first name is 'Victor'", func() {
				gomega.Expect(theBook.FirstName()).To(gomega.Equal("Victor"))
			})
			ginkgo.It("author middle name is 'Marie'", func() {
				gomega.Expect(theBook.MiddleName()).To(gomega.Equal("Marie"))
			})
		})
		ginkgo.Context("author name is empty", func() {
			ginkgo.BeforeEach(func() {
				theBook.Author = ""
			})
			ginkgo.It("author first name is ''", func() {
				gomega.Expect(theBook.FirstName()).To(gomega.Equal(""))
			})
			ginkgo.It("author middle name is ''", func() {
				gomega.Expect(theBook.MiddleName()).To(gomega.Equal(""))
			})
			ginkgo.It("author last name is ''", func() {
				gomega.Expect(theBook.LastName()).To(gomega.Equal(""))
			})
			ginkgo.It("model status is invalid", func() {
				gomega.Expect(theBook.IsValid()).To(gomega.Equal(false))
			})
		})
	})

	ginkgo.Describe("decode by json string", func() {
		//多了一层，为了测试JustBeforeEach，确实不太好，好的约定优于功能齐全的组件
		ginkgo.Describe("bad cases", func() {
			var theBook *model.Book
			var err error
			var json string
			ginkgo.JustBeforeEach(func() {
				theBook, err = model.NewBookFromJSON(json)
				gomega.Expect(theBook).To(gomega.BeNil()) //should: "the returned model should be nil"
			})
			ginkgo.Context("given a json failed to parse", func() {
				ginkgo.BeforeEach(func() {
					json = `{
						"title":"Les Miserables",
						"author":"Victor Hugo",
						"pages":2783oops
					  }`
				})
				ginkgo.It("it should return error", func() {
					gomega.Expect(err).NotTo(gomega.BeNil())
				})
			})
			ginkgo.Context("given a json is incomplete", func() {
				ginkgo.BeforeEach(func() {
					json = `{
					"title":"Les Miserables",
					"author":"Victor Hugo",
				  }`
				})
				ginkgo.It("it should return error", func() {
					gomega.Expect(err).NotTo(gomega.BeNil())
				})
			})

		})

		ginkgo.Context("given a good json string", func() {
			var theBook *model.Book
			var err error
			ginkgo.BeforeEach(func() {
				json := `{
						"title":"Les Miserables",
						"author":"Victor Hugo",
						"pages":2783
					  }`
				theBook, err = model.NewBookFromJSON(json)
			})
			ginkgo.It("it should return no error", func() {
				gomega.Expect(err).To(gomega.BeNil())
			})
			ginkgo.It("object returned should be equal to json struct", func() {
				gomega.Expect(*theBook).To(gomega.Equal(model.Book{
					Title:  "Les Miserables",
					Author: "Victor Hugo",
					Pages:  2783,
				}))
			})
		})

	})

	ginkgo.Describe("report model weight", func() {
		var b *model.Book
		var originalWeightUnits string
		ginkgo.BeforeEach(func() {
			b = &model.Book{
				Title:  "Les Miserables",
				Author: "Victor Hugo",
				Pages:  2783,
				Weight: 500,
			}
			originalWeightUnits = os.Getenv(model.WeightEnvName)
			ginkgo.DeferCleanup(func() {
				err := os.Setenv(model.WeightEnvName, originalWeightUnits)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})

		})

		ginkgo.Context("has not set WEIGHT_UNITS environment", func() {
			ginkgo.BeforeEach(func() {
				err := os.Unsetenv(model.WeightEnvName)
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})
			ginkgo.It("use g as unit", func() {
				result, err := b.HumanReadableWeight()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(result).To(gomega.Equal("500g"))
			})
		})
		ginkgo.Context("has set WEIGHT_UNITS environment as 'kg'", func() {
			ginkgo.BeforeEach(func() {
				err := os.Setenv(model.WeightEnvName, "kg")
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
			})
			ginkgo.It("use kg as unit", func() {
				result, err := b.HumanReadableWeight()
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				gomega.Expect(result).To(gomega.Equal("0.500kg"))
			})
		})
		ginkgo.Context("has set WEIGHT_UNITS environment to invalid value", func() {
			var result string
			var err error
			ginkgo.BeforeEach(func() {
				err = os.Setenv(model.WeightEnvName, "xxx")
				gomega.Expect(err).NotTo(gomega.HaveOccurred())
				result, err = b.HumanReadableWeight()
			})
			ginkgo.It("return a error", func() {
				gomega.Expect(err).To(gomega.HaveOccurred())
			})
			ginkgo.It("return weight with zero", func() {
				gomega.Expect(result).To(gomega.Equal(""))
			})
		})

	})

})

var _ = ginkgo.Describe("table testing", func() {
	ginkgo.Context("test", func() {
		ginkgo.DescribeTable("Extracting the author's first and last name",
			func(author string, isValid bool, firstName string, lastName string) {
				book := model.Book{
					Title:  "My Book",
					Author: author,
					Pages:  10,
				}
				gomega.Expect(book.IsValid()).To(gomega.Equal(isValid))
				gomega.Expect(book.FirstName()).To(gomega.Equal(firstName))
				gomega.Expect(book.LastName()).To(gomega.Equal(lastName))
			},
			ginkgo.Entry("When author has both names", "Victor Hugo", true, "Victor", "Hugo"),
			ginkgo.Entry("When author has one name", "Hugo", true, "", "Hugo"),
			ginkgo.Entry("When author has a middle name", "Victor Marie Hugo", true, "Victor", "Hugo"),
			ginkgo.Entry("When author has no name", "", false, "", ""),
		)
	},
	)
})

var _ = ginkgo.Describe("Math", func() {
	ginkgo.DescribeTable("addition",
		func(a, b, c int) {
			gomega.Expect(a + b).To(gomega.Equal(c))
		},
		func(a, b, c int) string {
			return fmt.Sprintf("%d + %d = %d", a, b, c)
		},
		ginkgo.Entry(nil, 1, 2, 3),
		ginkgo.Entry(nil, -1, 2, 1),
		ginkgo.Entry(nil, 0, 0, 0),
		ginkgo.Entry(nil, 10, 100, 110),
	)
})
