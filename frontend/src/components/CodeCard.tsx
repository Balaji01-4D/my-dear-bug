import Prism from 'prismjs'
import 'prismjs/components/prism-go'
import 'prismjs/components/prism-javascript'
import 'prismjs/components/prism-typescript'
import 'prismjs/components/prism-python'
import { useEffect, useState } from 'react'
import type { Confession } from '@/types'
import clsx from 'clsx'
import { navigate } from '@/navigation'

export function CodeCard({ item, className, minHeight = 220, showDescription = true, showTags = true }: { item: Confession; className?: string; minHeight?: number; showDescription?: boolean; showTags?: boolean }) {
  useEffect(() => { Prism.highlightAllUnder(document.body) }, [item.snippet])

  const lang = (item.language || 'go').toLowerCase()
  const languageClass = `language-${lang === 'js' ? 'javascript' : lang}`
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
  }
  const dotColor = colorMap[lang] ?? '#111827'
  const [upvotes, setUpvotes] = useState<number>(item.upvotes ?? 0)
  const [voting, setVoting] = useState(false)

  async function handleUpvote() {
    if (voting) return
    setVoting(true)
    try {
      const res = await fetch(`/confessions/${item.id}/upvote`, { method: 'POST' })
      if (res.ok) setUpvotes(v => v + 1)
    } catch (e) {
      // no-op
    } finally {
      setVoting(false)
    }
  }

  return (
    <article
      className={clsx('code-card group rounded-2xl bg-white border border-neutral-200 hover:border-neutral-300 transition-colors shadow-soft cursor-pointer', className)}
      onClick={() => navigate(`/c/confessions/${item.id}`)}
      role="button"
      tabIndex={0}
      onKeyDown={(e) => { if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); navigate(`/c/confessions/${item.id}`) } }}
      aria-label={`Open confession ${item.title}`}
    >
      <div className="p-5">
        <div className="relative">
          {/* Language badge */}
          <div className="absolute left-3 top-3 z-10 inline-flex items-center gap-2 rounded-full bg-black/55 px-3 py-1 text-xs font-medium text-white backdrop-blur border border-white/10 shadow-sm">
            <span className="h-2.5 w-2.5 rounded-full" style={{ backgroundColor: dotColor }} />
            <span className="uppercase tracking-wide">{item.language}</span>
          </div>
          {/* Upvote count moved near title */}
          {/* Code area (fixed height + fade) */}
          <div className="rounded-xl bg-neutral-200 p-3 relative">
            <pre
              className={clsx(languageClass)}
              style={{ minHeight, maxHeight: minHeight, overflow: 'hidden', background: 'transparent' }}
            >
              <code className={languageClass}>{item.snippet}</code>
            </pre>
            {/* Fade + ellipsis indicator */}
            <div
              className="pointer-events-none absolute inset-x-3 bottom-3 h-8"
              style={{
                background: 'linear-gradient(to bottom, rgba(229,231,235,0), rgba(229,231,235,1))'
              }}
            />
            <div className="pointer-events-none absolute right-4 bottom-2 text-neutral-500 text-sm">…</div>
          </div>
        </div>
        {/* Title + upvote button + description */}
        <div className="mt-4 flex items-start justify-between gap-3">
          <h3 className="text-black font-semibold leading-snug line-clamp-2 flex-1">{item.title}</h3>
          <button
            onClick={(e) => { e.stopPropagation(); handleUpvote() }}
            disabled={voting}
            className="shrink-0 inline-flex items-center gap-1 rounded-md border border-neutral-300 px-2.5 py-1 text-xs text-neutral-800 hover:bg-neutral-100 disabled:opacity-60"
            aria-label="Upvote confession"
            title="Upvote"
          >
            <span>▲</span>
            <span>{upvotes}</span>
          </button>
        </div>
        {showDescription && item.description && (
          <p className="text-neutral-600 text-sm line-clamp-2 mt-1">{item.description}</p>
        )}
        {/* Tags */}
        {showTags && (
          <div className="mt-3 flex flex-wrap gap-2 text-xs text-neutral-600">
            {item.tags?.slice(0, 3).map(t => (
              <span key={t.id} className="px-2 py-0.5 rounded bg-neutral-100 border border-neutral-200">#{t.name}</span>
            ))}
          </div>
        )}
      </div>
    </article>
  )
}
