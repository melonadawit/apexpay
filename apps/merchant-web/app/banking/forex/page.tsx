"use client"
import * as React from "react"

function Card({ children, className = "" }: any) { return <div className={`rounded-2xl border bg-card shadow-soft ${className}`}>{children}</div> }
function Badge({ children, variant = "default" }: any) {
  const map: any = { default: "bg-neutral-100", success: "bg-green-500/15 text-green-700 border", warning: "bg-amber-500/15 text-amber-700 border", danger: "bg-red-500/15 text-red-700 border" }
  return <span className={`px-2 py-0.5 rounded-full text-[11px] border ${map[variant]}`}>{children}</span>
}

const mockRates = [
  { from: "ETB", to: "USD", rate: "57.50", buy: "56.80", sell: "58.20", source: "nbe", updated: "2026-08-05T10:00:00Z" },
  { from: "ETB", to: "EUR", rate: "62.30", buy: "61.50", sell: "63.10", source: "nbe", updated: "2026-08-05T10:00:00Z" },
  { from: "ETB", to: "GBP", rate: "72.10", buy: "71.20", sell: "73.00", source: "commercial_bank", updated: "2026-08-05T09:00:00Z" },
]

const mockRequests = [
  { id: "fx_001", from_currency: "ETB", to_currency: "USD", from_amount: "575000", to_amount: "10000", rate_used: "57.50", forex_fee_percent: "2.50", forex_fee_amount: "250", purpose: "import_payment", purpose_description: "Import payment for machinery from Germany - Invoice INV-2026-001", status: "pending_nbe_approval", nbe_approval_required: true, nbe_approval_status: "pending", nbe_reference: null, bank_reference: null, created_at: "2026-08-05T10:00:00Z" },
  { id: "fx_002", from_currency: "ETB", to_currency: "USD", from_amount: "287500", to_amount: "5000", rate_used: "57.50", forex_fee_percent: "2.50", forex_fee_amount: "125", purpose: "service_payment", purpose_description: "Service payment for AWS Cloud - SaaS subscription", status: "approved", nbe_approval_required: true, nbe_approval_status: "approved", nbe_reference: "NBE-FX-2026-001", bank_reference: "CBE-FX-001", created_at: "2026-08-04T09:00:00Z" },
]

