import type { Confession, ConfessionRequest, Tag } from './types'

const withParams = (base: string, params?: Record<string, string | number | undefined>) => {
  const url = new URL(base, window.location.origin)
  if (params) Object.entries(params).forEach(([k, v]) => {
    if (v !== undefined && v !== '') url.searchParams.set(k, String(v))
  })
  return url.toString().replace(window.location.origin, '')
}

// Simple in-memory GET cache with TTL and in-flight de-duplication
type CacheEntry<T> = { data?: T; expiry: number; promise?: Promise<T> }
const cache = new Map<string, CacheEntry<any>>()

function cacheKey(method: string, path: string) {
  return `${method} ${path}`
}

function now() { return Date.now() }

function getCached<T>(path: string, ttlMs: number, fetcher: () => Promise<T>): Promise<T> {
  const key = cacheKey('GET', path)
  const hit = cache.get(key)
  if (hit && hit.data !== undefined && hit.expiry > now()) {
    return Promise.resolve(hit.data as T)
  }
  if (hit?.promise) return hit.promise as Promise<T>
  const p = fetcher().then((data) => {
    cache.set(key, { data, expiry: now() + ttlMs })
    return data
  }).finally(() => {
    const cur = cache.get(key)
    if (cur) delete cur.promise
  })
  cache.set(key, { promise: p, expiry: now() + ttlMs })
  return p
}

export function invalidateCache(prefix?: string) {
  if (!prefix) { cache.clear(); return }
  const pref = prefix.trim()
  for (const k of Array.from(cache.keys())) {
    if (k.includes(pref)) cache.delete(k)
  }
}

export async function fetchConfessions(params?: { limit?: number; offset?: number; language?: string }) {
  const path = params?.language
    ? withParams(`/confessions/language/${encodeURIComponent(params.language)}`, { limit: params.limit ?? 12, offset: params.offset ?? 0 })
    : withParams('/confessions', { limit: params?.limit ?? 12, offset: params?.offset ?? 0 })
  return getCached<Confession[]>(path, 30_000, async () => {
    const res = await fetch(path)
    if (!res.ok) throw new Error('Failed to fetch confessions')
    return res.json()
  })
}

export async function fetchTopConfessions(limit = 6) {
  const path = withParams('/confessions/top', { limit })
  return getCached<Confession[]>(path, 60_000, async () => {
    const res = await fetch(path)
    if (!res.ok) throw new Error('Failed to fetch top confessions')
    return res.json()
  })
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
  const created = (await res.json()) as Confession
  // Invalidate relevant cached lists
  invalidateCache('/confessions')
  return created
}

export async function fetchConfession(id: number) {
  const path = `/confessions/${id}`
  return getCached<Confession>(path, 60_000, async () => {
    const res = await fetch(path)
    if (res.status === 404) throw new Error('Confession not found')
    if (!res.ok) throw new Error('Failed to fetch confession')
    return res.json()
  })
}

export async function fetchTags() : Promise<Tag[]> {
  const path = '/tags'
  return getCached<Tag[]>(path, 300_000, async () => {
    const res = await fetch(path)
    if (!res.ok) throw new Error('Failed to load tags')
    return res.json()
  })
}

export async function postUpvote(id: number) {
  const res = await fetch(`/confessions/${id}/upvote`, { method: 'POST' })
  if (!res.ok) throw new Error('Failed to upvote')
  // Invalidate affected entries: detail and lists/top
  invalidateCache(`/confessions/${id}`)
  invalidateCache('/confessions/top')
  invalidateCache('/confessions')
}
