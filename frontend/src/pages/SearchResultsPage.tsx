import { useEffect, useMemo, useState } from 'react'
import { Navbar } from '@/components/Navbar'
import { CodeCard } from '@/components/CodeCard'
import { ConfessionPopup } from '@/components/ConfessionPopup'
import { navigate } from '@/navigation'
import { fetchSearch } from '@/api'

export function SearchResultsPage() {
  const [locSearch, setLocSearch] = useState<string>(() => window.location.search)
  useEffect(() => {
    const onPop = () => setLocSearch(window.location.search)
    window.addEventListener('popstate', onPop)
    return () => window.removeEventListener('popstate', onPop)
  }, [])

  const query = useMemo(() => {
    const p = new URLSearchParams(locSearch)
    return {
      q: p.get('q') || undefined,
      language: p.get('language') || undefined,
      tag: p.get('tag') || undefined,
    }
  }, [locSearch])
  const [data, setData] = useState<any[] | null>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [selectedId, setSelectedId] = useState<number | null>(null)
  const isDesktop = typeof window !== 'undefined' && window.matchMedia('(min-width: 1024px)').matches

  useEffect(() => {
    let alive = true
    setLoading(true)
    fetchSearch({ ...query, limit: 30, offset: 0 })
      .then(list => { if (alive) { setData(list); setError(null) } })
      .catch(err => { if (alive) setError(err?.message || 'Failed to search') })
      .finally(() => { if (alive) setLoading(false) })
    return () => { alive = false }
  }, [query.q, query.language, query.tag])

  useEffect(() => { document.title = (query.q || query.tag || query.language ? `${query.q || query.tag || query.language} – Search` : 'Search') + ' – My Dear Bug' }, [query])

  const subtitle = useMemo(() => {
    const parts = [] as string[]
    if (query.q) parts.push(`for “${query.q}”`)
    if (query.language) parts.push(`in ${query.language}`)
    if (query.tag) parts.push(`#${query.tag}`)
    return parts.length ? parts.join(' · ') : 'Explore recent results'
  }, [query])

  return (
    <div className="min-h-screen bg-neutral-100 text-black">
      <Navbar />
      <div className="mx-auto max-w-7xl px-3 sm:px-6 lg:px-8 py-6 sm:py-8">
        <div className="mb-5 sm:mb-6">
          <h1 className="text-2xl sm:text-3xl font-bold">Search results</h1>
          <p className="text-neutral-600 mt-1">{subtitle}</p>
        </div>

        {loading && (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-5 sm:gap-6">
            {Array.from({ length: 6 }).map((_, i) => (
              <div key={i} className="h-64 rounded-2xl bg-white border border-neutral-200 animate-pulse" />
            ))}
          </div>
        )}
        {error && <div className="text-red-600">{error}</div>}
        {!loading && !error && (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-5 sm:gap-6">
            {(data ?? []).map((item: any) => (
              <CodeCard
                key={item.id}
                item={item}
                onOpen={(id) => {
                  if (isDesktop) {
                    setSelectedId(id)
                    // Update URL to shareable detail without rerouting the SPA
                    history.pushState({}, '', `/c/confessions/${id}`)
                  } else {
                    navigate(`/c/confessions/${id}`)
                  }
                }}
              />
            ))}
            {(!data || data.length === 0) && (
              <div className="text-neutral-500">No results found.</div>
            )}
          </div>
        )}
        {selectedId !== null && (
          <ConfessionPopup
            id={selectedId}
            onClose={() => {
              setSelectedId(null)
              // Restore search URL without triggering router change
              const base = '/search' + (window.location.search || '')
              history.pushState({}, '', base)
            }}
          />
        )}
      </div>
    </div>
  )
}
