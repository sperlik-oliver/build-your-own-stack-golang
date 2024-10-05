package main

import (
	"os"
	"os/exec"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/creack/pty"
)

func main() {
	a := app.New()
	w := a.NewWindow("sleight")

	ui := widget.NewTextGrid() // Create a new TextGrid

	bash := exec.Command("/bin/bash")
	pty, err := pty.Start(bash)

	if err != nil {
		fyne.LogError("Failed to open pty", err)
		os.Exit(1)
	}

	defer bash.Process.Kill()

	pty.Write([]byte("ls\r"))
	time.Sleep(1 * time.Second)
	b := make([]byte, 1024)
	_, err = pty.Read(b)
	if err != nil {
		fyne.LogError("Failed to read pty", err)
	}
	// s := fmt.Sprintf("read bytes from pty.\nContent:%s",  string(b))
	ui.SetText(string(b))

	// Create a new container with a wrapped layout
	// set the layout width to 420, height to 200
	w.SetContent(
		container.New(
			layout.NewGridWrapLayout(fyne.NewSize(420, 200)),
			ui,
		),
	)

	w.ShowAndRun()

}
