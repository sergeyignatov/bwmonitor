package main

import (
	//"flag"
	//"fmt"
	"github.com/sergeyignatov/bwmonitor/api"
	"github.com/urfave/cli"
	"net/http"
	"os"
)

var revision string

func run(c *cli.Context) error {

	//listen := flag.String("listen", ":5312", "host:port for HTTP listening")
	//flag.Parse()
	err := http.ListenAndServe(c.String("listen"), api.Router())
	return err
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
