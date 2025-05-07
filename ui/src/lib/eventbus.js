class EventBus extends EventTarget {
  publish(topic, data) {
    this.dispatchEvent(new CustomEvent(topic, { detail: data }))
  }

  subscribe(topic, callback) {
    this.addEventListener(topic, callback)
  }

  unsubscribe(topic, callback) {
    this.removeEventListener(topic, callback)
  }
}

export const eventBus = new EventBus()

var socket = new WebSocket(
  `${window.CONFIG.api.websocketScheme}://${window.CONFIG.api.websocketHost}${window.CONFIG.api.path}/events`,
)
socket.onmessage = (event) => {
  const data = JSON.parse(event.data)
  eventBus.publish('server:status', data)
}
