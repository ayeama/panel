<script setup>
import { RouterView } from 'vue-router'
</script>

<template>
  <header>
    <nav class="navbar bg-body-tertiary">
      <div class="container">
        <RouterLink class="navbar-brand" to="/">Panel</RouterLink>
      </div>
    </nav>
  </header>

  <main>
    <div class="container">
      <RouterView />
    </div>
  </main>
</template>
