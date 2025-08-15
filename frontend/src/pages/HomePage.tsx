import { Navbar } from '@/components/Navbar'
import { Section } from '@/components/Section'
import { CodeCard } from '@/components/CodeCard'
import { useAsync } from '@/hooks/useAsync'
import { fetchConfessions, fetchTopConfessions } from '@/api'
import { useCallback, useEffect, useMemo, useState } from 'react'
import Prism from 'prismjs'
import { ConfessionPopup } from '@/components/ConfessionPopup'
import { navigate } from '@/navigation'

export function HomePage() {
  const [filters, setFilters] = useState<{ q?: string; language?: string; tag?: string; collection?: string }>({})
  const [selectedId, setSelectedId] = useState<number | null>(null)
  const isDesktop = useMemo(() => typeof window !== 'undefined' && window.matchMedia('(min-width: 1024px)').matches, [])

  // Close popup on route pop (Back)
  useEffect(() => {
    const onPop = () => { setSelectedId(null) }
    window.addEventListener('popstate', onPop)
    return () => window.removeEventListener('popstate', onPop)
  }, [])
  const latest = useAsync(() => {
    if (filters.collection) {
      const col = filters.collection
  if (col === 'top') return fetchTopConfessions(4)
  if (col === 'trending-weekly') return fetch(`/confessions/trending/weekly?limit=4`).then(r=>{ if(!r.ok) throw new Error('failed'); return r.json() })
  if (col === 'trending-monthly') return fetch(`/confessions/trending/monthly?limit=4`).then(r=>{ if(!r.ok) throw new Error('failed'); return r.json() })
  if (col === 'hall-of-fame') return fetch(`/confessions/hall-of-fame?limit=4`).then(r=>{ if(!r.ok) throw new Error('failed'); return r.json() })
      if (col === 'random') return fetch(`/confessions/random`).then(r=>{ if(!r.ok) throw new Error('failed'); return r.json() }).then(x=>[x])
    }
  if (filters.language) return fetchConfessions({ limit: 4, offset: 0, language: filters.language })
    if (filters.q || filters.tag) {
      const params = new URLSearchParams()
      if (filters.q) params.set('q', filters.q)
      if (filters.tag) params.set('tag', filters.tag)
  params.set('limit', '4')
      params.set('offset', '0')
      return fetch(`/confessions/search?${params.toString()}`).then(r=>{ if(!r.ok) throw new Error('search failed'); return r.json() })
    }
  return fetchConfessions({ limit: 4, offset: 0 })
  }, [filters])

  const top = useAsync(() => fetchTopConfessions(4), [])

  const handleSelect = useCallback((sel: { q?: string; language?: string; tag?: string; collection?: string }) => {
    setFilters(sel)
  }, [])

  return (
    <div className="min-h-screen bg-neutral-100 text-black">
  <Navbar onSearchSelect={handleSelect} />

      {/* Hero (black, classic, large slogan) */}
      <section className="bg-black text-white border-b border-neutral-900">
  <div className="mx-auto max-w-7xl px-3 sm:px-6 lg:px-8 py-6 sm:py-8 lg:py-12">
          <h1 className="font-bold leading-[0.95] tracking-[-0.02em]">
            <span className="block text-[9.2vw] sm:text-[7.6vw] lg:text-[6.2vw] fade-up">
              CONFESS IT <span className="opacity-80 arrow-nudge" aria-hidden="true">→</span>
            </span>
            <span className="block text-[9.2vw] sm:text-[7.6vw] lg:text-[6.2vw] pl-[2vw] sm:pl-[1.5vw] fade-up fade-up-delay-1">
              HELP OTHERS <span className="opacity-80 arrow-nudge" aria-hidden="true">→</span>
            </span>
            <span className="block text-[9.2vw] sm:text-[7.6vw] lg:text-[6.2vw] pl-[4vw] sm:pl-[3vw] fade-up fade-up-delay-2">
              MOVE ON
            </span>
          </h1>
          <div className="mt-5 sm:mt-6 lg:mt-7">
            <a href="#discover" className="inline-flex items-center gap-2 px-6 py-3 rounded-xl bg-white text-black font-semibold shadow-lg hover:shadow-xl transition-shadow">
              Read Confessions
              <span aria-hidden="true" className="arrow-nudge">→</span>
            </a>
          </div>
        </div>
      </section>

      {/* Nominees (limit 2) */}
  <section className="py-10 sm:py-12 bg-neutral-100">
  <div className="mx-auto max-w-7xl px-3 sm:px-6 lg:px-8">
          <div className="text-center mb-8">
            <div className="text-sm text-neutral-600 mb-3">Latest</div>
    <h2 className="text-[36px] sm:text-[48px] md:text-[80px] leading-none font-extrabold tracking-tight">NOMINEES</h2>
    <p className="text-neutral-600 mt-2 sm:mt-4">Vote for the latest confessions on My Dear Bug</p>
          </div>
          <Nominees />
        </div>
      </section>

      {/* Mobile tabs for Recent/Top */}
      <div className="sm:hidden mx-auto max-w-7xl px-3 sm:px-6 lg:px-8 mt-6">
        <MobileTabs
          tabs={[{ key: 'recent', label: 'Recent' }, { key: 'top', label: 'Top' }]}
          render={(active) => (
            <div className="mt-4 grid grid-cols-1 gap-5">
              {active === 'recent' ? (
                (latest.loading ? Array.from({ length: 4 }).map((_, i) => (<div key={i} className="h-64 rounded-2xl bg-white border border-neutral-200 animate-pulse" />)) :
                  latest.data?.map((item: any) => (
                    <CodeCard key={item.id} item={item} onOpen={(id)=>{ if (isDesktop) { setSelectedId(id); history.pushState({}, '', `/c/confessions/${id}`) } else { navigate(`/c/confessions/${id}`) } }} />
                  )))
              ) : (
                (top.loading ? Array.from({ length: 3 }).map((_, i) => (<div key={i} className="h-64 rounded-2xl bg-white border border-neutral-200 animate-pulse" />)) :
                  top.data?.map((item: any) => (
                    <CodeCard key={item.id} item={item} onOpen={(id)=>{ if (isDesktop) { setSelectedId(id); history.pushState({}, '', `/c/confessions/${id}`) } else { navigate(`/c/confessions/${id}`) } }} />
                  )))
              )}
            </div>
          )}
        />
      </div>

      {/* Discover section (desktop) */}
      <Section id="discover" title="Recent confessions of the day" subtitle="Fresh community confessions across languages.">
  <div className="hidden sm:grid grid-cols-1 md:grid-cols-2 gap-5 sm:gap-6">
          {latest.loading && Array.from({ length: 6 }).map((_, i) => (
            <div key={i} className="h-64 rounded-2xl bg-white border border-neutral-200 animate-pulse" />
          ))}
          {latest.error && <p className="text-red-600">Failed to load. Refresh to retry.</p>}
          {!latest.loading && !latest.error && latest.data && latest.data.length === 0 && (
            <p className="text-neutral-400">No confessions yet. Be the first to share!</p>
          )}
          {latest.data?.map((item: any) => (
            <CodeCard
              key={item.id}
              item={item}
              onOpen={(id) => {
                if (isDesktop) {
                  setSelectedId(id)
                  // reflect in URL for shareability
                  history.pushState({}, '', `/c/confessions/${id}`)
                } else {
                  navigate(`/c/confessions/${id}`)
                }
              }}
            />
          ))}
        </div>
      </Section>

      {/* Top section */}
    <Section id="top" title="Top confessions" subtitle="Most upvoted this week.">
  <div className="hidden sm:grid grid-cols-1 md:grid-cols-2 gap-5 sm:gap-6">
          {top.loading && Array.from({ length: 3 }).map((_, i) => (
            <div key={i} className="h-64 rounded-2xl bg-white border border-neutral-200 animate-pulse" />
          ))}
          {!top.loading && top.data && top.data.length === 0 && (
            <p className="text-neutral-400">No upvoted confessions yet. Come back soon!</p>
          )}
          {top.data?.map((item: any) => (
            <CodeCard
              key={item.id}
              item={item}
              onOpen={(id) => {
                if (isDesktop) {
                  setSelectedId(id)
                  history.pushState({}, '', `/c/confessions/${id}`)
                } else {
                  navigate(`/c/confessions/${id}`)
                }
              }}
            />
          ))}
        </div>
      </Section>

      {/* Footer cta */}
  <section className="py-12 sm:py-16 border-t border-neutral-200 bg-neutral-50">
  <div className="mx-auto max-w-7xl px-3 sm:px-6 lg:px-8 grid grid-cols-1 md:grid-cols-2 gap-5 sm:gap-6">
          <div className="rounded-2xl p-8 bg-neutral-50 border border-neutral-200">
    <h3 className="text-xl font-semibold">Submit your confession for visibility and recognition</h3>
            <p className="text-neutral-600 mt-2">Get featured on the home page and help others squash bugs faster.</p>
    <button onClick={()=>navigate('/submit')} className="mt-4 px-4 py-2 rounded-lg bg-black text-white font-medium shadow-sm">Share a confession</button>
          </div>
          <div className="rounded-2xl p-8 bg-neutral-50 border border-neutral-200">
            <h3 className="text-xl font-semibold">Get access to curated pro fix-sets</h3>
            <p className="text-neutral-600 mt-2">Deep dives and production-ready recipes from expert builders.</p>
    <button className="mt-4 px-4 py-2 rounded-lg border border-neutral-300 hover:bg-neutral-100">Join waitlist</button>
          </div>
        </div>
      </section>

      {/* Desktop corner popup for confession detail */}
      {selectedId !== null && (
        <ConfessionPopup
          id={selectedId}
          onClose={() => {
            setSelectedId(null)
            history.pushState({}, '', '/')
          }}
        />
      )}
    </div>
  )
}

