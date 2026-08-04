"use client"
import * as React from "react"
import { motion } from "framer-motion"

export default function CheckoutPage() {
  const [status, setStatus] = React.useState<"open"|"processing"|"success"|"failed">("open")
  const [method, setMethod] = React.useState("telebirr")
  const [show2FA, setShow2FA] = React.useState(false)
  const methods = [
    { id:"telebirr", name:"Telebirr", icon:"📱", latency:210, success:0.96, fee:"2.9%" },
    { id:"cbe_birr", name:"CBE Birr", icon:"🏦", latency:260, success:0.89, fee:"2.9%" },
    { id:"bank", name:"Bank Transfer", icon:"🏧", latency:180, success:0.92, fee:"2.5%" },
    { id:"card", name:"Card", icon:"💳", latency:320, success:0.85, fee:"3.0%" },
    { id:"qr", name:"EthSwitch QR", icon:"🔗", latency:150, success:0.90, fee:"2.0%" },
  ]

  const pay = () => {
    setStatus("processing")
    // Simulate 2FA >5000 ETB per NBE ONPS/10/2025
    setTimeout(()=> {
      if (true) { setShow2FA(false); setStatus("success") }
    }, 2000)
  }

  if (status==="processing") return (
    <div className="min-h-screen flex items-center justify-center bg-neutral-50 p-6">
      <div className="w-full max-w-[420px] text-center space-y-4">
        <div className="mx-auto h-16 w-16 rounded-full border-4 border-primary border-t-transparent animate-spin" />
        <p className="font-semibold">Processing payment… • ክፍያ በመፈጸም ላይ</p>
        <p className="text-xs text-muted-foreground">Polling verify every 2s • routed via best connector</p>
      </div>
    </div>
  )

  if (status==="success") return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-green-50 to-emerald-100 p-6">
      <motion.div initial={{ scale:0.8, opacity:0 }} animate={{ scale:1, opacity:1 }} className="w-full max-w-[420px] bg-white rounded-2xl shadow-soft p-8 text-center space-y-4">
        <div className="mx-auto h-20 w-20 rounded-full bg-green-100 flex items-center justify-center text-4xl">✓</div>
        <h2 className="text-2xl font-bold">Payment Successful • ክፍያ ተሳክቷል</h2>
        <p className="text-sm">ETB 500.00 to ApexPay — receipt sent</p>
        <div className="flex gap-2">
          <button className="flex-1 rounded-xl border h-12">Download PDF</button>
          <button className="flex-1 rounded-xl bg-primary text-white h-12">Done • ጨርስ</button>
        </div>
      </motion.div>
    </div>
  )

  return (
    <div className="min-h-screen flex items-center justify-center bg-neutral-50 p-4">
      <div className="w-full max-w-[420px] rounded-2xl border border-black/10 bg-white shadow-soft overflow-hidden">
        <div className="p-6 space-y-6">
          <div className="text-center">
            <div className="mx-auto h-12 w-12 rounded-xl bg-primary/10 flex items-center justify-center font-bold text-primary">A</div>
            <h1 className="mt-3 font-bold">ApexPay Demo Store</h1>
            <p className="text-sm text-muted-foreground">Secure checkout • NBE Licensed Gateway</p>
            <p className="text-3xl font-bold mt-3">ETB 500.00</p>
          </div>

          <div className="space-y-2">
            <p className="text-sm font-semibold">Payment Method • የክፍያ ዘዴ</p>
            <p className="text-xs text-muted-foreground flex items-center gap-1">Using best route: <span className="bg-primary/10 text-primary px-2 py-0.5 rounded-full">Telebirr (2% faster today)</span> <span className="text-[10px]">via smart routing success_rate</span></p>
            {methods.map(m=> (
              <button key={m.id} onClick={()=> setMethod(m.id)} className={`w-full rounded-xl border p-3 flex items-center gap-3 text-left ${method===m.id ? "border-primary bg-primary/5" : "border-black/10 hover:bg-neutral-50"}`}>
                <span className="text-xl">{m.icon}</span>
                <div className="flex-1"><p className="text-sm font-medium">{m.name}</p><p className="text-[11px] text-muted-foreground">{m.latency}ms • {Math.round(m.success*100)}% success • fee {m.fee}</p></div>
                {method===m.id && <span className="text-primary">✓</span>}
              </button>
            ))}
          </div>

          {show2FA && (
            <div className="rounded-xl bg-amber-50 border border-amber-200 p-3">
              <p className="text-sm font-semibold">2FA Required >5000 ETB per ONPS/10/2025</p>
              <input placeholder="OTP 6-digit" className="mt-2 w-full rounded-xl border h-12 px-3" />
            </div>
          )}

          <button onClick={pay} className="w-full h-14 rounded-xl bg-primary text-white font-semibold shadow-soft hover:shadow-medium active:scale-[0.98] transition-all">Pay ETB 500.00 • ክፍያ</button>

          <p className="text-[11px] text-center text-muted-foreground">🔒 Secure by ApexPay • FIN never logged • Encrypted • NBE compliant</p>
        </div>
      </div>
    </div>
  )
}
