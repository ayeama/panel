<script setup>
import { computed } from 'vue'

const props = defineProps({
  status: String,
})

const label = computed(() => {
  return props.status
})

const color = computed(() => {
  switch (label.value) {
    case 'running':
      return 'text-bg-success'
    default:
      return 'text-bg-secondary'
  }
})
</script>

<template>
  <span class="badge" :class="color">{{ label }}</span>
</template>
