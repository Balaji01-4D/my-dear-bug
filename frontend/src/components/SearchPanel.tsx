import { useEffect, useMemo, useState } from 'react'
import { useAsync } from '@/hooks/useAsync'
import { fetchConfessions, fetchTopConfessions, fetchTags } from '@/api'
import type { Confession } from '@/types'

type Props = {
  open: boolean
  onClose: () => void
  onSelect: (sel: { q?: string; language?: string; tag?: string; collection?: string }) => void
}

export function SearchPanel({ open, onClose, onSelect }: Props) {
  const [active, setActive] = useState<'trending' | 'language' | 'tags' | 'collections'>('trending')
  const latest = useAsync(() => fetchConfessions({ limit: 50, offset: 0 }), [])
  const top = useAsync(() => fetchTopConfessions(10), [])
  const tags = useAsync(fetchTags, [])

  useEffect(() => {
    function onEsc(e: KeyboardEvent) { if (e.key === 'Escape') onClose() }
    document.addEventListener('keydown', onEsc)
    return () => document.removeEventListener('keydown', onEsc)
  }, [onClose])

  // derive languages from latest list
  const languages = useMemo(() => {
    const map = new Map<string, number>()
    ;(latest.data ?? []).forEach((c: Confession) => {
      const lang = (c.language || '').trim().toLowerCase()
      if (!lang) return
      map.set(lang, (map.get(lang) ?? 0) + 1)
    })
    return Array.from(map.entries()).sort((a,b)=>b[1]-a[1])
  }, [latest.data])

  if (!open) return null

  return (
    <div className="fixed inset-x-0 top-16 z-50">
  <div className="mx-auto max-w-6xl px-3 sm:px-6 lg:px-8">
        <div className="rounded-2xl bg-white/95 backdrop-blur border border-neutral-200 shadow-2xl">
          <div className="grid grid-cols-1 sm:grid-cols-12 gap-3 sm:gap-4 p-2 sm:p-4 max-h-[70vh] sm:max-h-[70vh] overflow-hidden">
            {/* Left rail */}
            <div className="sm:col-span-3 sm:pr-2 sm:border-r sm:border-neutral-200">
              <RailButton active={active==='trending'} onClick={()=>setActive('trending')} icon={<span className="text-neutral-700">↻</span>}>
                Trending
              </RailButton>
              <RailButton active={active==='language'} onClick={()=>setActive('language')} icon={<span className="text-neutral-700">⌘</span>}>
                By Language
              </RailButton>
              <RailButton active={active==='tags'} onClick={()=>setActive('tags')} icon={<span className="text-neutral-700">#</span>}>
                By Tags
              </RailButton>
              <RailButton active={active==='collections'} onClick={()=>setActive('collections')} icon={<span className="text-neutral-700">☰</span>}>
                Collections
              </RailButton>
            </div>

            {/* Content */}
            <div className="sm:col-span-9 min-h-0 overflow-y-auto pr-1">
              {active === 'trending' && (
                <TwoColList>
                  {(top.data ?? []).map(item => (
                    <Row key={item.id} label={item.title} count={item.upvotes} onClick={()=>{ onSelect({ q: item.title }); const s = new URLSearchParams({ q: item.title }); window.history.pushState({}, '', `/search?${s.toString()}`); window.dispatchEvent(new PopStateEvent('popstate')); onClose() }} />
                  ))}
                </TwoColList>
              )}

              {active === 'language' && (
                <TwoColList>
                  {languages.map(([name, count]) => (
                    <Row key={name} label={name} count={count} onClick={()=>{ onSelect({ language: name }); const s = new URLSearchParams({ language: name }); window.history.pushState({}, '', `/search?${s.toString()}`); window.dispatchEvent(new PopStateEvent('popstate')); onClose() }} />
                  ))}
                </TwoColList>
              )}

              {active === 'tags' && (
                <TwoColList>
                  {(tags.data ?? []).map((t) => (
                    <Row key={t.id} label={t.name} onClick={()=>{ onSelect({ tag: t.name }); const s = new URLSearchParams({ tag: t.name }); window.history.pushState({}, '', `/search?${s.toString()}`); window.dispatchEvent(new PopStateEvent('popstate')); onClose() }} />
                  ))}
                </TwoColList>
              )}

              {active === 'collections' && (
                <TwoColList>
                  <Row label="Top" onClick={()=>{ onSelect({ collection: 'top' }); const s = new URLSearchParams({ q: 'top' }); window.history.pushState({}, '', `/search?${s.toString()}`); window.dispatchEvent(new PopStateEvent('popstate')); onClose() }} />
                  <Row label="Trending (Weekly)" onClick={()=>{ onSelect({ collection: 'trending-weekly' }); const s = new URLSearchParams({ q: 'trending weekly' }); window.history.pushState({}, '', `/search?${s.toString()}`); window.dispatchEvent(new PopStateEvent('popstate')); onClose() }} />
                  <Row label="Trending (Monthly)" onClick={()=>{ onSelect({ collection: 'trending-monthly' }); const s = new URLSearchParams({ q: 'trending monthly' }); window.history.pushState({}, '', `/search?${s.toString()}`); window.dispatchEvent(new PopStateEvent('popstate')); onClose() }} />
                  <Row label="Hall of Fame" onClick={()=>{ onSelect({ collection: 'hall-of-fame' }); const s = new URLSearchParams({ q: 'hall of fame' }); window.history.pushState({}, '', `/search?${s.toString()}`); window.dispatchEvent(new PopStateEvent('popstate')); onClose() }} />
                  <Row label="Random" onClick={()=>{ onSelect({ collection: 'random' }); const s = new URLSearchParams({ q: 'random' }); window.history.pushState({}, '', `/search?${s.toString()}`); window.dispatchEvent(new PopStateEvent('popstate')); onClose() }} />
                </TwoColList>
              )}
            </div>
          </div>
        </div>
      </div>
  {/* backdrop */}
  <div className="fixed inset-0 z-40 bg-black/20" onClick={onClose} />
    </div>
  )
}

function RailButton({ active, onClick, icon, children }: { active?: boolean; onClick: () => void; icon?: React.ReactNode; children: React.ReactNode }) {
  return (
  <button onClick={onClick} className={`w-full flex items-center gap-2 px-3 py-2 rounded-lg border text-left transition ${active ? 'bg-white border-neutral-300 shadow-sm' : 'border-transparent hover:bg-neutral-50'}`}>
      <span className="inline-flex h-6 w-6 items-center justify-center rounded-md bg-neutral-200 text-xs">{icon}</span>
      <span className="text-[15px]">{children}</span>
    </button>
  )
}

function TwoColList({ children }: { children: React.ReactNode }) {
  return (
  <div className="grid grid-cols-1 md:grid-cols-2 gap-x-8 gap-y-3 p-1">
      {children}
    </div>
  )
}

function Row({ label, count, onClick }: { label: string; count?: number; onClick?: () => void }) {
  return (
  <button onClick={onClick} className="flex items-center justify-between px-2 py-2 rounded hover:bg-neutral-50 text-[15px]">
      <span className="truncate text-neutral-800">{label}</span>
      {typeof count === 'number' && <span className="text-neutral-500">{count}</span>}
    </button>
  )
}
