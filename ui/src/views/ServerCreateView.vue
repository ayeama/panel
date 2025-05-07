<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const images = ref([])
const image = ref('')

onMounted(async () => {
  getImages()
})

async function getImages() {
  var limit = 10
  var offset = 0

  try {
    const response = await fetch(
      `${window.CONFIG.api.scheme}://${window.CONFIG.api.host}${window.CONFIG.api.path}/images?limit=${limit}&offset=${offset}`,
    )
    const data = await response.json()
    images.value = data.items

    if (images.value.length > 0) {
      image.value = data.items[0].image
    }
  } catch (error) {
    console.log('failed to fetch images', error)
  }
}

async function createServer() {
  try {
    const response = await fetch(
      `${window.CONFIG.api.scheme}://${window.CONFIG.api.host}${window.CONFIG.api.path}/servers`,
      {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          image: image.value,
        }),
      },
    )
  } catch (error) {
    console.log('Failed creating server', error)
  } finally {
    router.push('/')
  }
}
</script>

<template>
  <div>
    <div class="row mb-2">
      <div class="col">
        <h2>Create Server</h2>
      </div>

      <div class="col my-auto">
        <div class="float-end">
          <button form="createServer" type="submit" class="btn btn-primary">Create</button>
        </div>
      </div>
    </div>

    <form v-on:submit.prevent="createServer" id="createServer">
      <div class="row">
        <div class="col">
          <div>
            <label for="imageSelect" class="form-label">Image</label>
            <select
              v-model="image"
              class="form-select"
              id="imageSelect"
              aria-label="Default select example"
            >
              <option v-for="image in images" v-bind:id="image.image">{{ image.image }}</option>
            </select>
          </div>
        </div>

        <!-- <div class="col">
          <div>
            <label for="nameInput" class="form-label">Name</label>
            <input type="text" class="form-control" id="nameInput" v-model="name" />
          </div>
        </div> -->
      </div>

      <!-- <div class="row mt-3">
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
      </div> -->
    </form>
  </div>
</template>
