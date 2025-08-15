import { useEffect, useRef, useState } from 'react'
import type { Confession } from '@/types'
import { fetchConfession, postUpvote } from '@/api'
import Prism from 'prismjs'
import 'prismjs/components/prism-go'
import 'prismjs/components/prism-javascript'
import 'prismjs/components/prism-typescript'
import 'prismjs/components/prism-python'

export function ConfessionPopup({ id, onClose }: { id: number; onClose: () => void }) {
  const [data, setData] = useState<Confession | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [upvotes, setUpvotes] = useState<number>(0)
  const [voting, setVoting] = useState(false)
  const codeRef = useRef<HTMLElement | null>(null)

  useEffect(() => {
    let alive = true
    setLoading(true)
    fetchConfession(id)
      .then(c => { if (alive) { setData(c); setUpvotes(c.upvotes ?? 0); setError(null) } })
      .catch(err => { if (alive) setError(err?.message || 'Failed to load') })
      .finally(() => { if (alive) setLoading(false) })
    return () => { alive = false }
  }, [id])

  useEffect(() => {
    if (data?.snippet && codeRef.current) {
      try { Prism.highlightElement(codeRef.current) } catch {}
    }
  }, [data?.snippet])

  const lang = (data?.language || 'go').toLowerCase()
  const languageClass = `language-${lang === 'js' ? 'javascript' : lang}`
  const colorMap: Record<string, string> = {
    go: '#06b6d4', python: '#4f46e5', javascript: '#f59e0b', typescript: '#2563eb',
    rust: '#ea580c', ruby: '#dc2626', php: '#9333ea', java: '#0ea5e9', kotlin: '#7c3aed', swift: '#f97316', shell: '#10b981'
  }
  const dotColor = colorMap[lang] ?? '#111827'

  async function handleUpvote() {
    if (voting || !data) return
    setVoting(true)
    try { await postUpvote(data.id); setUpvotes(v => v + 1) } finally { setVoting(false) }
  }

  // Close on Esc and lock body scroll
  useEffect(() => {
    const onKey = (e: KeyboardEvent) => { if (e.key === 'Escape') onClose() }
    window.addEventListener('keydown', onKey)
    const prev = document.body.style.overflow
    document.body.style.overflow = 'hidden'
    return () => { window.removeEventListener('keydown', onKey); document.body.style.overflow = prev }
  }, [onClose])

  return (
    <div className="fixed inset-0 z-50 hidden lg:flex items-center justify-center">
      {/* Backdrop */}
      <div className="absolute inset-0 bg-black/70" onClick={onClose} />
      {/* Modal Shell (Instagram-like) */}
      <div className="relative z-10 w-[min(1100px,96vw)] h-[82vh] bg-white rounded-2xl overflow-hidden shadow-2xl flex">
        {/* Left: media/code area */}
        <div className="flex-1 bg-[#0a0a0a] flex items-center justify-center overflow-hidden">
      <div className="w-full h-full overflow-auto p-6">
            {data?.snippet ? (
        <pre className={`${languageClass} m-0`} style={{ background: 'transparent', fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace', fontVariantLigatures: 'none', lineHeight: 1.5, whiteSpace: 'pre' }}>
                <code ref={el => { codeRef.current = el }} className={`${languageClass} text-neutral-100`}>{data.snippet}</code>
              </pre>
            ) : (
              <div className="text-neutral-400 text-sm">{loading ? 'Loading…' : 'No snippet'}</div>
            )}
          </div>
        </div>
        {/* Right: meta + comments */}
  <div className="w-[420px] max-w-[45%] border-l border-neutral-200 flex flex-col bg-white">
          {/* Header */}
          <div className="flex items-center justify-between gap-3 p-4 border-b border-neutral-200">
            <div className="min-w-0">
              <div className="flex items-center gap-2 text-xs text-neutral-600">
                <span className="inline-flex items-center gap-2 rounded-full px-2.5 py-1 border bg-white" style={{ borderColor: dotColor }}>
                  <span className="h-2.5 w-2.5 rounded-full" style={{ backgroundColor: dotColor }} />
                  <span className="uppercase tracking-wide font-medium" style={{ color: dotColor }}>{data?.language}</span>
                </span>
                {!!upvotes && <span>• {upvotes} upvotes</span>}
              </div>
              <h3 className="mt-1 text-base font-semibold truncate" title={data?.title}>{data?.title || (loading ? 'Loading…' : 'Not found')}</h3>
            </div>
            <div className="flex items-center gap-2">
              <button onClick={handleUpvote} disabled={voting} className="rounded-md border border-neutral-300 px-3 py-1.5 text-sm hover:bg-neutral-100 disabled:opacity-60">▲</button>
              <button onClick={onClose} className="rounded-md border border-neutral-300 px-3 py-1.5 text-sm hover:bg-neutral-100">✕</button>
            </div>
          </div>
          {/* Body: description only */}
      <div className="flex-1 overflow-auto p-4">
            {error && <div className="text-red-600 text-sm mb-2">{error}</div>}
            {data?.description && (
        <div className="text-[15px] leading-relaxed text-neutral-800 whitespace-pre-wrap break-words max-w-prose">
                {data.description}
              </div>
            )}
          </div>

          {/* Footer: tags at bottom */}
          <div className="p-4 border-t border-neutral-200">
            {data?.tags?.length ? (
              <div className="flex flex-wrap gap-1.5 text-xs">
                {data.tags.map(t => (
                  <span key={t.id} className="px-2 py-0.5 rounded bg-neutral-100 border border-neutral-200">#{t.name}</span>
                ))}
              </div>
            ) : (
              <div className="text-xs text-neutral-400">No tags</div>
            )}
          </div>
        </div>
      </div>
    </div>
  )
}
