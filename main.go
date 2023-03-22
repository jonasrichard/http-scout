package main

import (
	"log"

	"github.com/gizak/termui/v3"
	"github.com/jonasrichard/httpscout/capture"
	"github.com/jonasrichard/httpscout/ui"
)

func main() {
	capture := capture.NewCapture()

    streamCh := make(chan ui.Stream)

    go capture.Run(streamCh)

	if err := termui.Init(); err != nil {
		log.Fatalf("Cannot initialize terminal %v", err)
	}
	defer termui.Close()

	ui.New(capture.Devices).Dashboard(streamCh)
}
