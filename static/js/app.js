import { Stream } from "./ws/util.js"
/**
* @param el {HTMLElement}
*/
function app(el) {
	const board = new Board()
	const ws = new Stream("ws://localhost:8000/ws", (data) => {
		console.log("Received message: ", data)
	})

}

class Board {
	/** @type {import("./types").BoardMetadata} */
	metadata

	/** @type {int[][]} */
	squares

	constructor() {
		this.metadata = {
			colors: {}
		}
	}

	to

	draw() {
		const parent = document.getElementById("board")
		assert(parent !== null, "Failed to find board element")
	}
}

window.onload = app
