<script setup>
import { toRef } from 'vue'

const props = defineProps({
  status: {
    type: String,
    required: true,
    validator: (value) => {
      ;['created', 'running', 'stopped', 'exited', 'dead'].includes(value)
    },
  },
})

const status = toRef(props, 'status')
</script>

<template>
  <span
    class="badge"
    v-bind:class="{
      'text-bg-secondary': status == 'created',
      'text-bg-success': status == 'running',
      'text-bg-secondary': status == 'stopped',
      'text-bg-secondary': status == 'exited',
      'text-bg-danger': status == 'dead',
    }"
    >{{ status }}</span
  >
</template>
