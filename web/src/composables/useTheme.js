import { onMounted, onUnmounted } from 'vue'

function setTheme(dark) {
  document.documentElement.setAttribute('data-bs-theme', dark ? 'dark' : 'light')
}

export function useTheme() {
  const media = window.matchMedia('(prefers-color-scheme: dark)')
  const handleTheme = (e) => setTheme(e.matches)

  onMounted(() => {
    setTheme(media.matches)
    media.addEventListener('change', handleTheme)
  })

  onUnmounted(() => {
    media.removeEventListener('change', handleTheme)
  })
}
