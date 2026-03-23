import type { ImageList, Image } from '~/types/image'

const api = useApi()

export function useImages() {
    const read = async (): Promise<Image[]> => {
      const response = await $fetch<ImageList>(api.http('/images'))
      return response.items
    }

  return {
    read,
  }
}
