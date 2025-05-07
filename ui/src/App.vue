<script setup>
import { RouterView } from 'vue-router'
import NavBar from './components/NavBar.vue'
</script>

<template>
  <div>
    <div class="mb-2">
      <NavBar />
    </div>

    <div class="container">
      <RouterView />
    </div>
  </div>
</template>
