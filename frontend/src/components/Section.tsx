import { ReactNode } from 'react'

export function Section({ id, title, subtitle, actions, children }: { id?: string, title: string, subtitle?: string, actions?: ReactNode, children: ReactNode }) {
  return (
  <section id={id} className="py-12">
  <div className="mx-auto max-w-7xl px-3 sm:px-6 lg:px-8">
        <div className="flex flex-col sm:flex-row sm:items-end sm:justify-between gap-4 mb-6">
          <div>
      <h2 className="text-2xl font-semibold tracking-tight">
            <span className="bg-gradient-to-r from-neutral-900 via-neutral-800 to-neutral-500 bg-clip-text text-transparent">
              {title}
            </span>
          </h2>
  {subtitle && <p className="text-neutral-700 mt-1 border-l-2 border-neutral-400 pl-3 font-medium">{subtitle}</p>}
          </div>
          <div className="sm:text-right">{actions}</div>
        </div>
        {children}
      </div>
    </section>
  )
}
