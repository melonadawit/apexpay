import { NextRequest, NextResponse } from "next/server"

// Logout: revokes the session server-side and clears the httpOnly cookie.

const API_BASE = process.env.APEXPAY_API_URL || "http://api:8080"
const COOKIE = "apexpay_session"

export async function POST(req: NextRequest) {
  const token = req.cookies.get(COOKIE)?.value
  if (token) {
    // Best-effort revocation of the backend session.
    try {
      await fetch(`${API_BASE}/v1/auth/logout`, {
        method: "POST",
        headers: { Authorization: `Bearer ${token}` },
      })
    } catch {
      /* ignore upstream errors; clear the cookie regardless */
    }
  }
  const res = NextResponse.json({ status: "logged_out" })
  res.cookies.set(COOKIE, "", { httpOnly: true, sameSite: "lax", path: "/", maxAge: 0 })
  return res
}
