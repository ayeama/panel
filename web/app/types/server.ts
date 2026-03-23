export type Server = {
    id: string
    name: string
    image: string
    status: string
    ports: string[]
}

export type ServerList = {
    items: Server[]
}

export type ServerCreateRequest = {
    image: string
    env: Map<string, string>
}
