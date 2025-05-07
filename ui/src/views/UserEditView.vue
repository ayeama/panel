<script setup>
import ArrowLeft from '@/components/icons/ArrowLeft.vue'
</script>

<template>
  <div>
    <RouterLink to="/servers" class="icon-link">
      <ArrowLeft />
      Users
    </RouterLink>

    <h1>Edit User</h1>
    <p>user id: {{ $route.params.id }}</p>
  </div>
</template>
