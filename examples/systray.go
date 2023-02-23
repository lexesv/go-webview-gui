package main

import "C"
import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/lexesv/go-webview-gui"
	"github.com/lexesv/go-webview-gui/dialog"
	"github.com/lexesv/go-webview-gui/systray"
)

var (
	iconData []byte
)

func main() {
	var err error

	// See NewInstance function
	new_window := flag.Bool("new_window", false, "")
	title := flag.String("title", "New Window", "")
	width := flag.Int("width", 200, "")
	height := flag.Int("height", 200, "")
	hint := flag.String("hint", "none", "")
	url := flag.String("url", "", "")
	flag.Parse()
	if *new_window {
		NewInstance(title, width, height, hint, url)
	}

	iconData, err = os.ReadFile("icon.png")
	if err != nil {
		panic(err)
	}

	w := webview.New(false, false)
	defer w.Destroy()

	w.SetEventsHandler(func(state webview.WindowState) {
		//fmt.Println(state)
		switch state {
		case webview.WindowClose:
			w.Hide() // Click "Show" after
		case webview.WindowResize:
			// Example: save window size for restore in next launch
		case webview.WindowMove:
			// Example: save window position for restore in next launch
		}
	})

	systray.Register(onReady(w))
	w.SetTitle("Systray Example")
	w.SetIconBites(iconData, len(iconData))
	w.SetSize(480, 320, webview.HintNone)
	w.SetHtml("Thanks for using Golang Webview GUI!")
	w.Run()
}

// NewInstance  - WebView does not support creating a new instance from the current application.
// Therefore, this is a possible option for creating a new window.
func NewInstance(title *string, width *int, height *int, hint *string, url *string) {
	w := webview.New(false, true)
	defer w.Destroy()
	w.SetTitle(*title)
	var _hint webview.Hint
	switch *hint {
	case "none":
		_hint = webview.HintNone
	case "fixed":
		_hint = webview.HintFixed
	}
	fmt.Println(*width, *height, _hint)
	w.SetSize(*width, *height, _hint)
	w.Navigate(*url)
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
			mFocus := systray.AddMenuItem("Focus", "")
			mMove := systray.AddMenuItem("Move", "")
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
			mNewWindow := systray.AddMenuItem("New Window", "")
			systray.AddSeparator()
			mQuit := systray.AddMenuItem("Quit                 ", "")

			for {
				select {
				case <-mQuit.ClickedCh:
					syscall.Kill(syscall.Getpid(), syscall.SIGINT) // kill child window/s
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

				case <-mFocus.ClickedCh:
					w.Dispatch(func() {
						w.Focus()
					})

				case <-mMove.ClickedCh:
					w.Dispatch(func() {
						w.Move(100, 100)
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
				case <-mNewWindow.ClickedCh:
					path, err := os.Executable()
					if err != nil {
						dialog.Message("%s", err.Error()).Error()
						break
					}
					dir, err := os.Getwd()
					if err != nil {
						dialog.Message("%s", err.Error()).Error()
						break
					}
					p := []string{
						"-new_window",
						"-title", "About",
						"-width", "300",
						"-height", "200",
						"-hint", "fixed",
						"-url", "file://" + dir + "/about.html",
					}
					cmd := exec.Command(path, p...)
					//fmt.Println(cmd.Args)
					c := make(chan os.Signal, 2)
					signal.Notify(c, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGINT)
					go func() {
						<-c
						cmd.Process.Kill()
						os.Exit(1)
					}()
					if err := cmd.Start(); err != nil {
						dialog.Message("%s", err.Error()).Error()
					}

				}
			}
		}()
	}
}
