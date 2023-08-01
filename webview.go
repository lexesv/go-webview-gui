package webview

/*
#cgo linux openbsd freebsd netbsd CXXFLAGS: -DWEBVIEW_GTK -std=c++11
#cgo linux openbsd freebsd netbsd pkg-config: gtk+-3.0 webkit2gtk-4.0

#cgo darwin CXXFLAGS: -DWEBVIEW_COCOA -std=c++17
#cgo darwin LDFLAGS: -framework WebKit -framework Cocoa

#cgo windows CXXFLAGS: -DWEBVIEW_EDGE -std=c++17 -I./windows/include
#cgo windows LDFLAGS: -static -static-libstdc++ -static-libgcc -ladvapi32 -lole32 -lshell32 -lshlwapi -luser32 -lversion -lGdiplus

#include "webview.h"

#include <stdlib.h>
#include <stdint.h>

void CgoWebViewDispatch(webview_t w, uintptr_t arg);
void CgoWebViewBind(webview_t w, const char *name, uintptr_t index);
void event_handler(int state);
typedef void (*closure)();
*/
import "C"
import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"sync"
	"unsafe"
)

func init() {
	// Ensure that main.main is called from the main thread
	runtime.LockOSThread()
}

// Hint are used to configure window sizing and resizing
type Hint int

// WindowState used when changing the native window status
type WindowState int

const (
	// Width and height are default size
	HintNone = C.WEBVIEW_HINT_NONE
	// Window size can not be changed by a user
	HintFixed = C.WEBVIEW_HINT_FIXED
	// Width and height are minimum bounds
	HintMin = C.WEBVIEW_HINT_MIN
	// Width and height are maximum bounds
	HintMax = C.WEBVIEW_HINT_MAX

	//--------------------------------------------------------- MacOS Linux Windows
	WindowClose          = C.WEBVIEW_WINDOW_CLOSE          // 0   +     +     +
	WindowFocus          = C.WEBVIEW_WINDOW_FOCUS          // 1   +     +     +
	WindowBlur           = C.WEBVIEW_WINDOW_BLUR           // 2   +     +     +
	WindowMove           = C.WEBVIEW_WINDOW_MOVE           // 3   +     +     +
	WindowResize         = C.WEBVIEW_WINDOW_RESIZE         // 4   +     +     +
	WindowFullScreen     = C.WEBVIEW_WINDOW_FULLSCREEN     // 5   +     +     +
	WindowExitFullScreen = C.WEBVIEW_WINDOW_EXITFULLSCREEN // 6   +     +     +
	WindowMaximize       = C.WEBVIEW_WINDOW_MAXIMIZE       // 7   -     +     +
	WindowUnmaximize     = C.WEBVIEW_WINDOW_UNMAXIMIZE     // 8   -     +     -
	WindowMinimize       = C.WEBVIEW_WINDOW_MINIMIZE       // 9   +     +     +
	WindowUnminimize     = C.WEBVIEW_WINDOW_UNMINIMIZE     // 10  +     +     -
	// WindowMinimize WindowUnminimize - Does not work with Stage Manager (MacOS)
)

