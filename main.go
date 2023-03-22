package main

import (
	"fmt"
	"log"

	"github.com/gizak/termui/v3"
	"github.com/jonasrichard/httpscout/capture"
	"github.com/jonasrichard/httpscout/ui"
)

func main2() {
	capture := capture.NewCapture()

	if err := capture.Run(); err != nil {
		fmt.Println(err)
	}
}

func main() {
	dashboard()
}

func dashboard() {
    if err := termui.Init(); err != nil {
        log.Fatalf("Cannot initialize terminal %v", err)
    }
    defer termui.Close()

    ui.New().Dashboard()
}
