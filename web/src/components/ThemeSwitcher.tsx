import { useTheme } from '../hooks/useTheme'

export function ThemeSwitcher() {
  const { theme, setTheme } = useTheme()

  return (
    <div className="join" role="group" aria-label="Theme">
      <button
        type="button"
        className={`btn btn-sm join-item ${theme === 'light' ? 'btn-active' : 'btn-ghost'}`}
        onClick={() => setTheme('light')}
        aria-pressed={theme === 'light'}
        aria-label="Light theme"
      >
        Light
      </button>
      <button
        type="button"
        className={`btn btn-sm join-item ${theme === 'dark' ? 'btn-active' : 'btn-ghost'}`}
        onClick={() => setTheme('dark')}
        aria-pressed={theme === 'dark'}
        aria-label="Dark theme"
      >
        Dark
      </button>
    </div>
  )
}
