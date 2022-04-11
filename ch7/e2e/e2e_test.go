package e2e_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
	"github.com/weenxin/ulitmate_go_notebook_reading/ch7/model"
	"io"
	"io/ioutil"
	"net/http"
)

var _ = ginkgo.Describe("Api", func() {
	ginkgo.Describe("Normal use", func() {
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
		ginkgo.It("smoking test", func() {
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
