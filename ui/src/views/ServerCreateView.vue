<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { HOST } from '@/config'
import ArrowLeft from '@/components/icons/ArrowLeft.vue'

const router = useRouter()

const id = ref('')
const name = ref('')
const loading = ref(false)

async function createServer() {
  loading.value = true

  try {
    const response = await fetch(`https://${HOST}/servers`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        id: id.value,
        name: name.value,
      }),
    })
  } catch (err) {
    console.log('Failed creating server', err)
  } finally {
    // await new Promise(r => setTimeout(r, 250))
    loading.value = false
    router.push('/servers')
  }
}
</script>

<template>
  <div>
    <div>
      <RouterLink to="/servers" class="icon-link">
        <ArrowLeft />
        Servers
      </RouterLink>
    </div>

    <h1>Create Server</h1>

    <form v-on:submit.prevent="createServer">
      <div class="mb-3">
        <label for="idInput" class="form-label">Identity</label>
        <input
          type="text"
          class="form-control"
          id="idInput"
          aria-describedby="idHelp"
          v-model="id"
        />
        <div id="idHelp" class="form-text">A sha256sum</div>
      </div>

      <div class="mb-3">
        <label for="nameInput" class="form-label">Name</label>
        <input type="text" class="form-control" id="nameInput" v-model="name" />
      </div>

      <button v-if="!loading" type="submit" class="btn btn-primary">Create</button>
      <button v-else type="button" class="btn btn-primary" disabled>
        <span class="spinner-border spinner-border-sm" aria-hidden="true"></span>
        <span class="ms-1" role="status">Creating...</span>
      </button>
    </form>
  </div>
</template>
