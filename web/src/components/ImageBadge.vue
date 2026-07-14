<script setup>
import { computed } from 'vue'

const props = defineProps({
  image: String,
})

const label = computed(() => {
  // TODO should be a useImages composable?
  const parts = props.image.split("/")
  return parts.at(-1)
})

const color = computed(() => {
  return 'text-bg-secondary'
})
</script>

<template>
  <span class="badge" :class="color">{{ label }}</span>
</template>
