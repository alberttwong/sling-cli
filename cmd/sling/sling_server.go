package main

import (
	"fmt"

	"github.com/flarco/g"
	"github.com/flarco/sling/core"
	"github.com/flarco/sling/core/server"
	"github.com/spf13/cast"
)

func slingServer(c *g.CliSC) (err error) {
	fmt.Println("sling - An Extract-Load tool")
	fmt.Println("Slings from a data source to a data target.\nVersion " + core.Version)

	host := "localhost"
	port := "9876"

	for k, v := range c.Vals {
		switch k {
		case "host":
			host = cast.ToString(v)
		case "port":
			port = cast.ToString(v)
		}
	}

	server := server.NewServer(host, port)

	server.Start()

	return
}
