import type { ServerList, Server, ServerCreateRequest } from '~/types/server'

const api = useApi()

export function useServers() {
  const create = async (server: ServerCreateRequest): Promise<Server> => {
    return await $fetch<Server>(api.http('/servers'), {method: "POST", body: server})
  }

  const read = async (): Promise<Server[]> => {
    const response = await $fetch<ServerList>(api.http('/servers'))
    return response.items
  }

  const readOne = async (id: string): Promise<Server> => {
    return await $fetch<Server>(api.http(`/servers/${id}`))
  }

  // NOTE 'delete' is an operator
  const deleteOne = async (id: string) => {
    await $fetch<Server>(api.http(`/servers/${id}`), {method: "DELETE"})
  }

  const start = async (id: string): Promise<Server> => {
    return await $fetch<Server>(api.http(`/servers/${id}/start`), {method: "POST"})
  }

  const stop = async (id: string): Promise<Server> => {
    return await $fetch<Server>(api.http(`/servers/${id}/stop`), {method: "POST"})
  }

  return {
    create,
    read,
    readOne,
    deleteOne,
    start,
    stop,
  }
}
