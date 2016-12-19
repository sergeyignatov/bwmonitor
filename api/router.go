package api

import (
	//"fmt"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
	"github.com/sergeyignatov/bwmonitor/common"
	"net/http"
	"time"
)

func Fail(c *gin.Context, err error) {
	c.Error(err)
	c.JSON(500, common.NewApiResponse(err))
}

func Router() http.Handler {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))
	router.Use(gin.ErrorLogger())
	router.Use(static.Serve("/", static.LocalFile("./static", true)))
	root := router.Group("/api/1.0")
	{
		root.POST("/bw", apiMeasureBW)
		root.GET("/bw/:name", apiServeFile)
		root.GET("/ping", apiPing)
	}
	return router
}
