<script setup>
import { onMounted, onUnmounted, ref } from 'vue';

import { API_WS } from '@/api';

const props = defineProps({
  id: String,
})

const cpu = ref(null)
const memory = ref(null)
const disk = ref(null)

let socket_stats = null

onMounted(() => {
  socket_stats = new WebSocket(`${API_WS}/instances/${props.id}/stats`)
  socket_stats.onmessage = (e) => {
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

onUnmounted(() => {
  if (socket_stats) {
    socket_stats.close()
    socket_stats = null
  }
})

function progressbar_width(v) {
  if (v) {
    return `${v}%`
  }
  return '0%'
}

function progressbar_color(v) {
  if (v >= 85) {
    return "text-bg-danger"
  } else if (v >= 70) {
    return "text-bg-warning"
  }
  return "text-bg-primary"
}
</script>

<template>
  <div class="row">
    <div class="col col-12">
      <span id="cpuProgressbar" class="form-text">CPU</span>

      <div class="progress" role="progressbar" aria-label="instance cpu" :aria-valuenow="cpu" aria-valuemin="0" aria-valuemax="100" aria-describedby="cpuProgressbar" style="height: 1rem">
        <div class="progress-bar" :class="progressbar_color(cpu)" :style="{ width: progressbar_width(cpu) }"></div>
      </div>
    </div>

    <div class="col col-12">
      <span id="memoryProgressbar" class="form-text">Memory</span>

      <div class="progress" role="progressbar" aria-label="instance memory" :aria-valuenow="memory" aria-valuemin="0" aria-valuemax="100" aria-describedby="memoryProgressbar" style="height: 1rem">
        <div class="progress-bar" :class="progressbar_color(memory)" :style="{ width: progressbar_width(memory) }"></div>
      </div>
    </div>

    <div class="col col-12">
      <span id="diskProgressbar" class="form-text">Disk</span>

      <div class="progress" role="progressbar" aria-label="instance disk" :aria-valuenow="disk" aria-valuemin="0" aria-valuemax="100" aria-describedby="diskProgressbar" style="height: 1rem">
        <div class="progress-bar" :class="progressbar_color(disk)" :style="{ width: progressbar_width(disk) }"></div>
      </div>
    </div>
  </div>
</template>
