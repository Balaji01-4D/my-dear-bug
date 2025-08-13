import Prism from 'prismjs'
import 'prismjs/components/prism-go'
import 'prismjs/components/prism-javascript'
import 'prismjs/components/prism-typescript'
import 'prismjs/components/prism-python'
import { useEffect } from 'react'
import type { Confession } from '@/types'
import clsx from 'clsx'

export function CodeCard({ item, className }: { item: Confession; className?: string }) {
  useEffect(() => { Prism.highlightAllUnder(document.body) }, [item.snippet])

  const lang = (item.language || 'go').toLowerCase()
  const languageClass = `language-${lang === 'js' ? 'javascript' : lang}`

  return (
    <article className={clsx('code-card group rounded-2xl bg-white border border-neutral-200 hover:border-neutral-300 transition-colors shadow-soft', className)}>
      <div className="p-4 flex items-start justify-between">
        <div>
          <h3 className="text-black font-semibold leading-tight line-clamp-1">{item.title}</h3>
          <p className="text-neutral-600 text-sm line-clamp-2 mt-1">{item.description}</p>
        </div>
        <span className="ml-3 text-xs px-2 py-1 rounded-md bg-neutral-100 text-neutral-700 border border-neutral-200">{item.language}</span>
      </div>
      <pre className={clsx('overflow-hidden m-4 mt-0 rounded-xl', languageClass)}>
        <code className={languageClass}>{item.snippet}</code>
      </pre>
      <div className="px-4 pb-4 flex items-center justify-between text-xs text-neutral-600">
        <div className="flex flex-wrap gap-2">
          {item.tags?.slice(0, 3).map(t => (
            <span key={t.id} className="px-2 py-0.5 rounded bg-neutral-100 border border-neutral-200">#{t.name}</span>
          ))}
        </div>
        <span className="ml-2">▲ {item.upvotes}</span>
      </div>
    </article>
  )
}
