<script setup>
import { RouterLink } from 'vue-router'

defineProps({
  instances: Array,
})
</script>

<template>
  <table class="table">
    <thead>
      <tr>
        <th scope="col">Name</th>
        <th scope="col">Image</th>
        <th scope="col">Status</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="item in instances" v-bind="item.id">
        <td>
          <RouterLink :to="`/instances/${item.id}`">{{ item.name }}</RouterLink>
        </td>
        <td>{{ item.image }}</td>
        <td>{{ item.status }}</td>
      </tr>
    </tbody>
  </table>
</template>
