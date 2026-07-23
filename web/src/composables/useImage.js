import { ref } from 'vue'

import { API_URL } from '@/api'

export function useImage() {
  const image = ref(null)
  const images = ref([])

  async function imageRead(id) {
    const response = await fetch(`${API_URL}/images/${id}`)
    image.value = await response.json()
  }

  async function imageReadMany() {
    const response = await fetch(`${API_URL}/images`)
    images.value = await response.json()
  }

  return {
    image,
    images,
    imageRead,
    imageReadMany,
  }
}
