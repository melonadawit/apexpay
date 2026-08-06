"use client"
import * as React from "react"

function Card({ children, className = "" }: any) { return <div className={`rounded-2xl border bg-card shadow-soft ${className}`}>{children}</div> }
function Badge({ children, variant = "default" }: any) {
  const map: any = { default: "bg-neutral-100", success: "bg-green-500/15 text-green-700 border", warning: "bg-amber-500/15 text-amber-700 border", danger: "bg-red-500/15 text-red-700 border" }
  return <span className={`px-2 py-0.5 rounded-full text-[11px] border ${map[variant]}`}>{children}</span>
}

const mockIntegrations = [
  { id: "int_001", provider: "tally", status: "connected", credentials_encrypted: "enc_tally_001", last_sync_at: "2026-08-05T10:00:00Z", last_sync_status: "success", last_sync_error: null, created_at: "2026-01-15" },
  { id: "int_002", provider: "zoho", status: "connected", credentials_encrypted: "enc_zoho_001", last_sync_at: "2026-08-05T09:00:00Z", last_sync_status: "success", last_sync_error: null, created_at: "2026-02-01" },
  { id: "int_003", provider: "quickbooks", status: "disconnected", credentials_encrypted: null, last_sync_at: "2026-07-20T10:00:00Z", last_sync_status: "failed", last_sync_error: "Invalid credentials - token expired", created_at: "2026-03-01" },
]

const mockSyncLogs = [
  { id: "log_001", integration_id: "int_001", sync_type: "payments", status: "success", payload: { count: 100 }, response: { synced: 100 }, records_synced: 100, error_message: null, created_at: "2026-08-05T10:00:00Z" },
  { id: "log_002", integration_id: "int_001", sync_type: "payouts", status: "success", payload: { count: 50 }, response: { synced: 50 }, records_synced: 50, error_message: null, created_at: "2026-08-05T09:30:00Z" },
  { id: "log_003", integration_id: "int_002", sync_type: "invoices", status: "failed", payload: { count: 20 }, response: { synced: 15 }, records_synced: 15, error_message: "5 invoices failed - GSTIN validation failed", created_at: "2026-08-05T09:00:00Z" },
]

