/** @param {boolean} truthy
  * @param {string} message
*/
export function assert(truthy, message) {
	if (!truthy) {
		console.error(message)
		debugger
	}

}
