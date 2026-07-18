<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'

import { useImage } from '@/composables/useImage'
import { useInstance } from '@/composables/useInstance'

const router = useRouter()

const { images, imageReadMany } = useImage()
const { instance, instanceCreate } = useInstance()

const formImage = ref(null)
const selectedFormImage = computed(() => images.value.find((i) => i.name === formImage.value))

const formCpu = ref(1.0)
const formMemory = ref(1.0)
const formDisk = ref(0.0)

onMounted(() => {
  imageReadMany()
})

async function instanceCreateRedirect() {
  if (!selectedFormImage.value) {
    return
  }

  const data = {
    image_id: selectedFormImage.value.id,
    resources: {
      cpu: formCpu.value,
      memory: formMemory.value,
      disk: formDisk.value,
    },
  }

  await instanceCreate(data)
  router.push(`/instances/${instance.value.id}`)
}
</script>

<template>
  <div class="row row-cols-1 g-2 pt-2">
    <div class="col">
      <h1 class="h4 mb-0">Instance create</h1>
    </div>

    <div class="col">
      <label for="datalistImages" class="form-label">Image</label>
      <input
        v-model="formImage"
        id="datalistImages"
        class="form-control"
        list="datalistOptionsImages"
        placeholder="Search..."
      />
      <datalist id="datalistOptionsImages">
        <option v-for="item in images" :key="item.id">{{ item.name }}</option>
      </datalist>
    </div>

    <div class="col">
      <div class="row row-cols-3">
        <div class="col">
          <label for="cpuInput" class="form-label">CPU</label>
          <input
            id="cpuInput"
            class="form-control"
            type="number"
            min="0"
            step="0.1"
            v-model.number="formCpu"
          />
        </div>

        <div class="col">
          <label for="memoryInput" class="form-label">Memory</label>
          <input
            id="memoryInput"
            class="form-control"
            type="number"
            min="0"
            step="0.1"
            v-model.number="formMemory"
          />
        </div>

        <div class="col">
          <label for="diskInput" class="form-label">Disk</label>
          <input
            id="diskInput"
            class="form-control"
            type="number"
            min="0"
            step="0.1"
            v-model.number="formDisk"
            readonly
            disabled
          />
        </div>
      </div>
    </div>

    <div class="col">
      <button type="button" class="btn btn-primary float-end" v-on:click="instanceCreateRedirect()">
        Create
      </button>
    </div>
  </div>
</template>
