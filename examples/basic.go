package main

import "C"
import (
	"github.com/lexesv/go-webview-gui"
	"github.com/lexesv/go-webview-gui/dialog"
	"github.com/lexesv/go-webview-gui/systray"
)

func main() {
	NewWebView()
}

func NewWebView() {
	w := webview.New(false)
	defer w.Destroy()
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

			mShowTitle := systray.AddMenuItem("GetTitle", "")
			mHide := systray.AddMenuItem("Hide", "")
			mShow := systray.AddMenuItem("Show", "")
			mMaximize := systray.AddMenuItem("Maximize", "")
			mUnmaximize := systray.AddMenuItem("Unmaximize", "")
			mMinimize := systray.AddMenuItem("Minimize", "")
			mUnminimize := systray.AddMenuItem("Unminimize", "")
			systray.AddSeparator()
			mQuit := systray.AddMenuItem("Quit                 ", "")

			for {
				select {
				case <-mQuit.ClickedCh:
					w.Terminate()
				case <-mShowTitle.ClickedCh:
					dialog.Message("%s", w.GetTitle()).Info()
				case <-mHide.ClickedCh:
					w.Dispatch(func() {
						w.Hide()
					})
				case <-mShow.ClickedCh:
					w.Dispatch(func() {
						w.Show()
					})
				case <-mMaximize.ClickedCh:
					w.Dispatch(func() {
						w.Maximize()
					})
				case <-mUnmaximize.ClickedCh:
					w.Dispatch(func() {
						w.Unmaximize()
					})
				case <-mMinimize.ClickedCh:
					w.Dispatch(func() {
						w.Minimize()
					})
				case <-mUnminimize.ClickedCh:
					w.Dispatch(func() {
						w.Unminimize()
					})

				}
			}
		}()
	}
}

func onTrayExit() {
	// clean up here
}
