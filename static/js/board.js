import { assert } from "./utils.js";

export const INIT_MSG = 0;
export const BOARD_MSG = 1;
export const PARTIAL_MSG = 2;

export class Board {
  /** @type {import("./types").BoardMetadata} */
  #metadata;

  /** @type {number[]} */
  #squares;

  /** @type {number} */
  #size;

  /** @type {HTMLDivElement} */
  #element;

  constructor() {
    const parent = /** @type{HTMLDivElement} */ (
      document.getElementById("board")
    );
    assert(parent !== null, "Failed to find board element");
    assert(parent?.children.length !== 0, "Board element is empty");
    this.#element = parent;
    this.#metadata = {
      colors: {},
    };
    this.#size = parent?.children.length ?? 0;
    this.#squares = [];
  }

  draw() {
    if (!this.#squares) return;
    const parent = this.#element;
    // Otherwise, diff the board and update only those that have changed
    for (let i = 0; i < parent.children.length; ++i) {
      const row = parent.children[i];
      for (let j = 0; j < row.children.length; ++j) {
        const div = /** @type{HTMLDivElement} */ (row.children[j]);
        const square = this.#squares[this.#index(i, j)];
        const color = this.#metadata.colors[square];
        if (div.style.backgroundColor !== color) {
          div.style.backgroundColor = color;
        }
      }
    }
  }

  redrawBoard() {
    if (!this.#squares) return;
    const parent = /** @type{HTMLDivElement} */ (
      document.getElementById("board")
    );
    assert(parent !== null, "Failed to find board element");

    if (parent.children.length === 0) {
      this.#initializeBoard(parent);
      return;
    }
    // Otherwise, diff the board and update only those that have changed
    for (let i = 0; i < parent.children.length; ++i) {
      const row = parent.children[i];
      for (let j = 0; j < row.children.length; ++j) {
        const div = /** @type{HTMLDivElement} */ (row.children[j]);
        const square = this.#squares[this.#index(i, j)];
        const color = this.#metadata.colors[square];
        if (div.style.backgroundColor !== color) {
          div.style.backgroundColor = color;
        }
      }
    }
  }

  /** @param {HTMLElement} parent */
  #initializeBoard(parent) {
    /** @returns {HTMLDivElement} */
    const newRow = () => {
      const row = document.createElement("div");
      row.style.flexGrow = "1";
      row.style.display = "flex";
      row.style.flexDirection = "row";
      return row;
    };

    /** @type {HTMLDivElement[]} */
    const elements = [];

    /** @type {HTMLDivElement} */
    let rowEl;
    for (let i = 0; i < this.#squares.length; ++i) {
      // initialize rowEl
      if (i % this.#size === 0) {
        rowEl = newRow();
        elements.push(rowEl);
      }

      const square = this.#squares[i];
      const div = document.createElement("div");
      const color = this.#metadata.colors[square];
      div.style.backgroundColor = color;
      div.style.flexGrow = "1";
      // @ts-expect-error
      rowEl.appendChild(div);
    }
    parent.replaceChildren(...elements);
  }

  /** @param {import("./types").Message} msg */
  handleMessage(msg) {
    switch (msg.type) {
      case INIT_MSG: {
        if (!msg.colors || !msg.board) break;
        this.#onInit(msg.colors, msg.board);
        break;
      }
      case BOARD_MSG: {
        // Pass
        break;
      }
      case PARTIAL_MSG: {
        if (!msg.diffs) break;
        this.#onDiffsReceived(msg.diffs);
        break;
      }
    }
  }

  /**
   * @param {number} i
   * @param {number} j
   * @returns {number}
   */
  #index(i, j) {
    return this.#size * i + j;
  }

  /**
   * @param {Object.<number, string>} colors
   * @param {number[]} squares
   * */
  #onInit(colors, squares) {
    this.#metadata.colors = colors;
    this.#size = Math.sqrt(squares.length);
    this.#squares = squares;
  }

  /** @param {import('./types.js').Diff[]} diffs */
  #onDiffsReceived(diffs) {
    if (!this.#squares) return;
    for (const diff of diffs) {
      this.#squares[this.#index(diff.row, diff.col)] = diff.team;
    }
  }
}
