"use client"

// Client helpers for dashboard auth. Login/logout hit the server routes that manage the
// httpOnly session cookie; /auth/me confirms an active session to the browser.

export type MerchantContext = { merchant_id: string; legal_name: string; role: string }
export type SessionUser = {
  user?: { id: string; email: string; name: string; status: string }
  merchant?: MerchantContext
  merchants?: MerchantContext[]
}

export async function login(email: string, password: string): Promise<SessionUser> {
  const res = await fetch("/api/auth/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  })
  if (!res.ok) {
    let detail = "Login failed"
    try {
      const body = await res.json()
      detail = body.message || body.error || detail
    } catch {
      /* ignore */
    }
    throw new Error(detail)
  }
  return res.json()
}

export async function logout(): Promise<void> {
  await fetch("/api/auth/logout", { method: "POST" })
}

export async function fetchMe(): Promise<SessionUser | null> {
  const res = await fetch("/api/proxy/auth/me", { cache: "no-store" })
  if (!res.ok) return null
  return res.json()
}
