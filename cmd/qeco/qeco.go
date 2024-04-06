package main

import (
	"fmt"
	"os"

	qeco "github.com/kawabatas/toy-interactive-filtering-tool"
	"github.com/nsf/termbox-go"
)

func main() {
	var err error

	// only support to receive in from Stdin
	// TODO: support either a file or Stdin
	in := os.Stdin // *os.File

	ctx := qeco.NewCtx()
	defer func() {
		if result := ctx.Result(); result != "" {
			os.Stdout.WriteString(result)
		}
	}()

	ctx.ReadBuffer(in)

	err = termbox.Init()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer termbox.Close()

	view := ctx.NewView()
	filter := ctx.NewFilter()
	input := ctx.NewInput()

	go view.Loop()
	go filter.Loop()
	go input.Loop()

	view.Refresh()

	ctx.WaitDone()
}
