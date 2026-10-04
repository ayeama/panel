import { ref } from 'vue'

import { API_ERROR_NOT_FOUND, API_URL } from '@/api'

export function useInstance() {
  const instance = ref(null)
  const instances = ref([])

  async function instanceCreate(data) {
    const response = await fetch(`${API_URL}/instances`, {
      method: 'POST',
      body: JSON.stringify(data),
    })

    instance.value = await response.json()
  }

  async function instanceRead(id) {
    const response = await fetch(`${API_URL}/instances/${id}`)

    if (response.status === 404) {
      throw new Error(API_ERROR_NOT_FOUND)
    }

    instance.value = await response.json()
  }

  async function instanceReadMany() {
    const response = await fetch(`${API_URL}/instances`)
    instances.value = await response.json()
  }

  async function instanceUpdate(id) {}

  async function instanceDelete(id) {
    const response = await fetch(`${API_URL}/instances/${id}`, {
      method: 'DELETE',
    })
  }

  async function instanceStart(id) {
    const response = await fetch(`${API_URL}/instances/${id}/start`, {
      method: 'POST',
    })
    await instanceRead(id)
  }

  async function instanceStop(id) {
    const response = await fetch(`${API_URL}/instances/${id}/stop`, {
      method: 'POST',
    })
    await instanceRead(id)
  }

  async function instanceRestore(id, data) {
    const response = await fetch(`${API_URL}/instances/${id}/restore`, {
      method: 'POST',
      body: data,
    })
  }

  return {
    instance,
    instances,
    instanceCreate,
    instanceRead,
    instanceReadMany,
    instanceUpdate,
    instanceDelete,
    instanceStart,
    instanceStop,
    instanceRestore,
  }
}
