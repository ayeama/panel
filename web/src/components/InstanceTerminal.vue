<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

import { Terminal } from '@xterm/xterm'
import { AttachAddon } from '@xterm/addon-attach'
import { FitAddon } from '@xterm/addon-fit'

import { API_WS } from '@/api'

const props = defineProps({
  id: String,
})

const terminal_element = ref(null)

let terminal = null
let terminal_fit = null
let terminal_attach = null

let socket_attach = null

function resize() {
  if (terminal_fit) {
    terminal_fit.fit()
  }
}

onMounted(() => {
  window.addEventListener('resize', resize)

  terminal = new Terminal({rows: 18})
  terminal_fit = new FitAddon()
  terminal.loadAddon(terminal_fit)
  
  socket_attach = new WebSocket(`${API_WS}/instances/${props.id}/attach`)
  socket_attach.onopen = () => {
    terminal_attach = new AttachAddon(socket_attach)
    terminal.loadAddon(terminal_attach)

    terminal.open(terminal_element.value)
    terminal_fit.fit()
  }
})

onUnmounted(() => {
  if (socket_attach) {
    socket_attach.close()
    socket_attach = null
  }

  if (terminal_attach) {
    terminal_attach.dispose()
    terminal_attach = null
  }

  if (terminal_fit) {
    terminal_fit.dispose()
    terminal_fit = null
  }

  if (terminal) {
    terminal.dispose()
    terminal = null
  }
})
</script>

<template>
  <div ref="terminal_element"></div>
</template>
