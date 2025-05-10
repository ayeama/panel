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

    <div class="row">
      <div class="col">
        <h1>Create Server</h1>
      </div>

      <div class="col">
        <div class="float-end">
          <button v-if="!loading" form="createServer" type="submit" class="btn btn-primary">
            Create
          </button>
        </div>
      </div>
    </div>

    <form v-on:submit.prevent="createServer" id="createServer">
      <div class="row">
        <div class="col">
          <div>
            <label for="gameSelect" class="form-label">Game</label>
            <select class="form-select" id="gameSelect" aria-label="Game select">
              <option selected>Minecraft</option>
              <option value="1">Valheim</option>
            </select>
          </div>
        </div>

        <div class="col">
          <div>
            <label for="variantSelect" class="form-label">Variant</label>
            <select class="form-select" id="variantSelect" aria-label="Variant select" disabled>
              <option selected>Vanilla</option>
              <option value="1">Paper</option>
            </select>
          </div>
        </div>

        <div class="col">
          <div>
            <label for="versionSelect" class="form-label">Version</label>
            <select class="form-select" id="versionSelect" aria-label="Version select" disabled>
              <option selected>1.21.5</option>
              <option value="1">1.21.5</option>
              <option value="2">1.21.4</option>
              <option value="3">1.21.3</option>
              <option value="4">1.21.2</option>
              <option value="5">1.21.1</option>
              <option value="6">1.21</option>
              <option value="7">1.20.1</option>
              <option value="8">1.20</option>
            </select>
          </div>
        </div>
      </div>

      <div class="row mt-3">
        <div class="col">
          <div>
            <label for="nameInput" class="form-label">Name</label>
            <input type="text" class="form-control" id="nameInput" v-model="name" />
          </div>
        </div>
      </div>

      <div class="row mt-3">
        <div class="col">
          <label for="cpuInput" class="form-label">CPU</label>
          <div class="input-group">
            <input type="number" class="form-control" id="cpuInput" value="100" />
            <span class="input-group-text">%</span>
          </div>
          <div class="form-text">400% available</div>
        </div>

        <div class="col">
          <label for="memoryInput" class="form-label">Memory</label>
          <div class="input-group">
            <input type="number" class="form-control" id="memoryInput" value="1000" />
            <span class="input-group-text">MB</span>
          </div>
          <div class="form-text">2000MB available</div>
        </div>

        <div class="col">
          <label for="diskInput" class="form-label">Disk</label>
          <div class="input-group">
            <input type="number" class="form-control" id="diskInput" value="20" />
            <span class="input-group-text">GB</span>
          </div>
          <div class="form-text">40GB available</div>
        </div>
      </div>
    </form>
  </div>
</template>
