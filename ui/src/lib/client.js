import { ref } from 'vue'

export var authenticated = ref(false)

function setAuthenticated(value) {
  authenticated.value = value
}

export async function request(path, data = {}) {
  data.credentials = 'include'
  const response = await fetch(
    `${window.CONFIG.api.scheme}://${window.CONFIG.api.host}${window.CONFIG.api.path}${path}`,
    data,
  )
  if (response.status == 401) {
    setAuthenticated(false)
    window.location.href = '/signin'
    return
  }
  setAuthenticated(true)
  return response
}
