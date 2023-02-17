package main

import "C"
import (
	"fmt"

	"github.com/lexesv/go-webview-gui"
	"github.com/lexesv/go-webview-gui/systray"
)

func main() {
	NewWebView()
}

func NewWebView() {
	w := webview.New(false)
	//defer w.Destroy()
	w.SetTitle("Basic Example")
	w.SetSize(480, 320, webview.HintNone)
	w.SetHtml("Thanks for using webview!")
	//w.Navigate("https://google.com")
	//w.SetBorderless()
	systray.Register(onReady(w))
	w.Run()
}

func onReady(w webview.WebView) func() {
	return func() {
		systray.SetTitle("Tray")
		go func() {
			mShowHide := systray.AddMenuItem("Show title", "")
			systray.AddSeparator()
			mQuit := systray.AddMenuItem("Quit                 ", "")

			for {
				select {
				case <-mQuit.ClickedCh:
					w.Terminate()
				case <-mShowHide.ClickedCh:
					//mShowHide.SetTitle("Hide")
					//w.Hide()
					fmt.Println(w.GetTitle())
				}
			}
		}()
	}
}

func onTrayExit() {
	// clean up here
}
