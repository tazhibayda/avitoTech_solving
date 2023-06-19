package handler

import "github.com/gin-gonic/gin"

func Routers(router *gin.Engine) {

	route := router.Group("")
	route.GET("/balance", balance)
	route.GET("/transaction/", transaction)
	route.POST("/top-up/", topUp)
	route.POST("/debit/", debit)
	route.POST("/transfer/", transfer)
}
