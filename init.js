document.onreadystatechange = function () {
    contentState(document.readyState);
}

document.addEventListener('DOMContentLoaded', () => {
    window.getHtml(window.document.getElementsByTagName('html')[0].innerHTML);
    window.getUrl(window.location.href);
    window.getPageTitle(window.document.title);

    Array.from(document.querySelectorAll('a[target="_blank"]'))
        .forEach(link => link.removeAttribute('target'));

    window.getDraggebleData().then(result => {
        if (result == null) return
        result.forEach(function (e) {
            if (e.v === true) {
                setDraggableRegion(e.id).then(result => {
                    console.log(result);
                }, error => {
                    console.log(error);
                });
            }
            let el = document.getElementById(e.id);

            let interval = setInterval(function () {
                window.getDraggebleElementValue(e.id).then(result => {
                    if (result === false) {
                        unsetDraggableRegion(e.id).then(result => {
                            window.delDraggebleElement(e.id);
                            console.log(result);
                        }, error => {
                            clearInterval(interval);
                            console.log(error, e.id);
                        });
                    }
                });
            }, 250);

            /*function DraggableListener(event) {
                window.getDraggebleElementValue(e.id).then(result => {
                    if (result === false) {
                        unsetDraggableRegion(e.id).then(result => {
                            window.delDraggebleElement(e.id);
                            el.removeEventListener("mouseover", DraggableListener, false);
                            console.log(result);
                        }, error => {
                            console.log(error, e.id);
                        });
                    }
                });
            }
            el.addEventListener("mouseover", DraggableListener);*/
        });
    });

});

document.addEventListener('keydown', (event) => {

        // Ctrl+C or Cmd+C pressed
        if ((event.ctrlKey || event.metaKey) && event.keyCode === 67) {
            console.log("copy");
            document.execCommand("copy");
        }

        // Ctrl+V or Cmd+V pressed
        if ((event.ctrlKey || event.metaKey) && event.keyCode === 86) {
            console.log("paste");
            document.execCommand("paste");
        }

        // Ctrl+X or Cmd+X pressed
        if ((event.ctrlKey || event.metaKey) && event.keyCode === 88) {
            console.log("cut");
            document.execCommand("cut");
        }


});

function __awaiter(thisArg, _arguments, P, generator) {
    function adopt(value) {
        return value instanceof P ? value : new P(function (resolve) {
            resolve(value);
        });
    }

    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) {
            try {
                step(generator.next(value));
            } catch (e) {
                reject(e);
            }
        }

        function rejected(value) {
            try {
                step(generator["throw"](value));
            } catch (e) {
                reject(e);
            }
        }

        function step(result) {
            result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected);
        }

        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
}

const draggableRegions = new WeakMap();

function setDraggableRegion(domElementOrId) {
    return new Promise((resolve, reject) => {

        const draggableRegion = domElementOrId instanceof Element ? domElementOrId : document.getElementById(domElementOrId);
        let initialClientX = 0;
        let initialClientY = 0;
        if (!draggableRegion) {
            return reject('Unable to find DOM element');
        }
        if (draggableRegions.has(draggableRegion)) {
            return reject('This DOM element is already an active draggable region');
        }
        draggableRegion.addEventListener('pointerdown', startPointerCapturing);
        draggableRegion.addEventListener('pointerup', endPointerCapturing);
        draggableRegions.set(draggableRegion, {pointerdown: startPointerCapturing, pointerup: endPointerCapturing});

        function onPointerMove(evt) {
            return __awaiter(this, void 0, void 0, function* () {
                yield move(evt.screenX - initialClientX, evt.screenY - initialClientY);
            });
        }

        function startPointerCapturing(evt) {
            if (evt.button !== 0) {
                return;
            }
            initialClientX = evt.clientX;
            initialClientY = evt.clientY;
            draggableRegion.addEventListener('pointermove', onPointerMove);
            draggableRegion.setPointerCapture(evt.pointerId);
        }

        function endPointerCapturing(evt) {
            draggableRegion.removeEventListener('pointermove', onPointerMove);
            draggableRegion.releasePointerCapture(evt.pointerId);
        }

        resolve('Draggable region was activated');
    });
}

function unsetDraggableRegion(domElementOrId) {
    return new Promise((resolve, reject) => {
        const draggableRegion = domElementOrId instanceof Element ? domElementOrId : document.getElementById(domElementOrId);
        if (!draggableRegion) {
            return reject('Unable to find DOM element');
        }
        if (!draggableRegions.has(draggableRegion)) {
            return reject('DOM element is not an active draggable region');
        }
        const {pointerdown, pointerup} = draggableRegions.get(draggableRegion);
        draggableRegion.removeEventListener('pointerdown', pointerdown);
        draggableRegion.removeEventListener('pointerup', pointerup);
        draggableRegions.delete(draggableRegion);
        resolve('Draggable region was deactivated');
    });
}

/*
const {setItem} = localStorage.__proto__;

localStorage.setItem = function (key, value) {
    console.log('localStorage.__proto__.setItem', key, value);
    //setItem.call(this, key, value);
};*/

/*
var localStorage = (function (exports) {
    'use strict';

    function setItem(k, v) {
        //
    }

    exports = {
        __proto__: null,
        setItem: setItem,
    };

    return exports;
})({});*/
