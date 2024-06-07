const StreamState = {
	connected: 0,
	connecting: 1,
	error: 2,
	closed: 3,
}

export class Stream {
	state = StreamState.connecting
	/** @type{string} */
	#url
	/** @type {(any) => void} */
	#handleMessage

	/** 
	 * @param {string} url 
	 * @param {(any) => void} handleMessage
	 * */
	constructor(url, handleMessage) {
		this.#url = url
		this.#handleMessage = handleMessage

		this.connect()
	}

	connect() {
		this.ws = new WebSocket(this.#url)
		this.state = StreamState.connecting
		this.ws.onopen = () => {
			this.state = StreamState.connected
		}

		this.ws.onclose = () => {
			this.state = StreamState.closed
		}

		this.ws.onmessage = (event) => {
			this.#handleMessage(event.data)
		}

		this.ws.onerror = () => {
			console.error("Websocket connection closed")
			this.state = StreamState.error
			this.connect()
		}
	}
}
