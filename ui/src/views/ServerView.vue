<script setup>
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import ArrowLeft from '@/components/icons/ArrowLeft.vue'
import ModalConfirm from '@/components/ModalConfirm.vue'
import Terminal from '@/components/Terminal.vue'
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
            <RouterLink to="{name: 'ServerEditView', params: {id: id}}" class="btn btn-secondary">
              Edit
            </RouterLink>
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
      <form>
        <div class="row gy-2">
          <div class="col-6">
            <label for="inputId" class="form-label">ID</label>
            <input type="text" class="form-control" id="inputId" readonly v-model="data.id">
          </div>

          <div class="col-6">
            <label for="inputName" class="form-label">Name</label>
            <input type="text" class="form-control" id="inputName" v-model="data.name">
          </div>

          <div class="col-6">
            <label for="inputStatus" class="form-label">Status</label>
            <input type="text" class="form-control" id="inputStatus" readonly v-model="data.status">
          </div>

        </div>
      </form>
    </div>

    <ModalConfirm
      id="modalConfirmDeleteServer"
      v-bind:title="'Delete server?'"
      v-on:confirm="handleDeleteConfirmModalConfirm"
    />
  </div>
</template>
