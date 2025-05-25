<script setup>
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'

const props = defineProps({
  url: {
    required: true,
  },
})

let socket = ref(null)
let cpu = ref(0.0)
let memory = ref(0.0)
// let disk = ref(12.0)

onMounted(() => {
  create()
})

onBeforeUnmount(() => {
  if (socket.value) {
    socket.value.close()
  }
})

function create() {
  socket.value = new WebSocket(props.url)
  socket.value.onopen = () => {}
  socket.value.onclose = () => {}
  socket.value.onerror = () => {}
  socket.value.onmessage = (event) => {
    const data = JSON.parse(event.data)
    cpu.value = data.cpu.toFixed(2)
    memory.value = data.memory.toFixed(2)
  }
}
</script>

<template>
  <div>
    <div class="row">
      <div class="col-2">
        <p>cpu</p>
      </div>
      <div class="col">
        <div
          class="progress"
          role="progressbar"
          v-bind:aria-valuenow="cpu"
          aria-valuemin="0"
          aria-valuemax="100"
          style="height: 25px"
        >
          <div
            class="progress-bar"
            v-bind:class="{
              'bg-success': cpu < 75,
              'bg-warning': cpu >= 75 && cpu < 85,
              'bg-danger': cpu >= 85,
            }"
            v-bind:style="{ width: cpu + '%' }"
          >
            {{ cpu }}%
          </div>
        </div>
      </div>
    </div>

    <div class="row">
      <div class="col-2">
        <p>memory</p>
      </div>
      <div class="col">
        <div
          class="progress"
          role="progressbar"
          v-bind:aria-valuenow="memory"
          aria-valuemin="0"
          aria-valuemax="100"
          style="height: 25px"
        >
          <div
            class="progress-bar"
            v-bind:class="{
              'bg-success': memory < 75,
              'bg-warning': memory >= 75 && memory < 85,
              'bg-danger': memory >= 85,
            }"
            v-bind:style="{ width: memory + '%' }"
          >
            {{ memory }}%
          </div>
        </div>
      </div>
    </div>

    <!-- <div class="row">
      <div class="col-2">
        <p>disk</p>
      </div>
      <div class="col">
        <div class="progress" role="progressbar" v-bind:aria-valuenow="disk" aria-valuemin="0" aria-valuemax="100"
          style="height: 25px">
          <div class="progress-bar" v-bind:class="{
            'bg-success': disk < 75,
            'bg-warning': disk >= 75 && disk < 85,
            'bg-danger': disk >= 85,
          }" v-bind:style="{width: disk+'%'}">{{ disk }}%</div>
        </div>
      </div>
    </div> -->
  </div>
</template>
