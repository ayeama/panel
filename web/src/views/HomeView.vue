<script setup>
import { onMounted } from 'vue'
import { RouterLink } from 'vue-router'

import { useInstance } from '@/composables/useInstance'
import InstanceTable from '@/components/InstanceTable.vue'

const { instances, instanceReadMany } = useInstance()

onMounted(() => {
  instanceReadMany()
})
</script>

<template>
  <div class="row row-cols-1 py-2">
    <div class="col">
      <div class="row pb-2">
        <div class="col">
          <h1 class="mb-0">Instances</h1>
        </div>

        <div class="col col-auto d-flex align-items-end">
          <RouterLink class="btn btn-primary" role="button" to="/instances/create"
            >Create</RouterLink
          >
        </div>
      </div>
    </div>

    <div class="col">
      <InstanceTable :instances="instances" />
    </div>
  </div>
</template>
