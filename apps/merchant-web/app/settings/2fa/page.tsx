"use client"
import * as React from "react"
import { useRequireAuth } from "@/lib/api/require-auth"
import { api, type TwoFAEnroll } from "@/lib/api/client"

export default function TwoFAPage() {
  const { checking } = useRequireAuth()
  const [enroll, setEnroll] = React.useState<TwoFAEnroll | null>(null)
  const [code, setCode] = React.useState("")
  const [result, setResult] = React.useState<string>("")
  const [err, setErr] = React.useState("")

  if (checking) return <Centered>Checking session…</Centered>

  const doEnroll = async () => {
    setErr("")
    try { setEnroll(await api.twofa.enroll("merchant@apexpay.et")) }
    catch (e) { setErr((e as Error).message) }
  }

  const doVerify = async () => {
    if (!enroll) return
    const r = await api.twofa.verify(enroll.secret, code)
    setResult(r.verified ? "✓ Verified" : "✗ Invalid code")
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-neutral-50 to-primary-50/20 p-6">
      <div className="max-w-lg mx-auto space-y-6">
        <h1 className="text-3xl font-bold">Two-Factor Authentication • 2FA</h1>
        <p className="text-sm text-muted-foreground">TOTP (RFC 6238) — add an authenticator app for strong 2FA on your account.</p>

        <div className="rounded-2xl border bg-card p-6 space-y-3">
          {!enroll ? (
            <button onClick={doEnroll} className="w-full rounded-xl bg-primary text-white h-11 text-sm font-semibold">Enroll TOTP</button>
          ) : (
            <>
              <p className="text-xs text-muted-foreground break-all">Secret: <span className="font-mono">{enroll.secret}</span></p>
              <p className="text-xs text-muted-foreground break-all">Scan URL: <span className="font-mono">{enroll.otpauth_url}</span></p>
              <input value={code} onChange={(e) => setCode(e.target.value)} maxLength={6} placeholder="6-digit code" className="w-full rounded-xl border h-11 px-3 text-sm text-center tracking-widest" />
              <button onClick={doVerify} className="w-full rounded-xl bg-primary text-white h-11 text-sm font-semibold">Verify</button>
              {result && <p className={`text-sm font-medium ${result.startsWith("✓") ? "text-green-600" : "text-red-600"}`}>{result}</p>}
            </>
          )}
          {err && <p className="text-sm text-red-600">{err}</p>}
        </div>
      </div>
    </div>
  )
}

function Centered({ children }: { children: React.ReactNode }) {
  return <div className="min-h-screen flex items-center justify-center bg-neutral-50 text-sm text-muted-foreground">{children}</div>
}
