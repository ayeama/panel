<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router';

import { useImage } from '@/composables/useImage';
import { useInstance } from '@/composables/useInstance';

const router = useRouter()

const { images, imageReadMany } = useImage()
const { instance, instanceCreate } = useInstance()

const formImage = ref(null)
const selectedFormImage = computed(() => images.value.find(i => i.name === formImage.value))

onMounted(() => {
  imageReadMany()
})

async function instanceCreateRedirect() {
  if (!selectedFormImage.value) {
    return
  }

  await instanceCreate({image_id: selectedFormImage.value.id})
  router.push(`/instances/${instance.value.id}`)
}
</script>

<template>
  <div class="row row-cols-1 py-2">
    <div class="col">
      <h1>Instance Create</h1>
    </div>

    <div class="col">
      <div class="mb-2">
        <label for="datalistImages" class="form-label">Image</label>
        <input v-model="formImage" id="datalistImages" class="form-control" list="datalistOptionsImages" placeholder="Search...">
        <datalist id="datalistOptionsImages">
          <option v-for="item in images" v-bind:key="item.id">{{ item.name }}</option>
        </datalist>
      </div>

      <div class="d-flex justify-content-end">
        <button type="button" class="btn btn-primary" v-on:click="instanceCreateRedirect()">Create</button>
      </div>
    </div>
  </div>
</template>
