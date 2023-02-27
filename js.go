package webview

import (
	"embed"
	_ "embed"
	"strings"

	"github.com/lexesv/go-webview-gui/dialog"
)

//go:embed init.js
var EF embed.FS

func (w *webview) initJSFunc() {

	js, _ := EF.ReadFile("init.js")
	w.Init(string(js))

	w.Bind("_alert", func(s ...string) {
		v := strings.Join(s, "")
		dialog.Message("%s", v).Info()
	})

	w.Bind("confirm", func(s interface{}) bool {
		yn := dialog.Message("%v", s).YesNo()
		return yn
	})

	w.Bind("move", func(x, y int) {
		w.Move(x, y)
	})

	w.Bind("getHtml", func(s string) {
		w.Html = s
	})

	w.Bind("contentState", func(s string) {
		w.ContentState = s
		if events.handle_cs != nil {
			events.handle_cs(s)
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

}

/*var js = `
'use strict';

document.onreadystatechange = function () {
	contentState(document.readyState);
}

function __awaiter(thisArg, _arguments, P, generator) {
        function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
        return new (P || (P = Promise))(function (resolve, reject) {
            function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
            function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
            function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
            step((generator = generator.apply(thisArg, _arguments || [])).next());
        });
}

const draggableRegions = new WeakMap();

function setDraggableRegion(domElementOrId) {
    return new Promise((resolve, reject) => {
        const draggableRegion = domElementOrId instanceof Element ?
            domElementOrId : document.getElementById(domElementOrId);
        let initialClientX = 0;
        let initialClientY = 0;
        if (!draggableRegion) {
            return reject();
        }
        if (draggableRegions.has(draggableRegion)) {
            return reject();
        }
        draggableRegion.addEventListener('pointerdown', startPointerCapturing);
        draggableRegion.addEventListener('pointerup', endPointerCapturing);
        draggableRegions.set(draggableRegion, { pointerdown: startPointerCapturing, pointerup: endPointerCapturing });
        function onPointerMove(evt) {
            return __awaiter(this, void 0, void 0, function* () {
                yield move(evt.screenX - initialClientX, evt.screenY - initialClientY);
            });
        }
        function startPointerCapturing(evt) {
            if (evt.button !== 0)
                return;
            initialClientX = evt.clientX;
            initialClientY = evt.clientY;
            draggableRegion.addEventListener('pointermove', onPointerMove);
            draggableRegion.setPointerCapture(evt.pointerId);
        }
        function endPointerCapturing(evt) {
            draggableRegion.removeEventListener('pointermove', onPointerMove);
            draggableRegion.releasePointerCapture(evt.pointerId);
        }
        resolve();
    });
}

function unsetDraggableRegion(domElementOrId) {
    return new Promise((resolve, reject) => {
        const draggableRegion = domElementOrId instanceof Element ?
            domElementOrId : document.getElementById(domElementOrId);
        if (!draggableRegion) {
            return reject();
        }
        if (!draggableRegions.has(draggableRegion)) {
            return reject();
        }
        const { pointerdown, pointerup } = draggableRegions.get(draggableRegion);
        draggableRegion.removeEventListener('pointerdown', pointerdown);
        draggableRegion.removeEventListener('pointerup', pointerup);
        draggableRegions.delete(draggableRegion);
        resolve();
    });
}
`*/
