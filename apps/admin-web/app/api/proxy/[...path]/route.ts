import { NextRequest, NextResponse } from "next/server"

// Same-origin proxy to the ApexPay API. The browser never sees the backend host or the
// admin credential; the server injects them. Backend URL + token are configurable via env.
const API = process.env.APEXPAY_API_URL || "http://localhost:8080"
const ADMIN_TOKEN = process.env.APEXPAY_ADMIN_TOKEN || ""

async function proxy(req: NextRequest, params: { path: string[] }): Promise<NextResponse> {
  const path = params.path.join("/")
  const url = `${API}/v1/${path}${req.nextUrl.search}`
  const headers: Record<string, string> = {}
  if (ADMIN_TOKEN) headers["Authorization"] = `Bearer ${ADMIN_TOKEN}`
  if (req.headers.get("content-type")) headers["Content-Type"] = req.headers.get("content-type")!

  const body = req.method === "GET" || req.method === "HEAD" ? undefined : await req.text()

  try {
    const upstream = await fetch(url, { method: req.method, headers, body })
    const text = await upstream.text()
    let data: unknown
    try {
      data = JSON.parse(text)
    } catch {
      data = text
    }
    return NextResponse.json(data, { status: upstream.status })
  } catch (e) {
    return NextResponse.json({ success: false, data: { message: (e as Error).message } }, { status: 502 })
  }
}

export async function GET(req: NextRequest, ctx: { params: { path: string[] } }) {
  return proxy(req, ctx.params)
}
export async function POST(req: NextRequest, ctx: { params: { path: string[] } }) {
  return proxy(req, ctx.params)
}
export async function PUT(req: NextRequest, ctx: { params: { path: string[] } }) {
  return proxy(req, ctx.params)
}
export async function DELETE(req: NextRequest, ctx: { params: { path: string[] } }) {
  return proxy(req, ctx.params)
}
