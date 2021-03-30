package main

import (
	_ "net/http/pprof"
	"treco/cmd"
)

func main() {
	_ = cmd.Execute()
}
