package main

import (
	"github.com/sergeyignatov/bwmonitor/api"
	"github.com/sergeyignatov/bwmonitor/common"
	"github.com/urfave/cli"
	"net/http"
	"os"
)

var revision string

func run(c *cli.Context) error {
	context := common.NewContext()
	return http.ListenAndServe(c.String("listen"), api.Router(&context))
}

func main() {
	app := cli.NewApp()
	app.Version = revision
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "listen",
			Value: ":5312",
			Usage: "host:port for HTTP listening",
		},
	}
	app.Action = run
	app.Run(os.Args)
}