type WebView interface {

	// Run runs the main loop until it's terminated. After this function exits -
	// you must destroy the webview.
	Run()

	// Terminate stops the main loop. It is safe to call this function from
	// a background thread.
	Terminate()

	// Dispatch posts a function to be executed on the main thread. You normally
	// do not need to call this function, unless you want to tweak the native
	// window.
	Dispatch(f func())

	// Destroy destroys a webview and closes the native window.
	Destroy()

	// Window returns a native window handle pointer. When using GTK backend the
	// pointer is GtkWindow pointer, when using Cocoa backend the pointer is
	// NSWindow pointer, when using Win32 backend the pointer is HWND pointer.
	Window() unsafe.Pointer

	// SetTitle updates the title of the native window. Must be called from the UI
	// thread.
	SetTitle(title string)

	// SetSize updates native window size. See Hint constants.
	SetSize(w int, h int, hint Hint)

	// Navigate navigates webview to the given URL. URL may be a properly encoded data.
	// URI. Examples:
	// w.Navigate("https://github.com/webview/webview")
	// w.Navigate("data:text/html,%3Ch1%3EHello%3C%2Fh1%3E")
	// w.Navigate("data:text/html;base64,PGgxPkhlbGxvPC9oMT4=")
	Navigate(url string)

	// SetHtml sets the webview HTML directly.
	// Example: w.SetHtml(w, "<h1>Hello</h1>");
	SetHtml(html string)

	// Init injects JavaScript code at the initialization of the new page. Every
	// time the webview will open a new page - this initialization code will
	// be executed. It is guaranteed that code is executed before window.onload.
	Init(js string)

	// Eval evaluates arbitrary JavaScript code. Evaluation happens asynchronously,
	// also the result of the expression is ignored. Use RPC bindings if you want
	// to receive notifications about the results of the evaluation.
	Eval(js string)

	// Bind binds a callback function so that it will appear under the given name
	// as a global JavaScript function. Internally it uses webview_init().
	// Callback receives a request string and a user-provided argument pointer.
	// Request string is a JSON array of all the arguments passed to the
	// JavaScript function.
	//
	// f must be a function
	// f must return either value and error or just error
	Bind(name string, f interface{}) error

	// SetUserAgent sets a custom user agent string for the webview.
	SetUserAgent(userAgent string)

	// SetWindowEventsHandler sets the window status change event handling function
	// Should be called before calling the "Run" method
	// Example:
	// w.SetWindowEventsHandler("test", func(state webview.WindowState) {
	//		switch state {
	//		case webview.WindowClose:
	//			w.Hide()
	//		case webview.WindowResize:
	//			// Example: save window size for restore in next launch
	//		case webview.WindowMove:
	//			// Example: save window position for restore in next launch
	//		}
	//	})
	SetWindowEventsHandler(key string, f func(state WindowState))

	// UnSetWindowEventsHandler unsets the window status change event handling function
	UnSetWindowEventsHandler(key string)

	// IsExistWindowEventsHandler checks if this event handler exists
	IsExistWindowEventsHandler(key string) bool

	// GetTitle gets the title of the native window.
	GetTitle() string

	// Hide hides the native window.
	Hide()

	// Show shows the native window.
	Show()

	// SetBorderless activates the borderless mode.
	SetBorderless()

	// SetBordered activates the bordered mode.
	SetBordered()

	// IsMaximized If the native window is maximized it returns true, otherwise false.
	IsMaximized() bool

	// Maximize maximizes the native window.
	Maximize()

	// Unmaximize unmaximizes the native window.
	Unmaximize()

	// IsMinimized If the native window is minimized it returns true, otherwise false.
	// IsMinimized Does not work with Stage Manager (MacOS)
	IsMinimized() bool

	// Minimize minimizes the native window.
	Minimize()

	// Unminimize unminimizes the native window.
	Unminimize()

	// IsVisible If the native window is visible it returns true, otherwise false.
	IsVisible() bool

	// SetFullScreen activates the full-screen mode.
	SetFullScreen()

	// ExitFullScreen exits full-screen mode.
	ExitFullScreen()

	// IsFullScreen If the native window is in full-screen mode it returns true, otherwise false.
	IsFullScreen() bool

	// SetIcon sets the application icon (from filename).
	// Example:
	// w.SetIcon("icon.png")
	SetIcon(icon string) error

	// SetIconBites sets the application icon (from []byte).
	// Example:
	// iconData, err = os.ReadFile("icon.png")
	// w.SetIconBites(iconData, len(iconData))
	SetIconBites(b []byte, size int)

	// SetAlwaysOnTop activates (=true) / deactivates(=false) the top-most mode.
	SetAlwaysOnTop(onTop bool)

	// GetSize gets the size of the native window.
	GetSize() (width int, height int, hint Hint)

	// GetPosition gets the coordinates of the native window.
	GetPosition() (x, y int)

	// Move moves the native window to the specified coordinates.
	Move(x, y int)

	// Focus set the focus on the native window.
	Focus()

	// SetContentStateHandler sets the document status change event handling function
	// Should be called before calling the "Run" method
	// Example:
	// w.SetContentStateHandler("test", func(state string) {
	//		fmt.Printf("document content state: %s\n", state)
	// })
	// Status of the document:
	// uninitialized - Has not started loading
	// loading - Is loading
	// loaded - Has been loaded
	// interactive - Has loaded enough to interact with
	// complete - Fully loaded
	SetContentStateHandler(key string, f func(state string))

	// UnSetContentStateHandler unsets the document status change event handling function
	UnSetContentStateHandler(key string)

	// IsExistContentStateHandler checks if this event handler exists
	IsExistContentStateHandler(key string) bool

	// SetDraggable converts a given DOM element to a draggable region. The user will be able to drag the native window by dragging the given DOM element.
	// This feature is suitable to make custom window bars along with the borderless mode.
	SetDraggable(id string)

	// UnSetDraggable converts a draggable region to a normal DOM elements by removing drag event handlers.
	UnSetDraggable(id string)

	// OpenUrlInBrowser Opens the specified link in the web browser
	OpenUrlInBrowser(url string) (err error)

	// Just for an example of using the Bind function. See js.go and init.js
	GetHtml() (s string)
	GetUrl() string
	GetPageTitle() string
}

