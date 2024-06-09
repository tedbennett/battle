import { Stream } from "./ws/util.js"
import { Board } from "./board.js"
import { assert } from "./utils.js"
/**
* @param el {HTMLElement}
*/
function app(el) {
	const board = new Board()
	const ws = new Stream("ws://localhost:8000/ws", (data) => {
		const msg = parseMessage(data)
		board.handleMessage(msg)
		board.draw()
	})

}

const VERSION = 1

const INIT_MSG = 0
const BOARD_MSG = 1
const PARTIAL_MSG = 2

/** 
	* @param {Uint8Array} buf
	* @returns {import("./types.js").Message}
*/
function parseMessage(buf) {
	assert(buf[0] === VERSION, "Invalid version")
	const type = buf[1]

	const len = read16(buf[2], buf[3])
	switch (type) {
		case INIT_MSG: {
			const colors = {}
			// 4 bytes per team/color pair
			for (let i = 4; i < len + 4; i += 4) {
				const team = buf[i]
				const color = `rgb(${buf[i + 1]}, ${buf[i + 2]}, ${buf[i + 3]})`
				for (let j = i + 1; j < i + 4; ++j) {
					assert(buf[j] <= 255, "parsed invalid color value")
				}
				colors[team] = color
			}
			return {
				type: INIT_MSG,
				colors,
			}
		}
		case BOARD_MSG: {
			let board = []
			// 2 bytes per count/char pair
			for (let i = 4; i < len + 4; i += 2) {
				const count = buf[i]
				const char = buf[i + 1]
				for (let j = i; j < i + 2; ++j) {
					assert(buf[j] <= 255, "parsed invalid rle encoded value")
				}
				board = board.concat(Array(count).fill(char))
			}
			return {
				type: BOARD_MSG,
				board,
			}
		}
	}
}

/**
* @param {number} byteA
* @param {number} byteB
* @returns {number}
*/
function read16(byteA, byteB) {
	return byteA << 8 | byteB
}


window.onload = app
