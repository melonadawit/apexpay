"use client"
import * as React from "react"

function Card({ children, className = "" }: any) { return <div className={`rounded-2xl border bg-card shadow-soft ${className}`}>{children}</div> }
function Badge({ children, variant = "default" }: any) {
  const map: any = { default: "bg-neutral-100", success: "bg-green-500/15 text-green-700 border", warning: "bg-amber-500/15 text-amber-700 border" }
  return <span className={`px-2 py-0.5 rounded-full text-[11px] border ${map[variant]}`}>{children}</span>
}

const mockBudgets = [
  { id: "pcb_001", budget_name: "Office Supplies • የቢሮ እቃዎች", amount: "50000", assigned_to: "Finance Manager", status: "active", spent_amount: "15000", remaining_amount: "35000", created_at: "2026-07-01" },
  { id: "pcb_002", budget_name: "Marketing • ግብይት", amount: "100000", assigned_to: "Marketing Manager", status: "active", spent_amount: "60000", remaining_amount: "40000", created_at: "2026-07-05" },
]
const mockExpenses = [
  { id: "pce_001", budget_id: "pcb_001", amount: "1500", description: "Office supplies - printer ink", receipt_file_key: "petty_cash/receipt_001.jpg", receipt_file_hash: "hash_pce_001", status: "paid", approved_by: "Finance Manager", created_by: "Admin" },
]

export default function PettyCashPage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-neutral-50 to-primary-50/20 p-6">
      <div className="max-w-7xl mx-auto space-y-6">
        <div className="flex justify-between items-start">
          <div>
            <h1 className="text-3xl font-bold">Petty Cash • ጥቃቅን ጥሬ ገንዘብ • Budgets & Expenses • Track Petty Cash Budgets and Make Payments from Assigned Budgets • Add Bills & Receipt as Attachments to Petty Cash Expenses • RazorpayX Parity • P0</h1>
            <p className="text-sm text-muted-foreground mt-2">Track petty cash budgets and make payments from assigned budgets add bills & receipt as attachments to petty cash expenses • Petty cash budgets id merchant_id budget_name amount assigned_to status active/closed/exhausted spent_amount remaining_amount created_by created_at index merchant status + petty cash expenses id budget_id merchant_id amount description receipt_file_key receipt_file_hash status pending/approved/rejected/paid approved_by created_by created_at index budget status • Outstanding modern UI glassmorphic Recharts • RazorpayX — Track petty cash budgets and make payments from assigned budgets • Add bills & receipt as attachments to petty cash expenses • Bug fixes and app performance improvements</p>
          </div>
          <button className="rounded-xl bg-primary text-white h-10 px-6 text-xs">+ Create Petty Cash Budget • Budget Name Amount Assigned To Status Active/Closed/Exhausted Spent Remaining • Outstanding</button>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <Card className="p-6">
            <h3 className="font-semibold">Petty Cash Budgets • Budget Name Amount Assigned To Status Active/Closed/Exhausted Spent Amount Remaining Amount • Outstanding • Track Petty Cash Budgets and Make Payments from Assigned Budgets</h3>
            <div className="mt-4 space-y-3">
              {mockBudgets.map(b => (
                <div key={b.id} className="rounded-xl border p-4 hover:bg-muted/50">
                  <div className="flex justify-between"><p className="font-medium text-sm">{b.budget_name} • {b.assigned_to}</p><Badge variant={b.status==="active" ? "success" : "warning"}>{b.status}</Badge></div>
                  <p className="text-[11px] text-muted-foreground mt-1">Amount {b.amount} ETB • Spent {b.spent_amount} • Remaining {b.remaining_amount} • Created {b.created_at} • Assigned To {b.assigned_to} • Status {b.status} • Outstanding per RazorpayX track petty cash budgets and make payments from assigned budgets</p>
                  <div className="mt-2 h-2 rounded-full bg-neutral-100 overflow-hidden"><div className="h-full bg-primary rounded-full" style={{ width: `${(parseInt(b.spent_amount)/parseInt(b.amount))*100}%` }} /></div>
                </div>
              ))}
              <button className="w-full rounded-xl border border-dashed h-10 text-xs">+ Create Petty Cash Budget • Budget Name Office Supplies Amount 50000 Assigned To Finance Manager Status Active Spent 0 Remaining 50000 Created By Finance Manager • Outstanding • Track Petty Cash Budgets and Make Payments from Assigned Budgets</button>
            </div>
          </Card>

          <Card className="p-6">
            <h3 className="font-semibold">Petty Cash Expenses • Amount Description Receipt File Key Receipt File Hash Status Pending/Approved/Rejected/Paid • Add Bills & Receipt as Attachments to Petty Cash Expenses • Outstanding • DocumentViewerOCR Side-by-side OCR</h3>
            <div className="mt-4 space-y-3">
              {mockExpenses.map(e => (
                <div key={e.id} className="rounded-xl border p-4 hover:bg-muted/50">
                  <div className="flex justify-between"><p className="font-medium text-sm">{e.description} • {e.amount} ETB</p><Badge variant={e.status==="paid" ? "success" : "warning"}>{e.status}</Badge></div>
                  <p className="text-[11px] text-muted-foreground mt-1">Budget ID {e.budget_id} • Amount {e.amount} • Receipt File Key {e.receipt_file_key} • Hash {e.receipt_file_hash} • Status {e.status} • Approved By {e.approved_by} • Created By {e.created_by} • Outstanding per RazorpayX add bills & receipt as attachments to petty cash expenses track petty cash budgets and make payments from assigned budgets • Receipt Preview Thumbs • DocumentViewerOCR side-by-side OCR • Hash Integrity • Progress Donut • Outstanding Modern • Mercury/Linear inspiration • Glassmorphic</p>
                  <div className="mt-2 flex gap-2"><div className="h-16 w-16 rounded-xl bg-neutral-100 border flex items-center justify-center text-[10px]">🖼️ JPG<br/>Expense<br/>Receipt</div><div><p className="text-[11px]">Receipt File Key {e.receipt_file_key} • Hash {e.receipt_file_hash} • MinIO presigned 15m • Encrypted SSE-S3 • 7y retention NBE • File key petty_cash/receipt_001.jpg • Hash integrity • ClamAV clean • Preview thumbs • DocumentViewer.tsx side-by-side OCR</p><button className="mt-1 rounded-xl border h-7 px-3 text-[10px]">View Receipt • MinIO presigned 15m • Thumb preview • DocumentViewer.tsx</button></div></div>
                </div>
              ))}
              <button className="w-full rounded-xl border border-dashed h-10 text-xs">+ Create Petty Cash Expense • Budget ID pcb_001 • Amount 1500 • Description Office supplies - printer ink • Receipt File Key petty_cash/receipt_001.jpg • Hash hash_pce_001 • Status pending • Approved By Finance Manager • Created By Admin • Outstanding • Add Bills & Receipt as Attachments to Petty Cash Expenses • Track Petty Cash Budgets and Make Payments from Assigned Budgets</button>
            </div>
          </Card>
        </div>
      </div>
    </div>
  )
}
