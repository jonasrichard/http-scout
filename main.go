package main

import (
	"log"

	"github.com/gizak/termui/v3"
	"github.com/jonasrichard/httpscout/capture"
	"github.com/jonasrichard/httpscout/ui"
)

func main() {
	capture := capture.NewCapture()

	dashboard(capture.Devices)

	//if err := capture.Run(); err != nil {
	//	fmt.Println(err)
	//}
}

func dashboard(devices []string) {
	if err := termui.Init(); err != nil {
		log.Fatalf("Cannot initialize terminal %v", err)
	}
	defer termui.Close()

	ui.New(devices).Dashboard()
}
