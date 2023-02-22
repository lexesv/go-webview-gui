package main

import (
	"fmt"

	"github.com/lexesv/go-webview-gui"
)

func main() {

	w := webview.New(true)
	defer w.Destroy()

	webview.Events.Handle = func(state webview.WindowState) {
		//fmt.Println(state)
		switch state {
		case webview.WindowClose:
			w.Terminate() // Click "Show" after
		case webview.WindowResize:
			// Example: save window size for restore in next launch
		case webview.WindowMove:
			// Example: save window position for restore in next launch
		}
	}
	w.SetTitle("Basic Example")
	err := w.SetIcon("icon.png")
	if err != nil {
		fmt.Println(err.Error())
	}
	w.SetSize(480, 320, webview.HintNone)
	w.SetHtml("Thanks for using Golang Webview GUI!")
	w.Run()
}
