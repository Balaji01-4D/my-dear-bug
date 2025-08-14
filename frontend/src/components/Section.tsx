import { ReactNode } from 'react'

export function Section({ id, title, subtitle, actions, children }: { id?: string, title: string, subtitle?: string, actions?: ReactNode, children: ReactNode }) {
  return (
  <section id={id} className="py-12">
  <div className="mx-auto max-w-7xl px-3 sm:px-6 lg:px-8">
        <div className="flex flex-col sm:flex-row sm:items-end sm:justify-between gap-4 mb-6">
          <div>
      <h2 className="text-2xl font-semibold tracking-tight">{title}</h2>
      {subtitle && <p className="text-neutral-600 mt-1">{subtitle}</p>}
          </div>
          <div className="sm:text-right">{actions}</div>
        </div>
        {children}
      </div>
    </section>
  )
}
