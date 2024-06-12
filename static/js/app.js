import { Stream } from "./ws/util.js";
import { Board, INIT_MSG, BOARD_MSG, PARTIAL_MSG } from "./board.js";
import { assert } from "./utils.js";

function app() {
  const board = new Board();
  new Stream("ws://localhost:8000/ws", (data) => {
    const msg = parseMessage(data);
    board.handleMessage(msg);
    board.draw();
  });
}

const VERSION = 1;

/**
 * @param {Uint8Array} buf
 * @returns {import("./types.js").Message}
 */
function parseMessage(buf) {
  assert(buf[0] === VERSION, "Invalid version");
  const type = buf[1];

  const len = read16(buf[2], buf[3]);
  switch (type) {
    case INIT_MSG: {
      /** @type{Object.<number, string>} */
      const colors = {};
      // 4 bytes per team/color pair
      for (let i = 4; i < len + 4; i += 4) {
        const team = buf[i];
        const color = `rgb(${buf[i + 1]}, ${buf[i + 2]}, ${buf[i + 3]})`;
        for (let j = i + 1; j < i + 4; ++j) {
          assert(buf[j] <= 255, "parsed invalid color value");
        }
        colors[team] = color;
      }
      const offset = 4 + len;
      const boardLen = read16(buf[offset], buf[offset + 1]);
      /** @type{number[][]} */
      let board = [];
      // 2 bytes per count/char pair
      for (let i = offset + 2; i < offset + boardLen + 2; i += 2) {
        const count = buf[i];
        const char = buf[i + 1];
        for (let j = i; j < i + 2; ++j) {
          assert(buf[j] <= 255, "parsed invalid rle encoded value");
        }
        board = board.concat(Array(count).fill(char));
      }
      return {
        type: INIT_MSG,
        colors,
        board,
      };
    }
    case BOARD_MSG: {
      /** @type{number[][]} */
      let board = [];
      // 2 bytes per count/char pair
      for (let i = 4; i < len + 4; i += 2) {
        const count = buf[i];
        const char = buf[i + 1];
        for (let j = i; j < i + 2; ++j) {
          assert(buf[j] <= 255, "parsed invalid rle encoded value");
        }
        board = board.concat(Array(count).fill(char));
      }
      return {
        type: BOARD_MSG,
        board,
      };
    }
    case PARTIAL_MSG: {
      /** @returns {import("./types.js").Diff[]} */
      let diffs = [];
      // 3 bytes per row/col/team set
      for (let i = 4; i < len + 4; i += 3) {
        const row = buf[i];
        const col = buf[i + 1];
        const team = buf[i + 2];
        for (let j = i; j < i + 3; ++j) {
          assert(buf[j] <= 255, "parsed invalid diff value");
        }
        diffs.push({ row, col, team });
      }
      return {
        type: PARTIAL_MSG,
        diffs,
      };
    }
  }
  throw new Error("Invalid message type received");
}

/**
 * @param {number} byteA
 * @param {number} byteB
 * @returns {number}
 */
function read16(byteA, byteB) {
  return (byteA << 8) | byteB;
}

window.onload = () => app();
