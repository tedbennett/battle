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
		const parent = document.getElementById("board")
		assert(parent !== null, "Failed to find board element")

		/** @type {HTMLDivElement} */
		const elements = []

		for (const row of this.#squares) {
			const rowEl = document.createElement("div");
			rowEl.style.flexGrow = 1;
			for (const square of row) {
				const div = document.createElement("div");
				div.style.backgroundColor = this.#metadata.colors[square];
				div.style.flexGrow = 1;
				div.appendChild(div)
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
			}
			case BOARD_MSG: {
				this.onBoardChange(msg.board)
			}
		}

	}
}
