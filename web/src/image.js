export function imageRegistry(image) {
  return image.split(':').at(0)
}

export function imageTag(image) {
  return image.split(':').at(1)
}

export function imageLabel(image) {
  return image.split('/').at(-1)
}
