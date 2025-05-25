<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import ArrowLeft from '@/components/icons/ArrowLeft.vue'
import ModalConfirm from '@/components/ModalConfirm.vue'
import Terminal from '@/components/Terminal.vue'
import ServerStats from '@/components/ServerStats.vue'
import { HOST } from '@/config'

const route = useRoute()
const router = useRouter()

const id = route.params.id
const data = ref('')

async function handleDeleteConfirmModalConfirm() {
  try {
    const response = await fetch(`https://${HOST}/servers/${id}`, {
      method: 'DELETE',
    })
  } catch (err) {
    console.log('Failed to delete server', err)
  } finally {
    router.push('/servers')
  }
}

async function getServer(id) {
  try {
    const response = await fetch(`https://${HOST}/servers/${id}`)
    data.value = await response.json()
  } catch (err) {
    console.log('Failed to get server', err)
  }
}

async function startServer() {
  try {
    const response = await fetch(`https://${HOST}/servers/${id}/start`, {
      method: 'POST',
    })
    // TODO handle?
  } catch (err) {
    console.log('Failed to start server')
  }
}

async function stopServer() {
  try {
    const response = await fetch(`https://${HOST}/servers/${id}/stop`, {
      method: 'POST',
    })
    // TODO handle?
  } catch (err) {
    console.log('Failed to stop server')
  }
}

onMounted(async () => {
  await getServer(id)
})
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
        <h1>Server</h1>
      </div>

      <div class="col my-auto">
        <div class="float-end">
          <div class="">
            <!-- TODO routing -->
            <button class="btn btn-secondary" v-on:click="startServer">Start</button>
            <button class="btn btn-secondary ms-2" v-on:click="stopServer">Stop</button>
            <button class="btn btn-secondary ms-2">Edit</button>
            <!-- <RouterLink to="{name: 'ServerEditView', params: {id: id}}" class="btn btn-secondary">
              Edit
            </RouterLink> -->
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
      <Terminal v-bind:url="`wss://${HOST}/servers/${id}/attach`" />
    </div>

    <div>
      <div class="row mt-3">
        <div class="col-6">
          <form>
            <div class="row gy-2">
              <div class="col-12">
                <div>
                  <label for="idInput" class="form-label">ID</label>
                  <input
                    type="text"
                    class="form-control"
                    id="idInput"
                    v-model="data.id"
                    disabled
                    readonly
                  />
                </div>
              </div>

              <div class="col-12">
                <div>
                  <label for="addressInput" class="form-label">Address</label>
                  <input
                    type="text"
                    class="form-control"
                    id="addressInput"
                    value="127.0.0.1:25565"
                    disabled
                    readonly
                  />
                </div>
              </div>

              <div class="col-12">
                <div>
                  <label for="nameInput" class="form-label">Name</label>
                  <input type="text" class="form-control" id="nameInput" v-model="data.name" />
                </div>
              </div>

              <div class="col-12">
                <div>
                  <label for="statusInput" class="form-label">Status</label>
                  <input
                    type="text"
                    class="form-control"
                    id="statusInput"
                    disabled
                    readonly
                    v-model="data.status"
                  />
                </div>
              </div>
            </div>
          </form>
        </div>

        <div class="col-6">
          <ServerStats v-bind:url="`wss://${HOST}/servers/${id}/stats`" />
        </div>
      </div>
    </div>

    <ModalConfirm
      id="modalConfirmDeleteServer"
      v-bind:title="'Delete server?'"
      v-on:confirm="handleDeleteConfirmModalConfirm"
    />
  </div>
</template>
