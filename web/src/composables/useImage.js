import { ref } from 'vue'

import { API_URL } from '@/api'

export function useImage() {
  const images = ref([])

  async function imageReadMany() {
    const response = await fetch(`${API_URL}/images`)
    images.value = await response.json()
  }

  return {
    images,
    imageReadMany,
  }
}
