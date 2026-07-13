<script setup>
import { onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { useInstance } from '@/composables/useInstance'
import InstanceStatistics from '@/components/InstanceStatistics.vue'
import InstanceStatusBadge from '@/components/InstanceStatusBadge.vue'
import InstanceTerminal from '@/components/InstanceTerminal.vue'

import { API_URL } from '@/api'

const route = useRoute()
const router = useRouter()

const id = route.params.id

const { instance, instanceRead, instanceDelete, instanceStart, instanceStop } = useInstance()

onMounted(() => {
  instanceRead(id)
})

async function instanceDeleteRedirect(id) {
  await instanceDelete(id)
  router.push('/')
}

function running(status) {
  switch (status) {
    case 'running':
      return true
    default:
      return false
  }
}
</script>

<!-- TODO add bootstrap placeholders for loading elements -->
<template>
  <div v-if="instance" class="row g-2 py-2">
    <div class="col">
      <div class="row">
        <div class="col d-flex">
          <h1 class="mb-0">{{ instance.name }}</h1>
        </div>

        <div class="col col-auto d-flex align-items-end">
          <div class="btn-group">
            <button
              v-if="!running(instance.status)"
              type="button"
              class="btn btn-secondary"
              v-on:click="instanceStart(id)"
            >
              Start
            </button>
            <button v-else type="button" class="btn btn-secondary" v-on:click="instanceStop(id)">
              Stop
            </button>

            <button
              type="button"
              class="btn btn-secondary dropdown-toggle dropdown-toggle-split"
              data-bs-toggle="dropdown"
              aria-expanded="false"
            >
              <span class="visually-hidden">Toggle Dropdown</span>
            </button>

            <ul class="dropdown-menu">
              <li><a class="dropdown-item" href="#todo">Backup</a></li>
              <li><a class="dropdown-item" href="#todo">Restore</a></li>
              <li>
                <a class="dropdown-item" :href="`${API_URL}/instances/${id}/logs`" download=""
                  >Logs</a
                >
              </li>
              <li><hr class="dropdown-divider" /></li>
              <li>
                <a
                  class="dropdown-item text-danger"
                  href="#todo"
                  v-on:click="instanceDeleteRedirect(id)"
                  >Delete</a
                >
              </li>
            </ul>
          </div>
        </div>
      </div>
    </div>

    <div class="col col-12">
      <InstanceTerminal v-if="instance" :id="instance.id" />
    </div>

    <div class="col">
      <div v-if="instance">
        <input
          id="instanceName"
          class="form-control-plaintext"
          type="text"
          :value="instance.name"
          readonly
        />
        <input
          id="instanceImage"
          class="form-control-plaintext"
          type="text"
          :value="instance.image"
          readonly
        />
        <InstanceStatusBadge :status="instance.status" />

        <input
          v-for="port in instance.ports"
          id="instancePort"
          class="form-control-plaintext"
          type="text"
          :value="port"
          readonly
        />

        <label for="instanceCpu" class="form-label">CPU</label>
        <input
          id="instanceCpu"
          class="form-control"
          type="text"
          :value="instance.resources.cpu"
          readonly
        />

        <label for="instanceMemory" class="form-label">Memory</label>
        <input
          id="instanceMemory"
          class="form-control"
          type="text"
          :value="instance.resources.memory"
          readonly
        />

        <label for="instanceMemory" class="form-label">Disk</label>
        <input
          id="instanceDisk"
          class="form-control"
          type="text"
          :value="instance.resources.disk"
          readonly
        />

        <label for="instanceWebhook" class="form-label">Webhook</label>
        <input
          id="instanceWebhook"
          class="form-control"
          type="text"
          :value="instance.webhook"
          readonly
        />
      </div>
    </div>

    <div class="col">
      <InstanceStatistics v-if="instance" :id="instance.id" />
    </div>
  </div>
</template>
