import { assert } from "./utils.js"

const INIT_MSG = 0
const BOARD_MSG = 1

export class Board {
	/** @type {import("./types").BoardMetadata} */
	#metadata

	/** @type {number[]} */
	#squares

	/** @type {number} */
	#size

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
				const square = this.#squares[this.#index(i, j)]
				const color = this.#metadata.colors[square]
				if (div.style.backgroundColor !== color) {
					div.style.backgroundColor = color
				}
			}
		}
	}

	/** @param {HTMLElement} parent */
	#initializeBoard(parent) {
		/** @returns {HTMLDivElement} */
		const newRow = () => {
			const row = document.createElement("div");
			row.style.flexGrow = 1;
			row.style.display = "flex";
			row.style.flexDirection = "row";
			return row;
		}


		/** @type {HTMLDivElement[]} */
		const elements = []

		/** @type {HTMLDivElement} */
		let rowEl;
		for (let i = 0; i < this.#squares.length; ++i) {
			// initialize rowEl
			if (i % this.#size === 0) {
				rowEl = newRow();
				elements.push(rowEl)
			}

			const square = this.#squares[i]
			const div = document.createElement("div")
			const color = this.#metadata.colors[square]
			div.style.backgroundColor = color
			div.style.flexGrow = 1
			rowEl.appendChild(div)
		}
		parent.replaceChildren(...elements)
	}

	/** @param {import("./types").Message} msg */
	handleMessage(msg) {
		switch (msg.type) {
			case INIT_MSG: {
				this.#onColorsChange(msg.colors)
			}
			case BOARD_MSG: {
				this.#onBoardChange(msg.board)
			}
		}

	}

	/** 
	 * @param {number} i
	 * @param {number} j
	 * @returns {number}
	 */
	#index(i, j) {
		return (this.#size * i) + j
	}

	/** @param {Object.<number, string>} */
	#onColorsChange(colors) {
		this.#metadata.colors = colors
	}

	/** @param {number[][]} squares */
	#onBoardChange(squares) {
		this.#size = Math.sqrt(squares.length)
		this.#squares = squares
	}

}
