package main

import (
	"context"
	"fmt"
	"log"
	"time"

    "github.com/jonasrichard/httpscout/capture"
	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/keyboard"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/tcell"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgets/text"
)

func main() {
    if err := capture.Run(); err != nil {
        fmt.Println(err)
    }
}

func dashboard() {
	t, err := tcell.New(tcell.ColorMode(terminalapi.ColorMode256))

	if err != nil {
		log.Fatalf("Failed to initialize terminal %v", err)
	}
	defer t.Close()

    var i int

	requests, _ := text.New()

    for i = 0; i < 10; i++ {
        requests.Write(fmt.Sprintf("Text #%v\n", i))
    }
    
    content, _ := text.New()

	c, err := container.New(
		t,
		container.SplitHorizontal(
			container.Top(
                container.Border(linestyle.Light),
                container.BorderTitle("HTTP requests"),
				container.PlaceWidget(requests),
			),
			container.Bottom(
                container.Border(linestyle.Light),
                container.BorderTitle("Captured content"),
                container.PlaceWidget(content),
            ),
		),
	)

	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	quitter := func(k *terminalapi.Keyboard) {
		if k.Key == keyboard.KeyEsc || k.Key == keyboard.KeyCtrlC {
			cancel()
		}
	}

    go func(ctx context.Context, requests *text.Text) {
        ticker := time.NewTicker(time.Second)
        defer ticker.Stop()

        for {
            select {
            case <-ticker.C:
                requests.Write("A new text\n")
            case <-ctx.Done():
                return
            }
        }
    }(ctx, requests)

	if err := termdash.Run(ctx, t, c, termdash.KeyboardSubscriber(quitter)); err != nil {
		panic(err)
	}
}
