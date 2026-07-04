import { ref } from "vue";

export function useInstance() {
  const instances = ref([])

  async function instanceCreate() { }

  async function instanceRead(id) {
    try {
      const response = await fetch("https://localhost:8000/instances")
      instances.value = await response.json()
    } catch (error) {
      console.log(error)
    }
  }

  async function instanceUpdate(id) { }

  async function instanceDelete(id) { }

  return {
    instances,
    instanceCreate,
    instanceRead,
    instanceUpdate,
    instanceDelete
  }
}