type webview struct {
	w                 C.webview_t
	Hint              Hint
	data              map[string]interface{}
	DraggableElements sync.Map
}

// EventHandler It is used to intercept changes in the status of the native window
type eventsHandler struct {
	handle_ws   map[string]func(state WindowState)
	handle_cs   map[string]func(state string)
	exitOnClose bool
	exitFunc    func()
}

var (
	m        sync.Mutex
	index    uintptr
	dispatch = map[uintptr]func(){}
	bindings = map[uintptr]func(id, req string) (interface{}, error){}
	events   = eventsHandler{}
)

func boolToInt(b bool) C.int {
	if b {
		return 1
	}
	return 0
}

//export event_handler
func event_handler(state C.int) {
	if events.handle_ws != nil {
		for _, f := range events.handle_ws {
			f(WindowState(state))
		}
	}
	if state == WindowClose && events.exitOnClose {
		events.exitFunc()
	}
}

// New calls NewWindow to create a new window and a new webview instance. If debug
// is non-zero - developer tools will be enabled (if the platform supports them).
func New(debug, exitOnClose bool) WebView {
	res := C.webview_create(boolToInt(debug), nil)
	if res == nil {
		return nil
	}
	w := &webview{w: res}
	w.data = make(map[string]interface{})
	events.handle_ws = make(map[string]func(state WindowState))
	events.handle_cs = make(map[string]func(state string))
	events.exitOnClose = exitOnClose
	events.exitFunc = w.Exit
	w.initJSFunc()
	return w
}

func (w *webview) Exit() {
	C.webview_destroy(w.w)
	os.Exit(1)
}

func (w *webview) Destroy() {
	C.webview_destroy(w.w)
}

func (w *webview) Run() {
	C.webview_set_event_handler(C.closure(C.event_handler))
	C.webview_run(w.w)
}

func (w *webview) SetWindowEventsHandler(key string, f func(state WindowState)) {
	events.handle_ws[key] = f
}

func (w *webview) UnSetWindowEventsHandler(key string) {
	delete(events.handle_ws, key)
}

func (w *webview) IsExistWindowEventsHandler(key string) bool {
	_, ok := events.handle_ws[key]
	if ok {
		return true
	}
	return false
}

func (w *webview) Terminate() {
	C.webview_terminate(w.w)
}

func (w *webview) Window() unsafe.Pointer {
	return C.webview_get_window(w.w)
}

func (w *webview) Navigate(url string) {
	s := C.CString(url)
	defer C.free(unsafe.Pointer(s))
	C.webview_navigate(w.w, s)
}

func (w *webview) SetHtml(html string) {
	s := C.CString(html)
	defer C.free(unsafe.Pointer(s))
	C.webview_set_html(w.w, s)
}

func (w *webview) SetTitle(title string) {
	s := C.CString(title)
	defer C.free(unsafe.Pointer(s))
	C.webview_set_title(w.w, s)
}

func (w *webview) SetSize(width int, height int, hint Hint) {
	C.webview_set_size(w.w, C.int(width), C.int(height), C.int(hint))
}

func (w *webview) Init(js string) {
	s := C.CString(js)
	defer C.free(unsafe.Pointer(s))
	C.webview_init(w.w, s)
}

