<script setup>
import { useRouter } from 'vue-router'

import ImageBadge from './ImageBadge.vue'
import InstanceStatusBadge from './InstanceStatusBadge.vue'

const router = useRouter()

const props = defineProps({
  instances: Array,
})
</script>

<template>
  <div class="card overflow-hidden">
    <div class="card-body p-0">
      <div class="table-responsive">
        <table class="table table-hover mb-0 clickable">
          <caption class="ps-2">{{ instances.length }} instances </caption>

          <thead>
            <tr>
              <th scope="col">Name</th>
              <th scope="col">Status</th>
              <th scope="col">Image</th>
            </tr>
          </thead>
    
          <tbody>
            <tr
              v-if="props.instances.length > 0"
              v-for="item in props.instances"
              :key="item.id"
              v-on:click="router.push(`/instances/${item.id}`)"
            >
              <td>{{ item.name }}</td>
              <td><InstanceStatusBadge :status="item.status" /></td>
              <td><ImageBadge :image="item.image" /></td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
  </div>
</template>

<style scoped>
.clickable :deep(tbody > tr:hover) {
  cursor: pointer;
}
</style>
