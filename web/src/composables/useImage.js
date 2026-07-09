import { ref } from "vue";

import { API_URL } from "@/api";

export function useImage() {
    const image = ref(null)
    const images = ref([])

    async function imageRead(id) {
        try {
            const response = await fetch(`${API_URL}/images/${id}`)
            image.value = await response.json()
        } catch (error) {
            console.error(error)
        }
    }

    async function imageReadMany() {
        try {
            const response = await fetch(`${API_URL}/images`)
            images.value = await response.json()
        } catch (error) {
            console.error(error)
        }
    }

    return {
        image,
        images,
        imageRead,
        imageReadMany,
    }
}
