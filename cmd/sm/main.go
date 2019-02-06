package main

import (
	"os"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	sm := kingpin.New("sm", "CLI to interact with the Services Marketplace")
	sm.Command("broker", "")
	sm.UsageWriter(os.Stdout)

	kingpin.MustParse(sm.Parse(os.Args[1:]))
}
