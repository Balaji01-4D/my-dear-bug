import type { Confession } from './types'

const withParams = (base: string, params?: Record<string, string | number | undefined>) => {
  const url = new URL(base, window.location.origin)
  if (params) Object.entries(params).forEach(([k, v]) => {
    if (v !== undefined && v !== '') url.searchParams.set(k, String(v))
  })
  return url.toString().replace(window.location.origin, '')
}

export async function fetchConfessions(params?: { limit?: number; offset?: number; language?: string }) {
  const path = params?.language
    ? withParams(`/confessions/language/${encodeURIComponent(params.language)}`, { limit: params.limit ?? 12, offset: params.offset ?? 0 })
    : withParams('/confessions', { limit: params?.limit ?? 12, offset: params?.offset ?? 0 })
  const res = await fetch(path)
  if (!res.ok) throw new Error('Failed to fetch confessions')
  return (await res.json()) as Confession[]
}

export async function fetchTopConfessions(limit = 6) {
  const res = await fetch(withParams('/confessions/top', { limit }))
  if (!res.ok) throw new Error('Failed to fetch top confessions')
  return (await res.json()) as Confession[]
}

export async function fetchRandom() {
  const res = await fetch('/confessions/random')
  if (!res.ok) throw new Error('Failed to fetch random confession')
  return (await res.json()) as Confession
}
