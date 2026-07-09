import { ref } from 'vue'

import { API_URL } from '@/api'

export function useInstance() {
  const instance = ref(null)
  const instances = ref([])

  async function instanceCreate(data) {
    try {
      const response = await fetch(`${API_URL}/instances`, {
        method: 'POST',
        body: JSON.stringify(data),
      })
      const response_data = await response.json()
      await instanceRead(response_data.instance_id)
    } catch (error) {
      console.error(error)
    }
  }

  async function instanceRead(id) {
    try {
      const response = await fetch(`${API_URL}/instances/${id}`)
      instance.value = await response.json()
    } catch (error) {
      console.error(error)
    }
  }

  async function instanceReadMany() {
    try {
      const response = await fetch(`${API_URL}/instances`)
      instances.value = await response.json()
    } catch (error) {
      console.error(error)
    }
  }

  async function instanceUpdate(id) {}

  async function instanceDelete(id) {
    try {
      const response = await fetch(`${API_URL}/instances/${id}`, {
        method: 'DELETE',
      })
    } catch (error) {
      console.error(error)
    }
  }

  async function instanceStart(id) {
    try {
      const response = await fetch(`${API_URL}/instances/${id}/start`, {
        method: 'POST',
      })
    } catch (error) {
      console.error(error)
    }

    await instanceRead(id)
  }

  async function instanceStop(id) {
    try {
      const response = await fetch(`${API_URL}/instances/${id}/stop`, {
        method: 'POST',
      })
    } catch (error) {
      console.error(error)
    }

    await instanceRead(id)
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
  }
}
