<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

import { useInstance } from '@/composables/useInstance'
import InstanceStatistics from '@/components/InstanceStatistics.vue'
import InstanceStatusBadge from '@/components/InstanceStatusBadge.vue'
import InstanceTerminal from '@/components/InstanceTerminal.vue'

import { API_URL } from '@/api'

const route = useRoute()
const router = useRouter()

const id = route.params.id

const { instance, instanceRead, instanceDelete, instanceStart, instanceStop, instanceRestore } =
  useInstance()

const restoreFileInput = ref(null)

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

function restore() {
  restoreFileInput.value.click()
}

async function restoreFileSelected(event) {
  const file = event.target.files[0]
  if (!file) {
    return
  }

  const data = new FormData()
  data.append('backup', file)

  await instanceRestore(id, data)

  event.target.value = ''
}
</script>

<!-- TODO add bootstrap placeholders for loading elements -->
<template>
  <div v-if="instance" class="row g-2 py-2">
    <div class="col">
      <div class="row">
        <div class="col d-flex">
          <div class="row row-cols-1">
            <div class="col">
              <h1 class="h4 mb-0">{{ instance.name }}</h1>
            </div>

            <div class="col">
              <InstanceStatusBadge :status="instance.status" />
            </div>
          </div>
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
              <li>
                <a class="dropdown-item" :href="`${API_URL}/instances/${id}/backup`">Backup</a>
              </li>
              <li>
                <a class="dropdown-item" href="#" v-on:click="restore">Restore</a>
                <input
                  class="d-none"
                  ref="restoreFileInput"
                  type="file"
                  accept=".zip"
                  v-on:change="restoreFileSelected"
                />
              </li>
              <li>
                <a class="dropdown-item" :href="`${API_URL}/instances/${id}/logs`" download=""
                  >Logs</a
                >
              </li>
              <li><hr class="dropdown-divider" /></li>
              <li>
                <a
                  class="dropdown-item text-danger"
                  href="#"
                  v-on:click="instanceDeleteRedirect(id)"
                  >Delete</a
                >
              </li>
            </ul>
          </div>
        </div>
      </div>
    </div>

    <div class="col-12">
      <div class="card overflow-hidden">
        <div class="card-body p-0">
          <InstanceTerminal v-if="instance" :instance="instance" />
        </div>
      </div>
    </div>

    <div class="col-md-12 col-lg-8">
      <div v-if="instance" class="card">
        <div class="card-body">
          <h5 class="card-title">Details</h5>

          <div class="row g-2">
            <div class="col-12">
              <label for="instanceID" class="form-label">ID</label>
              <input
                id="instanceID"
                class="form-control"
                type="text"
                :value="instance.id"
                readonly
              />
            </div>

            <div class="col-12">
              <label for="instanceImage" class="form-label">Image</label>
              <input
                id="instanceImage"
                class="form-control"
                type="text"
                :value="instance.image"
                readonly
              />
            </div>

            <div class="col-12">
              <label for="instanceDisk" class="form-label">Ports</label>
              <div v-if="Object.keys(instance.ports).length === 0" class="text-muted small">
                No ports
              </div>

              <div v-else class="d-flex flex-column gap-2">
                <div v-for="(port, _) in instance.ports" :key="port" class="col-4">
                  <input
                    id="instanceDisk"
                    class="form-control"
                    type="number"
                    :value="port"
                    readonly
                  />
                </div>
              </div>
            </div>

            <div class="row mt-0 g-2">
              <div class="col-4">
                <label for="instanceCPU" class="form-label">CPU</label>
                <input
                  id="instanceCPU"
                  class="form-control"
                  type="number"
                  min="0"
                  step="0.1"
                  :value="instance.resources.cpu"
                  readonly
                />
              </div>

              <div class="col-4">
                <label for="instanceMemory" class="form-label">Memory</label>
                <input
                  id="instanceMemory"
                  class="form-control"
                  type="number"
                  min="0"
                  step="0.1"
                  :value="instance.resources.memory"
                  readonly
                />
              </div>

              <div class="col-4">
                <label for="instanceDisk" class="form-label">Disk</label>
                <input
                  id="instanceDisk"
                  class="form-control"
                  type="number"
                  min="0"
                  step="0.1"
                  :value="instance.resources.disk"
                  readonly
                  disabled
                />
              </div>
            </div>

            <div class="col-12">
              <label for="instanceWebhook" class="form-label">Webhooks</label>
              <div v-if="instance.webhooks.length === 0" class="text-muted small">No webhooks</div>

              <div v-else class="d-flex flex-column gap-2">
                <input
                  v-for="(webhook, _) in instance.webhooks"
                  :key="webhook"
                  id="instanceWebhook"
                  class="form-control"
                  type="text"
                  :value="webhook"
                  readonly
                />
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="col-md-12 col-lg-4">
      <div class="card">
        <div class="card-body">
          <h5 class="card-title">Stats</h5>

          <InstanceStatistics v-if="instance" :id="instance.id" />
        </div>
      </div>
    </div>
  </div>
</template>
