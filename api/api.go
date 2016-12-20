package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sergeyignatov/bwmonitor/client"
	"github.com/sergeyignatov/bwmonitor/common"
	"io"
	"os"
	"strconv"
)

func apiPing(c *gin.Context) {
	c.JSON(200, "ok")
}
func apiMeasureBWM(c *gin.Context) {
	dest := c.Params.ByName("dest")
	cc := client.NewClient(dest)
	t, err := cc.DownloadSpeed()
	if err != nil {
		Fail(c, err)
		return
	}
	c.JSON(200, common.NewApiResponse(t))
}
func apiMeasureBW(c *gin.Context) {
	dest := c.PostForm("dest")
	if dest == "" {
		Fail(c, fmt.Errorf("no dest"))
		return
	}
	cc := client.NewClient(dest)
	t, err := cc.DownloadSpeed()
	if err != nil {
		Fail(c, err)
		return
	}
	c.JSON(200, common.NewApiResponse(t))
}

func apiServeFile(c *gin.Context) {
	_ = c.Params.ByName("name")
	get_params := c.Request.URL.Query().Get("size")
	size := 128 * 1024
	if len(get_params) > 0 {
		if i, err := strconv.Atoi(string(get_params)); err == nil {
			size = i
		}
	}
	devzero, err := os.Open("/dev/zero") // For read access.
	if err != nil {
		return
	}

	defer devzero.Close()
	io.CopyN(c.Writer, devzero, int64(size))
	return
}
