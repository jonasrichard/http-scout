package ui

import (
	"github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

type ActiveView int

const (
	RequestList ActiveView = iota
	ContentView
	DeviceList
	Help
)

const TopFrame = "top-frame"
const BottomFrame = "bottom-frame"

type Stream struct {
	Timestamp string
	Host      string
	Path      string
	Request   string
	Response  string
}

type State struct {
	active  ActiveView
	streams []Stream
}

func New() *State {
	return &State{
		streams: make([]Stream, 0),
	}
}

func (s *State) Dashboard() (err error) {
	requests := widgets.NewList()
	requests.Title = "HTTP Requests"
	requests.SetRect(0, 0, 40, 10)
	requests.Border = true

	termui.Render(requests)

	for evt := range termui.PollEvents() {
		switch evt.ID {
        case "q":
            return
		}
	}

	return
}