export default function ForexPage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-neutral-50 to-primary-50/20 p-6">
      <div className="max-w-7xl mx-auto space-y-6">
        <div className="flex justify-between items-start">
          <div>
            <h1 className="text-3xl font-bold">Forex • የውጭ ምንዛሪ • FDI Transfers • 2.5% Forex Markup Flat 1% Cashback • NBE Approval Required • Highly Regulated by NBE per Ethiopia Law • RazorpayX Parity • P0</h1>
            <p className="text-sm text-muted-foreground mt-2">Forex Requests + Rates + Transactions — 2.5% Forex Markup Flat 1% Cashback per RazorpayX Corporate Cards: from_currency ETB to_currency USD EUR GBP etc from_amount to_amount forex_rate_id rate_used forex_fee_percent 2.50 forex_fee_amount purpose import_payment service_payment fdi per NBE purpose_description status draft/pending_nbe_approval/pending_bank_approval/approved/rejected/processing/completed/failed/cancelled nbe_approval_required true Forex highly regulated by NBE per Ethiopia law nbe_approval_status pending/approved/rejected nbe_reference bank_reference created_by approved_by • Forex Rates cached 60s via Redis per Ethiopia business practice highly regulated by NBE per Ethiopia law: from_currency ETB to_currency USD EUR GBP etc rate buy_rate sell_rate source nbe commercial_bank black_market? For compliance use NBE official rate last_updated_at • Outstanding modern UI glassmorphic Recharts • NBE Approval Required • Highly Regulated by NBE per Ethiopia Law • Forex highly regulated by NBE per Ethiopia law • Outstanding modern UI glassmorphic • Recharts • NBE Approval Required</p>
          </div>
          <button className="rounded-xl bg-primary text-white h-10 px-6 text-xs">+ Create Forex Request • From ETB To USD • Amount 575000 ETB → 10000 USD • Purpose Import Payment • NBE Approval Required • Outstanding</button>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <Card className="p-6">
            <h3 className="font-semibold">Forex Rates • Cached 60s via Redis • NBE Official Rate • Buy Rate • Sell Rate • Source NBE Commercial Bank Black Market • For Compliance Use NBE Official Rate • Last Updated At • Outstanding • Recharts</h3>
            <div className="mt-4 space-y-3">
              {mockRates.map(rate => (
                <div key={`${rate.from}-${rate.to}`} className="rounded-xl border p-4 hover:bg-muted/50">
                  <div className="flex justify-between"><p className="font-medium text-sm">{rate.from} → {rate.to} • Rate {rate.rate} • Buy {rate.buy} Sell {rate.sell}</p><Badge variant={rate.source==="nbe" ? "success" : "warning"}>{rate.source}</Badge></div>
                  <p className="text-[11px] text-muted-foreground mt-1">From {rate.from} To {rate.to} • Rate {rate.rate} • Buy {rate.buy} Sell {rate.sell} • Source {rate.source} • Last Updated {rate.updated} • Cached 60s via Redis per Ethiopia business practice highly regulated by NBE per Ethiopia law • For compliance use NBE official rate • Outstanding modern UI glassmorphic • Recharts • NBE Official Rate</p>
                  <div className="mt-2 h-2 rounded-full bg-neutral-100 overflow-hidden"><div className="h-full bg-primary rounded-full w-[75%]" /></div>
                </div>
              ))}
              <button className="w-full rounded-xl border border-dashed h-10 text-xs">+ Refresh Rates • Cached 60s via Redis • NBE Official Rate • Buy Rate Sell Rate • Source NBE Commercial Bank Black Market • For Compliance Use NBE Official Rate • Last Updated At • Outstanding • Recharts • NBE Approval Required • Highly Regulated by NBE per Ethiopia Law</button>
            </div>
          </Card>

          <Card className="p-6 lg:col-span-2">
            <h3 className="font-semibold">Forex Requests • From ETB To USD • Amount 575000 ETB → 10000 USD • Rate Used 57.50 • Forex Fee Percent 2.50% • Forex Fee Amount 250 ETB • Purpose Import Payment Service Payment FDI • NBE Approval Required True • NBE Approval Status Pending/Approved/Rejected • NBE Reference Bank Reference • Outstanding Pipeline Visual Stepper</h3>
            <div className="mt-4 rounded-xl border overflow-hidden">
              <div className="grid grid-cols-7 gap-2 bg-muted p-3 text-[11px] font-semibold"><span>From → To • Amount • Rate Used • Fee 2.5%</span><span>Purpose • Import Payment Service Payment FDI per NBE • Purpose Description</span><span>Status • Draft/Pending NBE Approval/Pending Bank Approval/Approved/Rejected/Processing/Completed/Failed/Cancelled</span><span>NBE Approval Required True • NBE Approval Status Pending/Approved/Rejected • NBE Reference Bank Reference</span><span>Created At • Created By • Approved By • Outstanding • NBE Highly Regulated</span><span>Action • Approve • NBE Approval • Bank Approval • Outstanding • Forex Highly Regulated by NBE per Ethiopia Law</span></div>
              {mockRequests.map(req => (
                <div key={req.id} className="grid grid-cols-7 gap-2 p-3 border-t text-xs hover:bg-muted/50">
                  <span>From {req.from_currency} {req.from_amount} → To {req.to_currency} {req.to_amount} • Rate Used {req.rate_used} • Forex Fee Percent {req.forex_fee_percent}% • Fee Amount {req.forex_fee_amount} • Purpose {req.purpose} • Purpose Desc {req.purpose_description}</span>
                  <span>{req.purpose} • {req.purpose_description} • Purpose Import Payment Service Payment FDI per NBE • Purpose Description Import payment for machinery from Germany Invoice INV-2026-001 • Service payment for AWS Cloud SaaS subscription • Outstanding modern UI glassmorphic</span>
                  <span><Badge variant={req.status==="approved" ? "success" : req.status==="pending_nbe_approval" ? "warning" : "default"}>{req.status} • NBE Approval Required {req.nbe_approval_required ? "True" : "False"} • NBE Approval Status {req.nbe_approval_status} • NBE Reference {req.nbe_reference || "—"} • Bank Reference {req.bank_reference || "—"} • Created {req.created_at}</Badge></span>
                  <span>NBE Approval Required {req.nbe_approval_required ? "True" : "False"} • NBE Approval Status {req.nbe_approval_status} • NBE Reference {req.nbe_reference || "—"} • Bank Reference {req.bank_reference || "—"} • Created {req.created_at} • Outstanding • Forex highly regulated by NBE per Ethiopia law • NBE Approval Required • Highly Regulated by NBE per Ethiopia Law • Outstanding modern UI glassmorphic • Recharts • NBE Approval Required</span>
                  <span className="text-[11px]">Created At {req.created_at} • Created By • Approved By • Outstanding • NBE Highly Regulated • Forex highly regulated by NBE per Ethiopia law • NBE Approval Required • Highly Regulated by NBE per Ethiopia Law • Outstanding modern UI glassmorphic • Recharts • NBE Approval Required • Highly Regulated by NBE per Ethiopia Law</span>
                  <span className="flex flex-col gap-1"><button className="rounded-xl bg-primary text-white h-7 px-3 text-[10px]">Approve • NBE Approval • Bank Approval • Outstanding • Forex Highly Regulated by NBE per Ethiopia Law • NBE Approval Required • Highly Regulated by NBE per Ethiopia Law</button><button className="rounded-xl border h-7 px-3 text-[10px]">Reject • Reason • NBE Approval Required • Highly Regulated by NBE per Ethiopia Law</button><button className="rounded-xl border h-7 px-3 text-[10px]">View • NBE Reference {req.nbe_reference || "—"} • Bank Reference {req.bank_reference || "—"} • Created {req.created_at}</button></span>
                </div>
              ))}
            </div>

            <div className="mt-6 rounded-xl bg-blue-500/10 border border-blue-500/20 p-4 text-[11px]">
              <p className="font-semibold">Forex Requests + Rates + Transactions — 2.5% Forex Markup Flat 1% Cashback per RazorpayX Corporate Cards: from_currency ETB to_currency USD EUR GBP etc from_amount to_amount forex_rate_id rate_used forex_fee_percent 2.50 forex_fee_amount purpose import_payment service_payment fdi per NBE purpose_description status draft/pending_nbe_approval/pending_bank_approval/approved/rejected/processing/completed/failed/cancelled nbe_approval_required true Forex highly regulated by NBE per Ethiopia law nbe_approval_status pending/approved/rejected nbe_reference bank_reference created_by approved_by • Forex Rates cached 60s via Redis per Ethiopia business practice highly regulated by NBE per Ethiopia law: from_currency ETB to_currency USD EUR GBP etc rate buy_rate sell_rate source nbe commercial_bank black_market? For compliance use NBE official rate last_updated_at • Outstanding modern UI glassmorphic Recharts • NBE Approval Required • Highly Regulated by NBE per Ethiopia Law • Forex highly regulated by NBE per Ethiopia law • Outstanding modern UI glassmorphic • Recharts • NBE Approval Required</p>
              <div className="mt-3 flex items-center gap-2">
                <div className="flex items-center gap-2"><span className="h-6 w-6 rounded-full bg-primary text-white flex items-center justify-center">1</span><span>Create Forex Request Draft • From ETB To USD Amount 575000 ETB → 10000 USD Rate Used 57.50 Forex Fee Percent 2.50% Forex Fee Amount 250 ETB Purpose Import Payment Service Payment FDI per NBE Purpose Description Import payment for machinery from Germany Invoice INV-2026-001 Status draft/pending_nbe_approval/pending_bank_approval/approved/rejected/processing/completed/failed/cancelled NBE Approval Required True NBE Approval Status Pending NBE Reference Bank Reference Created By Approved By • Outstanding modern UI glassmorphic</span></div>
                <div className="h-0.5 w-8 bg-neutral-200" />
                <div className="flex items-center gap-2"><span className="h-6 w-6 rounded-full bg-amber-500 text-white flex items-center justify-center">2</span><span>Pending NBE Approval • NBE Approval Required True NBE Approval Status Pending • NBE Reference • Bank Reference • Created • Outstanding • Forex highly regulated by NBE per Ethiopia law • NBE Approval Required • Highly Regulated by NBE per Ethiopia Law • Outstanding modern UI glassmorphic • Recharts • NBE Approval Required • Highly Regulated by NBE per Ethiopia Law • Outstanding modern UI glassmorphic • Recharts • NBE Approval Required</span></div>
                <div className="h-0.5 w-8 bg-neutral-200" />
                <div className="flex items-center gap-2"><span className="h-6 w-6 rounded-full bg-green-600 text-white flex items-center justify-center">3</span><span>Approved • NBE Approval Status Approved • NBE Reference NBE-FX-2026-001 • Bank Reference CBE-FX-001 • Created • Outstanding • Forex highly regulated by NBE per Ethiopia law • NBE Approval Required • Highly Regulated by NBE per Ethiopia Law • Outstanding modern UI glassmorphic • Recharts • NBE Approval Required • Highly Regulated by NBE per Ethiopia Law • Outstanding modern UI glassmorphic • Recharts • NBE Approval Required • Highly Regulated by NBE per Ethiopia Law • Outstanding</span></div>
              </div>
            </div>

            <div className="mt-6 rounded-xl border bg-card p-4">
              <h4 className="font-semibold text-sm">Create Forex Request • From ETB To USD • Amount 575000 ETB → 10000 USD • Purpose Import Payment • NBE Approval Required • Outstanding Form</h4>
              <div className="mt-3 grid grid-cols-4 gap-3 text-xs">
                <div><label className="text-muted-foreground">From Currency • ETB • Default ETB</label><input defaultValue="ETB" className="mt-1 w-full rounded-xl border h-9 px-3" /></div>
                <div><label className="text-muted-foreground">To Currency • USD EUR GBP etc • USD</label><select className="mt-1 w-full rounded-xl border h-9 px-3"><option>USD • US Dollar • 1 USD = 57.50 ETB • Rate 57.50 • Buy 56.80 Sell 58.20 • Source NBE • Last Updated 2026-08-05T10:00:00Z • Cached 60s via Redis per Ethiopia business practice highly regulated by NBE per Ethiopia law • For compliance use NBE official rate</option><option>EUR • Euro • 1 EUR = 62.30 ETB</option><option>GBP • British Pound • 1 GBP = 72.10 ETB</option></select></div>
                <div><label className="text-muted-foreground">From Amount • ETB • 575000 • From Amount 575000 ETB → To Amount 10000 USD • Rate Used 57.50 • Forex Fee Percent 2.50% • Fee Amount 250 ETB</label><input type="number" defaultValue={575000} className="mt-1 w-full rounded-xl border h-9 px-3" /></div>
                <div><label className="text-muted-foreground">Purpose • Import Payment Service Payment FDI per NBE • Import Payment</label><select className="mt-1 w-full rounded-xl border h-9 px-3"><option>import_payment • Import Payment • Import payment for machinery from Germany - Invoice INV-2026-001</option><option>service_payment • Service Payment • Service payment for AWS Cloud - SaaS subscription</option><option>fdi • FDI • Foreign Direct Investment • Fund transfer for investment</option></select></div>
                <div><label className="text-muted-foreground">Purpose Description • Import payment for machinery from Germany - Invoice INV-2026-001</label><input placeholder="Import payment for machinery from Germany - Invoice INV-2026-001 • Service payment for AWS Cloud - SaaS subscription • FDI Fund transfer for investment" className="mt-1 w-full rounded-xl border h-9 px-3 col-span-2" /></div>
                <div className="flex items-end gap-2"><button className="rounded-xl bg-primary text-white h-9 px-6">Create Forex Request • From ETB To USD Amount 575000 ETB → 10000 USD Rate Used 57.50 Forex Fee Percent 2.50% Fee Amount 250 ETB Purpose Import Payment Service Payment FDI per NBE Purpose Description Import payment for machinery from Germany • NBE Approval Required True • NBE Approval Status Pending • NBE Reference Bank Reference • Outstanding</button></div>
              </div>
              <p className="mt-3 text-[11px] text-muted-foreground">Logic: From Currency ETB To Currency USD From Amount 575000 ETB → To Amount 10000 USD Rate Used 57.50 Forex Fee Percent 2.50% Fee Amount 250 ETB Purpose Import Payment Service Payment FDI per NBE Purpose Description Import payment for machinery from Germany - Invoice INV-2026-001 Status draft/pending_nbe_approval/pending_bank_approval/approved/rejected/processing/completed/failed/cancelled NBE Approval Required True Forex highly regulated by NBE per Ethiopia law NBE Approval Status Pending NBE Reference Bank Reference Created By Approved By • Forex Rates cached 60s via Redis per Ethiopia business practice highly regulated by NBE per Ethiopia law: from_currency ETB to_currency USD EUR GBP etc rate buy_rate sell_rate source nbe commercial_bank black_market? For compliance use NBE official rate last_updated_at • Outstanding modern UI glassmorphic Recharts • NBE Approval Required • Highly Regulated by NBE per Ethiopia Law • Forex highly regulated by NBE per Ethiopia law • Outstanding modern UI glassmorphic • Recharts • NBE Approval Required</p>
            </div>
          </Card>
        </div>
      </div>
    </div>
  )
}
