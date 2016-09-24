package main

//go:generate go-bindata -pkg cmd -o cmd/templates.go templates/

import "github.com/Luzifer/worktime/cmd"

func main() {
	cmd.Execute()
}
