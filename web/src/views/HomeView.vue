<script setup>
import { onMounted } from 'vue'

import { useInstance } from '@/composables/useInstance'
import InstanceTable from '@/components/InstanceTable.vue'

const { instances, instanceReadMany } = useInstance()

onMounted(() => {
  instanceReadMany()
})
</script>

<template>
  <div class="row">
    <div class="col">
      <InstanceTable v-bind:instances="instances" />
    </div>
  </div>
</template>
