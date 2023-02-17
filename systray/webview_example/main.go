package main

import (
	"github.com/webview/webview"

	"github.com/ghostiam/systray"
	"github.com/ghostiam/systray/example/icon"
)

func main() {
	debug := true
	w := webview.New(debug)
	defer w.Destroy()
	w.SetTitle("Minimal webview example")
	w.SetSize(800, 600, webview.HintNone)
	w.Navigate("https://en.m.wikipedia.org/wiki/Main_Page")

	systray.Register(onReady(w))

	w.Run()
}

func onReady(w webview.WebView) func() {
	return func() {
		systray.SetTemplateIcon(icon.Data, icon.Data)
		systray.SetTitle("Webview example")
		mShowLantern := systray.AddMenuItem("Show Lantern", "")
		mShowWikipedia := systray.AddMenuItem("Show Wikipedia", "")
		mQuit := systray.AddMenuItem("Quit", "Quit the whole app")
		go func() {
			for {
				select {
				case <-mShowLantern.ClickedCh:
					w.Dispatch(func() { w.Navigate("https://www.getlantern.org") })
				case <-mShowWikipedia.ClickedCh:
					w.Dispatch(func() { w.Navigate("https://www.wikipedia.org") })
				case <-mQuit.ClickedCh:
					w.Terminate()
				}
			}
		}()
	}
}
