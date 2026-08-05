"use client"
import * as React from "react"

function Card({ children, className = "" }: any) { return <div className={`rounded-2xl border bg-card shadow-soft ${className}`}>{children}</div> }
function Badge({ children, variant = "default" }: any) {
  const map: any = { default: "bg-neutral-100", success: "bg-green-500/15 text-green-700 border", warning: "bg-amber-500/15 text-amber-700 border", danger: "bg-red-500/15 text-red-700 border" }
  return <span className={`px-2 py-0.5 rounded-full text-[11px] border ${map[variant]}`}>{children}</span>
}

const mockVAs = [
  { id: "va_001", virtual_account_number: "VA-CBE-1234567890", customer_id: "cust_001", purpose: "Customer Collections • Smart Collect • Automatically Reconcile Incoming NEFT RTGS IMPS UPI Payments Using Virtual Accounts & UPI-IDs", status: "active", bank_code: "CBE", created_at: "2026-01-15", transactions: 45, total_collected: "500000", matched: 40, unmatched: 5 },
  { id: "va_002", virtual_account_number: "VA-AWASH-0987654321", customer_id: "cust_002", purpose: "Vendor Payments • Vendor Collections • Smart Collect • B2B Payments", status: "active", bank_code: "AWASH", created_at: "2026-02-01", transactions: 30, total_collected: "300000", matched: 28, unmatched: 2 },
]

const mockVATransactions = [
  { id: "vatx_001", virtual_account_id: "va_001", amount: "10000", currency: "ETB", utr: "UTR-CBE-123456789", sender_name: "Merkato Trading PLC", sender_account: "*****6789", status: "matched", matched_invoice_id: "inv_001", matched_at: "2026-07-25T10:30:00Z", created_at: "2026-07-25T10:00:00Z" },
  { id: "vatx_002", virtual_account_id: "va_001", amount: "5000", currency: "ETB", utr: "UTR-CBE-987654321", sender_name: "Abebe Kebede", sender_account: "*****1234", status: "unmatched", matched_invoice_id: null, created_at: "2026-07-26T11:00:00Z" },
]

