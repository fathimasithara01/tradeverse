package main

import (
	"github.com/fathimasithara01/tradeverse/intenal/graph"
	"github.com/fathimasithara01/tradeverse/intenal/server"
)

func main() {
	graph.InitSchema()
	r := server.SetupRoutes()
	r.Run(":8080")
}
