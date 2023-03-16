package webview

import (
	"embed"
	_ "embed"
	"fmt"
	"net/url"
	"strings"
)

//go:embed init.js
var EF embed.FS

func (w *webview) initJSFunc() {

	js, _ := EF.ReadFile("init.js")
	w.Init(string(js))

	w.Bind("open", func(s string) {
		if !strings.Contains(s, "http") {
			u, err := url.Parse(w.GetUrl())
			if err == nil {
				s = fmt.Sprintf("%s://%s/%s", u.Scheme, u.Host, s)
			}
		}
		w.Navigate(s)
	})

	w.Bind("getUrl", func(s string) {
		w.data["url"] = s
	})

	w.Bind("getHtml", func(s string) {
		w.data["html"] = s
	})

	w.Bind("getPageTitle", func(s string) {
		w.data["title"] = s
	})

	w.Bind("move", func(x, y int) {
		w.Move(x, y)
	})

	w.Bind("contentState", func(s string) {
		if events.handle_cs != nil {
			for _, f := range events.handle_cs {
				f(s)
			}
		}
	})

	type DResult struct {
		Id string `json:"id,omitempty"`
		V  bool   `json:"v,omitempty"`
	}

	w.Bind("getDraggebleData", func() (res []DResult) {
		w.DraggableElements.Range(func(key, value any) bool {
			r := DResult{Id: key.(string), V: value.(bool)}
			res = append(res, r)
			return true
		})
		return res
	})

	w.Bind("delDraggebleElement", func(id string) {
		w.DraggableElements.Delete(id)
	})

	w.Bind("getDraggebleElementValue", func(id string) bool {
		if v, ok := w.DraggableElements.Load(id); !ok {
			return false
		} else {
			return v.(bool)
		}
	})

	//w.Bind("setCookie")

}
