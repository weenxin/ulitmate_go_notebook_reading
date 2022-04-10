package api_test

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/onsi/gomega"
	"github.com/weenxin/ulitmate_go_notebook_reading/ch7/api"
	"github.com/weenxin/ulitmate_go_notebook_reading/ch7/service"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/onsi/ginkgo/v2"
)

var address = "0.0.0.0:8080"
var server *http.Server

func TestApi(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Api Suite")
}

var _ = ginkgo.BeforeSuite(func() {

	ginkgo.By("initializing database")
	dsn := os.Getenv("DATABASE_DSN")
	gomega.Expect(dsn).NotTo(gomega.BeEmpty())
	err := service.InitManagerFromDsn(dsn)
	gomega.Expect(err).NotTo(gomega.HaveOccurred())

	ginkgo.By("start server")
	r := gin.Default()
	group := r.Group("/books")
	api.InitRoute(group)

	server = &http.Server{
		Addr:    address,
		Handler: r,
	}

	go func() {
		// service connections
		server.ListenAndServe()
	}()

})

var _ = ginkgo.AfterSuite(func() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("Server Shutdown:", err)
	}
})
