<script setup>
import { ref, nextTick, onMounted, computed } from 'vue'
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

const formWebhooks = ref([])

const spinnerCreating = ref(false)

async function addWebhook(webhook) {
  formWebhooks.value.push({ url: webhook })

  const i = formWebhooks.value.length - 1
  await nextTick()
  document.getElementById(`webhook${i}`)?.focus()
}

function removeWebhook(index) {
  formWebhooks.value.splice(index, 1)
}

onMounted(() => {
  imageReadMany()
})

async function instanceCreateRedirect() {
  if (!selectedFormImage.value) {
    return
  }

  const webhooks = [
    ...new Set(
      formWebhooks.value
        .filter((webhook) => {
          if (!webhook.url) {
            return false
          }
          try {
            new URL(webhook.url)
          } catch {
            return false
          }
          return true
        })
        .map((webhook) => webhook.url),
    ),
  ]

  const data = {
    image: selectedFormImage.value.name,
    resources: {
      cpu: formCpu.value,
      memory: formMemory.value,
      disk: formDisk.value,
    },
    webhooks: webhooks,
  }

  spinnerCreating.value = true
  try {
    await instanceCreate(data)
  } finally {
    spinnerCreating.value = false
  }

  router.push(`/instances/${instance.value.id}`)
}
</script>

<template>
  <div class="row g-2 pt-2">
    <div class="col">
      <div class="row align-items-end">
        <div class="col">
          <h1 class="h2 mb-0">Instance create</h1>
        </div>

        <div class="col col-auto d-flex">
          <!-- TODO add spinner in the button -->
          <button v-if="!spinnerCreating" type="button" class="btn btn-primary" v-on:click="instanceCreateRedirect()">
            Create
          </button>

          <button v-else class="btn btn-primary" type="button" disabled>
            <span class="spinner-border spinner-border-sm" aria-hidden="true"></span>
            <span role="status"> Creating</span>
          </button>
        </div>
      </div>
    </div>

    <div class="col-12">
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

    <div class="col-12">
      <div class="row">
        <div class="col-4">
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

        <div class="col-4">
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

        <div class="col-4">
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

    <div class="col-12">
      <div class="d-flex justify-content-between align-items-end">
        <label for="webhookInput" class="form-label">Webhooks</label>
        <button class="btn btn-sm btn-secondary mb-2" v-on:click="addWebhook('')">Add</button>
      </div>

      <div v-if="formWebhooks.length === 0" class="text-muted small">No webhooks</div>
      <div v-else class="d-flex flex-column gap-2">
        <div v-for="(webhook, i) in formWebhooks" :key="`${webhook[i]}`" class="input-group">
          <input
            class="form-control"
            type="url"
            v-model="webhook.url"
            :id="`webhook${i}`"
            :aria-describedby="`webhook${i}-remove`"
          />
          <button
            class="btn btn-outline-danger"
            type="button"
            :id="`webhook${i}-remove`"
            v-on:click="removeWebhook(i)"
          >
            Remove
          </button>
        </div>
      </div>
    </div>
  </div>
</template>
