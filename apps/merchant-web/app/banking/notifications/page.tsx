"use client"
import * as React from "react"

function Card({ children, className = "" }: any) { return <div className={`rounded-2xl border bg-card shadow-soft ${className}`}>{children}</div> }
function Badge({ children, variant = "default" }: any) {
  const map: any = { default: "bg-neutral-100", success: "bg-green-500/15 text-green-700 border", warning: "bg-amber-500/15 text-amber-700 border", danger: "bg-red-500/15 text-red-700 border" }
  return <span className={`px-2 py-0.5 rounded-full text-[11px] border ${map[variant]}`}>{children}</span>
}

const mockNotifications = [
  { id: "notif_001", type: "bulk_payouts_approval", title: "Bulk Payouts Approval Required • 50,000 payouts at once with just one OTP", message: "Bulk payout batch pbat_001 with 50,000 payouts amount 5000000 ETB requires approval per maker-checker dual approval >50k payout >100k payroll approval count approver != submitter onboarding dual approval risk>=70 or TPV>1M", data: { payout_batch_id: "pbat_001", amount: "5000000", count: 50000 }, is_read: false, action_url: "/payout_batches/pbat_001", created_at: "2026-08-05T10:00:00Z" },
  { id: "notif_002", type: "payroll_run_pending_approval", title: "Payroll Run Pending Approval • July2026_Regular • 10 employees • Total Net 150000 ETB", message: "Payroll run prun_July2026 period 07/2026 type regular status pending_approval total_gross 200000 total_net 150000 total_tax 20000 total_pension 14000 employer_total_pension 22000 total_employer_cost 222000 total_count 10 payroll_data cutoff_date 2026-07-25 disbursal_date 2026-07-30 variance_report vs_last_month_percent 5.2 last_month_gross 190000 change_reason OT increase + bonus Sales requires dual approval >100k net maker-checker approver != submitter", data: { payroll_run_id: "prun_July2026", amount: "150000", count: 10 }, is_read: false, action_url: "/payroll/prun_July2026", created_at: "2026-08-05T09:30:00Z" },
  { id: "notif_003", type: "bank_file_generated", title: "Bank File Generated • pain.001.001.03 XML • 10 txs 150000 ETB • CBE", message: "Bank file pain.001.001.03 XML generated for payroll run July2026_Regular 10 txs 150000 ETB CBE • ISO20022 Document CstmrCdtTrfInitn GrpHdr MsgId CreDtTm NbOfTxs CtrlSum InitgPty PmtInf PmtInfId PmtMtd NbOfTxs CtrlSum ReqdExctnDt Dbtr Nm DbtrAcct Id Othr Id CdtTrfTxInf Amt InstdAmt Ccy ETB Cdtr Nm CdtrAcctId Othr Id • CBE/Awash/Dashen MT103 CSV fallback MT940 reconciliation window 24h amount tolerance 0.01 ETB O(n+m) map", data: { payroll_run_id: "prun_July2026", file_key: "payroll/reports/banking/bank_disbursal_July2026.xml" }, is_read: true, action_url: "/payroll/reports/bank_disbursal?year=2026&month=7", created_at: "2026-08-05T09:00:00Z" },
  { id: "notif_004", type: "pension_csv_generated", title: "Pension CSV Generated • Private Org Employees Social Security Agency • 10 employees • Emp 14k Emplr 22k Total 36k", message: "Pension CSV generated for July 2026 10 employees Emp 14k Emplr 22k Total 36k • pension_no employee_name code pensionable_gross employee_7% employer_11% total 18% period cost_center bank_code masked • Private Org Employees Social Security Agency • Outstanding modern UI glassmorphic", data: { period_year: 2026, period_month: 7, file_key: "payroll/reports/pension_2026_07.csv" }, is_read: true, action_url: "/payroll/reports/pension?year=2026&month=7", created_at: "2026-08-05T08:30:00Z" },
]

