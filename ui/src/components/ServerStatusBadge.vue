<script setup>
const props = defineProps({
  status: {
    type: String,
    required: true,
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
      ].includes(value)
    },
  },
})

const statusClass = {
  created: 'text-bg-secondary',
  running: 'text-bg-success',
  stopped: 'text-bg-secondary',
  error: 'text-bg-danger',
  unknown: 'text-bg-warning',
}

function normalise(status) {
  switch (status) {
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
</script>

<template>
  <span class="badge" v-bind:class="statusClass[normalise(status)]">{{ normalise(status) }}</span>
</template>
