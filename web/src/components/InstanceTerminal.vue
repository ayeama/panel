<script setup>
import { ref, computed, onMounted, onUnmounted, useTemplateRef, watch } from 'vue'

import { Terminal } from '@xterm/xterm'
import { FitAddon } from '@xterm/addon-fit'

import { API_WS } from '@/api'

const props = defineProps({
  instance: Object,
})

const terminal_element = useTemplateRef('terminal_element')
const terminal = ref(null)
const terminal_fit = ref(null)

let socket = null

const reconnectPending = ref(false)
const status = ref('disconnected')
const status_color = computed(() => {
  switch (status.value) {
    case 'connected':
      return 'text-success border-success'
    case 'connecting':
    case 'disconnecting':
    case 'disconnected':
      return 'text-secondary border-secondary'
  }
})

// const connected = computed(() => {
//   if (status.value == 'connected') {
//     return true
//   }

//   return false
// })

function connect() {
  if (socket != null) {
    return
  }

  // TODO needed?
  if (terminal.value != null) {
    terminal.value.clear()
    terminal.value.focus()
  }

  status.value = 'connecting'

  socket = new WebSocket(`${API_WS}/instances/${props.instance.id}/attach`)

  socket.onopen = () => {
    status.value = 'connected'
  }

  socket.onclose = () => {
    socket = null
    status.value = 'disconnected'

    if (reconnectPending.value) {
      reconnectPending.value = false
      connect()
    }
  }

  socket.onmessage = (e) => {
    terminal.value.write(e.data)
  }

  socket.onerror = () => {}
}

function disconnect() {
  if (socket != null) {
    status.value = 'disconnecting'
    socket.close()
  }
}

function reconnect() {
  if (socket != null) {
    socket.close()
  }

  reconnectPending.value = true
}

function resize() {
  if (terminal_fit.value != null) {
    terminal_fit.value.fit()
  }
}

function create() {
  if (terminal.value == null) {
    terminal.value = new Terminal({ rows: 20 })

    // TODO ignore signals (CTL-C etc)?
    terminal.value.onData((data) => {
      if (socket != null) {
        socket.send(data)
      }
    })

    terminal_fit.value = new FitAddon()
    terminal.value.loadAddon(terminal_fit.value)

    terminal.value.open(terminal_element.value)

    resize()
  }
}

function destroy() {
  disconnect()

  if (terminal.value != null) {
    try {
      terminal.value.dispose()
    } catch (error) {
      if (error.message === 'Could not dispose an addon that has not been loaded') {
        return
      }
      console.error(error)
    } finally {
      terminal.value = null
    }
  }
}

watch(
  () => props.instance.status,
  (status) => {
    status === 'running' ? connect() : disconnect()
  },
  { immediate: true },
)

onMounted(() => {
  window.addEventListener('resize', resize)

  create()
})

onUnmounted(() => {
  window.removeEventListener('resize', resize)

  destroy()
})
</script>

<template>
  <div ref="terminal_element"></div>
</template>
