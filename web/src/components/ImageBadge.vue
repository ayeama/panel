<script setup>
import { computed } from 'vue'

import { imageLabel } from '@/image'

const props = defineProps({
  image: String,
})

const label = computed(() => imageLabel(props.image))
const color = computed(() => 'text-bg-secondary')
</script>

<template>
  <span class="badge" :class="color">{{ label }}</span>
</template>