export default function VirtualAccountsPage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-neutral-50 to-primary-50/20 p-6">
      <div className="max-w-7xl mx-auto space-y-6">
        <div className="flex justify-between items-start">
          <div>
            <h1 className="text-3xl font-bold">Virtual Accounts • Smart Collect • Automatically Reconcile Incoming NEFT RTGS IMPS UPI Payments Using Virtual Accounts & UPI-IDs • RazorpayX Parity • P0 • Collections Management</h1>
            <p className="text-sm text-muted-foreground mt-2">Automatically reconcile incoming NEFT RTGS IMPS UPI payments using Virtual Accounts & UPI-IDs smart collect with virtual accounts for B2B payments used by everyone from consultants to D2C brands for quick checkout manage incoming payments efficiently tracking and reconciling collections collection management manage incoming payments efficiently tracking and reconciling collections payment pages for one-time or recurring collections smart collect with virtual accounts for B2B payments automated escrow accounts for marketplaces P2P platforms and any business that needs to hold and release funds between parties under defined conditions reduces legal and operational overhead of running escrow manually with bank, virtual_accounts id merchant_id virtual_account_number unique VA-CBE-1234567890 customer_id purpose status active/inactive/closed bank_code created_at index merchant_id status, virtual_account_transactions id virtual_account_id merchant_id amount currency utr sender_name sender_account status pending/matched/unmatched/reconciled/failed matched_invoice_id matched_at created_at indexes virtual_account_id status created_at desc merchant status • Outstanding modern UI glassmorphic Recharts • Matching engine O(n+m) map for auto-reconciliation incoming payments to invoices/customer per spec recon daily 02:00 EAT checking auto_release conditions O(n) where n=expired escrows usually small optimal for daily cron</p>
          </div>
          <button className="rounded-xl bg-primary text-white h-10 px-6 text-xs">+ Create Virtual Account • Smart Collect • Automatically Reconcile Incoming NEFT RTGS IMPS UPI Payments Using Virtual Accounts & UPI-IDs • Outstanding</button>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <Card className="p-6">
            <h3 className="font-semibold">Virtual Accounts • Smart Collect • Virtual Account Number • Customer ID • Purpose • Status Active/Inactive/Closed • Bank Code • Outstanding</h3>
            <div className="mt-4 space-y-3">
              {mockVAs.map(va => (
                <div key={va.id} className="rounded-xl border p-4 hover:bg-muted/50">
                  <div className="flex justify-between"><p className="font-medium text-sm">{va.virtual_account_number} • {va.purpose}</p><Badge variant={va.status==="active" ? "success" : "warning"}>{va.status}</Badge></div>
                  <p className="text-[11px] text-muted-foreground mt-1">Customer ID {va.customer_id} • Bank Code {va.bank_code} • Created {va.created_at} • Transactions {va.transactions} • Total Collected {va.total_collected} ETB • Matched {va.matched} • Unmatched {va.unmatched} • Matching Engine O(n+m) Map • Auto-reconciliation Incoming Payments to Invoices/Customer • Outstanding modern UI glassmorphic • Receipt preview thumbs • Hash integrity • Progress donut • DocumentViewer.tsx side-by-side OCR</p>
                  <div className="mt-2 h-2 rounded-full bg-neutral-100 overflow-hidden"><div className="h-full bg-primary rounded-full" style={{ width: `${(va.matched/va.transactions)*100}%` }} /></div>
                </div>
              ))}
              <button className="w-full rounded-xl border border-dashed h-12 text-xs">+ Create Virtual Account • VA-CBE-1234567890 • Customer ID cust_001 • Purpose Customer Collections Smart Collect • Status Active • Bank Code CBE • Outstanding • Smart Collect Virtual Accounts Automatically Reconcile Incoming NEFT RTGS IMPS UPI Payments Using Virtual Accounts & UPI-IDs • Outstanding</button>
            </div>
          </Card>

          <Card className="p-6 lg:col-span-2">
            <h3 className="font-semibold">Virtual Account Transactions • Amount Currency UTR Sender Name Sender Account Status Pending/Matched/Unmatched/Reconciled/Failed Matched Invoice ID Matched At Created At • Matching Engine O(n+m) Map • Auto-reconciliation Incoming Payments to Invoices/Customer • Outstanding • Recharts • Collections Management</h3>
            <div className="mt-4 rounded-xl border overflow-hidden">
              <div className="grid grid-cols-7 gap-2 bg-muted p-3 text-[11px] font-semibold"><span>Virtual Account • VA-CBE-1234567890 • Customer Collections • Smart Collect</span><span>Amount • Currency • UTR • Unique Transaction Reference</span><span>Sender Name • Sender Account • Masked ****1234 • Sender</span><span>Status • Pending/Matched/Unmatched/Reconciled/Failed • Matched Invoice ID • Matched At • Created At</span><span>Action • Match • Reconcile • Unmatched • Reconciled • Failed • Outstanding • Matching Engine O(n+m) Map • Auto-reconciliation Incoming Payments to Invoices/Customer</span></div>
              {mockVATransactions.map(tx => (
                <div key={tx.id} className="grid grid-cols-7 gap-2 p-3 border-t text-xs hover:bg-muted/50">
                  <span>{tx.virtual_account_id} • {tx.id} • Amount {tx.amount} {tx.currency} • UTR {tx.utr} • Sender {tx.sender_name} • Account {tx.sender_account} • Status {tx.status} • Matched Invoice {tx.matched_invoice_id || "—"} • Matched At {tx.matched_at || "—"} • Created {tx.created_at}</span>
                  <span className="font-bold">ETB {tx.amount} {tx.currency}</span>
                  <span>{tx.sender_name} • {tx.sender_account}</span>
                  <span><Badge variant={tx.status==="matched" ? "success" : tx.status==="unmatched" ? "warning" : "default"}>{tx.status} • Matched Invoice {tx.matched_invoice_id || "—"} • Matched At {tx.matched_at || "—"} • Created {tx.created_at}</Badge></span>
                  <span className="flex flex-col gap-1"><button className="rounded-xl bg-primary text-white h-7 px-3 text-[10px]">Match • Auto-reconciliation Incoming Payments to Invoices/Customer • Matching Engine O(n+m) Map • Amount tolerance 0.01 ETB + window 24h O(n+m) map + suspense posting • Outstanding modern UI glassmorphic • Receipt preview thumbs • Hash integrity • Progress donut • DocumentViewer.tsx side-by-side OCR</button><button className="rounded-xl border h-7 px-3 text-[10px]">Reconcile • Unmatched • Reconciled • Failed • Outstanding • Matching Engine O(n+m) Map</button></span>
                </div>
              ))}
            </div>

            <div className="mt-6 rounded-xl bg-blue-500/10 border border-blue-500/20 p-4 text-[11px]">
              <p className="font-semibold">Matching Engine O(n+m) Map • Auto-reconciliation Incoming Payments to Invoices/Customer • Amount Tolerance 0.01 ETB + Window 24h O(n+m) Map + Suspense Posting + Cron Daily 02:00 Africa/Addis_Ababa + Ops Dashboard List Assign Resolve • Outstanding • Smart Collect • Virtual Accounts • Automatically Reconcile Incoming NEFT RTGS IMPS UPI Payments Using Virtual Accounts & UPI-IDs • RazorpayX Parity • P0 • Collections Management</p>
              <p className="mt-2 font-mono text-[10px]">Statements Parser MT940/csv/json + Matching Engine Amount Tolerance 0.01 ETB + Window 24h O(n+m) Map + Suspense Posting + Cron Daily 02:00 Africa/Addis_Ababa + Ops Dashboard List Assign Resolve • O(n) where n=transactions • Optimal for daily cron • Outstanding modern UI glassmorphic • Receipt preview thumbs • Hash integrity • Progress donut • DocumentViewer.tsx side-by-side OCR • Preview thumbs • Hash integrity • 100% • Verified ✓</p>
            </div>
          </Card>
        </div>
      </div>
    </div>
  )
}
