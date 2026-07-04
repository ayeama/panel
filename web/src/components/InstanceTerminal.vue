<script setup>
import { onMounted } from 'vue'

import { Terminal } from '@xterm/xterm'
import { AttachAddon } from '@xterm/addon-attach'
import { FitAddon } from '@xterm/addon-fit'

import { API_WS } from '@/api'

const props = defineProps({
  id: String,
})

onMounted(() => {
  const terminal = new Terminal()
  
  const socket = new WebSocket(`${API_WS}/instances/${props.id}/attach`)
  socket.onopen = () => {
    const fitAddon = new FitAddon()
    terminal.loadAddon(fitAddon)

    const attachAddon = new AttachAddon(socket)
    terminal.loadAddon(attachAddon)

    terminal.open(document.getElementById('terminal'))
    fitAddon.fit()
  }
})
</script>

<template>
  <div id="terminal"></div>
</template>
