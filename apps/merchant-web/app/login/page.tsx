"use client"
import * as React from "react"
import { useRouter } from "next/navigation"
import { login } from "@/lib/api/auth"
import { useLanguage } from "@/components/providers/language-provider"

export default function LoginPage() {
  const { t } = useLanguage()
  const router = useRouter()
  const [email, setEmail] = React.useState("")
  const [password, setPassword] = React.useState("")
  const [error, setError] = React.useState("")
  const [loading, setLoading] = React.useState(false)

  const submit = async (e: React.FormEvent) => {
    e.preventDefault()
    setLoading(true)
    setError("")
    try {
      await login(email, password)
      router.push("/dashboard")
      router.refresh()
    } catch (err) {
      setError((err as Error).message)
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-primary to-primary-light p-6">
      <div className="w-full max-w-md rounded-2xl border border-black/10 bg-white shadow-soft p-8 space-y-6">
        <div className="text-center">
          <div className="mx-auto h-14 w-14 rounded-2xl bg-primary/10 flex items-center justify-center font-bold text-primary text-xl">A</div>
          <h1 className="mt-4 text-2xl font-bold">{t("Sign in","ይግቡ")}</h1>
          <p className="text-sm text-muted-foreground">ApexPay Merchant Dashboard</p>
        </div>

        <form onSubmit={submit} className="space-y-4">
          <div>
            <label className="text-sm font-medium">Email</label>
            <input
              type="email"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="you@company.et"
              className="mt-1 w-full rounded-xl border h-12 px-3"
              required
            />
          </div>
          <div>
            <label className="text-sm font-medium">Password</label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="••••••••"
              className="mt-1 w-full rounded-xl border h-12 px-3"
              required
            />
          </div>

          {error && <p className="text-sm text-red-600">{error}</p>}

          <button
            type="submit"
            disabled={loading}
            className="w-full h-12 rounded-xl bg-primary text-white font-semibold disabled:opacity-50"
          >
            {loading ? "Signing in…" : "Sign in • ግባ"}
          </button>
        </form>

        <p className="text-[11px] text-center text-muted-foreground">
          Demo seed: demo@apexpay.et / Admin@12345
        </p>
      </div>
    </div>
  )
}
