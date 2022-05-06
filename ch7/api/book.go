package api

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/weenxin/ulitmate_go_notebook_reading/ch7/model"
	"github.com/weenxin/ulitmate_go_notebook_reading/ch7/service"
	"net/http"
)

func getBook(ctx *gin.Context) {
	bookId, err := cast.ToUintE(ctx.Param("book_id"))
	if err != nil {
		makeResponse(ctx, http.StatusBadRequest, "failed", "invalid book id", nil)
		return
	}
	book, err := service.GetManager().GetBook(bookId)
	if err != nil {
		makeResponse(ctx, http.StatusBadRequest, "failed", err.Error(), nil)
		return
	}
	makeResponse(ctx, http.StatusOK, "success", "", book)
}

func CreateBook(ctx *gin.Context) {
	var book model.Book
	err := ctx.Bind(&book)
	if err != nil {
		makeResponse(ctx, http.StatusBadRequest, "failed", err.Error(), nil)
	}

	err = service.GetManager().AddBook(&book)
	if err != nil {
		makeResponse(ctx, http.StatusBadRequest, "failed", err.Error(), nil)
		return
	}
	makeResponse(ctx, http.StatusOK, "success", "", book)
}

func makeResponse(ctx *gin.Context, code int, status, msg string, data interface{}) {
	ctx.JSON(code, gin.H{
		"status":  status,
		"message": msg,
		"data":    data,
	})
}
