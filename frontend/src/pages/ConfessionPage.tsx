import { useEffect, useMemo, useRef, useState } from 'react'
import { Navbar } from '@/components/Navbar'
import { fetchConfession, postUpvote } from '@/api'
import type { Confession } from '@/types'
import { navigate } from '@/navigation'
import Prism from 'prismjs'
import 'prismjs/components/prism-go'
import 'prismjs/components/prism-javascript'
import 'prismjs/components/prism-typescript'
import 'prismjs/components/prism-python'

export function ConfessionPage({ id }: { id: number }) {
  const [data, setData] = useState<Confession | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [upvotes, setUpvotes] = useState<number>(0)
  const [voting, setVoting] = useState(false)
  const [copied, setCopied] = useState(false)
  const codeRef = useRef<HTMLElement | null>(null)

  useEffect(() => {
    let mounted = true
    setLoading(true)
    fetchConfession(id)
      .then(c => { if (mounted) { setData(c); setUpvotes(c.upvotes ?? 0); setError(null) } })
      .catch(err => { if (mounted) setError(err?.message || 'Failed to load') })
      .finally(() => { if (mounted) setLoading(false) })
    return () => { mounted = false }
  }, [id])

  useEffect(() => {
    if (data?.title) document.title = `${data.title} – My Dear Bug`
  }, [data?.title])

  useEffect(() => {
    if (data?.snippet && codeRef.current) {
      try { Prism.highlightElement(codeRef.current) } catch {}
    }
  }, [data?.snippet])

  const lang = (data?.language || 'go').toLowerCase()
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

  async function handleUpvote() {
    if (voting || !data) return
    setVoting(true)
    try {
      await postUpvote(data.id)
      setUpvotes(v => v + 1)
    } finally {
      setVoting(false)
    }
  }

  async function handleCopyCode() {
    if (!data?.snippet) return
    try {
      await navigator.clipboard.writeText(data.snippet)
      setCopied(true)
      setTimeout(() => setCopied(false), 1200)
    } catch {}
  }

  async function handleCopyLink() {
    try {
      const url = new URL(window.location.href)
      if (!url.pathname.startsWith('/c/')) {
        url.pathname = `/c${url.pathname}`
      }
      await navigator.clipboard.writeText(url.toString())
    } catch {}
  }

  const created = useMemo(() => {
    if (!data?.createdAt) return ''
    try { return new Date(data.createdAt).toLocaleString() } catch { return data.createdAt }
  }, [data?.createdAt])

  return (
    <div className="min-h-screen bg-neutral-100 text-black">
      <Navbar />
  <div className="mx-auto max-w-3xl sm:max-w-5xl lg:max-w-6xl xl:max-w-7xl px-4 sm:px-6 lg:px-8 py-8">
        <button onClick={() => navigate('/')} className="text-sm text-neutral-600 hover:text-neutral-800">← Back to home</button>

        {loading && (
          <div className="mt-6 space-y-4">
            <div className="h-8 w-3/4 bg-neutral-200 rounded animate-pulse" />
            <div className="h-5 w-1/3 bg-neutral-200 rounded animate-pulse" />
            <div className="h-64 bg-white border border-neutral-200 rounded-2xl animate-pulse" />
          </div>
        )}
        {!loading && error && (
          <div className="mt-6 rounded-xl border border-red-200 bg-red-50 text-red-700 p-4">{error}</div>
        )}
        {!loading && data && (
          <article className="mt-6 bg-white border border-neutral-200 rounded-2xl shadow-soft">
            <header className="p-6 sm:p-8 border-b border-neutral-200">
              <div className="flex items-center justify-between gap-3 flex-wrap">
                <div className="flex items-center gap-3 text-xs text-neutral-600">
                  <span className="inline-flex items-center gap-2 rounded-full px-3 py-1 border bg-white"
                        style={{ borderColor: dotColor }}>
                    <span className="h-2.5 w-2.5 rounded-full" style={{ backgroundColor: dotColor }} />
                    <span className="uppercase tracking-wide font-medium" style={{ color: dotColor }}>{data.language}</span>
                  </span>
                  <span>•</span>
                  <time dateTime={data.createdAt}>{created}</time>
                  <span>•</span>
                  <span>{upvotes} upvotes</span>
                </div>
                <div className="flex items-center gap-2">
                  <button
                    onClick={handleCopyLink}
                    className="inline-flex items-center gap-1 rounded-md border border-neutral-300 px-3 py-1.5 text-sm text-neutral-800 hover:bg-neutral-100"
                    aria-label="Copy link"
                  >
                    Copy link
                  </button>
                  <button
                    onClick={handleUpvote}
                    disabled={voting}
                    className="inline-flex items-center gap-1 rounded-md border border-neutral-300 px-3 py-1.5 text-sm text-neutral-800 hover:bg-neutral-100 disabled:opacity-60"
                    aria-label="Upvote confession"
                    title="Upvote"
                  >
                    ▲ Upvote
                  </button>
                </div>
              </div>
              <h1 className="mt-3 text-3xl sm:text-4xl font-extrabold tracking-tight leading-tight">{data.title}</h1>
              {data.tags?.length ? (
                <div className="mt-3 flex flex-wrap gap-2 text-xs">
                  {data.tags.map(t => (
                    <button
                      key={t.id}
                      onClick={() => navigate(`/?tag=${encodeURIComponent(t.name)}`)}
                      className="px-2 py-0.5 rounded bg-neutral-100 border border-neutral-300 text-neutral-800 shadow-sm hover:bg-neutral-200"
                    >
                      #{t.name}
                    </button>
                  ))}
                </div>
              ) : null}
            </header>

            {data.snippet && (
              <div className="p-6 sm:p-8 border-b border-neutral-200 bg-neutral-50">
                <div className="rounded-xl bg-neutral-200 p-3 relative">
                  <button
                    onClick={handleCopyCode}
                    className="absolute right-3 top-3 inline-flex items-center gap-1 rounded-md border border-neutral-300 bg-white/70 backdrop-blur px-2.5 py-1 text-xs text-neutral-800 hover:bg-white"
                    aria-label="Copy code"
                  >
                    {copied ? 'Copied' : 'Copy'}
                  </button>
                  <div className="overflow-x-auto">
                    <pre className={languageClass} style={{ background: 'transparent' }}>
                      <code ref={el => { codeRef.current = el }} className={languageClass}>{data.snippet}</code>
                    </pre>
                  </div>
                </div>
              </div>
            )}

            <section className="p-6 sm:p-8 prose prose-neutral max-w-none">
              <p className="whitespace-pre-wrap text-neutral-800 leading-relaxed">{data.description}</p>
            </section>
          </article>
        )}
      </div>
    </div>
  )
}
