package main

import (
	"context"
	"os"
	"runtime/trace"
)

func main() {
	// setup tracing
	f, err := os.Create("trace.out")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	trace.Start(f)

	// new FmtSubscriber w/ desired output writer
	sub := NewFmtSubscriber(os.Stdout)

	wm := New()
	ctx, task := trace.NewTask(context.Background(), "wm.Run")
	defer task.End()
	wm.Run(ctx, sub)

	trace.Stop()
}
