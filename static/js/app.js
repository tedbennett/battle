import { Stream } from "./ws/util.js"
import { Board } from "./board.js"
/**
* @param el {HTMLElement}
*/
function app(el) {
	const board = new Board()
	const ws = new Stream("ws://localhost:8000/ws", (data) => {
		const msg = JSON.parse(data)
		board.handleMessage(msg)
		board.draw()
	})

}


window.onload = app
