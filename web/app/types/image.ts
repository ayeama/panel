export type Image = {
    id: string
    reference: string
    repository: string
    tag: string
    env: Map<string, string>
}

export type ImageList = {
    items: Image[]
}
