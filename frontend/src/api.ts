import type { Confession, ConfessionRequest } from './types'

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

export async function createConfession(payload: ConfessionRequest) {
  const res = await fetch('/confessions', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  })
  if (!res.ok) {
    let msg = 'Failed to create confession'
    try { const data = await res.json(); if (data?.error) msg = data.error } catch {}
    throw new Error(msg)
  }
  return (await res.json()) as Confession
}

export async function fetchConfession(id: number) {
  const res = await fetch(`/confessions/${id}`)
  if (res.status === 404) throw new Error('Confession not found')
  if (!res.ok) throw new Error('Failed to fetch confession')
  return (await res.json()) as Confession
}
