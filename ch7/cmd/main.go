package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/weenxin/ulitmate_go_notebook_reading/ch7/api"
	"github.com/weenxin/ulitmate_go_notebook_reading/ch7/service"
)

var (
	dsn     = flag.String("dsn", "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local", "database dsn")
	address = flag.String("address", "0.0.0.0:8080", "server bind address")
)

func main() {
	flag.Parse()
	err := service.InitManagerFromDsn(*dsn)
	if err != nil {
		fmt.Println("init database failed")
	}
	r := gin.Default()
	group := r.Group("/books")
	api.InitRoute(group)
	err = r.Run(*address)
	if err != nil {
		panic("failed to start server")
	}
}
