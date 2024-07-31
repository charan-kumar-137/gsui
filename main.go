package main

import (
	"log"

	"github.com/charan-kumar-137/gsui/display"
	"github.com/charan-kumar-137/gsui/gcs"
)

func main() {
	err := gcs.Init("test")

	if err != nil {
		display.Run()
	} else {
		log.Fatalln(err)
	}
	// display.Run()

}
