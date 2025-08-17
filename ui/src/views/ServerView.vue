<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { request } from '@/lib/client'
import ModalConfirm from '@/components/ModalConfirm.vue'
import Terminal from '@/components/Terminal.vue'
import ServerStats from '@/components/ServerStats.vue'
import ServerStatusBadge from '@/components/ServerStatusBadge.vue'
import Clipboard from '@/components/icons/Clipboard.vue'
import Check2 from '@/components/icons/Check2.vue'
import ThreeDotsVertical from '@/components/icons/ThreeDotsVertical.vue'

const route = useRoute()
const router = useRouter()

const id = route.params.id
const data = ref({})

const command = computed(() => (data.value.status === 'running' ? 'connect' : 'disconnect'))

async function handleDeleteConfirmModalConfirm() {
  try {
    const response = await request(`/servers/${id}`, {
      method: 'DELETE',
    })
  } catch (error) {
    console.log('Failed to delete server', error)
  } finally {
    router.push('/')
  }
}

async function getServer(id) {
  try {
    const response = await request(`/servers/${id}`)
    data.value = await response.json()
  } catch (error) {
    console.log('Failed to get server', error)
  }
}

async function startServer() {
  try {
    const response = await request(`/servers/${id}/start`, {
      method: 'POST',
    })
    // TODO handle?
  } catch (error) {
    console.log('Failed to start server')
  }
}

async function stopServer() {
  try {
    const response = await request(`/servers/${id}/stop`, {
      method: 'POST',
    })
    // TODO handle?
  } catch (error) {
    console.log('Failed to stop server')
  }
}

async function backupServer() {
  try {
    const response = await request(`/servers/${id}/backup`, {
      method: 'POST',
    })
  } catch (error) {
    console.log('Failed to backup server')
  }
}

onMounted(async () => {
  await getServer(id)
})

function copyAddress() {
  navigator.clipboard.writeText(data.value.address)
}

function copySFTPAddress(i) {
  navigator.clipboard.writeText('sftp://root@' + data.value.sidecar_addresses[i])
}
</script>

<template>
  <div>
    <div class="row mb-2">
      <div class="col col-auto my-auto">
        <h2>
          {{ data.name }}
        </h2>
      </div>

      <div class="col">
        <div class="row">
          <div>
            <ServerStatusBadge
              v-bind:server_id="data.id"
              v-bind:status="data.status"
              v-on:status="data.status = $event"
              class="align-top"
            />
          </div>
        </div>

        <div class="row">
          <div>
            <span class="badge text-bg-secondary">{{ data.image }}</span>
          </div>
        </div>
      </div>

      <div class="col my-auto">
        <div class="float-end">
          <div>
            <button
              class="btn btn-secondary"
              v-if="data.status === 'running'"
              v-on:click="stopServer"
            >
              Stop
            </button>
            <button class="btn btn-secondary" v-else v-on:click="startServer">Start</button>

            <button
              class="btn dropdown-toggle"
              type="button"
              data-bs-toggle="dropdown"
              aria-expanded="false"
            >
              <ThreeDotsVertical />
            </button>
            <ul class="dropdown-menu dropdown-menu-end">
              <li>
                <button class="dropdown-item" type="button" v-on:click="startServer">Start</button>
              </li>
              <li>
                <button class="dropdown-item" type="button" v-on:click="stopServer">Stop</button>
              </li>
              <li>
                <button class="dropdown-item" type="button" v-on:click="backupServer">
                  Backup
                </button>
              </li>
              <li><hr class="dropdown-divider" /></li>
              <li>
                <button
                  class="dropdown-item btn-danger"
                  type="button"
                  data-bs-toggle="modal"
                  data-bs-target="#modalConfirmDeleteServer"
                >
                  Delete
                </button>
              </li>
            </ul>
          </div>
        </div>
      </div>
    </div>

    <div class="mb-2">
      <Terminal v-if="data.id" v-bind:server_id="data.id" v-bind:command="command" />
    </div>

    <div>
      <div class="row mt-3">
        <div class="col-12 col-lg-6">
          <form>
            <div class="row gy-2">
              <div class="col-12 col-xl-6">
                <div>
                  <label for="addressInput" class="form-label">Server Address</label>
                  <div class="input-group">
                    <input
                      v-bind:value="data.address"
                      type="text"
                      class="form-control"
                      id="addressInput"
                      disabled
                      readonly
                      aria-describedby="addressCopyIcon"
                    />
                    <span
                      class="input-group-text"
                      id="addressCopyIcon"
                      style="cursor: pointer"
                      v-on:click="copyAddress()"
                      ><Clipboard
                    /></span>
                  </div>
                </div>
              </div>

              <div class="col-12 col-xl-6" v-for="(sidecar_address, i) in data.sidecar_addresses">
                <div>
                  <label for="addressInput" class="form-label">SFTP Address</label>
                  <div class="input-group">
                    <input
                      v-bind:value="'sftp://root@' + sidecar_address"
                      type="text"
                      class="form-control"
                      id="addressInput"
                      disabled
                      readonly
                      aria-describedby="addressCopyIcon"
                    />
                    <span
                      class="input-group-text"
                      id="addressCopyIcon"
                      style="cursor: pointer"
                      v-on:click="copySFTPAddress(i)"
                      ><Clipboard
                    /></span>
                  </div>
                </div>
              </div>
            </div>
          </form>
        </div>

        <!-- <div class="col-6">
          <ServerStats v-bind:url="`wss://${HOST}/servers/${id}/stats`" />
        </div> -->
      </div>
    </div>

    <ModalConfirm
      id="modalConfirmDeleteServer"
      v-bind:title="'Delete server?'"
      v-on:confirm="handleDeleteConfirmModalConfirm"
    />
  </div>
</template>

<style>
/* TODO does this effect other dropdowns? */
.dropdown-toggle::after {
  display: none !important;
}
</style>
