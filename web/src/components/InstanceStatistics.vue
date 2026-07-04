<script setup>
import { onMounted, ref } from 'vue';

import { API_WS } from '@/api';

const props = defineProps({
  id: String,
})

const cpu = ref(null)
const memory = ref(null)
const disk = ref(null)

onMounted(() => {
  const socket = new WebSocket(`${API_WS}/instances/${props.id}/stats`)
  socket.onmessage = (e) => {
    try {
      const data = JSON.parse(e.data)
      cpu.value = data.cpu
      memory.value = data.memory
      disk.value = data.disk
    } catch (error) {
      console.error(error)
    }
  }
})

function progressbar_width(v) {
  return `${v}%`
}

function progressbar_color(v) {
  // if (v >= 70) {
  //   return "text-bg-warning"
  // } else if (v >= 85) {
  //   return "text-bg-danger"
  // }
  return "text-bg-secondary"
}
</script>

<template>
  <div class="row"><div class="col">
      <div class="progress" role="progressbar" aria-label="instance cpu" :aria-valuenow="cpu" aria-valuemin="0" aria-valuemax="100">
        <div class="progress-bar" :class="progressbar_color(cpu)" :style="{ width: progressbar_width(cpu) }"></div>
      </div>
    </div>
  </div>
  
  <div class="row"><div class="col">
      <div class="progress" role="progressbar" aria-label="instance memory" :aria-valuenow="memory" aria-valuemin="0" aria-valuemax="100">
        <div class="progress-bar" :class="progressbar_color(memory)" :style="{ width: progressbar_width(memory) }"></div>
      </div>
    </div>
  </div>

  <div class="row"><div class="col">
      <div class="progress" role="progressbar" aria-label="instance disk" :aria-valuenow="disk" aria-valuemin="0" aria-valuemax="100">
        <div class="progress-bar" :class="progressbar_color(disk)" :style="{ width: progressbar_width(disk) }"></div>
      </div>
    </div>
  </div>
</template>