export default function AccountingIntegrationsPage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-neutral-50 to-primary-50/20 p-6">
      <div className="max-w-7xl mx-auto space-y-6">
        <div className="flex justify-between items-start">
          <div>
            <h1 className="text-3xl font-bold">Accounting Integrations • የሂሳብ ውህደት • Two-way Sync Tally Zoho QuickBooks Xero Sage Other • CA Access Controls • Accounting Payouts • Create Razorpay Readable Payout Files • Payout Result Files Uploadable for Reconciliation • RazorpayX Parity • P0</h1>
            <p className="text-sm text-muted-foreground mt-2">Stress-free integrations with Tally, Zoho and QuickBooks to eliminate data entry & reconciliation, two-way sync between RazorpayX payments and accounting software, CA access controls, accounting payouts, create Razorpay readable payout files from accounting software, these files can be imported in dashboard and payouts can be generated, payout result files downloaded from dashboard can directly be uploaded to supported accounting software for reconciliation, you agree and authorise Razorpay to access import use or process your data from accounting software, Razorpay makes no warranties express implied or representations as to accuracy of details extracted, seamless integrations with Tally Zoho QuickBooks to eliminate data entry & reconciliation, generate financial reports in minutes and view real-time financial insights at a glance, real-time cash flow insights intuitive dashboard immediate visibility into cash flow and expenses, seamless accounting integrations easily integrate with popular accounting software like Zoho Books and Tally for efficient bookkeeping, accounting integrations with Zoho Books and Tally, accounting integrations with popular accounting software simplifying bookkeeping and reconciliation, with features like instant payouts and dedicated support it aims to enhance overall banking experience, manage receivables and payables in one place, get instant loans without collaterals, automate vendor and tax payments, get in-depth reporting into cash flow trends, automate your vendor and tax payments, get in-depth reporting into cash flow trends, real-time dashboard with payment success rates by method bank and device, cohort analysis for subscription businesses, revenue analytics, settlement tracking, and webhook-powered event streams for custom reporting, settlement reports need reconciliation against GSTR-1 data accounting for T+2 settlement cycles and gateway fees, section 194-O requires 1% TDS on e-commerce seller payments verify this in Form 26AS regularly, Tally Zoho QuickBooks integrations enable automated reconciliation but settlement timing differences need manual attention, generate financial reports in minutes, outstanding modern UI glassmorphic Recharts</p>
          </div>
          <button className="rounded-xl bg-primary text-white h-10 px-6 text-xs">+ Connect Accounting Integration • Tally Zoho QuickBooks Xero Sage Other • Two-way Sync • CA Access Controls • Outstanding</button>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          <Card className="p-6">
            <h3 className="font-semibold">Accounting Integrations • Provider Tally Zoho QuickBooks Xero Sage Other • Status Connected/Disconnected/Error/Pending • Credentials Encrypted AES-GCM • Last Sync At • Last Sync Status Success/Failed/Partial • Last Sync Error • Created By • Created At • Updated At • Unique Merchant Provider • Outstanding</h3>
            <div className="mt-4 space-y-3">
              {mockIntegrations.map(intg => (
                <div key={intg.id} className="rounded-xl border p-4 hover:bg-muted/50">
                  <div className="flex justify-between"><p className="font-medium text-sm">{intg.provider} • {intg.id}</p><Badge variant={intg.status==="connected" ? "success" : intg.status==="disconnected" ? "danger" : "warning"}>{intg.status}</Badge></div>
                  <p className="text-[11px] text-muted-foreground mt-1">Provider {intg.provider} • Status {intg.status} • Credentials Encrypted {intg.credentials_encrypted ? "enc_*" : "None"} • Last Sync At {intg.last_sync_at} • Last Sync Status {intg.last_sync_status} • Last Sync Error {intg.last_sync_error || "—"} • Created At {intg.created_at} • Outstanding per RazorpayX stress-free integrations with Tally Zoho QuickBooks to eliminate data entry & reconciliation two-way sync between RazorpayX payments and accounting software CA access controls accounting payouts create Razorpay readable payout files from accounting software these files can be imported in dashboard and payouts can be generated payout result files downloaded from dashboard can directly be uploaded to supported accounting software for reconciliation</p>
                  <div className="mt-2 flex gap-2">
                    <button className="rounded-xl bg-primary text-white h-7 px-3 text-[10px]">Sync Now • Two-way Sync Tally Zoho QuickBooks • Outstanding • Create Razorpay Readable Payout Files • Payout Result Files Uploadable for Reconciliation</button>
                    <button className="rounded-xl border h-7 px-3 text-[10px]">Disconnect • Credentials Encrypted AES-GCM • Last Sync At • Last Sync Status • Last Sync Error • Created By</button>
                  </div>
                </div>
              ))}
              <button className="w-full rounded-xl border border-dashed h-10 text-xs">+ Connect Accounting Integration • Provider Tally Zoho QuickBooks Xero Sage Other • Status Connected/Disconnected/Error/Pending • Credentials Encrypted AES-GCM via CONNECTOR_ENCRYPTION_KEY • Last Sync At • Last Sync Status Success/Failed/Partial • Last Sync Error • Created By • Outstanding • Two-way Sync Tally Zoho QuickBooks • CA Access Controls • Accounting Payouts • Create Razorpay Readable Payout Files • Payout Result Files Uploadable for Reconciliation • Outstanding modern UI glassmorphic</button>
            </div>
          </Card>

          <Card className="p-6 lg:col-span-2">
            <h3 className="font-semibold">Accounting Sync Logs • Integration ID • Sync Type Payments/Payouts/Invoices/Bills/Journal Entries/Contacts/Full Sync • Status Success/Failed/Partial/Pending • Payload Response Records Synced Error Message • Created At • Outstanding • Recharts • Financial Reports • Cash Flow Insights • Real-time Dashboard</h3>
            <div className="mt-4 rounded-xl border overflow-hidden">
              <div className="grid grid-cols-6 gap-2 bg-muted p-3 text-[11px] font-semibold"><span>Integration ID • Provider Tally Zoho QuickBooks</span><span>Sync Type • Payments/Payouts/Invoices/Bills/Journal Entries/Contacts/Full Sync</span><span>Status • Success/Failed/Partial/Pending • Payload Response Records Synced Error Message</span><span>Payload • Count 100 • Response • Synced 100 • Records Synced 100 • Error Message • Created At</span><span>Action • Sync Now • View Logs • Re-sync • Outstanding • Two-way Sync • Recharts • Financial Reports</span></div>
              {mockSyncLogs.map(log => (
                <div key={log.id} className="grid grid-cols-6 gap-2 p-3 border-t text-xs hover:bg-muted/50">
                  <span>{log.integration_id} • {log.id} • Sync Type {log.sync_type} • Status {log.status} • Payload {JSON.stringify(log.payload)} • Response {JSON.stringify(log.response)} • Records Synced {log.records_synced} • Error {log.error_message || "—"} • Created {log.created_at}</span>
                  <span>{log.sync_type} • Sync Type Payments/Payouts/Invoices/Bills/Journal Entries/Contacts/Full Sync • Status {log.status} • Payload {JSON.stringify(log.payload)} • Response {JSON.stringify(log.response)} • Records Synced {log.records_synced} • Error {log.error_message || "—"} • Created {log.created_at}</span>
                  <span><Badge variant={log.status==="success" ? "success" : log.status==="failed" ? "danger" : "warning"}>{log.status} • Payload Count {log.payload.count} • Response Synced {log.response.synced} • Records Synced {log.records_synced} • Error {log.error_message || "—"} • Created {log.created_at}</Badge></span>
                  <span className="text-[11px]">Payload {JSON.stringify(log.payload)} • Response {JSON.stringify(log.response)} • Records Synced {log.records_synced} • Error {log.error_message || "—"} • Created {log.created_at} • Outstanding per RazorpayX stress-free integrations with Tally Zoho QuickBooks to eliminate data entry & reconciliation two-way sync between RazorpayX payments and accounting software CA access controls accounting payouts create Razorpay readable payout files from accounting software these files can be imported in dashboard and payouts can be generated payout result files downloaded from dashboard can directly be uploaded to supported accounting software for reconciliation</span>
                  <span className="flex flex-col gap-1"><button className="rounded-xl bg-primary text-white h-7 px-3 text-[10px]">View Logs • Sync Type Payments/Payouts/Invoices/Bills/Journal Entries/Contacts/Full Sync • Status Success/Failed/Partial/Pending • Payload Response Records Synced Error Message • Created At • Outstanding • Recharts • Financial Reports • Cash Flow Insights • Real-time Dashboard • Cohort Analysis for Subscription Businesses Revenue Analytics Settlement Tracking Webhook-powered Event Streams for Custom Reporting</button><button className="rounded-xl border h-7 px-3 text-[10px]">Re-sync • Two-way Sync Tally Zoho QuickBooks • Outstanding • Create Razorpay Readable Payout Files • Payout Result Files Uploadable for Reconciliation</button></span>
                </div>
              ))}
            </div>

            <div className="mt-6 rounded-xl bg-blue-500/10 border border-blue-500/20 p-4 text-[11px]">
              <p className="font-semibold">Accounting Integrations Two-way Sync Tally Zoho QuickBooks CA Access Controls per RazorpayX: provider tally/zoho/quickbooks/xero/sage/other status connected/disconnected/error/pending credentials_encrypted via AES-GCM CONNECTOR_ENCRYPTION_KEY last_sync_at last_sync_status success/failed/partial last_sync_error created_by + accounting_sync_logs integration_id merchant_id sync_type payments/payouts/invoices/bills/journal_entries/contacts/full_sync status success/failed/partial/pending payload response records_synced error_message + outstanding modern UI glassmorphic Recharts • Financial Reports • Cash Flow Insights • Real-time Dashboard • Cohort Analysis for Subscription Businesses Revenue Analytics Settlement Tracking Webhook-powered Event Streams for Custom Reporting • Settlement Reports Need Reconciliation Against GSTR-1 Data Accounting for T+2 Settlement Cycles and Gateway Fees • Section 194-O Requires 1% TDS on E-commerce Seller Payments Verify in Form 26AS Regularly • Tally Zoho QuickBooks Integrations Enable Automated Reconciliation but Settlement Timing Differences Need Manual Attention • Outstanding modern UI glassmorphic • Recharts • Financial Reports • Cash Flow Insights • Real-time Dashboard • Cohort Analysis • Revenue Analytics • Settlement Tracking • Webhook-powered Event Streams</p>
            </div>
          </Card>
        </div>
      </div>
    </div>
  )
}