func (w *webview) Eval(js string) {
	s := C.CString(js)
	defer C.free(unsafe.Pointer(s))
	C.webview_eval(w.w, s)
}

func (w *webview) SetUserAgent(userAgent string) {
	ua := C.CString(userAgent)
	defer C.free(unsafe.Pointer(ua))
	C.webview_set_user_agent(w.w, ua)
}

func (w *webview) GetSize() (width int, height int, hint Hint) {
	if !w.IsVisible() {
		return
	}
	wc := C.int(width)
	hc := C.int(height)
	C.webview_get_size(w.w, (*C.int)(&wc), (*C.int)(&hc))
	width = int(wc)
	height = int(hc)

	hint = w.Hint
	return width, height, hint
}

func (w *webview) GetPosition() (x, y int) {
	if !w.IsVisible() {
		return
	}
	xc := C.int(x)
	yc := C.int(y)
	C.webview_get_position(w.w, (*C.int)(&xc), (*C.int)(&yc))
	x = int(xc)
	y = int(yc)
	return x, y
}

func (w *webview) Move(x, y int) {
	C.webview_move(w.w, C.int(x), C.int(y))
}

func (w *webview) Focus() {
	C.webview_focus(w.w)
}

func (w *webview) Hide() {
	if w.IsVisible() {
		C.webview_hide(w.w)
	}
}

func (w *webview) Show() {
	if !w.IsVisible() {
		C.webview_show(w.w)
	}
}

func (w *webview) IsVisible() bool {
	if C.webview_is_visible(w.w) != 0 {
		return true
	}
	return false
}

func (w *webview) GetTitle() string {
	s := C.webview_get_title(w.w)
	return C.GoString(s)
}

func (w *webview) SetBorderless() {
	C.webview_set_borderless(w.w)
}

func (w *webview) SetBordered() {
	C.webview_set_bordered(w.w, C.int(w.Hint))
}

func (w *webview) IsMaximized() bool {
	if C.webview_is_maximized(w.w) != 0 {
		return true
	}
	return false
}

func (w *webview) Maximize() {
	if w.IsMaximized() {
		return
	}
	C.webview_maximize(w.w)
}

func (w *webview) Unmaximize() {
	if !w.IsMaximized() {
		return
	}
	C.webview_unmaximize(w.w)
}

func (w *webview) IsMinimized() bool {
	if C.webview_is_minimized(w.w) != 0 {
		return true
	}
	return false
}

func (w *webview) Minimize() {
	if w.IsMinimized() {
		return
	}
	C.webview_minimize(w.w)
}

func (w *webview) Unminimize() {
	if !w.IsMinimized() {
		return
	}
	C.webview_unminimize(w.w)
}

func (w *webview) SetFullScreen() {
	if w.IsFullScreen() {
		return
	}
	C.webview_set_full_screen(w.w)
}

func (w *webview) ExitFullScreen() {
	if !w.IsFullScreen() {
		return
	}
	C.webview_exit_full_screen(w.w)
}

func (w *webview) IsFullScreen() bool {
	if C.webview_is_full_screen(w.w) != 0 {
		return true
	}
	return false
}

func (w *webview) SetIcon(icon string) error {
	b, err := os.ReadFile(icon)
	if err != nil {
		return err
	}
	C.webview_set_icon(w.w, C.CString(string(b)), C.long(len(b)))
	return nil
}

func (w *webview) SetIconBites(b []byte, size int) {
	C.webview_set_icon(w.w, C.CString(string(b)), C.long(size))
}

func (w *webview) SetAlwaysOnTop(onTop bool) {
	C.webview_set_always_ontop(w.w, boolToInt(onTop))
}

func (w *webview) SetContentStateHandler(key string, f func(state string)) {
	events.handle_cs[key] = f
}

func (w *webview) UnSetContentStateHandler(key string) {
	delete(events.handle_cs, key)
}

func (w *webview) IsExistContentStateHandler(key string) bool {
	_, ok := events.handle_cs[key]
	if ok {
		return true
	}
	return false
}

func (w *webview) SetDraggable(id string) {
	w.DraggableElements.Store(id, true)
}

