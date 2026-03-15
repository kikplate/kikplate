package main

import (
	"github.com/kickplate/api/bootstrap"
)

func main() {
	err := bootstrap.RootApp.Execute()
	if err != nil {
		return
	}
}
