package main

import (
	"github.com/lexesv/go-webview-gui"
	"github.com/ncruces/zenity"
)

const html = `<button id="increment">Tap me</button>
<div>You tapped <span id="count">0</span> time(s).</div>
<br>
<button id="openFileDialog">OpenFile</button>
<div>Selected file: <span id="selectedFile"></span></div>
<br>

<script>
  const [incrementElement, countElement, openFileDialogElement, selectedFileElement] =
    document.querySelectorAll("#increment, #count, #openFileDialog, #selectedFile");

  document.addEventListener("DOMContentLoaded", () => {

    incrementElement.addEventListener("click", () => {
      window.increment().then(result => {
        countElement.textContent = result.count;
      });
 	});

	openFileDialogElement.addEventListener("click", () => {
      window.openFileDialog().then(result => {
        selectedFileElement.textContent = result.file;
      });
    });

  });

localStorage.setItem("", "");

</script>`

type Result struct {
	Count uint   `json:"count,omitempty"`
	File  string `json:"file,omitempty"`
}

func main() {
	var count uint = 0
	w := webview.New(true, true)
	defer w.Destroy()

	w.SetTitle("Bind Example")
	w.SetSize(480, 320, webview.HintNone)
	w.Bind("increment", func() Result {
		count++
		return Result{Count: count}
	})
	w.Bind("openFileDialog", func() Result {
		file, err := zenity.SelectFile(zenity.FileFilters{
			{"Go files", []string{"*.go"}, false},
			{"Web files", []string{"*.html", "*.js", "*.css"}, true},
			{"Image files", []string{"*.png", "*.gif", "*.ico", "*.jpg", "*.webp"}, true},
		})
		if err != nil {
			zenity.Error(err.Error())
		}
		return Result{File: file}
	})
	w.SetHtml(html)
	w.Run()
}