func (w *webview) UnSetDraggable(id string) {
	w.DraggableElements.Store(id, false)
}

func (w *webview) GetHtml() string {
	if _, ok := w.data["html"]; !ok {
		return ""
	}
	return w.data["html"].(string)
}

func (w *webview) GetPageTitle() string {
	if _, ok := w.data["title"]; !ok {
		return ""
	}
	return w.data["title"].(string)
}

func (w *webview) GetUrl() string {
	if _, ok := w.data["url"]; !ok {
		return ""
	}
	return w.data["url"].(string)
}

func (w *webview) OpenUrlInBrowser(url string) (err error) {
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return err
}

func (w *webview) Dispatch(f func()) {
	m.Lock()
	for ; dispatch[index] != nil; index++ {
	}
	dispatch[index] = f
	m.Unlock()
	C.CgoWebViewDispatch(w.w, C.uintptr_t(index))
}

//export _webviewDispatchGoCallback
func _webviewDispatchGoCallback(index unsafe.Pointer) {
	m.Lock()
	f := dispatch[uintptr(index)]
	delete(dispatch, uintptr(index))
	m.Unlock()
	f()
}

//export _webviewBindingGoCallback
func _webviewBindingGoCallback(w C.webview_t, id *C.char, req *C.char, index uintptr) {
	m.Lock()
	f := bindings[uintptr(index)]
	m.Unlock()
	jsString := func(v interface{}) string { b, _ := json.Marshal(v); return string(b) }
	status, result := 0, ""
	if res, err := f(C.GoString(id), C.GoString(req)); err != nil {
		status = -1
		result = jsString(err.Error())
	} else if b, err := json.Marshal(res); err != nil {
		status = -1
		result = jsString(err.Error())
	} else {
		status = 0
		result = string(b)
	}
	s := C.CString(result)
	defer C.free(unsafe.Pointer(s))
	C.webview_return(w, id, C.int(status), s)
}

func (w *webview) Bind(name string, f interface{}) error {
	v := reflect.ValueOf(f)
	// f must be a function
	if v.Kind() != reflect.Func {
		return errors.New("only functions can be bound")
	}
	// f must return either value and error or just error
	if n := v.Type().NumOut(); n > 2 {
		return errors.New("function may only return a value or a value+error")
	}

	binding := func(id, req string) (interface{}, error) {
		raw := []json.RawMessage{}
		if err := json.Unmarshal([]byte(req), &raw); err != nil {
			return nil, err
		}

		isVariadic := v.Type().IsVariadic()
		numIn := v.Type().NumIn()
		if (isVariadic && len(raw) < numIn-1) || (!isVariadic && len(raw) != numIn) {
			return nil, errors.New("function arguments mismatch")
		}
		args := []reflect.Value{}
		for i := range raw {
			var arg reflect.Value
			if isVariadic && i >= numIn-1 {
				arg = reflect.New(v.Type().In(numIn - 1).Elem())
			} else {
				arg = reflect.New(v.Type().In(i))
			}
			if err := json.Unmarshal(raw[i], arg.Interface()); err != nil {
				return nil, err
			}
			args = append(args, arg.Elem())
		}
		errorType := reflect.TypeOf((*error)(nil)).Elem()
		res := v.Call(args)
		switch len(res) {
		case 0:
			// No results from the function, just return nil
			return nil, nil
		case 1:
			// One result may be a value, or an error
			if res[0].Type().Implements(errorType) {
				if res[0].Interface() != nil {
					return nil, res[0].Interface().(error)
				}
				return nil, nil
			}
			return res[0].Interface(), nil
		case 2:
			// Two results: first one is value, second is error
			if !res[1].Type().Implements(errorType) {
				return nil, errors.New("second return value must be an error")
			}
			if res[1].Interface() == nil {
				return res[0].Interface(), nil
			}
			return res[0].Interface(), res[1].Interface().(error)
		default:
			return nil, errors.New("unexpected number of return values")
		}
	}

	m.Lock()
	for ; bindings[index] != nil; index++ {
	}
	bindings[index] = binding
	m.Unlock()
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	C.CgoWebViewBind(w.w, cname, C.uintptr_t(index))
	return nil
}