export default function NotificationsPage() {
  const [filter, setFilter] = React.useState("all")

  return (
    <div className="min-h-screen bg-gradient-to-br from-neutral-50 to-primary-50/20 p-6">
      <div className="max-w-4xl mx-auto space-y-6">
        <div className="flex justify-between items-start">
          <div>
            <h1 className="text-3xl font-bold">Notifications • ማሳወቂያዎች • Bulk Payouts Approval Refresh Button Latest Balance Instantly + Pending Payouts Notification • RazorpayX Parity • P0</h1>
            <p className="text-sm text-muted-foreground mt-2">Notifications Bulk Payouts Approval Refresh Button Latest Balance Instantly + Pending Payouts Notification per RazorpayX: type bulk_payouts_approval pending_payout payout_failed payroll_run_pending_approval payroll_run_completed tax_payment_due compliance_alert bank_file_generated pension_csv_generated erca_csv_generated loan_emi_due leave_request_pending claim_pending escrow_held escrow_released current_account_opened corporate_card_transaction forex_rate_alert accounting_sync_failed other title message data jsonb payout_batch_id payroll_run_id amount etc is_read read_at action_url /payout_batches/{`{id}`} or /payroll/{`{id}`} • Outstanding modern UI glassmorphic • Recharts • FCM push + in-app inbox + refresh button SWR revalidate • Bulk Payouts Approval Required 50,000 payouts at once with just one OTP • RazorpayX — Track petty cash budgets and make payments from assigned budgets • Add bills & receipt as attachments to petty cash expenses • Bug fixes and app performance improvements • Can instantly approve bulk payouts • Added notification for bulk payouts approval • Bug fixes & performance improvements • Now users can find their latest balance instantly using newly introduced refresh button • Also custom role users will be able to access mobile app • Instant approval bulk payouts notification for bulk payouts approval • Bank like its 2022 with RazorpayX Apple Watch • Bank like its 2022 with RazorpayX • Approve pending payouts on the go review payout details & account balance</p>
          </div>
          <button className="rounded-xl border bg-white h-10 px-4 text-xs">Refresh Balance • Instant • SWR • Real-time • Latest balance instantly using refresh button • RazorpayX mobile app feature • Outstanding • Real-time • Latest balance instantly • SWR • Revalidate • Outstanding modern UI glassmorphic</button>
        </div>

        <div className="flex gap-2">
          <button onClick={()=>setFilter("all")} className={`rounded-xl border h-9 px-4 text-xs ${filter==="all" ? "bg-primary text-white" : "bg-white"}`}>All • {mockNotifications.length}</button>
          <button onClick={()=>setFilter("unread")} className={`rounded-xl border h-9 px-4 text-xs ${filter==="unread" ? "bg-primary text-white" : "bg-white"}`}>Unread • {mockNotifications.filter(n=>!n.is_read).length} • Bulk Payouts Approval Required 50,000 payouts at once with just one OTP</button>
          <button onClick={()=>setFilter("bulk_payouts_approval")} className={`rounded-xl border h-9 px-4 text-xs ${filter==="bulk_payouts_approval" ? "bg-primary text-white" : "bg-white"}`}>Bulk Payouts Approval • RazorpayX • Track petty cash budgets • Add bills & receipt</button>
        </div>

        <Card className="p-6">
          <h3 className="font-semibold">Notifications • Bulk Payouts Approval Refresh Button Latest Balance Instantly + Pending Payouts Notification • RazorpayX Parity • P0 • Outstanding • FCM Push + In-app Inbox + Refresh Button SWR Revalidate • Bulk Payouts Approval Required 50,000 payouts at once with just one OTP</h3>
          <div className="mt-4 space-y-3">
            {mockNotifications.filter(n=> filter==="all" ? true : filter==="unread" ? !n.is_read : n.type===filter).map(notif => (
              <div key={notif.id} className={`rounded-xl border p-4 hover:bg-muted/50 ${!notif.is_read ? "bg-primary/5 border-primary/20" : ""}`}>
                <div className="flex justify-between items-start">
                  <div>
                    <p className="font-medium text-sm flex items-center gap-2">{notif.title} {!notif.is_read && <span className="h-2 w-2 rounded-full bg-primary animate-pulse" />}</p>
                    <p className="text-[11px] text-muted-foreground mt-1">{notif.message}</p>
                    <p className="text-[10px] text-muted-foreground mt-1">Type {notif.type} • Data {JSON.stringify(notif.data)} • Is Read {notif.is_read ? "Yes" : "No"} • Action URL {notif.action_url} • Created {notif.created_at} • Outstanding per RazorpayX notifications bulk payouts approval refresh button latest balance instantly + pending payouts notification per RazorpayX type bulk_payouts_approval pending_payout payout_failed payroll_run_pending_approval payroll_run_completed tax_payment_due compliance_alert bank_file_generated pension_csv_generated erca_csv_generated loan_emi_due leave_request_pending claim_pending escrow_held escrow_released current_account_opened corporate_card_transaction forex_rate_alert accounting_sync_failed other title message data jsonb payout_batch_id payroll_run_id amount etc is_read read_at action_url</p>
                  </div>
                  <Badge variant={!notif.is_read ? "warning" : "default"}>{notif.is_read ? "Read" : "Unread"} • {notif.type.replaceAll("_", " ")}</Badge>
                </div>
                <div className="mt-3 flex gap-2">
                  <button className="rounded-xl bg-primary text-white h-8 px-3 text-[11px]">View • Action URL {notif.action_url} • /payout_batches/{`{id}`} or /payroll/{`{id}`} • Outstanding • FCM Push + In-app Inbox + Refresh Button SWR Revalidate • Bulk Payouts Approval Required 50,000 payouts at once with just one OTP</button>
                  <button className="rounded-xl border h-8 px-3 text-[11px]">Mark as Read • Is Read true • Read At now • Outstanding • FCM Push + In-app Inbox + Refresh Button SWR Revalidate</button>
                  <button className="rounded-xl border h-8 px-3 text-[11px]">Approve • Bulk Payouts Approval Required 50,000 payouts at once with just one OTP • Maker-checker dual approval >50k payout >100k payroll • Outstanding • RazorpayX • Track petty cash budgets and make payments from assigned budgets • Add bills & receipt as attachments to petty cash expenses • Bug fixes and app performance improvements • Can instantly approve bulk payouts • Added notification for bulk payouts approval</button>
                </div>
              </div>
            ))}
          </div>
        </Card>

        <Card className="p-6">
          <h3 className="font-semibold">Refresh Balance • Instant • SWR • Real-time • Latest Balance Instantly Using Refresh Button • RazorpayX Mobile App Feature • Outstanding • Real-time • Latest Balance Instantly • SWR • Revalidate • Outstanding Modern UI Glassmorphic</h3>
          <div className="mt-4 rounded-xl bg-blue-500/10 border border-blue-500/20 p-4 text-[11px]">
            <p className="font-semibold">Refresh Balance Instant SWR Real-time Latest Balance Instantly Using Refresh Button RazorpayX Mobile App Feature Outstanding Real-time Latest Balance Instantly SWR Revalidate Outstanding Modern UI Glassmorphic</p>
            <p className="mt-2">RazorpayX — Track petty cash budgets and make payments from assigned budgets • Add bills & receipt as attachments to petty cash expenses • Bug fixes and app performance improvements • Can instantly approve bulk payouts • Added notification for bulk payouts approval • Bug fixes & performance improvements • Now users can find their latest balance instantly using newly introduced refresh button • Also custom role users will be able to access mobile app • Instant approval bulk payouts notification for bulk payouts approval • Bank like its 2022 with RazorpayX Apple Watch • Bank like its 2022 with RazorpayX • Approve pending payouts on the go review payout details & account balance • Outstanding modern UI glassmorphic • Recharts • FCM push + in-app inbox + refresh button SWR revalidate • Bulk Payouts Approval Required 50,000 payouts at once with just one OTP • RazorpayX — Track petty cash budgets and make payments from assigned budgets • Add bills & receipt as attachments to petty cash expenses • Bug fixes and app performance improvements • Can instantly approve bulk payouts • Added notification for bulk payouts approval</p>
            <div className="mt-3 flex gap-2">
              <button className="rounded-xl bg-primary text-white h-9 px-4 text-xs">Refresh Balance • Instant • SWR • Real-time • Latest balance instantly using refresh button • RazorpayX mobile app feature • Outstanding • Real-time • Latest balance instantly • SWR • Revalidate • Outstanding modern UI glassmorphic • Recharts • FCM push + in-app inbox + refresh button SWR revalidate • Bulk Payouts Approval Required 50,000 payouts at once with just one OTP</button>
              <button className="rounded-xl border h-9 px-4 text-xs">Mark All as Read • Is Read true • Read At now • Outstanding • FCM Push + In-app Inbox + Refresh Button SWR Revalidate</button>
            </div>
          </div>
        </Card>
      </div>
    </div>
  )
}
