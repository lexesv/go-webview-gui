package main

import "C"
import (
	"os"

	"github.com/lexesv/go-webview-gui"
	"github.com/lexesv/go-webview-gui/dialog"
	"github.com/lexesv/go-webview-gui/systray"
)

var (
	iconData []byte
)

func main() {
	iconData, _ = os.ReadFile("icon.png")
	w := webview.New(false)
	defer w.Destroy()
	systray.Register(onReady(w))
	w.SetTitle("Systray Example")
	w.SetSize(480, 320, webview.HintNone)
	w.SetHtml("Thanks for using Golang Webview GUI!")
	w.Run()
}

func onReady(w webview.WebView) func() {
	return func() {
		systray.SetTitle("Tray")
		systray.SetIcon(iconData)
		go func() {

			mShowTitle := systray.AddMenuItem("GetTitle", "")
			mGetSizePosition := systray.AddMenuItem("Get Size & Position", "")
			mHide := systray.AddMenuItem("Hide", "")
			mShow := systray.AddMenuItem("Show", "")
			mMaximize := systray.AddMenuItem("Maximize", "")
			mUnmaximize := systray.AddMenuItem("Unmaximize", "")
			mMinimize := systray.AddMenuItem("Minimize", "")
			mUnminimize := systray.AddMenuItem("Unminimize", "")
			mSetFullScreen := systray.AddMenuItem("SetFullScreen", "")
			mExitFullScreen := systray.AddMenuItem("ExitFullScreen", "")
			mSetAlwaysOnTopTrue := systray.AddMenuItem("SetAlwaysOnTop True", "")
			mSetAlwaysOnTopFalse := systray.AddMenuItem("SetAlwaysOnTop False", "")
			mSetBorderless := systray.AddMenuItem("SetBorderless", "")
			mSetBordered := systray.AddMenuItem("SetBordered", "")
			systray.AddSeparator()
			mNewWebView := systray.AddMenuItem("NewWebView", "")
			systray.AddSeparator()
			mQuit := systray.AddMenuItem("Quit                 ", "")

			for {
				select {
				case <-mQuit.ClickedCh:
					w.Terminate()
				case <-mShowTitle.ClickedCh:
					dialog.Message("%s", w.GetTitle()).Info()
				case <-mGetSizePosition.ClickedCh:
					width, height, hint := w.GetSize()
					x, y := w.GetPosition()
					dialog.Message("Size:%dx%d %v. \nPosition: X:%d Y:%d", width, height, hint, x, y).Info()
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
				case <-mSetBorderless.ClickedCh:
					w.Dispatch(func() {
						w.SetBorderless()
					})
				case <-mSetBordered.ClickedCh:
					w.Dispatch(func() {
						w.SetBordered()
					})
				case <-mSetAlwaysOnTopTrue.ClickedCh:
					w.Dispatch(func() {
						w.SetAlwaysOnTop(true)
					})
				case <-mSetAlwaysOnTopFalse.ClickedCh:
					w.Dispatch(func() {
						w.SetAlwaysOnTop(false)
					})
				case <-mSetFullScreen.ClickedCh:
					w.Dispatch(func() {
						w.SetFullScreen()
					})
				case <-mExitFullScreen.ClickedCh:
					w.Dispatch(func() {
						w.ExitFullScreen()
					})
				case <-mNewWebView.ClickedCh:
					//NewWebView()

				}
			}
		}()
	}
}

func onTrayExit() {
	// clean up here
}
