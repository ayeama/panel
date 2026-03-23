<script setup lang="ts">
const props = defineProps<{
  status: string
}>()

const label = computed(() => {
    const label = props.status.toLowerCase()
    switch (label) {
    case 'initialized':
        return 'stopped'
    default:
        return label
    }
})

const color = computed(() => {
  switch (label.value) {
    case 'running':
        return 'secondary'
    case 'created':
    case 'stopped':
        return 'neutral'
    case 'exited':
    default:
        return 'error'
  }
})
</script>

<template>
    <UBadge variant="outline" :color="color">{{ label }}</UBadge>
</template>
