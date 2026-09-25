<script setup>
import { computed } from 'vue'
import { RouterView } from 'vue-router'

import { Theme, useTheme } from './composables/useTheme'
import SunFill from './components/icons/SunFill.vue'
import MoonFill from './components/icons/MoonFill.vue'
import CircleHalfFill from './components/icons/CircleHalfFill.vue'

const { getTheme, setTheme } = useTheme()

const themeIcon = computed(() => {
  switch (getTheme()) {
    case Theme.LIGHT:
      return SunFill
    case Theme.DARK:
      return MoonFill
    default:
      return CircleHalfFill
  }
})
</script>

<template>
  <header>
    <nav class="navbar bg-body-tertiary">
      <div class="container">
        <RouterLink class="navbar-brand" to="/">Panel</RouterLink>

        <div class="dropdown">
          <button class="btn dropdown-toggle" type="button" data-bs-toggle="dropdown">
            <component :is="themeIcon" />
          </button>
          <ul class="dropdown-menu dropdown-menu-end">
            <li>
              <a
                class="dropdown-item"
                :class="{ active: getTheme() === Theme.LIGHT }"
                href="#"
                @click.prevent="setTheme(Theme.LIGHT)"
                ><SunFill /> Light</a
              >
            </li>
            <li>
              <a
                class="dropdown-item"
                :class="{ active: getTheme() === Theme.DARK }"
                href="#"
                @click.prevent="setTheme(Theme.DARK)"
                ><MoonFill /> Dark</a
              >
            </li>
            <li>
              <a
                class="dropdown-item"
                :class="{ active: getTheme() === Theme.SYSTEM }"
                href="#"
                @click.prevent="setTheme(Theme.SYSTEM)"
                ><CircleHalfFill /> Auto</a
              >
            </li>
          </ul>
        </div>
      </div>
    </nav>
  </header>

  <main>
    <div class="container">
      <RouterView />
    </div>
  </main>
</template>
