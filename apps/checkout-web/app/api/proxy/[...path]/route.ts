import { NextRequest, NextResponse } from "next/server"

// Server-side proxy from the checkout UI to the ApexPay Go API.
//
// The browser talks only to this same-origin route (`/api/proxy/...`) so it never needs
// CORS or the backend host. The Go API base is read server-side, so it is not exposed to
// the browser.
//
// Env: APEXPAY_API_URL (defaults to the compose service name `http://api:8080`). In local
// dev outside compose, set it to `http://localhost:8080`.

const API_BASE = process.env.APEXPAY_API_URL || "http://api:8080"

export async function GET(req: NextRequest, { params }: { params: { path: string[] } }) {
  return proxy(req, params.path, "GET")
}

export async function POST(req: NextRequest, { params }: { params: { path: string[] } }) {
  return proxy(req, params.path, "POST")
}

async function proxy(req: NextRequest, path: string[], method: string) {
  const query = req.nextUrl.searchParams.toString()
  const url = `${API_BASE}/${path.join("/")}${query ? `?${query}` : ""}`

  const headers: Record<string, string> = {
    // Forward the caller's idempotency key if present; the public checkout routes need
    // no authorization header (the payment-link token is the capability).
  }
  if (req.headers.get("Idempotency-Key")) headers["Idempotency-Key"] = req.headers.get("Idempotency-Key")!

  const body = method === "POST" ? await req.text() : undefined

  try {
    const upstream = await fetch(url, {
      method,
      headers,
      body,
      cache: "no-store",
    })
    const text = await upstream.text()
    return new NextResponse(text, {
      status: upstream.status,
      headers: { "Content-Type": "application/json" },
    })
  } catch (err) {
    return NextResponse.json(
      { error: "checkout_gateway_unreachable", message: "Payment service is unreachable." },
      { status: 503 }
    )
  }
}
