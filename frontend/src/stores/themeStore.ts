import { create } from 'zustand'

interface ThemeState {
  isDarkMode: boolean
  toggleDarkMode: () => void
  setDarkMode: (isDark: boolean) => void
}

export const useThemeStore = create<ThemeState>((set) => ({
  isDarkMode: localStorage.getItem('theme') === 'dark',

  toggleDarkMode: () => {
    set((state) => {
      const newDarkMode = !state.isDarkMode
      localStorage.setItem('theme', newDarkMode ? 'dark' : 'light')
      applyTheme(newDarkMode)
      return { isDarkMode: newDarkMode }
    })
  },

  setDarkMode: (isDark: boolean) => {
    localStorage.setItem('theme', isDark ? 'dark' : 'light')
    applyTheme(isDark)
    set({ isDarkMode: isDark })
  },
}))

export function applyTheme(isDark: boolean) {
  if (isDark) {
    document.documentElement.classList.add('dark-mode')
  } else {
    document.documentElement.classList.remove('dark-mode')
  }
}

// 初始化主题
export function initTheme() {
  const isDarkMode = localStorage.getItem('theme') === 'dark'
  applyTheme(isDarkMode)
}
