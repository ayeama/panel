<script setup>
import { ref, watch, onBeforeUnmount } from 'vue'

const props = defineProps({
  server_id: {
    type: String,
    required: true,
  },
  command: {
    type: String,
    default: 'disconnect',
  },
})

const emit = defineEmits(['status'])

let socket = ref(null)
const status = ref({ message: 'disconnected', class: 'text-bg-secondary' })
const input = ref('')
const output = ref('')
const textarea = ref(null)
const follow = ref(true)

// onMounted(() => {
//   updateStatus()
//   create()
// })

onBeforeUnmount(() => {
  disconnect()
})

function connect() {
  if (!props.server_id) {
    return
  }

  socket.value = new WebSocket(
    `${window.CONFIG.api.websocketScheme}://${window.CONFIG.api.websocketHost}${window.CONFIG.api.path}/servers/${props.server_id}/attach`,
  )
  socket.value.onopen = () => {
    updateStatus()
  }
  socket.value.onclose = () => {
    disconnect()
  }
  socket.value.onerror = () => {
    updateStatus()
  }
  socket.value.onmessage = (event) => {
    output.value += event.data
    if (follow.value) {
      requestAnimationFrame(() => {
        textarea.value.scrollTop = textarea.value.scrollHeight
      })
    }
  }
}

function disconnect() {
  if (socket.value) {
    socket.value.close()
    socket.value = null
    updateStatus()
  }
}

async function send() {
  if (socket.value) {
    socket.value.send(input.value + '\n')
    input.value = ''
  }
}

function updateStatus() {
  if (!socket.value) {
    status.value = { message: 'disconnected', class: 'text-bg-secondary' }
    emit('status', status.value.message)
    return
  }

  switch (socket.value.readyState) {
    case WebSocket.CONNECTING:
      status.value = { message: 'connecting', class: 'text-bg-secondary' }
      break
    case WebSocket.OPEN:
      status.value = { message: 'connected', class: 'text-bg-success' }
      break
    case WebSocket.CLOSING:
      status.value = { message: 'disconnecting', class: 'text-bg-secondary' }
      break
    case WebSocket.CLOSED:
      status.value = { message: 'disconnected', class: 'text-bg-secondary' }
      break
    default:
      status.value = { message: 'unknown', class: 'text-bg-warning' }
  }
  emit('status', status.value.message)
}

function toggle() {
  if (socket.value) {
    disconnect()
  } else {
    connect()
  }
  updateStatus()
}

watch(
  () => props.command,
  (cmd) => {
    if (cmd === 'connect') {
      connect()
    } else if (cmd === 'disconnect') {
      disconnect()
    }
  },
  { immediate: true },
)
</script>

<template>
  <div>
    <div class="card">
      <div
        class="card-header"
        data-bs-toggle="collapse"
        data-bs-target="#card-body"
        style="cursor: pointer"
      >
        <div class="row">
          <div class="col">
            <h1 class="fs-5 mb-0">Terminal</h1>
          </div>

          <div class="col">
            <div class="float-end">
              <span
                class="badge"
                v-bind:class="status.class"
                v-on:click="toggle"
                style="cursor: pointer"
              >
                {{ status.message }}
                <span class="visually-hidden">terminal connection status</span></span
              >
            </div>
          </div>
        </div>
      </div>

      <div class="card-body p-0">
        <form @submit.prevent="send">
          <div class="row">
            <div class="col">
              <textarea
                ref="textarea"
                v-model="output"
                class="form-control px-1 border-0 rounded-0"
                rows="18"
                style="resize: none; font-family: monospace"
                disabled
                readonly
              ></textarea>
            </div>
          </div>

          <div class="row">
            <div class="col">
              <div class="input-group">
                <span class="input-group-text rounded-top-0 border-0 border-top border-right"
                  >$</span
                >

                <input
                  v-model="input"
                  class="form-control px-2 border-0 border-top border-left border-right"
                  style="font-family: monospace"
                />
                <button
                  type="submit"
                  class="btn btn-secondary border-0 border-top border-left border-right"
                >
                  send
                </button>
                <a
                  class="btn btn-secondary rounded-top-0 border-0 border-top border-left"
                  v-bind:class="{ active: follow }"
                  v-bind:aria-disabled="!follow"
                  v-on:click="follow = !follow"
                  >follow</a
                >
              </div>
            </div>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>
