package webview

import (
	"github.com/lexesv/go-webview-gui/dialog"
)

func (w *webview) initJSFunc() {
	w.Init(js)

	w.Bind("alert", func(s interface{}) {
		dialog.Message("%v", s).Info()
	})
	w.Bind("confirm", func(s interface{}) bool {
		yn := dialog.Message("%v", s).YesNo()
		return yn
	})
	w.Bind("move", func(x, y int) {
		w.Move(x, y)
	})

	w.Bind("contentState", func(s string) {
		w.ContentState = s
	})
}

var js = `
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
`
