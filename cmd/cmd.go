package cmd

import (
	"flag"
	//"fmt"
	"github.com/sergeyignatov/bwmonitor/api"
	"net/http"
)

func Run() error {
	listen := flag.String("listen", ":5312", "host:port for HTTP listening")
	flag.Parse()
	err := http.ListenAndServe(*listen, api.Router())
	return err
}
