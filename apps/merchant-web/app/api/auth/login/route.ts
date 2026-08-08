import { NextRequest, NextResponse } from "next/server"

// Login: calls the Go API /v1/auth/login with the submitted credentials. On success,
// stores the opaque session token in an httpOnly cookie so the browser never sees it.

const API_BASE = process.env.APEXPAY_API_URL || "http://api:8080"
const COOKIE = "apexpay_session"

export async function POST(req: NextRequest) {
  let body: { email?: string; password?: string }
  try {
    body = await req.json()
  } catch {
    return NextResponse.json({ error: "invalid json" }, { status: 400 })
  }
  const email = body.email?.trim() || ""
  const password = body.password || ""
  if (!email || !password) {
    return NextResponse.json({ error: "email and password required" }, { status: 400 })
  }

  const upstream = await fetch(`${API_BASE}/v1/auth/login`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ email, password }),
  })
  const text = await upstream.text()
  const res = new NextResponse(text, {
    status: upstream.status,
    headers: { "Content-Type": "application/json" },
  })

  if (upstream.ok) {
    let data: { token?: string; expires_at?: string } | null = null
    try {
      data = JSON.parse(text)
    } catch {
      /* ignore */
    }
    if (data?.token) {
      // httpOnly + SameSite=Lax + Secure in production. expires matches the API session.
      res.cookies.set(COOKIE, data.token, {
        httpOnly: true,
        sameSite: "lax",
        path: "/",
        secure: process.env.NODE_ENV === "production",
        expires: data.expires_at ? new Date(data.expires_at) : undefined,
      })
    }
  }
  return res
}
