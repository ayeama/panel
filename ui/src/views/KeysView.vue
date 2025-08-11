<script setup>
import { ref, onMounted } from 'vue'
import ModalConfirm from '@/components/ModalConfirm.vue'

const data = ref({})

const newKeyComment = ref('')
const newKeyKey = ref('')

async function refresh() {
  newKeyComment.value = ""
  newKeyKey.value = ""
  await getKeys()
}

async function getKeys() {
  try {
    const response = await fetch(
      `${window.CONFIG.api.scheme}://${window.CONFIG.api.host}${window.CONFIG.api.path}/keys`,
    )
    data.value = await response.json()
  } catch (error) {
    console.log('Failed to get keys', error)
  }
}

async function createKey() {
  try {
    const response = await fetch(
      `${window.CONFIG.api.scheme}://${window.CONFIG.api.host}${window.CONFIG.api.path}/keys`,
      {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          comment: newKeyComment.value,
          key: newKeyKey.value,
        }),
      },
    )
  } catch (error) {
    console.log('Failed creating key', error)
  } finally {
    await refresh()
  }
}

const selectedKeyId = ref(null)

function setSelectedKey(id) {
  selectedKeyId.value = id
}

async function handleDeleteConfirmModalConfirm() {
  if (!selectedKeyId.value) {
    console.error('No key was selected before attempting to delete')
    return
  }

  try {
    const response = await fetch(
      `${window.CONFIG.api.scheme}://${window.CONFIG.api.host}${window.CONFIG.api.path}/keys/${selectedKeyId.value}`,
      {
        method: 'DELETE',
      },
    )
  } catch (error) {
    console.log('Failed to delete key', error)
  } finally {
    await refresh()
  }
}

onMounted(async () => {
  await getKeys()
})
</script>

<template>
  <div>
    <div>
      <h2>keys</h2>
    </div>

    <div>
      <form v-on:submit.prevent="createKey" id="createKey">
        
        <div class="mb-3">
          <label for="commentInput" class="form-label">Comment</label>
          <input v-model="newKeyComment" type="text" class="form-control" id="commentInput" placeholder="user@localhost">
        </div>

        <div class="mb-3">
          <label for="keyInput" class="form-label">Key</label>
          <textarea v-model="newKeyKey" class="form-control" id="keyInput" rows="6" placeholder="Begins with 'ssh-rsa', 'ecdsa-sha2-nistp256', 'ecdsa-sha2-nistp384', 'ecdsa-sha2-nistp521', 'ssh-ed25519', 'sk-ecdsa-sha2-nistp256@openssh.com', or 'sk-ssh-ed25519@openssh.com'"></textarea>
        </div>

        <div class="mb-3">
          <button class="btn btn-primary" form="createKey" type="submit">Create</button>
        </div>

      </form>
    </div>

    <div class="table-responsive">
      <table class="table table-hover m-0">
        <thead>
          <tr>
            <th scope="col">Key</th>
            <th>Comment</th>
            <th></th>
          </tr>
        </thead>

        <tbody>
          <tr
            v-for="key in data.items"
            v-bind:key="key.key"
            style="cursor: pointer"
          >
            <td scope="row" class="text-truncate">SHA256:{{ key.key }}</td>
            <td>{{ key.comment }}</td>
            <td><button class="btn btn-danger" type="button" data-bs-toggle="modal" data-bs-target="#modalConfirmDeleteServer" v-on:click="setSelectedKey(key.id)">Delete</button></td>
          </tr>
        </tbody>
      </table>
    </div>

    <ModalConfirm
      id="modalConfirmDeleteServer"
      v-bind:title="'Delete key?'"
      v-on:confirm="handleDeleteConfirmModalConfirm"
    />
  </div>
</template>
