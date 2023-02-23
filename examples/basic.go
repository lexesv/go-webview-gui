package main

import (
	"fmt"

	"github.com/lexesv/go-webview-gui"
)

func main() {

	w := webview.New(true, false)
	defer w.Destroy()

	w.SetTitle("Basic Example")
	err := w.SetIcon("icon.png")
	if err != nil {
		fmt.Println(err.Error())
	}
	w.SetSize(480, 320, webview.HintNone)
	w.SetHtml("Thanks for using Golang Webview GUI!")
	w.Run()
}
