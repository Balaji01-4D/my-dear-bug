import { Navbar } from '@/components/Navbar'
import { Section } from '@/components/Section'
import { CodeCard } from '@/components/CodeCard'
import { useAsync } from '@/hooks/useAsync'
import { fetchConfessions, fetchTopConfessions } from '@/api'
import { useCallback, useState } from 'react'
import Prism from 'prismjs'

export function HomePage() {
  const [filters, setFilters] = useState<{ q?: string; language?: string; tag?: string; collection?: string }>({})
  const latest = useAsync(() => {
    if (filters.collection) {
      const col = filters.collection
      if (col === 'top') return fetchTopConfessions(12)
      if (col === 'trending-weekly') return fetch(`/confessions/trending/weekly?limit=12`).then(r=>{ if(!r.ok) throw new Error('failed'); return r.json() })
      if (col === 'trending-monthly') return fetch(`/confessions/trending/monthly?limit=12`).then(r=>{ if(!r.ok) throw new Error('failed'); return r.json() })
      if (col === 'hall-of-fame') return fetch(`/confessions/hall-of-fame?limit=12`).then(r=>{ if(!r.ok) throw new Error('failed'); return r.json() })
      if (col === 'random') return fetch(`/confessions/random`).then(r=>{ if(!r.ok) throw new Error('failed'); return r.json() }).then(x=>[x])
    }
    if (filters.language) return fetchConfessions({ limit: 12, offset: 0, language: filters.language })
    if (filters.q || filters.tag) {
      const params = new URLSearchParams()
      if (filters.q) params.set('q', filters.q)
      if (filters.tag) params.set('tag', filters.tag)
      params.set('limit', '12')
      params.set('offset', '0')
      return fetch(`/confessions/search?${params.toString()}`).then(r=>{ if(!r.ok) throw new Error('search failed'); return r.json() })
    }
    return fetchConfessions({ limit: 12, offset: 0 })
  }, [filters])

  const top = useAsync(() => fetchTopConfessions(6), [])

  const handleSelect = useCallback((sel: { q?: string; language?: string; tag?: string; collection?: string }) => {
    setFilters(sel)
  }, [])

  return (
    <div className="min-h-screen bg-neutral-100 text-black">
  <Navbar onSearchSelect={handleSelect} />

      {/* Hero */}
      <section className="hero-gradient border-b border-neutral-200 bg-white">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 py-16 lg:py-24 grid lg:grid-cols-2 gap-10 items-center">
          <div>
            <span className="inline-flex items-center text-xs uppercase tracking-widest text-neutral-500">Site of the Day</span>
            <h1 className="mt-3 text-6xl sm:text-7xl font-extrabold leading-[1.05]">
              METAMASK
            </h1>
            <p className="mt-4 text-neutral-700 text-lg max-w-prose">Antinomy Studio • Aug 13, 2025 • Score 7.39 of 10</p>
            <div className="mt-6 flex gap-3">
              <a href="#discover" className="px-4 py-2 rounded-lg bg-black text-white font-medium">Explore</a>
              <a href="#" className="px-4 py-2 rounded-lg border border-neutral-300">Submit a snippet</a>
            </div>
          </div>
          <div className="hidden lg:block">
            <div className="rounded-2xl bg-neutral-50 border border-neutral-200 p-4 shadow-soft">
              <div className="grid grid-cols-1 gap-4">
                {/* Preview skeletons replaced when data loads */}
                {(latest.data ?? Array.from({ length: 3 }).map((_,i)=>({id:i,title:'',description:'',language:'go',snippet:'package main\nfunc main(){}',tags:[],upvotes:0} as any))).slice(0,3).map((c:any)=> (
                  <div key={c.id} className="rounded-xl bg-white border border-neutral-200 overflow-hidden">
                    <pre className="m-0 p-4 text-xs text-neutral-700"><code>{(c.snippet ?? '').split('\n').slice(0,6).join('\n')}</code></pre>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Nominees (limit 2) */}
      <section className="py-12 bg-neutral-100">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between mb-6">
            <h2 className="text-4xl md:text-5xl font-extrabold tracking-tight">NOMINEES</h2>
            <p className="text-neutral-600 hidden md:block">Vote for the latest snippets on My Dear Bug</p>
          </div>
          <Nominees />
        </div>
      </section>

      {/* Discover section */}
      <Section id="discover" title="Recent bites of the day" subtitle="Fresh community snippets across languages.">
        <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-6">
          {latest.loading && Array.from({ length: 6 }).map((_, i) => (
            <div key={i} className="h-64 rounded-2xl bg-white border border-neutral-200 animate-pulse" />
          ))}
          {latest.error && <p className="text-red-600">Failed to load. Refresh to retry.</p>}
          {!latest.loading && !latest.error && latest.data && latest.data.length === 0 && (
            <p className="text-neutral-400">No snippets yet. Be the first to share!</p>
          )}
          {latest.data?.map(item => (
            <CodeCard key={item.id} item={item} />
          ))}
        </div>
      </Section>

      {/* Top section */}
      <Section id="top" title="Top snippets" subtitle="Most upvoted this week.">
        <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-6">
          {top.loading && Array.from({ length: 3 }).map((_, i) => (
            <div key={i} className="h-64 rounded-2xl bg-white border border-neutral-200 animate-pulse" />
          ))}
          {!top.loading && top.data && top.data.length === 0 && (
            <p className="text-neutral-400">No upvoted snippets yet. Come back soon!</p>
          )}
          {top.data?.map(item => (
            <CodeCard key={item.id} item={item} />
          ))}
        </div>
      </Section>

      {/* Footer cta */}
      <section className="py-16 border-t border-neutral-200 bg-white">
        <div className="mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 grid md:grid-cols-2 gap-6">
          <div className="rounded-2xl p-8 bg-neutral-50 border border-neutral-200">
            <h3 className="text-xl font-semibold">Submit your snippet for visibility and recognition</h3>
            <p className="text-neutral-600 mt-2">Get featured on the home page and help others squash bugs faster.</p>
            <button className="mt-4 px-4 py-2 rounded-lg bg-black text-white font-medium">Share a snippet</button>
          </div>
          <div className="rounded-2xl p-8 bg-neutral-50 border border-neutral-200">
            <h3 className="text-xl font-semibold">Get access to curated pro fix-sets</h3>
            <p className="text-neutral-600 mt-2">Deep dives and production-ready recipes from expert builders.</p>
            <button className="mt-4 px-4 py-2 rounded-lg border border-neutral-300">Join waitlist</button>
          </div>
        </div>
      </section>
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
    <div className="grid md:grid-cols-2 gap-8">
      {(nominees.data ?? Array.from({ length: 2 }).map((_,i)=>({id:i,title:'Loading…',description:'',language:'go',snippet:'// Loading...',tags:[],upvotes:0}))).map((c:any)=> {
        const lang = (c.language || 'code').toLowerCase();
        const colorMap: Record<string, string> = {
          go: '#06b6d4',
          python: '#4f46e5',
          javascript: '#f59e0b',
          typescript: '#2563eb',
          rust: '#ea580c',
          ruby: '#dc2626',
          php: '#9333ea',
          java: '#0ea5e9',
          kotlin: '#7c3aed',
          swift: '#f97316',
          shell: '#10b981'
        };
        const dotColor = colorMap[lang] ?? '#111827';
        return (
          <article key={c.id} className="rounded-2xl bg-white border border-neutral-200 p-4 theme-monokai">
            <div className="relative">
              {/* Language badge */}
              <div className="absolute left-3 top-3 z-10 inline-flex items-center gap-2 rounded-full bg-black/55 px-3 py-1 text-xs font-medium text-white backdrop-blur border border-white/10 shadow-sm">
                <span className="h-2.5 w-2.5 rounded-full" style={{ backgroundColor: dotColor }} />
                <span className="uppercase tracking-wide">{c.language}</span>
              </div>
              {/* Upvotes badge */}
              <div className="absolute right-3 top-3 z-10 rounded-full bg-black/45 px-2.5 py-1 text-xs text-white border border-white/10 shadow-sm">▲ {c.upvotes}</div>
              <pre className={`language-${lang}`} style={{minHeight: '320px'}}>
                <code className={`language-${lang}`}>{c.snippet}</code>
              </pre>
            </div>
            {/* Bottom title */}
            <div className="mt-4 flex items-start justify-between gap-3">
              <h3 className="text-lg md:text-xl font-semibold leading-snug">{c.title}</h3>
            </div>
          </article>
        );
      })}
    </div>
  )
}
