<script setup lang="ts">
const props = defineProps<{
  image: string
}>()

function short(image: string) {
    return image.split('/').pop() || image
}
</script>

<template>
    <UBadge variant="outline" color="neutral">{{ image }}</UBadge>
</template>
