package ui

import (
	"fmt"

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
	streams       []Stream
	devices       []string
	active        ActiveView
	tw            int
	th            int
	requests      *widgets.List
	content       *widgets.Paragraph
	help          *widgets.Paragraph
	deviceChooser *widgets.List
}

func New(devices []string) *State {
	return &State{
		streams: make([]Stream, 0),
		devices: devices,
	}
}

func (s *State) Dashboard(ch chan Stream) (err error) {
	s.tw, s.th = termui.TerminalDimensions()

	requests := widgets.NewList()
	requests.Title = "HTTP Requests"
	requests.SetRect(0, 0, s.tw, s.th/2)
	requests.Border = true
	requests.SelectedRowStyle = termui.NewStyle(termui.ColorBlack, termui.ColorWhite)

	s.requests = requests

	content := widgets.NewParagraph()
	content.Title = "Content"
	content.SetRect(0, s.th/2, s.tw, s.th)
	content.Border = true

	s.content = content

	s.deviceChooser = s.ChooseDevice()

	termui.Render(requests, content)

	events := termui.PollEvents()

	for {
		select {
		case evt := <-events:
			switch evt.ID {
			case "d":
				s.active = DeviceList
			case "q", "<Escape>", "<C-c>":
				return
			}

			switch s.active {
			// Handle key events for requests list
			case RequestList:
				switch evt.ID {
				case "j", "<Down>":
					s.requests.ScrollDown()

                    // TODO update content in a more optimal way
                    i := s.requests.SelectedRow
                    s.content.Text = s.streams[i].Request + "\n\n" + s.streams[i].Response
				case "k", "<Up>":
					s.requests.ScrollUp()

                    // TODO update content in a more optimal way
                    i := s.requests.SelectedRow
                    s.content.Text = s.streams[i].Request + "\n\n" + s.streams[i].Response
				}

			case DeviceList:
				switch evt.ID {
				case "j", "<Down>":
					s.deviceChooser.ScrollDown()
				case "k", "<Up>":
					s.deviceChooser.ScrollUp()
				case "<Enter>":
					deviceName := s.devices[s.deviceChooser.SelectedRow]
					s.requests.Title = "HTTP Requests - " + deviceName
					s.requests.BorderStyle = termui.NewStyle(termui.ColorGreen)
					s.active = RequestList
				}
			}

		case stream := <-ch:
			s.streams = append(s.streams, stream)

			line := fmt.Sprintf("%s %s %s", stream.Timestamp, stream.Host, stream.Path)

			s.requests.Rows = append(s.requests.Rows, line)
		}

		// Redraw UI
		switch s.active {
		case DeviceList:
			termui.Render(s.requests, s.content, s.deviceChooser)
		default:
			termui.Render(s.requests, s.content)
		}
	}
}

func (s *State) Help() *widgets.Paragraph {
	help := widgets.NewParagraph()
	help.Title = "Help"
	help.SetRect(s.tw/2-15, s.th/2-10, s.tw/2+15, s.th/2+10)
	help.Border = true

	help.Text = "S    Start/stop capture"

	return help
}

func (s *State) ChooseDevice() *widgets.List {
	devs := widgets.NewList()
	devs.Title = "Devices"
	devs.SetRect(s.tw/2-10, s.th/2-5, s.tw/2+10, s.th/2+5)
	devs.Border = true
	devs.Rows = s.devices
	devs.SelectedRowStyle = termui.NewStyle(termui.ColorBlack, termui.ColorWhite)

	return devs
}
