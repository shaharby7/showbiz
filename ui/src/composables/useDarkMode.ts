import { ref, watch, onMounted } from 'vue'

const STORAGE_KEY = 'showbiz-dark-mode'

type ThemeMode = 'light' | 'dark' | 'system'

const mode = ref<ThemeMode>('system')
const isDark = ref(false)

function applyTheme() {
  const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
  isDark.value = mode.value === 'dark' || (mode.value === 'system' && prefersDark)
  document.documentElement.classList.toggle('dark', isDark.value)
}

export function useDarkMode() {
  onMounted(() => {
    const stored = localStorage.getItem(STORAGE_KEY) as ThemeMode | null
    if (stored && ['light', 'dark', 'system'].includes(stored)) {
      mode.value = stored
    }
    applyTheme()

    window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', applyTheme)
  })

  watch(mode, (val) => {
    localStorage.setItem(STORAGE_KEY, val)
    applyTheme()
  })

  function toggle() {
    mode.value = isDark.value ? 'light' : 'dark'
  }

  return { mode, isDark, toggle }
}
