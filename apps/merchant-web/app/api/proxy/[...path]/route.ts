import { NextRequest, NextResponse } from "next/server"

// Server-side proxy from the merchant dashboard to the ApexPay Go API.
//
// The merchant's API key lives ONLY here (server-side, from env), never in the browser.
// The browser talks to this same-origin route (/api/proxy/...), which attaches the
// Authorization header and forwards to the Go API — so no CORS and no secret leakage.
//
// Env:
//   APEXPAY_API_URL  Go API base (default http://api:8080 for compose; use
//                    http://localhost:8080 for local dev outside compose)
//   APEXPAY_API_KEY  the merchant test/live API key that authenticates requests

const API_BASE = process.env.APEXPAY_API_URL || "http://api:8080"
const API_KEY = process.env.APEXPAY_API_KEY || ""

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
  // Merchant API key is injected server-side only.
  if (API_KEY) headers["Authorization"] = `Bearer ${API_KEY}`
  // Forward idempotency keys so mutations are safe to retry.
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
  } catch (err) {
    return NextResponse.json(
      { error: "gateway_unreachable", message: "ApexPay API is unreachable." },
      { status: 503 }
    )
  }
}
