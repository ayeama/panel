<script setup>
import { ref, computed } from 'vue'
import { eventBus } from '@/lib/eventbus.js'

const props = defineProps({
  server_id: {
    type: String,
    required: false,
  },
  status: {
    type: String,
    required: false,
    default: 'unknown',
    validator: (value) => {
      return [
        'configured',
        'created',
        'dead',
        'exited',
        'healthy',
        'initialized',
        'paused',
        'removing',
        'running',
        'stopped',
        'stopping',
        'unhealthy',
        'unknown',
      ].includes(value)
    },
  },
})

const emit = defineEmits(['status'])

const statusClass = {
  created: 'text-bg-secondary',
  running: 'text-bg-success',
  stopped: 'text-bg-secondary',
  error: 'text-bg-danger',
  unknown: 'text-bg-warning',
}

function normalise(s) {
  switch (s) {
    case 'configured':
      return 'created'
    case 'created':
      return 'created'
    case 'dead':
      return 'error'
    case 'exited':
      return 'stopped'
    case 'healthy':
      return 'running'
    case 'initialized':
      return 'created'
    case 'paused':
      return 'stopped'
    case 'removing':
      return 'stopped'
    case 'running':
      return 'running'
    case 'stopped':
      return 'stopped'
    case 'stopping':
      return 'stopped'
    case 'unhealthy':
      return 'running'
    default:
      return 'unknown'
  }
}

const statusOverride = ref(null)
const status = computed(() => {
  if (statusOverride.value) {
    emit('status', normalise(statusOverride.value))
    return normalise(statusOverride.value)
  }
  emit('status', normalise(props.status))
  return normalise(props.status)
})

eventBus.subscribe("server:status", (event) => {
  if (props.server_id && props.server_id === event.detail.id) {
    statusOverride.value = event.detail.status
  }
})
</script>

<template>
  <span class="badge" v-bind:class="statusClass[status]">{{ status }}</span>
</template>
