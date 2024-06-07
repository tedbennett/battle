/** @param {boolean} truthy
  * @param {string} message
*/
export const assert = (truthy, message) => {
	if (!truthy) {
		console.error(message)
		throw new Error(message)
	}

}
