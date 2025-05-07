<script setup>
import ServersTable from '@/components/ServersTable.vue'
</script>

<template>
  <div>
    <div class="row mb-2">
      <div class="col">
        <h2>Servers</h2>
      </div>

      <div class="col my-auto">
        <div class="float-end">
          <RouterLink to="/servers/create" class="btn btn-primary">Create</RouterLink>
        </div>
      </div>
    </div>

    <div>
      <ServersTable />
    </div>
  </div>
</template>
