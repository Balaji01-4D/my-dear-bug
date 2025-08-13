import { useEffect, useState } from 'react'

export function useAsync<T>(fn: () => Promise<T>, deps: any[] = []) {
  const [data, setData] = useState<T | null>(null)
  const [error, setError] = useState<Error | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    let active = true
    setLoading(true)
    setError(null)
    fn()
      .then((d) => active && setData(d))
      .catch((e) => active && setError(e as Error))
      .finally(() => active && setLoading(false))
    return () => { active = false }
  // eslint-disable-next-line react-hooks/exhaustive-deps
  }, deps)

  return { data, error, loading }
}
