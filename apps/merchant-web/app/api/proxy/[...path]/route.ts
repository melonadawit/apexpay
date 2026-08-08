import { NextRequest, NextResponse } from "next/server"

// Server-side proxy from the merchant dashboard to the ApexPay Go API.
//
// Auth model:
//   - A logged-in dashboard user holds an opaque session token in the httpOnly
//     `apexpay_session` cookie (set by app/api/auth/login). The proxy reads that
//     cookie and injects it as a Bearer session token — the API key never lives in,
//     or reaches, the browser.
//   - If no session cookie is present, the proxy falls back to APEXPAY_API_KEY
//     (useful for local dev / server-to-server), so the app still works without login.
//
// Env:
//   APEXPAY_API_URL  Go API base (default http://api:8080)
//   APEXPAY_API_KEY  optional fallback merchant API key for dev without login

const API_BASE = process.env.APEXPAY_API_URL || "http://api:8080"
const API_KEY = process.env.APEXPAY_API_KEY || ""
const COOKIE = "apexpay_session"

export async function GET(req: NextRequest, { params }: { params: { path: string[] } }) {
  return proxy(req, params.path, "GET")
}
export async function POST(req: NextRequest, { params }: { params: { path: string[] } }) {
  return proxy(req, params.path, "POST")
}
export async function PUT(req: NextRequest, { params }: { params: { path: string[] } }) {
  return proxy(req, params.path, "PUT")
}
export async function DELETE(req: NextRequest, { params }: { params: { path: string[] } }) {
  return proxy(req, params.path, "DELETE")
}

async function proxy(req: NextRequest, path: string[], method: string) {
  const query = req.nextUrl.searchParams.toString()
  const url = `${API_BASE}/v1/${path.join("/")}${query ? `?${query}` : ""}`

  const headers: Record<string, string> = { "Content-Type": "application/json" }
  // 1) Prefer the httpOnly dashboard session cookie.
  const session = req.cookies.get(COOKIE)?.value
  if (session) {
    headers["Authorization"] = `Bearer ${session}`
  } else if (API_KEY) {
    headers["Authorization"] = `Bearer ${API_KEY}`
  }
  const idem = req.headers.get("Idempotency-Key") || req.headers.get("X-Idempotency-Key")
  if (idem) headers["Idempotency-Key"] = idem

  const body = method === "GET" || method === "DELETE" ? undefined : await req.text()

  try {
    const upstream = await fetch(url, { method, headers, body, cache: "no-store" })
    const text = await upstream.text()
    return new NextResponse(text, {
      status: upstream.status,
      headers: { "Content-Type": "application/json" },
    })
  } catch {
    return NextResponse.json(
      { error: "gateway_unreachable", message: "ApexPay API is unreachable." },
      { status: 503 }
    )
  }
}
