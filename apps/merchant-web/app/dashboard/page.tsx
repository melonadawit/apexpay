"use client"
import * as React from "react"
import { motion } from "framer-motion"
import { Card, GlassCard } from "@/components/ui/card"
import { DonutProgress } from "@/components/ui/progress"
import { TPVRecharts, HealthRecharts } from "./recharts"

export default function DashboardPage() {
  const [tpv] = React.useState(125430)
  return (
    <div className="min-h-screen bg-gradient-to-br from-neutral-50 to-primary-50 p-6">
      <div className="max-w-6xl mx-auto space-y-6">
        <div className="flex justify-between items-center">
          <div><h1 className="text-3xl font-bold">Dashboard • ዳሽቦርድ</h1><p className="text-sm text-muted-foreground">Welcome Meron, TPV today + AI chat panel glassmorphic</p></div>
          <DonutProgress value={78} />
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <motion.div initial={{ opacity:0, y:10 }} animate={{ opacity:1, y:0 }} className="rounded-2xl bg-gradient-to-br from-primary to-primary-light p-6 text-foreground shadow-medium">
            <p className="text-foreground/70 text-sm">TPV Today • ዛሬ — Real Recharts TPV + merchant_tpv_daily materialized view refreshed hourly worker + SWR polling 5s</p><p className="text-3xl font-bold mt-2">ETB {tpv.toLocaleString()}</p>
            <p className="text-xs mt-2 bg-card/20 inline-block px-2 py-0.5 rounded-full">+12% vs yesterday • ትናንት • GET /v1/payments?merchant_id + /v1/admin/connectors/health Recharts AreaChart</p>
            <div className="mt-4">
              <TPVRecharts />
            </div>
          </motion.div>
          <Card className="p-6"><p className="text-sm text-muted-foreground">Success Rate • ስኬት — Health Sampler 30s + Circuit Breaker</p><p className="text-2xl font-bold mt-2">96.2%</p><p className="text-xs text-green-600 mt-1">▲ 2% routing fallback used 3 times • fallback_used false • Health snapshot telebirr 0.96 210ms</p><div className="mt-2"><HealthRecharts /></div></Card>
          <Card className="p-6"><p className="text-sm text-muted-foreground">Active Links • ሊንኮች</p><p className="text-2xl font-bold mt-2">12</p><p className="text-xs">QR • Telegram share • 3 paid today</p></Card>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <div className="lg:col-span-2 space-y-4">
            <Card className="p-4"><h3 className="font-semibold">Recent Payments • የቅርብ ክፍያዎች</h3>
              <div className="mt-3 space-y-2">
                {[
                  {tx:"txr_01H", amt:"ETB 500", method:"Telebirr", status:"succeeded", ago:"2 min"},
                  {tx:"txr_02H", amt:"ETB 1,000", method:"CBE Birr", status:"succeeded", ago:"5 min"},
                  {tx:"txr_03H", amt:"ETB 6,000", method:"Bank", status:"succeeded", ago:"10 min", twoFA:true},
                ].map(p=>(
                  <div key={p.tx} className="flex items-center justify-between rounded-xl border p-3 hover:bg-muted">
                    <div className="flex items-center gap-3"><div className="h-8 w-8 rounded-full bg-green-500/20 text-green-700 flex items-center justify-center text-xs">✓</div><div><p className="text-sm font-medium">{p.amt} • {p.method} {p.twoFA && "• 2FA"}</p><p className="text-xs text-muted-foreground">{p.tx} • {p.ago} ago • routed via best</p></div></div>
                    <span className="text-xs px-2 py-0.5 rounded-full bg-green-500/20 text-green-700">{p.status}</span>
                  </div>
                ))}
              </div>
            </Card>

            <Card className="p-4"><h3 className="font-semibold">Quick Actions • ፈጣን እርምጃዎች</h3>
              <div className="mt-3 grid grid-cols-3 gap-2">
                <button className="rounded-xl border p-3 text-sm font-medium hover:bg-muted">Create Link • ሊንክ</button>
                <button className="rounded-xl border p-3 text-sm font-medium hover:bg-muted">Pay Vendor • አቅራቢ ክፍያ</button>
                <button className="rounded-xl border p-3 text-sm font-medium hover:bg-muted">Run Payroll • ደሞዝ</button>
              </div>
            </Card>
          </div>

          <GlassCard className="p-4 space-y-3">
            <h3 className="font-semibold flex items-center gap-2">AI Chat • Swarm 🤖</h3>
            <div className="rounded-xl bg-card border p-3 text-sm">
              <p className="text-muted-foreground">Goal: Create link 100 ETB for coffee if today TPV&gt;0</p>
              <div className="mt-2 space-y-2">
                <div className="flex gap-2"><span className="h-6 w-6 rounded-full bg-primary text-foreground flex items-center justify-center text-xs">1</span><span className="text-xs">Tool get_tpv → ETB 125,430</span><span className="ml-auto text-green-600 text-xs">✓</span></div>
                <div className="flex gap-2"><span className="h-6 w-6 rounded-full bg-primary text-foreground flex items-center justify-center text-xs">2</span><span className="text-xs">Tool create_payment_link 100 ETB</span><span className="ml-auto text-green-600 text-xs">✓</span></div>
              </div>
              <p className="mt-2 text-xs font-medium">Final: Created link https://checkout.apexpay.et/c/coffee100</p>
            </div>
            <div className="rounded-xl bg-muted p-3 text-xs">
              <p className="font-semibold">RAG Ask • ተገዢነት</p>
              <p className="mt-1">Q: When is 2FA required? A: Transactions above 5000 ETB require 2FA per ONPS/10/2025 [1] score 0.92</p>
            </div>
          </GlassCard>
        </div>
      </div>
    </div>
  )
}
