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
	err := w.SetIcon("icon.png")
	if err != nil {
		fmt.Println(err.Error())
	}
	w.SetSize(480, 320, webview.HintNone)
	w.SetUserAgent("My User Agent")
	w.SetHtml(`<div style="width:100%; height: 25px; border: 1px solid blue; padding: 5px; text-align: center;" id="drg">Drag me</div>`)
	w.Init("console.log('init 1');")
	w.Init("console.log('init 2');")
	w.SetBorderless()
	w.SetDraggable("drg")
	//w.SetHtml("Thanks for using Golang Webview GUI!")
	//w.Navigate("https://httpbin.org/headers")
	w.Run()
}

func EventHandler(state webview.WindowState) {
	if state == webview.WindowClose {
		dialog.Message("%s", "Window Closed").Title("Info").Info()
	}
}
