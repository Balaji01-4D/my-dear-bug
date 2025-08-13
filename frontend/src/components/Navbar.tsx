import { useRef, useState } from 'react'
import { SearchPanel } from '@/components/SearchPanel'

type Props = {
  onSearchSelect?: (sel: { q?: string; language?: string; tag?: string; collection?: string }) => void
}

export function Navbar({ onSearchSelect }: Props) {
  const [open, setOpen] = useState(false)
  const inputRef = useRef<HTMLInputElement>(null)
  return (
    <header className="sticky top-0 z-40 bg-white/90 backdrop-blur border-b border-neutral-200">
      <div className="relative mx-auto max-w-7xl px-4 sm:px-6 lg:px-8 h-16 flex items-center gap-4">
        {/* Logo */}
        <a href="#" className="text-2xl font-extrabold tracking-tight text-black mr-1">W.</a>

        {/* Primary nav */}
        <nav className="hidden lg:flex items-center gap-6 text-[15px] text-neutral-700">
          <button
            type="button"
            aria-expanded={open}
            onClick={() => { const next=!open; setOpen(next); if (next) queueMicrotask(()=>inputRef.current?.focus()) }}
            className="group inline-flex items-center gap-1 text-neutral-800 hover:text-black"
          >
            <span className="font-medium">Explore</span>
            <svg className={`transition-transform ${open ? 'rotate-180' : ''}`} width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><polyline points="6 9 12 15 18 9"></polyline></svg>
          </button>
          <span className="inline-block w-24" aria-hidden="true"></span>
          <span className="inline-block w-24" aria-hidden="true"></span>
          <span className="inline-block w-16" aria-hidden="true"></span>
          <span className="inline-block w-16" aria-hidden="true"></span>
        </nav>

        {/* Centered Search */}
        <div className="absolute left-1/2 -translate-x-1/2 w-full max-w-xl">
          <div className="relative">
            <input
              ref={inputRef}
              onFocus={()=>setOpen(true)}
              placeholder="Search by Inspiration"
              className="w-full h-11 rounded-xl bg-neutral-200 text-neutral-900 placeholder-neutral-700 px-10 pr-4 border border-transparent focus:outline-none focus:ring-0 focus:border-neutral-300"
            />
            <svg className="absolute left-3 top-1/2 -translate-y-1/2 text-neutral-600" width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><circle cx="11" cy="11" r="8"></circle><path d="m21 21-4.3-4.3"></path></svg>
          </div>
        </div>

        {/* Right actions */}
        <div className="ml-auto hidden sm:flex items-center gap-3 text-[15px]">
          <a
            href="https://github.com/balaji01-4d/my-dear-bug"
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center gap-2 px-3 py-2 rounded-lg bg-black text-white hover:bg-black/90"
            aria-label="Open GitHub repository"
          >
            <svg viewBox="0 0 16 16" width="16" height="16" aria-hidden="true" fill="currentColor">
              <path d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.013 8.013 0 0 0 16 8c0-4.42-3.58-8-8-8Z"/>
            </svg>
            <span>GitHub</span>
          </a>
          <button className="px-3 py-2 rounded-lg border border-neutral-300 hover:bg-neutral-100">Sumbit Confession</button>
        </div>
      </div>
      <SearchPanel
        open={open}
        onClose={()=>setOpen(false)}
        onSelect={(sel)=>{ onSearchSelect?.(sel); }}
      />
    </header>
  )
}
