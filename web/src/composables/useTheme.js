import { ref, onMounted, onUnmounted } from 'vue'

export const Theme = Object.freeze({
  LIGHT: 'light',
  DARK: 'dark',
  SYSTEM: 'system',
})

export function useTheme() {
  const theme = ref(Theme.SYSTEM)

  function getTheme() {
    return theme.value
  }

  function setTheme(value) {
    theme.value = value
    applyTheme()
  }

  const media = window.matchMedia('(prefers-color-scheme: dark)')

  function applyTheme() {
    const dark = theme.value == Theme.DARK || (theme.value == Theme.SYSTEM && media.matches)
    document.documentElement.setAttribute('data-bs-theme', dark ? Theme.DARK : Theme.LIGHT)
  }

  function handleTheme() {
    if (theme.value == Theme.SYSTEM) {
      applyTheme()
    }
  }

  onMounted(() => {
    applyTheme()
    media.addEventListener('change', handleTheme)
  })

  onUnmounted(() => {
    media.removeEventListener('change', handleTheme)
  })

  return {
    getTheme,
    setTheme,
  }
}
