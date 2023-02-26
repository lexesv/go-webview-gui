package main

import (
	"fmt"

	"github.com/lexesv/go-webview-gui"
	"github.com/lexesv/go-webview-gui/dialog"
)

func main() {

	w := webview.New(true, true)
	defer w.Destroy()

	w.SetWindowEventsHandler(EventHandler)

	w.SetTitle("Basic Example")
	err := w.SetIcon("../asset/icon.png")
	if err != nil {
		fmt.Println(err.Error())
	}
	w.SetSize(480, 320, webview.HintNone)
	w.SetHtml("Thanks for using Golang Webview GUI!")
	w.Run()
}

func EventHandler(state webview.WindowState) {
	if state == webview.WindowClose {
		dialog.Message("%s", "Window Closed").Title("Info").Info()
	}
}