function Nominees() {
  const nominees = useAsync(async () => {
    const res = await fetch('/confessions/trending/weekly?limit=2')
    if (!res.ok) throw new Error('failed to load nominees')
    return res.json()
  }, [])

  // highlight after load
  if (nominees.data) setTimeout(() => Prism.highlightAllUnder(document.body), 0)

  return (
    <div className="grid grid-cols-1 md:grid-cols-2 gap-5 sm:gap-6">
      {(nominees.data ?? Array.from({ length: 2 }).map((_,i)=>({id:i,title:'Loading…',description:'',language:'go',snippet:'// Loading...',tags:[],upvotes:0}))).map((c:any)=> (
  <CodeCard key={c.id} item={c} className="w-full nominee-gold" />
      ))}
    </div>
  )
}

function MobileTabs({ tabs, render }: { tabs: { key: string; label: string }[]; render: (active: string) => React.ReactNode }) {
  const [active, setActive] = useState(tabs[0]?.key || 'recent')
  return (
    <div>
      <div className="inline-flex rounded-xl border border-neutral-300 bg-white p-1">
        {tabs.map(t => (
          <button key={t.key} onClick={()=>setActive(t.key)} className={`px-4 py-2 rounded-lg text-sm font-medium ${active===t.key ? 'bg-black text-white' : 'text-neutral-700 hover:bg-neutral-100'}`}>{t.label}</button>
        ))}
      </div>
      {render(active)}
    </div>
  )
}
