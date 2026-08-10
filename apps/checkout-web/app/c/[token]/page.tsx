"use client"
import * as React from "react"
import {
  fetchLink,
  initializePayment,
  fetchStatus,
  verify2FA,
  type PaymentLink,
  type InitResult,
} from "../../../lib/checkout"

type Phase = "loading" | "open" | "processing" | "2fa" | "success" | "failed"

const METHODS = [
  { id: "telebirr", name: "Telebirr", icon: "📱" },
  { id: "cbe_birr", name: "CBE Birr", icon: "🏦" },
  { id: "bank", name: "Bank Transfer", icon: "🏧" },
  { id: "card", name: "Card", icon: "💳" },
  { id: "qr", name: "EthSwitch QR", icon: "🔗" },
]

export default function CheckoutTokenPage({ params }: { params: { token: string } }) {
  const token = params.token
  const [phase, setPhase] = React.useState<Phase>("loading")
  const [link, setLink] = React.useState<PaymentLink | null>(null)
  const [method, setMethod] = React.useState("telebirr")
  const [init, setInit] = React.useState<InitResult | null>(null)
  const [otp, setOtp] = React.useState("")
  const [error, setError] = React.useState("")

  // Load the payment link from the real public API.
  React.useEffect(() => {
    let cancelled = false
    fetchLink(token)
      .then((pl) => {
        if (cancelled) return
        setLink(pl)
        setPhase(pl.status === "active" || pl.status === "open" ? "open" : "failed")
      })
      .catch((e) => {
        if (cancelled) return
        setError(e.message || "Could not load payment")
        setPhase("failed")
      })
    return () => {
      cancelled = true
    }
  }, [token])

  const pay = async () => {
    setPhase("processing")
    setError("")
    try {
      const res = await initializePayment(token, method)
      setInit(res)
      if (res.requires_2fa) {
        setPhase("2fa")
      } else {
        poll(res.tx_ref)
      }
    } catch (e) {
      setError((e as Error).message)
      setPhase("open")
    }
  }

  const submit2FA = async () => {
    if (!init) return
    try {
      await verify2FA(token, init.id, otp)
      setOtp("")
      poll(init.tx_ref)
    } catch (e) {
      setError((e as Error).message || "Invalid OTP")
    }
  }

  // Poll the real status endpoint every 2s until settled.
  const poll = (txRef: string) => {
    const interval = setInterval(async () => {
      try {
        const s = await fetchStatus(token, txRef)

        if (s.status === "succeeded") {
          clearInterval(interval)
          setPhase("success")
        } else if (s.status === "failed") {
          clearInterval(interval)
          setError("Payment failed")
          setPhase("failed")
        }
      } catch {
        // transient poll error — keep retrying
      }
    }, 2000)
  }

  if (phase === "loading") {
    return (
      <Centered>
        <div className="mx-auto h-12 w-12 rounded-full border-4 border-primary border-t-transparent animate-spin" />
        <p className="mt-4 text-sm text-muted-foreground">Loading secure checkout…</p>
      </Centered>
    )
  }

  if (phase === "processing") {
    return (
      <Centered>
        <div className="mx-auto h-12 w-12 rounded-full border-4 border-primary border-t-transparent animate-spin" />
        <p className="mt-4 font-semibold">Processing payment…</p>
      </Centered>
    )
  }

  if (phase === "success") {
    return (
      <Centered>
        <div className="w-full max-w-[420px] bg-white rounded-2xl shadow-soft p-8 text-center space-y-4">
          <div className="mx-auto h-20 w-20 rounded-full bg-green-100 flex items-center justify-center text-4xl">✓</div>
          <h2 className="text-2xl font-bold">Payment Successful</h2>
          <p className="text-sm text-muted-foreground">
            {init ? `${init.amount} ${init.currency} • Ref ${init.tx_ref}` : "Receipt sent to your email"}
          </p>
          <div className="flex gap-2">
            <button className="flex-1 rounded-xl border h-12 text-sm" onClick={() => window.print()}>Download PDF</button>
            <button className="flex-1 rounded-xl bg-primary text-white h-12 text-sm">Done</button>
          </div>
        </div>
      </Centered>
    )
  }

  const amount = link ? link.amount : "—"

  return (
    <Centered>
      <div className="w-full max-w-[420px] rounded-2xl border border-black/10 bg-white shadow-soft overflow-hidden">
        <div className="p-6 space-y-6">
          <div className="text-center">
            <div className="mx-auto h-12 w-12 rounded-xl bg-primary/10 flex items-center justify-center font-bold text-primary">A</div>
            <h1 className="mt-3 font-bold">{link?.description || "Secure Checkout"}</h1>
            <p className="text-sm text-muted-foreground">🔒 Secure checkout</p>
            <p className="text-3xl font-bold mt-3">ETB {amount}</p>
          </div>

          <div className="space-y-2">
            <p className="text-sm font-semibold">Payment Method</p>
            {METHODS.map((m) => (
              <button
                key={m.id}
                onClick={() => setMethod(m.id)}
                className={`w-full rounded-xl border p-3 flex items-center gap-3 text-left ${
                  method === m.id ? "border-primary bg-primary/5" : "border-black/10 hover:bg-neutral-50"
                }`}
              >
                <span className="text-xl">{m.icon}</span>
                <div className="flex-1">
                  <p className="text-sm font-medium">{m.name}</p>
                </div>
                {method === m.id && <span className="text-primary">✓</span>}
              </button>
            ))}
          </div>

          {phase === "2fa" && (
            <div className="rounded-xl bg-amber-50 border border-amber-200 p-3 space-y-2">
              <p className="text-sm font-semibold">Enter the 6-digit code sent to your phone</p>
              <input
                value={otp}
                onChange={(e) => setOtp(e.target.value)}
                placeholder="6-digit OTP"
                maxLength={6}
                inputMode="numeric"
                className="w-full rounded-xl border h-12 px-3"
              />
              <button onClick={submit2FA} className="w-full rounded-xl bg-amber-600 text-white h-10 text-sm">
                Verify
              </button>
            </div>
          )}

          {error && <p className="text-sm text-red-600">{error}</p>}

          <button
            onClick={pay}
            disabled={phase === "failed"}
            className="w-full h-14 rounded-xl bg-primary text-white font-semibold shadow-soft hover:shadow-medium active:scale-[0.98] transition-all disabled:opacity-50"
          >
            Pay ETB {amount}
          </button>

          <p className="text-[11px] text-center text-muted-foreground">🔒 Secure checkout by ApexPay</p>
        </div>
      </div>
    </Centered>
  )
}

function Centered({ children }: { children: React.ReactNode }) {
  return (
    <div className="min-h-screen flex items-center justify-center bg-neutral-50 p-4">
      {children}
    </div>
  )
}
