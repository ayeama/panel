import { ref } from 'vue'

import { API_URL } from '@/api'

export function useInstance() {
  const instance = ref(null)
  const instances = ref([])

  async function instanceCreate() {}

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

  async function instanceDelete(id) {}

  return {
    instance,
    instances,
    instanceCreate,
    instanceRead,
    instanceReadMany,
    instanceUpdate,
    instanceDelete,
  }
}
