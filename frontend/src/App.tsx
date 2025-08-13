import { useEffect, useState } from 'react'
import { HomePage } from './pages/HomePage'
import { SubmitPage } from '@/pages/SubmitPage'
import { ConfessionPage } from '@/pages/ConfessionPage'
import { normalizeConfessionPath } from '@/navigation'

export default function App() {
  const [path, setPath] = useState<string>(window.location.pathname)

  useEffect(() => {
    // Normalize legacy direct hits: /confessions/:id -> /c/confessions/:id
    normalizeConfessionPath()
    const onPop = () => setPath(window.location.pathname)
    window.addEventListener('popstate', onPop)
    return () => window.removeEventListener('popstate', onPop)
  }, [])

  if (path === '/submit') return <SubmitPage />

  // Preferred SPA detail path
  const matchC = path.match(/^\/c\/confessions\/(\d+)$/)
  if (matchC) {
    const id = parseInt(matchC[1], 10)
    if (!Number.isNaN(id)) return <ConfessionPage id={id} />
  }
  // Legacy path support
  const match = path.match(/^\/confessions\/(\d+)$/)
  if (match) {
    const id = parseInt(match[1], 10)
    if (!Number.isNaN(id)) return <ConfessionPage id={id} />
  }

  return <HomePage />
}
