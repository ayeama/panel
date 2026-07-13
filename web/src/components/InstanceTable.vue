<script setup>
import { useRouter } from 'vue-router'

import InstanceStatusBadge from './InstanceStatusBadge.vue'

const router = useRouter()

const props = defineProps({
  instances: Array,
})
</script>

<template>
  <div class="table-responsive">
    <table class="table table-hover clickable">
      <thead>
        <tr>
          <th scope="col">Name</th>
          <th scope="col">Image</th>
          <th scope="col">Status</th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="item in props.instances"
          :key="item.id"
          v-on:click="router.push(`/instances/${item.id}`)"
        >
          <td>{{ item.name }}</td>
          <td>{{ item.image }}</td>
          <td><InstanceStatusBadge :status="item.status" /></td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<style scoped>
.clickable :deep(tbody > tr:hover) {
  cursor: pointer;
}
</style>
