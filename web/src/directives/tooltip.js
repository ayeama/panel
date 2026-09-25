export default {
  mounted(el, binding) {
    if (binding.value) {
      el.setAttribute('data-bs-title', binding.value)
    }
    el._tooltip = new bootstrap.Tooltip(el)
  },

  unmounted(el) {
    el._tooltip?.dispose()
  },
}
