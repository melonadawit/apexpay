"use client"
import * as React from "react"
import { useRouter } from "next/navigation"
import { fetchMe, type SessionUser } from "./auth"

// useRequireAuth checks for an active dashboard session and redirects to /login if absent.
// Returns the session user (or null while checking) so pages can show a loading state.
export function useRequireAuth(): { user: SessionUser | null; checking: boolean } {
  const router = useRouter()
  const [user, setUser] = React.useState<SessionUser | null>(null)
  const [checking, setChecking] = React.useState(true)

  React.useEffect(() => {
    let cancelled = false
    ;(async () => {
      const me = await fetchMe()
      if (cancelled) return
      if (!me) {
        router.replace("/login")
        return
      }
      setUser(me)
      setChecking(false)
    })()
    return () => {
      cancelled = true
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  return { user, checking }
}
