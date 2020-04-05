package main

import (
	"github.com/wchy1001/docker-images/cmd"
	"log"
	"os"
)

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	cmd.Execute()
}
