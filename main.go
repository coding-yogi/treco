package main

import (
	"github.com/pkg/profile"
	"net/http"
	_"net/http/pprof"
	"treco/cmd"
)

func main() {
	defer profile.Start(profile.MemProfile).Stop()

	go func() {
		http.ListenAndServe(":8081", nil)
	}()

	_ = cmd.Execute()
}
