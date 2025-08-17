<script setup>
import { ref } from 'vue'
import DiscordSymbol from '@/components/icons/DiscordSymbol.vue'

const oauth2URL = ref(window.CONFIG.oauth2.url)
</script>

<template>
  <div>
    <div class="row d-flex justify-content-center">
      <div class="col-sm-4">
        <h2>Sign in</h2>
      </div>
    </div>

    <div class="row d-flex justify-content-center">
      <div class="col-sm-4">
        <form>
          <div class="mb-2">
            <label for="emailInput" class="form-label">Email address</label>
            <input type="email" class="form-control" id="emailInput" disabled />
          </div>

          <div class="mb-2">
            <label for="passwordInput" class="form-label">Password</label>
            <input type="password" class="form-control" id="passwordInput" disabled />
          </div>

          <div class="mb-2 form-check">
            <div class="float-end">
              <label class="form-check-label" for="gridCheck"> Remember me </label>
              <input class="form-check-input" type="checkbox" id="gridCheck" disabled />
            </div>
          </div>

          <div class="d-grid gap-2">
            <button type="submit" class="btn btn-primary" disabled>Sign in</button>
          </div>
        </form>
      </div>
    </div>

    <div class="row mb-5 d-flex justify-content-center">
      <div class="col-sm-4"></div>
    </div>

    <div class="row mb-2 d-flex justify-content-center">
      <div class="col-sm-4 d-flex justify-content-center">
        <a class="btn btn-outline-secondary" v-bind:href="oauth2URL">
          <!-- <span><DiscordSymbol class="icon-sm" /></span> -->
          Continue with Discord
        </a>
      </div>
    </div>
  </div>
</template>

<style>
.icon-sm {
  width: 1rem;
  height: 1rem;
}
</style>
