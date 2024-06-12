/** @param {boolean} truthy
 * @param {string} message
 */
export function assert(truthy, message) {
  if (!truthy) {
    console.error(message);
    debugger;
  }
}

/**
 * @template {any} T
 * @param {T} maybe
 * @param {string} message
 * @returns {asserts maybe is NotNullable<T>}
 */
export function assertNotNull(maybe, message) {
  if (maybe === null) {
    console.error(message);
    debugger;
  }
}
