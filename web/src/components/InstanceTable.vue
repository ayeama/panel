<script setup>
defineProps({
  instances: Array
})
</script>

<template>
  <table class="table table-hover">
    <thead>
      <tr>
        <th scope="col">ID</th>
        <th scope="col">Name</th>
        <th scope="col">Status</th>
      </tr>
    </thead>
    <tbody>
      <tr v-for="item in instances" v-bind="item.id">
        <td scope="row">{{ item.id }}</td>
        <td>{{ item.name }}</td>
        <td>{{ item.status }}</td>
      </tr>
    </tbody>
  </table>
</template>
