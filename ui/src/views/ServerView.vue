<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import ModalConfirm from '@/components/ModalConfirm.vue'
import Terminal from '@/components/Terminal.vue'
import ServerStats from '@/components/ServerStats.vue'
import ServerStatusBadge from '@/components/ServerStatusBadge.vue'
import { HOST } from '@/config'

const route = useRoute()
const router = useRouter()

const id = route.params.id
const data = ref({})

const command = computed(() => (data.value.status === 'running' ? 'connect' : 'disconnect'))

async function handleDeleteConfirmModalConfirm() {
  try {
    const response = await fetch(`https://${HOST}/servers/${id}`, {
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
    const response = await fetch(`https://${HOST}/servers/${id}`)
    data.value = await response.json()
  } catch (error) {
    console.log('Failed to get server', error)
  }
}

async function startServer() {
  try {
    const response = await fetch(`https://${HOST}/servers/${id}/start`, {
      method: 'POST',
    })
    // TODO handle?
  } catch (error) {
    console.log('Failed to start server')
  }
}

async function stopServer() {
  try {
    const response = await fetch(`https://${HOST}/servers/${id}/stop`, {
      method: 'POST',
    })
    // TODO handle?
  } catch (error) {
    console.log('Failed to stop server')
  }
}

onMounted(async () => {
  await getServer(id)
})
</script>

<template>
  <div>
    <div class="row mb-2">
      <div class="col">
        <h2>
          {{ data.name }} <ServerStatusBadge v-bind:server_id="data.id" v-bind:status="data.status" v-on:status="data.status=$event" class="fs-6 align-top" />
        </h2>
      </div>

      <div class="col my-auto">
        <div class="float-end">
          <div class="">
            <!-- TODO routing -->
            <button class="btn btn-secondary" v-on:click="startServer">Start</button>
            <button class="btn btn-secondary ms-2" v-on:click="stopServer">Stop</button>
            <button
              type="button"
              class="btn btn-danger ms-2"
              data-bs-toggle="modal"
              data-bs-target="#modalConfirmDeleteServer"
            >
              Delete
            </button>
          </div>
        </div>
      </div>
    </div>

    <div class="mb-2">
      <Terminal v-if="data.id" v-bind:server_id="data.id" v-bind:command="command" />
    </div>

    <div>
      <div class="row mt-3">
        <div class="col-6">
          <form>
            <div class="row gy-2">
              <div class="col-6">
                <div>
                  <label for="addressInput" class="form-label">Address</label>
                  <input
                    v-bind:value="data.address"
                    type="text"
                    class="form-control"
                    id="addressInput"
                    disabled
                    readonly
                  />
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
