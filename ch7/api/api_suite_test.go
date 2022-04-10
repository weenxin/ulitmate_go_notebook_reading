package api_test

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/onsi/gomega"
	"github.com/weenxin/ulitmate_go_notebook_reading/ch7/api"
	"github.com/weenxin/ulitmate_go_notebook_reading/ch7/service"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/onsi/ginkgo/v2"
)

var address = "127.0.0.1:8080"
var server *http.Server

func TestApi(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Api Suite")
}

var _ = ginkgo.BeforeSuite(func() {

	ginkgo.By("initialzing tables")
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
	err = service.InitManagerFromDsn(dsn)
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
	//wait for server start
	time.Sleep(1 * time.Second)

})

var _ = ginkgo.AfterSuite(func() {
	ginkgo.By("stop server")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		fmt.Println("Server Shutdown:", err)
	}

	ginkgo.By("clear database")
	rootPassoword := os.Getenv("ROOT_DATABASE_PWD")
	cmd := exec.Command("mysql", "-uroot", "-h127.0.0.1", fmt.Sprintf("-p%s", rootPassoword))
	cmd.Stdin = bytes.NewBuffer([]byte("drop database if exists test;"))
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	gomega.Expect(err).NotTo(gomega.HaveOccurred())
})
