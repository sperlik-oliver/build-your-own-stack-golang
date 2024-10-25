package main

import (
	"bufio"
	"os"
	"os/exec"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/creack/pty"
)

// MaxBufferSize sets the size limit
// for our command output buffer.
const MaxBufferSize = 16

func main() {
	// Initialize application
	application := app.New()
	// Initialize window
	window := application.NewWindow("germ")

	// Create a new TextGrid
	ui := widget.NewTextGrid()

	// Spawn a bash subprocess
	os.Setenv("TERM", "dumb")
	cmd := exec.Command("/bin/bash")

	// Start pseudoterminal
	PTY, _ := pty.Start(cmd)

	// Close the bash process when we close our process
	defer cmd.Process.Kill()

	// On key event
	window.Canvas().SetOnTypedKey(func(event *fyne.KeyEvent) {
		if event.Name == fyne.KeyEnter || event.Name == fyne.KeyReturn {
			PTY.Write([]byte{'\r'})
		}
	})

	// On rune event
	window.Canvas().SetOnTypedRune(func(r rune) {
		PTY.WriteString(string(r))
	})

	buffer := [][]rune{}
	reader := bufio.NewReader(PTY)

	// Goroutine that reads from PTY
	go func() {
		line := []rune{}
		buffer = append(buffer, line)
		for {
			r, _, _ := reader.ReadRune()

			line = append(line, r)
			buffer[len(buffer)-1] = line
			if r == '\n' {
				if len(buffer) > MaxBufferSize {
					buffer = buffer[1:]
				}

				line = []rune{}
				buffer = append(buffer, line)
			}

		}
	}()

	// Goroutine that renders to UI
	go func() {
		for {
			time.Sleep(10 * time.Millisecond)

			var lines string
			for _, line := range buffer {
				lines = lines + string(line)
			}
			ui.SetText(string(lines))
		}
	}()

	// Create a new container with a wrapped layout and show itd
	window.SetContent(
		fyne.NewContainerWithLayout(
			layout.NewGridWrapLayout(fyne.NewSize(420, 200)),
			ui,
		),
	)

	window.ShowAndRun()

}
