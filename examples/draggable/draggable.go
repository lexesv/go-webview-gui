package main

import (
	"fmt"
	"time"

	"github.com/lexesv/go-webview-gui"
)

const html = `
<div style="width:99%; height: 20px; border: 1px solid blue; padding: 5px; text-align: center;" id="drg">Drag me</div>
<br> <input type="button" value="Close" id="close">

<script>
  const closeEl = document.getElementById("close");
  document.addEventListener("DOMContentLoaded", () => {
    closeEl.addEventListener("click", () => {
      window.closeWindow();
 	});
  });

</script>`

func main() {

	w := webview.New(true, true)
	defer w.Destroy()

	w.SetTitle("Basic Example")
	err := w.SetIcon("../asset/icon.png")
	if err != nil {
		fmt.Println(err.Error())
	}
	w.SetSize(480, 320, webview.HintNone)
	w.SetHtml(html)
	w.Bind("closeWindow", func() {
		w.Terminate()
	})
	w.SetBorderless()
	w.SetDraggable("drg")
	go func() {
		time.Sleep(time.Second * 20)
		w.UnSetDraggable("drg")
	}()
	w.Run()
}
