import { assert } from "./utils.js"

const INIT_MSG = 0
const BOARD_MSG = 1

export class Board {
	/** @type {import("./types").BoardMetadata} */
	#metadata

	/** @type {int[][]} */
	#squares

	constructor() {
		this.#metadata = {
			colors: {}
		}
	}

	draw() {
		if (!this.#squares) return;
		const parent = document.getElementById("board")
		assert(parent !== null, "Failed to find board element")

		if (parent.children.length === 0) {
			this.#initializeBoard(parent)
			return
		}
		// Otherwise, diff the board and update only those that have changed
		for (let i = 0; i < parent.children.length; ++i) {
			const row = parent.children[i]
			for (let j = 0; j < row.children.length; ++j) {
				const div = row.children[j]
				const square = this.#squares[i][j]
				const color = this.#metadata.colors[square]
				if (div.style.backgroundColor !== color) {
					div.style.backgroundColor = color
				}
			}
		}
	}

	/** @param {HTMLElement} parent */
	#initializeBoard(parent) {
		/** @type {HTMLDivElement} */
		const elements = []

		for (const row of this.#squares) {
			const rowEl = document.createElement("div");
			rowEl.style.flexGrow = 1;
			rowEl.style.display = "flex";
			rowEl.style.flexDirection = "row";
			for (const square of row) {
				const div = document.createElement("div")
				const color = this.#metadata.colors[square]
				div.style.backgroundColor = color
				div.style.flexGrow = 1
				rowEl.appendChild(div)
			}
			elements.push(rowEl)
		}
		parent.replaceChildren(...elements)
	}

	/** @param {Object.<number, string>} */
	onColorsChange(colors) {
		this.#metadata.colors = colors
	}

	/** @param {number[][]} squares */
	onBoardChange(squares) {
		this.#squares = squares
	}

	/** @param {import("./types").Message} msg */
	handleMessage(msg) {
		switch (msg.type) {
			case INIT_MSG: {
				this.onColorsChange(msg.colors)
				this.onBoardChange(msg.board)
			}
			case BOARD_MSG: {
				this.onBoardChange(msg.board)
			}
		}

	}
}
