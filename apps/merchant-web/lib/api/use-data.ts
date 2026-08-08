"use client"
import * as React from "react"

// Minimal data-fetching hook with loading/error state and optional polling.
// Polling (e.g. every 5s for the dashboard) keeps live numbers fresh without a
// heavier SWR dependency. Callers pass a stable loader reference to avoid loops.

export function useData<T>(loader: () => Promise<T>, deps: unknown[], pollMs?: number) {
  const [data, setData] = React.useState<T | null>(null)
  const [error, setError] = React.useState<string>("")
  const [loading, setLoading] = React.useState(true)

  const runRef = React.useRef<() => Promise<void>>(async () => {})
  runRef.current = async () => {
    try {
      const result = await loader()
      setData(result)
      setError("")
    } catch (e) {
      setError((e as Error).message || "Failed to load")
    } finally {
      setLoading(false)
    }
  }

  React.useEffect(() => {
    let cancelled = false
    const run = async () => {
      try {
        const result = await loader()
        if (!cancelled) {
          setData(result)
          setError("")
        }
      } catch (e) {
        if (!cancelled) setError((e as Error).message || "Failed to load")
      } finally {
        if (!cancelled) setLoading(false)
      }
    }

    void run()
    if (pollMs && pollMs > 0) {
      const interval = setInterval(run, pollMs)
      return () => {
        cancelled = true
        clearInterval(interval)
      }
    }
    return () => {
      cancelled = true
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, deps)

  const refetch = React.useCallback(() => {
    void runRef.current()
  }, [])

  return { data, error, loading, refetch }
}

export function formatETB(value: string | number | undefined): string {
  if (value === undefined || value === null || value === "") return "0"
  const n = typeof value === "number" ? value : Number(value)
  if (Number.isNaN(n)) return "0"
  return n.toLocaleString("en-US", { maximumFractionDigits: 2 })
}
