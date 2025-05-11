<script setup>
import { ENVIRONMENTS } from '@/config'
import Moon from '@/components/icons/Moon.vue'
</script>

<template>
  <header>
    <nav class="navbar navbar-expand-lg bg-body-tertiary">
      <div class="container">
        <span class="navbar-brand mb-0 h1">Panel</span>

        <div class="d-flex">
          <div>
            <button class="btn" type="button"><Moon /></button>
          </div>

          <div>
            <select class="form-select" aria-label="Environment select">
              <option
                v-for="environment in ENVIRONMENTS"
                v-bind:key="environment.name"
                v-bind:value="environment.name"
              >
                {{ environment.name }}
              </option>
            </select>
          </div>
        </div>
      </div>
    </nav>
  </header>
</template>
